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

            <v-data-iterator :items="users">
              <template v-slot:header="{ page, pageCount, prevPage, nextPage }">
                <div class="text-truncate">
                  Most popular mice
                </div>

                <div class="d-flex align-center">
                  <v-btn class="me-8" variant="text">
                    <span class="text-decoration-underline text-none">See all</span>
                  </v-btn>

                  <div class="d-inline-flex">
                    <v-btn
                        :disabled="page === 1"
                        class="me-2"
                        icon="mdi-arrow-left"
                        size="small"
                        variant="tonal"
                        @click="prevPage"
                    ></v-btn>

                    <v-btn
                        :disabled="page === pageCount"
                        icon="mdi-arrow-right"
                        size="small"
                        variant="tonal"
                        @click="nextPage"
                    ></v-btn>
                  </div>
                </div>
              </template>

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
                            <v-icon v-if="!user.raw.is_locked" color="success" end size="x-small">
                              mdi-lock-open
                            </v-icon>
                            <v-icon v-if="user.raw.is_locked" color="error" end size="x-small">
                              mdi-lock
                            </v-icon>
                            <v-icon v-if="user.raw.is_online" color="success" end size="x-small">
                              mdi-lan-connect
                            </v-icon>
                            <v-icon v-if="!user.raw.is_online" color="error" end size="x-small">
                              mdi-lan-disconnect
                            </v-icon>
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
                        </tbody>
                      </v-table>
                    </v-sheet>
                  </v-col>
                </v-row>
              </template>

              <!--              <template v-slot:footer="{ page, pageCount }">-->
              <!--                <v-footer-->
              <!--                    class="justify-space-between text-body-2 mt-4"-->
              <!--                    color="surface-variant"-->
              <!--                >-->
              <!--                  Total mice: {{ users.length }}-->

              <!--                  <div>-->
              <!--                    Page {{ page }} of {{ pageCount }}-->
              <!--                  </div>-->
              <!--                </v-footer>-->
              <!--              </template>-->

            </v-data-iterator>

            <!--            <v-row align="center" justify="center">-->
            <!--              <v-col class="mt-3" cols="12" md="12">-->
            <!--                <v-row align="center" justify="start">-->
            <!--                  <v-col-->
            <!--                      v-for="(item, index) in users"-->
            <!--                      :key="`ocserv-users-${index}`"-->
            <!--                      cols="12"-->
            <!--                      lg="4"-->
            <!--                      md="3"-->
            <!--                      sm="12"-->
            <!--                  >-->
            <!--                    <v-card class="mt-1" elevation="6">-->
            <!--                      <v-card-title class="text-subtitle-1 bg-primary py-4">-->
            <!--                        <v-row align="center" justify="center">-->

            <!--                                      <v-col cols="12" md="10">-->
            <!--                                        {{ item.username }}-->
            <!--                                        <span class="text-capitalize">({{ item.is_online ? t("ONLINE") : t("OFFLINE") }})</span>-->
            <!--                                      </v-col>-->

            <!--                          <v-col cols="12" md="1">-->
            <!--                                        <v-icon v-if="item.is_locked" end>mdi-lock</v-icon>-->
            <!--                                        <v-icon v-if="!item.is_locked" end>mdi-lock-open</v-icon>-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="1">-->
            <!--                                        <v-menu>-->
            <!--                                          <template v-slot:activator="{ props }">-->
            <!--                                            <v-icon start v-bind="props">-->
            <!--                                              mdi-dots-vertical-->
            <!--                                            </v-icon>-->
            <!--                                          </template>-->

            <!--                                          <v-list color="info">-->
            <!--                                            <v-list-item @click="objHandler(item);editDialog = true">-->
            <!--                                              <v-list-item-title class="text-info text-capitalize me-5">-->
            <!--                                                {{ t("EDIT") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="info">mdi-pencil</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                            <v-list-item v-if="item.is_online" @click="objHandler(item);disconnectDialog = true">-->
            <!--                                              <v-list-item-title class="text-info text-capitalize me-5">-->
            <!--                                                {{ t("DISCONNECT") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="info">mdi-lan-disconnect</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                            <v-list-item v-if="!item.is_locked" @click="objHandler(item);lockDialog = true">-->
            <!--                                              <v-list-item-title class="text-info text-capitalize me-5">-->
            <!--                                                {{ t("LOCK") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="info">mdi-lock</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                            <v-list-item v-if="item.is_locked" @click="objHandler(item);unlockDialog = true">-->
            <!--                                              <v-list-item-title class="text-info text-capitalize me-5">-->
            <!--                                                {{ t("UNLOCK") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="info">mdi-lock-open</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                            <v-list-item @click="objHandler(item);statisticsDialog = true">-->
            <!--                                              <v-list-item-title class="text-info text-capitalize me-5">-->
            <!--                                                {{ t("STATISTICS") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="info">mdi-chart-bar-stacked</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                            <v-list-item @click="objHandler(item);deleteDialog=true">-->
            <!--                                              <v-list-item-title class="text-error  text-capitalize me-5">-->
            <!--                                                {{ t("DELETE") }}-->
            <!--                                              </v-list-item-title>-->
            <!--                                              <template v-slot:prepend>-->
            <!--                                                <v-icon class="ms-2" color="error">mdi-delete</v-icon>-->
            <!--                                              </template>-->
            <!--                                            </v-list-item>-->

            <!--                                          </v-list>-->
            <!--                                        </v-menu>-->
            <!--                          </v-col>-->
            <!--                        </v-row>-->
            <!--                      </v-card-title>-->

            <!--                      <v-card-text class="pa-5">-->
            <!--                        <v-row align="center" justify="start">-->
            <!--                          <v-col cols="12" md="6">-->
            <!--                            <span class="text-grey text-capitalize">-->
            <!--                              {{ t('GROUP') }}:-->
            <!--                            </span>-->
            <!--                            {{ item.group }}-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="6">-->
            <!--                                        <span class="text-grey text-capitalize">-->
            <!--                                          {{ t('PASSWORD') }}:-->
            <!--                                        </span>-->

            <!--                                        <span v-if="showPasswords[item.username]">-->
            <!--                                          {{ item.password }}-->
            <!--                                        </span>-->
            <!--                                        <span v-else>-->
            <!--                                          {{ '*'.repeat(item.password?.length || 0) }}-->
            <!--                                        </span>-->
            <!--                                        <v-icon-->
            <!--                                            v-if="!showPasswords[item.username]"-->
            <!--                                            class="ms-2"-->
            <!--                                            color="grey"-->
            <!--                                            icon="mdi-eye"-->
            <!--                                            @click="togglePassword(item.username)"-->
            <!--                                        />-->
            <!--                                        <v-icon-->
            <!--                                            v-else-->
            <!--                                            class="ms-2"-->
            <!--                                            color="grey"-->
            <!--                                            icon="mdi-eye-off"-->
            <!--                                            @click="togglePassword(item.username)"-->
            <!--                                        />-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="6">-->
            <!--                            <span class="text-grey text-capitalize">-->
            <!--                              {{ t('TRAFFIC_TYPE') }}:-->
            <!--                            </span>-->
            <!--                            <span class="text-capitalize">{{ trafficTypesTransformer(item.traffic_type) }}</span>-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="6">-->
            <!--                            <span class="text-grey text-capitalize">-->
            <!--                              {{ t('TRAFFIC_SIZE') }}:-->
            <!--                            </span>-->
            <!--                            {{ item.traffic_size }} GB-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="6">-->
            <!--                            <span class="text-grey text-capitalize">-->
            <!--                              TX:-->
            <!--                            </span>-->
            <!--                            {{ Math.round((item.tx / (1024 ** 3)) * 1000) / 1000 }} GB-->
            <!--                          </v-col>-->

            <!--                          <v-col cols="12" md="6">-->
            <!--                            <span class="text-grey">-->
            <!--                              RX:-->
            <!--                            </span>-->
            <!--                            {{ Math.round((item.rx / (1024 ** 3)) * 1000) / 1000 }} GB-->
            <!--                          </v-col>-->

            <!--                        </v-row>-->
            <!--                      </v-card-text>-->
            <!--                    </v-card>-->
            <!--                  </v-col>-->
            <!--                </v-row>-->
            <!--              </v-col>-->
            <!--            </v-row>-->
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