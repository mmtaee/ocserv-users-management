<script lang="ts" setup>
import {dummyGroupList} from "@/utils/dummy.ts";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import type {ModelsOcservGroup} from "@/api";
import {useLocale} from "vuetify/framework";

const Create = defineAsyncComponent(() => import('@/components/ocserv_group/Create.vue'));

const {t} = useLocale()


const otherGroups = reactive<ModelsOcservGroup[]>([])
const createDialog = ref(false)
const editDialog = ref(false)
const deleteDialog = ref(false)


const getGroups = async () => {
  console.log("getGroups")
}

const complete = () => {
  createDialog.value = false
  console.log("complete")
  getGroups()
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
      <v-row align="center" dense justify="start">
        <v-col cols="12" md="12">
          <v-row>
            <v-col cols="12" md="11">
              <h3 class="text-capitalize">{{ t("OTHER") }} {{ t("GROUPS") }}</h3>
            </v-col>
            <v-col class="ma-0 pa-0" cols="12" md="auto">
              <v-btn
                  color="primary"
                  variant="outlined"
                  @click="createDialog = true"
              >
                {{ t("CREATE") }}
              </v-btn>
            </v-col>
          </v-row>
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
                xl="2"
            >
              <v-card elevation="6">
                <v-card-title class="text-subtitle-1 bg-info">
                  {{ item.name }}
                </v-card-title>
                <v-card-actions>
                  <v-spacer/>
                  <v-icon color="odd" size="small">mdi-pencil</v-icon>
                  <v-icon color="error" size="small">mdi-delete</v-icon>
                </v-card-actions>
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <Create
      v-if="createDialog"
      v-model="createDialog"
      @complete="complete"
  />
</template>
