<script lang="ts" setup>
import {type ModelsUser, SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import type {Meta} from "@/utils/interfaces.ts";
import {useLocale} from "vuetify/framework";
import {formatDateTimeWithRelative} from "@/utils/convertors.ts";
import router from "@/plugins/router.ts";

const ReusablePagination = defineAsyncComponent(() => import("@/components/reusable/ReusablePagination.vue"))
const Delete = defineAsyncComponent(() => import("@/components/account/Delete.vue"))
const ChangePassword = defineAsyncComponent(() => import("@/components/account/ChangePassword.vue"))

const {t} = useLocale()
const adminUsers = reactive<ModelsUser[]>([])
const meta = reactive<Meta>({
  page: 1,
  size: 25,
  sort: "ASC",
  total_records: 0
})

// const logDialog = ref(false)
const changePasswordDialog = ref(false)
const deleteDialog = ref(false)
const selectedObj = ref<ModelsUser>({is_admin: false, last_login: "", uid: "", username: ""})

const api = new SystemUsersApi()

const getAdmins = () => {
  const api = new SystemUsersApi()
  api.systemUsersGet({
    ...getAuthorization(),
    ...meta
  }).then(res => {
    adminUsers.splice(0, adminUsers.length, ...(res.data.result ?? []))
    Object.assign(meta, res.data.meta)
  })
}


const objHandler = (obj: ModelsUser) => {
  selectedObj.value = JSON.parse(JSON.stringify(obj))
}

const deleteUser = () => {
  api.systemUsersUidDelete({
    ...getAuthorization(),
    uid: selectedObj.value.uid,
  }).then(() => {
    let index = adminUsers.findIndex(i => i.uid === selectedObj.value.uid)
    if (index > -1) {
      adminUsers.splice(index, 1)
    }
  }).finally(() => {
    deleteDialog.value = false
  })
}

const changePassword = (uid: string, password: string) => {
  api.systemUsersUidPasswordPost({
    ...getAuthorization(),
    uid: uid,
    request: {
      password: password
    }
  }).then(() => {
  }).finally(() => {
    changePasswordDialog.value = false
  })
}


const showLogs = async (uid: string) => {
  await router.push({path: '/logs', query: {uid: uid}})
}

onMounted(() => {
  getAdmins()
})

</script>

<template>
  <v-row align="start" justify="center">
    <v-col>
      <v-card min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title>
            {{ t('ADMIN_USERS') }}
          </v-toolbar-title>
        </v-toolbar>
        <v-card flat min-height="850">
          <v-card-text>
            <v-data-iterator :items="adminUsers" :items-per-page="meta.size">
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
                            <v-icon color="white">mdi-account</v-icon>
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
                              <v-list-item @click="showLogs(user.raw.uid)">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("LOGS") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                                </template>
                              </v-list-item>

                              <v-list-item @click="objHandler(user.raw);changePasswordDialog=true">
                                <v-list-item-title class="text-info text-capitalize me-5">
                                  {{ t("CHANGE_PASSWORD") }}
                                </v-list-item-title>
                                <template v-slot:prepend>
                                  <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
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
                          <th>UID:</th>
                          <td>
                            {{ user.raw.uid }}
                          </td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("CREATED_AT") }}:</th>
                          <td>
                            {{ formatDateTimeWithRelative(user.raw.created_at, "") }}
                          </td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("UPDATED_AT") }}:</th>
                          <td>
                            {{ formatDateTimeWithRelative(user.raw.updated_at, "") }}
                          </td>
                        </tr>
                        <tr style="text-align: right;">
                          <th>{{ t("LAST_LOGIN") }}:</th>
                          <td>
                            {{ formatDateTimeWithRelative(user.raw.last_login, t("NO_LOGIN_YET")) }}
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
                      @update:modelValue="getAdmins"
                  />
                </v-footer>
              </template>
            </v-data-iterator>

          </v-card-text>
        </v-card>

        <!-- Delete Dialog -->
        <Delete
            v-model="deleteDialog"
            :user="selectedObj"
            @done="deleteUser"
        />

        <!--  ChangePassword Dialog -->
        <ChangePassword
            v-model="changePasswordDialog"
            :user="selectedObj"
            @save="changePassword"
        />
      </v-card>
    </v-col>
  </v-row>
</template>
