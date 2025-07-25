<script lang="ts" setup>
import {type ModelsUser, SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {defineAsyncComponent, onMounted, reactive} from "vue";
import type {Meta} from "@/utils/interfaces.ts";
import {useLocale} from "vuetify/framework";
import {formatDateTimeWithRelative} from "@/utils/convertors.ts";

const ReusablePagination = defineAsyncComponent(() => import("@/components/reusable/ReusablePagination.vue"))

const {t} = useLocale()
const adminUsers = reactive<ModelsUser[]>([])
const meta = reactive<Meta>({
  page: 1,
  size: 25,
  sort: "ASC",
  total_records: 0
})

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

onMounted(() => {
  getAdmins()
})

</script>

<template>
  <v-card flat>
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
                        <v-list-item @click="">
                          <v-list-item-title class="text-info text-capitalize me-5">
                            {{ t("LOGS") }}
                          </v-list-item-title>
                          <template v-slot:prepend>
                            <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                          </template>
                        </v-list-item>

                        <v-list-item @click="">
                          <v-list-item-title class="text-info text-capitalize me-5">
                            {{ t("CHANGE_PASSWORD") }}
                          </v-list-item-title>
                          <template v-slot:prepend>
                            <v-icon class="ms-2" color="info">mdi-pencil</v-icon>
                          </template>
                        </v-list-item>

                        <v-list-item @click="">
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
</template>

<style scoped>

</style>