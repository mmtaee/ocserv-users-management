<script lang="ts" setup>
import {defineAsyncComponent, onBeforeMount, onMounted, reactive, ref} from "vue";
import {useLocale} from "vuetify/framework";
import {
  type ModelsOcservUser,
  type ModelsOcservUserConfig,
  ModelsOcservUserTrafficTypeEnum,
  OcservGroupsApi,
  type OcservUserCreateOcservUserData,
  OcservUsersApi,
  type OcservUserUpdateOcservUserData
} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import type {Meta} from "@/utils/interfaces.ts";
import {formatDateTimeWithRelative} from "../utils/convertors.ts";

const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_user/CreateOrUpdate.vue'));
const Delete = defineAsyncComponent(() => import('@/components/ocserv_user/Delete.vue'));
const Disconnect = defineAsyncComponent(() => import('@/components/ocserv_user/Disconnect.vue'));
const Lock = defineAsyncComponent(() => import('@/components/ocserv_user/Lock.vue'));
const Unlock = defineAsyncComponent(() => import('@/components/ocserv_user/Unlock.vue'));
const Statistics = defineAsyncComponent(() => import('@/components/ocserv_user/Statistics.vue'));
const ReusablePagination = defineAsyncComponent(() => import("@/components/reusable/ReusablePagination.vue"))

const api = new OcservUsersApi()

const {t} = useLocale()
const createDialog = ref(false)
const users = reactive<ModelsOcservUser[]>([])
const meta = reactive<Meta>({
  page: 1,
  size: 10,
  sort: "ASC",
  total_records: 0
})
const editDialog = ref(false)
const deleteDialog = ref(false)
const lockDialog = ref(false)
const unlockDialog = ref(false)
const disconnectDialog = ref(false)
const statisticsDialog = ref(false)

const loading = ref(false)

const showPasswords = reactive<Record<string, boolean>>({});

const createData = reactive<OcservUserCreateOcservUserData>({
  group: "",
  password: "",
  traffic_size: 0,
  traffic_type: ModelsOcservUserTrafficTypeEnum.TOTALLY_TRANSMIT,
  username: ""
})
const editData = reactive<OcservUserUpdateOcservUserData>({})
const selectedObj = ref<ModelsOcservUser>({
  created_at: "",
  group: "",
  is_locked: false,
  is_online: false,
  password: "",
  rx: 0,
  traffic_size: 10,
  traffic_type: ModelsOcservUserTrafficTypeEnum.TOTALLY_TRANSMIT,
  tx: 0,
  uid: "",
  username: ""
})

const groups = ref<string[]>([])

const trafficTypesTransformer = (item: ModelsOcservUserTrafficTypeEnum): string => {
  switch (item) {
    case ModelsOcservUserTrafficTypeEnum.FREE:
      return t('FREE');
    case ModelsOcservUserTrafficTypeEnum.MONTHLY_TRANSMIT:
      return t('MONTHLY_TRANSMIT');
    case ModelsOcservUserTrafficTypeEnum.MONTHLY_RECEIVE:
      return t('MONTHLY_RECEIVE');
    case ModelsOcservUserTrafficTypeEnum.TOTALLY_RECEIVE:
      return t('TOTALLY_RECEIVE');
    case ModelsOcservUserTrafficTypeEnum.TOTALLY_TRANSMIT:
      return t('TOTALLY_TRANSMIT');
    default:
      return item;
  }
}


const togglePassword = (username: string) => {
  showPasswords[username] = !showPasswords[username]
}

const completeCreate = (data: ModelsOcservUser, config: ModelsOcservUserConfig) => {
  Object.assign(createData, data)
  createData.config = config
  api.ocservUsersPost({
    ...getAuthorization(),
    request: createData,
  }).then((res) => {
    users.unshift(res.data)
    createDialog.value = false
  })
}

const completeEdit = (data: ModelsOcservUser, config: ModelsOcservUserConfig) => {
  Object.assign(editData, data)
  editData.config = config

  api.ocservUsersUidPatch({
    ...getAuthorization(),
    uid: selectedObj.value.uid,
    request: editData,
  }).then((res) => {
    let index = users.findIndex(i => i.uid = selectedObj.value.uid)
    if (index > -1) {
      users.splice(index, 1, res.data)
    }
    editDialog.value = false
  })
}

