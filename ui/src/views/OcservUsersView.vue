<script lang="ts" setup>

import {defineAsyncComponent, onBeforeMount, onMounted, reactive, ref} from "vue";
import {useLocale} from "vuetify/framework";
import {
  type ModelsOcservUser,
  type ModelsOcservUserConfig,
  ModelsOcservUserTrafficTypeEnum, OcservGroupsApi,
  type OcservUserCreateOcservUserData,
  OcservUsersApi,
  type OcservUserUpdateOcservUserData
} from "@/api";
import {getAuthorization} from "@/utils/request.ts";


const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_user/CreateOrUpdate.vue'));
const Delete = defineAsyncComponent(() => import('@/components/ocserv_user/Delete.vue'));
const Disconnect = defineAsyncComponent(() => import('@/components/ocserv_user/Disconnect.vue'));
const Lock = defineAsyncComponent(() => import('@/components/ocserv_user/Lock.vue'));
const Unlock = defineAsyncComponent(() => import('@/components/ocserv_user/Unlock.vue'));
const Statistics = defineAsyncComponent(() => import('@/components/ocserv_user/Statistics.vue'));

const api = new OcservUsersApi()

const {t} = useLocale()
const createDialog = ref(false)
const users = reactive<ModelsOcservUser[]>([])
const editDialog = ref(false)
const deleteDialog = ref(false)
const lockDialog = ref(false)
const unlockDialog = ref(false)
const disconnectDialog = ref(false)
const statisticsDialog = ref(false)

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

onMounted(() => {
  api.ocservUsersGet({
    ...getAuthorization(),
  }).then((res) => {
    Object.assign(users, res.data.result)
  })
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
        <v-toolbar color="primary">
          <v-toolbar-title>
            {{ t('OCSERV_USERS') }}
          </v-toolbar-title>

          <template v-slot:append>
            <v-btn
                class="ma-5"
                color="white"
                variant="outlined"
                @click="createDialog = true"
            >
              {{ t("CREATE") }}
            </v-btn>
          </template>

        </v-toolbar>
        <v-card flat>
          <v-card-text>
            <v-row align="center" justify="center">
              <v-col class="mt-3" cols="12" md="12">
                <v-row align="center" justify="start">
                  <v-col
                      v-for="(item, index) in users"
                      :key="`ocserv-users-${index}`"
                      cols="12"
                      lg="4"
                      md="3"
                      sm="12"
                  >
                    <v-card class="mt-1" elevation="6">
                      <v-card-title class="text-subtitle-1 bg-primary pa-4">
                        <v-row align="center" justify="start">

                          <v-col cols="12" md="11" sm="8">
                            <v-icon v-if="item.is_locked" start>mdi-lock</v-icon>
                            <v-icon v-if="!item.is_locked" start>mdi-lock-open</v-icon>
                            {{ item.username }}
                            <span class="text-capitalize">({{ item.is_online ? t("ONLINE") : t("OFFLINE") }})</span>
                          </v-col>
                          <v-col cols="12" md="1">
                            <v-menu>
                              <template v-slot:activator="{ props }">
                                <v-icon v-bind="props">
                                  mdi-dots-vertical
                                </v-icon>
                              </template>

                              <v-list color="primary">

                                <v-list-item @click="objHandler(item);editDialog = true">
                                  <v-list-item-title class="text-info text-capitalize me-5">
                                    {{ t("EDIT") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                                  </template>
                                </v-list-item>

                                <v-list-item v-if="item.is_online" @click="objHandler(item);disconnectDialog = true">
                                  <v-list-item-title class="text-info text-capitalize me-5">
                                    {{ t("DISCONNECT") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="info">mdi-lan-disconnect</v-icon>
                                  </template>
                                </v-list-item>

                                <v-list-item v-if="!item.is_locked" @click="objHandler(item);lockDialog = true">
                                  <v-list-item-title class="text-odd text-capitalize me-5">
                                    {{ t("LOCK") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="odd">mdi-lock</v-icon>
                                  </template>
                                </v-list-item>

                                <v-list-item v-if="item.is_locked" @click="objHandler(item);unlockDialog = true">
                                  <v-list-item-title class="text-info text-capitalize me-5">
                                    {{ t("UNLOCK") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="info">mdi-lock-open</v-icon>
                                  </template>
                                </v-list-item>

                                <v-list-item @click="objHandler(item);statisticsDialog = true">
                                  <v-list-item-title class="text-success text-capitalize me-5">
                                    {{ t("STATISTICS") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="success">mdi-chart-bar-stacked</v-icon>
                                  </template>
                                </v-list-item>

                                <v-list-item @click="objHandler(item);deleteDialog=true">
                                  <v-list-item-title class="text-error  text-capitalize me-5">
                                    {{ t("DELETE") }}
                                  </v-list-item-title>
                                  <template v-slot:prepend>
                                    <v-icon class="ms-2" color="error">mdi-delete</v-icon>
                                  </template>
                                </v-list-item>

                              </v-list>
                            </v-menu>
                          </v-col>
                        </v-row>
                      </v-card-title>

                      <v-card-text class="pa-5">
                        <v-row align="center" justify="start">
                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              {{ t('GROUP') }}:
                            </span>
                            {{ item.group }}
                          </v-col>

                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              {{ t('PASSWORD') }}:
                            </span>

                            <span v-if="showPasswords[item.username]">
                              {{ item.password }}
                            </span>
                            <span v-else>
                              {{ '*'.repeat(item.password?.length || 0) }}
                            </span>
                            <v-icon
                                v-if="!showPasswords[item.username]"
                                class="ms-2"
                                color="grey"
                                icon="mdi-eye"
                                @click="togglePassword(item.username)"
                            />
                            <v-icon
                                v-else
                                class="ms-2"
                                color="grey"
                                icon="mdi-eye-off"
                                @click="togglePassword(item.username)"
                            />
                          </v-col>

                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              {{ t('TRAFFIC_TYPE') }}:
                            </span>
                            <span class="text-capitalize">{{ trafficTypesTransformer(item.traffic_type) }}</span>
                          </v-col>

                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              {{ t('TRAFFIC_SIZE') }}:
                            </span>
                            {{ item.traffic_size }} GB
                          </v-col>

                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              TX:
                            </span>
                            {{ Math.round((item.tx / (1024 ** 3)) * 1000) / 1000 }} GB
                          </v-col>

                          <v-col cols="12" md="6">
                            <span class="text-grey">
                              RX:
                            </span>
                            {{ Math.round((item.rx / (1024 ** 3)) * 1000) / 1000 }} GB
                          </v-col>

                        </v-row>
                      </v-card-text>
                    </v-card>
                  </v-col>
                </v-row>
              </v-col>
            </v-row>
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