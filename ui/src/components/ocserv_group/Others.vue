<script lang="ts" setup>
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {
  type ModelsOcservGroup,
  type ModelsOcservGroupConfig,
  type OcservGroupCreateOcservGroupData,
  OcservGroupsApi, type OcservGroupUpdateOcservGroupData
} from "@/api";
import {useLocale} from "vuetify/framework";
import {getAuthorization} from "@/utils/request.ts";

const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_group/CreateOrUpdate.vue'));
const ReusableDialog = defineAsyncComponent(() => import('@/components/reusable/ReusableDialog.vue'));


const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits(["update:modelValue"]);


const {t} = useLocale()

const api = new OcservGroupsApi()

const otherGroups = reactive<ModelsOcservGroup[]>([])

const editDialog = ref(false)
const deleteDialog = ref(false)

const createData = reactive<OcservGroupCreateOcservGroupData>({config: {}, name: ""})
const editData = reactive<OcservGroupUpdateOcservGroupData>({config: {}})
const selectedObj = ref<ModelsOcservGroup>({config: undefined, id: 0, name: ""})

const completeCreate = (data: ModelsOcservGroupConfig, groupName: string) => {
  createData.name = groupName
  createData.config = data

  api.ocservGroupsPost({
    ...getAuthorization(),
    request: createData
  }).then((res) => {
    console.log("resp: ", res.data)
    otherGroups.unshift(res.data)
    emit("update:modelValue", false)
  })
}

const deleteGroup = () => {
  if (selectedObj.value.id !== undefined) {
    api.ocservGroupsIdDelete({
      ...getAuthorization(),
      id: selectedObj.value.id,
    }).then(res => {
      console.log("resp: ", res.data)
      let index = otherGroups.findIndex(i => i.id === selectedObj.value.id)
      if (index > -1) {
        otherGroups.splice(index, 1)
      }
    }).finally(() => {
      deleteDialog.value = false
    })
  }
}

const completeEdit = (data: ModelsOcservGroupConfig) => {
  Object.assign(editData.config, data)
  if (selectedObj.value.id !== undefined) {
    api.ocservGroupsIdPatch({
      ...getAuthorization(),
      id: selectedObj.value.id,
      request: {
        config: editData.config
      },
    }).then((res) => {
      const index = otherGroups.findIndex(group => group.id === selectedObj.value.id)
      otherGroups.splice(index, 1, res.data)
      editDialog.value = false
    })
  }
}

const objHandler = (obj: ModelsOcservGroup) => {
  selectedObj.value = JSON.parse(JSON.stringify(obj))
}

onMounted(() => {
  api.ocservGroupsGet({
    ...getAuthorization(),
  }).then((res) => {
    Object.assign(otherGroups, res.data.result)
  })
})

</script>

<template>
  <v-card flat>
    <v-card-text>
      <v-row align="center" justify="center">
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
                      <v-icon color="info" @click="objHandler(item);editDialog=true">mdi-pencil</v-icon>
                    </v-col>

                    <v-col cols="12" md="1" sm="1">
                      <v-icon color="error" @click="objHandler(item);deleteDialog = true">mdi-delete</v-icon>
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

  <!-- Create Dialog -->
  <CreateOrEdit
      v-model="props.modelValue"
      @complete="completeCreate"
      @update:modelValue="val => $emit('update:modelValue', val)"
  />

  <!-- Edit Dialog -->
  <CreateOrEdit
      v-model="editDialog"
      :initValue="selectedObj"
      @complete="completeEdit"
  />

  <!-- Delete dialog -->
  <ReusableDialog
      v-model="deleteDialog"
      color="error"
      transition="dialog-top-transition"
      width="500"
  >
    <template #dialogTitle>
      <v-icon class="mb-1">mdi-delete</v-icon>
      <span class="text-capitalize">{{ t("DELETE_GROUP_TITLE") }} ({{ selectedObj.name }})</span>
    </template>

    <template #dialogText>
      {{ t("DELETE_OCSERV_GROUP_MESSAGE") }} <br/><br/>
      <v-icon class="mb-1 ma-0" color="error">
        mdi-bullhorn
      </v-icon>
      <span class="text-subtitle-2">{{ t("DELETE_GROUP_MESSAGE_SUB") }}</span>
    </template>

    <template #dialogAction>
      <v-btn
          color="black"
          variant="outlined"
          @click="deleteDialog = false"
      >
        {{ t("CANCEL") }}
      </v-btn>

      <v-btn
          color="error"
          variant="outlined"
          @click="deleteGroup"
      >
        {{ t("DELETE") }}
      </v-btn>
    </template>
  </ReusableDialog>

</template>
