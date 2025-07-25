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
import type {Meta} from "@/utils/interfaces.ts";

const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_group/CreateOrUpdate.vue'));
const ReusableDialog = defineAsyncComponent(() => import('@/components/reusable/ReusableDialog.vue'));
const ReusablePagination = defineAsyncComponent(() => import('@/components/reusable/ReusablePagination.vue'));


const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits(["update:modelValue"]);

const api = new OcservGroupsApi()

const {t} = useLocale()
const loading = ref(false)
const meta = reactive<Meta>({
  page: 1,
  size: 10,
  sort: "ASC",
  total_records: 0
})
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

const getGroups = () => {
  loading.value = true
  api.ocservGroupsGet({
    ...getAuthorization(),
    ...meta
  }).then((res) => {
    otherGroups.splice(0, otherGroups.length, ...(res.data.result ?? []))
    Object.assign(meta, res.data.meta)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(() => {
  getGroups()
})

</script>

<template>
  <v-card flat>
    <v-card-text>
      <v-data-iterator :items="otherGroups" :items-per-page="meta.size">
        <template v-slot:default="{ items }">
          <v-row align="center" justify="start">
            <v-col
                v-for="(group, i) in items"
                :key="i"
                cols="12"
                sm="6"
                xl="3"
            >
              <v-sheet border>
                <v-list-item
                    :title="group.raw.name"
                    class="bg-primary"
                    density="comfortable"
                    lines="two"
                >
                  <template v-slot:prepend>
                    <v-avatar color="grey-lighten-1">
                      <v-icon color="white">mdi-account-network</v-icon>
                    </v-avatar>
                  </template>

                  <template v-slot:append>
                    <v-menu>
                      <template v-slot:activator="{ props }">
                        <v-icon start v-bind="props">
                          mdi-dots-vertical
                        </v-icon>
                      </template>

                      <v-list color="info">
                        <v-list-item @click="objHandler(group.raw);editDialog = true">
                          <v-list-item-title class="text-info text-capitalize me-5">
                            {{ t("EDIT") }}
                          </v-list-item-title>
                          <template v-slot:prepend>
                            <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                          </template>
                        </v-list-item>

                        <v-list-item @click="objHandler(group.raw);deleteDialog=true">
                          <v-list-item-title class="text-error  text-capitalize me-5">
                            {{ t("DELETE") }}
                          </v-list-item-title>
                          <template v-slot:prepend>
                            <v-icon class="ms-2" color="error">mdi-delete</v-icon>
                          </template>
                        </v-list-item>

                      </v-list>
                    </v-menu>
                  </template>

                </v-list-item>

                <v-table class="text-caption text-capitalize" density="compact">
                  <tbody>
                  <tr style="text-align: right;">
                    <th>IPv4 Network:</th>
                    <td>
                      {{ group.raw?.config?.["ipv4-network"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>
                  <tr style="text-align: right;">
                    <th>Max Same Clients:</th>
                    <td>
                      {{ group.raw?.config?.["max-same-clients"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>
                  <tr style="text-align: right;">
                    <th>MTU:</th>
                    <td>
                      {{ group.raw?.config?.["mtu"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>

                  <tr style="text-align: right;">
                    <th>TX Data Per Sec:</th>
                    <td>
                      {{ group.raw?.config?.["tx-data-per-sec"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>

                  <tr style="text-align: right;">
                    <th>RX Data Per Sec:</th>
                    <td>
                      {{ group.raw?.config?.["rx-data-per-sec"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>

                  <tr style="text-align: right;">
                    <th>Keepalive:</th>
                    <td>
                      {{ group.raw?.config?.["keepalive"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>

                  <tr style="text-align: right;">
                    <th>Net Priority:</th>
                    <td>
                      {{ group.raw?.config?.["net-priority"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>

                  <tr style="text-align: right;">
                    <th>Session Timeout:</th>
                    <td>
                      {{ group.raw?.config?.["session-timeout"] ?? t("NOT_CONFIGURED") }}
                    </td>
                  </tr>
                  </tbody>
                </v-table>
              </v-sheet>
            </v-col>
          </v-row>
        </template>

        <template v-slot:footer="{}">
          <v-footer class="justify-space-between text-body-2 mt-4">
            <ReusablePagination
                v-model="meta"
                @update:modelValue="getGroups"
            />
          </v-footer>
        </template>

      </v-data-iterator>


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
