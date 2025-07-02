<script lang="ts" setup>
import {dummyGroupList} from "@/utils/dummy.ts";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import type {ModelsOcservGroup, ModelsOcservGroupConfig} from "@/api";
import {useLocale} from "vuetify/framework";

const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_group/CreateOrUpdate.vue'));


const {t} = useLocale()


const otherGroups = reactive<ModelsOcservGroup[]>([])
const createDialog = ref(false)
const editDialog = ref(false)
const deleteDialog = ref(false)

const createData = reactive({})

const defaultGroup: ModelsOcservGroup = {name: "", id: 0}
const editObject = ref<ModelsOcservGroup>({config: undefined, id: 0, name: ""})
const editData = reactive({})


const getGroups = async () => {
  console.log("getGroups")
}

const complete = (data: ModelsOcservGroupConfig, groupName: string) => {
  createDialog.value = false
  console.log("groupName: ", groupName)
  console.log("complete: ", data)
  otherGroups.push({
    id: 99,
    name: groupName,
    config: data
  })
}

const completeEdit = (data: ModelsOcservGroupConfig) => {
  const index = otherGroups.findIndex(group => group.id === editObject.value.id)
  if (index !== -1) {
    otherGroups.splice(index, 1, editObject.value)
  }
  editDialog.value = false
}

const editHandler = (obj: ModelsOcservGroup) => {
  editObject.value = JSON.parse(JSON.stringify(obj))
  editDialog.value = true
}

const deleteGroup = (id: number) => {

}

onMounted(() => {
  // Fetch and assign group list
  // TODO: call group api list
  getGroups()
  Object.assign(otherGroups, dummyGroupList)
})


</script>

<template>
  <v-card flat>
    <v-card-text>
      <v-row align="center" justify="center">
        <v-col cols="12" lg="10" md="10" sm="10">
          <span class="text-capitalize text-subtitle-1">{{ t("OTHER") }} {{ t("GROUPS") }}</span>
        </v-col>
        <v-col cols="12" lg="auto" md="2" sm="2">
          <v-btn
              color="primary"
              variant="outlined"
              @click="createDialog = true"
          >
            {{ t("CREATE") }}
          </v-btn>
        </v-col>

        <v-divider/>
        <v-col class="mt-3" cols="12" md="12">
          <v-row>
            <v-col
                v-for="(item, index) in otherGroups"
                :key="`other-groups-${index}`"
                cols="12"
                lg="3"
                md="4"
                sm="6"
            >
              <v-card class="py-3" elevation="6">
                <v-card-title class="text-subtitle-1">
                  <v-row align="center" justify="center">
                    <v-col cols="12" md="9" sm="8">{{ item.name }}</v-col>

                    <v-col cols="12" md="1" sm="1">
                      <v-icon color="info" @click="editHandler(item)">mdi-pencil</v-icon>
                    </v-col>

                    <v-col cols="12" md="1" sm="1">
                      <v-icon color="error">mdi-delete</v-icon>
                    </v-col>
                  </v-row>
                </v-card-title>
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <CreateOrEdit
      v-model="createDialog"
      @complete="complete"
  />

  <CreateOrEdit
      v-model="editDialog"
      :initValue="editObject"
      @complete="completeEdit"
  />


</template>