const deleteUser = () => {
  api.ocservUsersUidDelete({
    ...getAuthorization(),
    uid: selectedObj.value.uid,
  }).then(() => {
    let index = users.findIndex(i => i.uid = selectedObj.value.uid)
    if (index > -1) {
      users.splice(index, 1)
    }
  }).finally(() => {
    deleteDialog.value = false
  })
}


const disconnectUser = () => {
  api.ocservUsersUsernameDisconnectPost({
    ...getAuthorization(),
    username: selectedObj.value.username,
  }).then(() => {
    let index = users.findIndex(i => i.uid = selectedObj.value.uid)
    if (index > -1) {
      users[index].is_online = false
    }
  }).finally(() => {
    disconnectDialog.value = false
  })
}


const lockUser = () => {
  api.ocservUsersUidLockPost({
    ...getAuthorization(),
    uid: selectedObj.value.uid,
  }).then(() => {
    let index = users.findIndex(i => i.uid = selectedObj.value.uid)
    if (index > -1) {
      users[index].is_locked = true
    }
  }).finally(() => {
    lockDialog.value = false
  })
}


const unlockUser = () => {
  api.ocservUsersUidUnlockPost({
    ...getAuthorization(),
    uid: selectedObj.value.uid,
  }).then(() => {
    let index = users.findIndex(i => i.uid = selectedObj.value.uid)
    if (index > -1) {
      users[index].is_locked = false
    }
  }).finally(() => {
    unlockDialog.value = false
  })
}

const objHandler = (obj: ModelsOcservUser) => {
  selectedObj.value = JSON.parse(JSON.stringify(obj))
}


const getUsers = () => {
  loading.value = true
  api.ocservUsersGet({
    ...getAuthorization(),
    ...meta,
  }).then((res) => {
    users.splice(0, users.length, ...(res.data.result ?? []))
    Object.assign(meta, res.data.meta)
  }).finally(() => {
    loading.value = false
  })
}

onMounted(() => {
  getUsers()
})

onBeforeMount(
    () => {
      const api = new OcservGroupsApi()
      api.ocservGroupsLookupGet({
        ...getAuthorization()
      }).then((res) => {
        groups.value = res.data
      })
    }
)


</script>

<template>
  <v-row>
    <v-col>
      <v-card min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title>
            {{ t('OCSERV_USERS') }}
          </v-toolbar-title>

          <template v-slot:append>
            <v-btn
                class="ma-5"
                color="primary"
                variant="elevated"
                @click="createDialog = true"
            >
              {{ t("CREATE") }}
            </v-btn>
          </template>

        </v-toolbar>
        <v-card flat>
          <v-card-text>

            <v-data-iterator :items="users" :items-per-page="meta.size">
              <template v-slot:default="{ items }">
                <v-row align="center" justify="start">
                  <v-col
                      v-for="(user, i) in items"
                      :key="i"
                      cols="12"
                      sm="6"
                      xl="3"
                  >
                    <v-sheet border>
                      <v-list-item
                          :title="user.raw.username"
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
                              <v-list-item @click="objHandler(user.raw);editDialog = true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("EDIT") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item v-if="user.raw.is_online"
                                           @click="objHandler(user.raw);disconnectDialog = true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("DISCONNECT") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-lan-disconnect</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item v-if="!user.raw.is_locked" @click="objHandler(user.raw);lockDialog = true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("LOCK") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-lock</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item v-if="user.raw.is_locked" @click="objHandler(user.raw);unlockDialog = true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("UNLOCK") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-lock-open</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item @click="objHandler(user.raw);statisticsDialog = true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("STATISTICS") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-chart-bar-stacked</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item @click="objHandler(user.raw);deleteDialog=true">
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
                          <th>{{ t("STATUS") }}:</th>
                          <td>
                            <v-tooltip v-if="!user.raw.is_locked" text="Unlocked">
                              <template #activator="{ props }">
                                <v-icon color="success" end v-bind="props">
                                  mdi-lock-open
                                </v-icon>
                              </template>
                            </v-tooltip>

                            <v-tooltip v-if="user.raw.is_locked" text="Locked">
                              <template #activator="{ props }">
                                <v-icon color="error" end v-bind="props">
                                  mdi-lock
                                </v-icon>
                              </template>
                            </v-tooltip>

                            <v-tooltip v-if="user.raw.is_online" text="Online">
                              <template #activator="{ props }">
                                <v-icon color="success" end v-bind="props">
                                  mdi-lan-connect
                                </v-icon>
                              </template>
                            </v-tooltip>

                            <v-tooltip v-if="!user.raw.is_online" text="Offline">
                              <template #activator="{ props }">
                                <v-icon color="error" end v-bind="props">
                                  mdi-lan-disconnect
                                </v-icon>
                              </template>
                            </v-tooltip>
                          </td>
                        </tr>

                        <tr style="text-align: right;">
                          <th>{{ t("GROUP") }}:</th>
                          <td>{{ user.raw.group }}</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("PASSWORD") }}:</th>
                          <td>
                            <span v-if="showPasswords[user.raw.username]">
                              {{ user.raw.password }}
                            </span>
                            <span v-else>
                              {{ '*'.repeat(user.raw.password?.length || 0) }}
                            </span>
                            <v-icon
                                v-if="!showPasswords[user.raw.username]"
                                class="ms-2 mb-2"
                                color="grey"
                                icon="mdi-eye"
                                @click="togglePassword(user.raw.username)"
                            />
                            <v-icon
                                v-else
                                class="ms-2 mb-2"
                                color="grey"
                                icon="mdi-eye-off"
                                @click="togglePassword(user.raw.username)"
                            />
                          </td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("TRAFFIC_TYPE") }}:</th>
                          <td>{{ trafficTypesTransformer(user.raw.traffic_type) }}</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("TRAFFIC_SIZE") }}:</th>
                          <td>{{ user.raw.traffic_size }} GB</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>RX:</th>
                          <td>{{ Math.round((user.raw.rx / (1024 ** 3)) * 1000) / 1000 }} GB</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>TX:</th>
                          <td>{{ Math.round((user.raw.tx / (1024 ** 3)) * 1000) / 1000 }} GB</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("DEACTIVATED_AT") }}:</th>
                          <td>{{ formatDateTimeWithRelative(user.raw.deactivated_at, t("USER_IS_ACTIVE")) }}</td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("EXPIRE_AT") }}:</th>
                          <td>{{ formatDateTimeWithRelative(user.raw.expire_at, t("USER_NO_EXPIRE_TIME_SET")) }}</td>
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
                      @update:modelValue="getUsers"
                  />
                </v-footer>
              </template>

            </v-data-iterator>

          </v-card-text>
        </v-card>
      </v-card>
    </v-col>
  </v-row>

  <!-- Create Dialog -->
  <CreateOrEdit
      v-model="createDialog"
      :groups="groups"
      @complete="completeCreate"
  />

  <!-- Edit Dialog -->
  <CreateOrEdit
      v-model="editDialog"
      :groups="groups"
      :initValue="selectedObj"
      @complete="completeEdit"
  />

  <!-- Delete Dialog -->
  <Delete
      v-model="deleteDialog"
      :user="selectedObj"
      @done="deleteUser"
  />

  <!-- Disconnect Dialog -->
  <Disconnect
      v-model="disconnectDialog"
      :user="selectedObj"
      @done="disconnectUser"
  />

  <!-- Lock Dialog -->
  <Lock
      v-model="lockDialog"
      :user="selectedObj"
      @done="lockUser"
  />

  <!-- Unlock Dialog -->
  <Unlock
      v-model="unlockDialog"
      :user="selectedObj"
      @done="unlockUser"
  />

  <!-- Statistics Dialog -->
  <Statistics
      v-model="statisticsDialog"
      :user="selectedObj"
  />

</template>
`