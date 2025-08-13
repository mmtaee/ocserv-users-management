<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {type AuditLogAuditLog, LogsApi, type ModelsUsersLookup, SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import type {Meta} from "@/utils/interfaces.ts";
import {formatDateTimeWithRelative} from "@/utils/convertors.ts";
import {useUserStore} from "@/stores/user.ts";
import {useRoute} from "vue-router";

const ReusablePagination = defineAsyncComponent(() => import("@/components/reusable/ReusablePagination.vue"))
const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const {t} = useLocale()
const route = useRoute()
const loading = ref(false)
const userLogs = reactive<AuditLogAuditLog[]>([])
const usersLookup = ref<ModelsUsersLookup[]>([])
const userSelected = ref("")
const meta = reactive<Meta>({
  page: 1,
  size: 25,
  sort: "DESC",
  total_records: 0
})

const changesResult = ref<any>({})
const showChangesDialog = ref(false)

const getLogs = () => {
  loading.value = true
  const api = new LogsApi()
  api.logsAuditGet({
    ...getAuthorization(),
    ...meta,
    uid: userSelected.value,
  }).then((res) => {
    userLogs.splice(0, userLogs.length, ...(res.data.result ?? []))
    Object.assign(meta, res.data.meta)
  }).finally(() => {
    loading.value = false
  })
}

const showChanges = (changes: string | undefined) => {
  if (changes && changes !== "null" && changes.length > 0) {
    changesResult.value = JSON.parse(changes) as any
    if (typeof changesResult.value != "object" && Object.keys(changesResult.value).length > 0) {
      changesResult.value = {
        result: changesResult.value
      }
    }
    showChangesDialog.value = true
  }
}

const getUsersLookup = () => {
  const api = new SystemUsersApi()
  api.systemUsersLookupGet({
    ...getAuthorization()
  }).then((res) => {
    usersLookup.value.splice(0, usersLookup.value.length, ...(res.data || []))
  })
}

const resetSearch = () => {
  userSelected.value = ''
  getLogs()
}


onMounted(() => {
  let uid = route.query.uid
  if (uid && typeof uid === 'string' && uid !== '') {
    userSelected.value = uid
  } else {
    const userStore = useUserStore()
    userSelected.value = userStore.uid
  }
  getUsersLookup()
  getLogs()
})


</script>

<template>
  <v-card class="mx-auto pa-4" flat>
    <v-card-title>
      <v-row align="start" justify="start">
        <v-col cols="12" lg="3" md="4">
          <v-autocomplete
              v-model="userSelected"
              :items="usersLookup"
              density="compact"
              item-title="username"
              item-value="uid"
              variant="underlined"
          >
          </v-autocomplete>
        </v-col>
        <v-col cols="12" md="auto">
          <v-btn
              :disabled="userSelected == ''"
              class="mt-3"
              color="primary"
              size="small"
              variant="flat"
              @click="getLogs"
          >
            {{ t("SEARCH") }}
          </v-btn>
          <v-btn
              :disabled="userSelected == ''"
              class="mt-3 ms-2"
              color="secondary"
              size="small"
              variant="outlined"
              @click="resetSearch"
          >
            {{ t("CLEAR") }}
          </v-btn>
        </v-col>
      </v-row>
    </v-card-title>

    <v-divider class="mb-2"/>

    <v-card-text>
      <v-data-iterator :items="userLogs" :items-per-page="meta.size" :loading="loading">
        <template v-slot:default="{ items }">
          <v-row align="center" justify="start">
            <v-col
                v-for="(log, i) in items"
                :key="i"
                cols="12"
                sm="6"
                xl="3"
            >
              <v-sheet border>
                <v-table class="text-caption text-capitalize" density="compact">
                  <tbody>
                  <tr class="text-capitalize" style="text-align: right;">
                    <th>{{ t("USERNAME") }}:</th>
                    <td>
                      {{ log.raw.username }}
                    </td>
                  </tr>
                  <tr class="text-capitalize" style="text-align: right;">
                    <th>{{ t("ACTION") }}:</th>
                    <td>
                      {{ JSON.parse(log.raw.action ?? "{}").action }}
                    </td>
                  </tr>
                  <tr class="text-capitalize" style="text-align: right;">
                    <th>{{ t("REASON") }}:</th>
                    <td>
                      {{ JSON.parse(log.raw.action ?? "{}").reason }}
                    </td>
                  </tr>
                  <tr style="text-align: right;">
                    <th>{{ t("CREATED_AT") }}:</th>
                    <td>
                      {{ formatDateTimeWithRelative(log.raw.created_at, "") }}
                    </td>
                  </tr>
                  <tr style="text-align: right;">
                    <th>{{ t("CHANGES") }}:</th>
                    <td>
                      <v-icon @click="showChanges(log.raw.changes)">mdi-eye</v-icon>
                    </td>
                  </tr>
                  </tbody>
                </v-table>
              </v-sheet>
            </v-col>
          </v-row>
        </template>

        <template v-slot:footer="{}">
          <v-footer v-if="userLogs.length == meta.size" class="justify-space-between text-body-2 mt-4">
            <ReusablePagination
                v-model="meta"
                @update:modelValue="getLogs"
            />
          </v-footer>
        </template>
      </v-data-iterator>
    </v-card-text>
  </v-card>

  <ReusableDialog
      v-model="showChangesDialog"
      color="white"
      transition="dialog-top-transition"
      width="auto"
  >

    <template #dialogText>
      <v-card elevation="4" variant="tonal">
        <v-sheet
            class="text-pre-wrap px-2"
            style="overflow: auto; max-height: 630px; white-space: pre;"
        >
              <pre style="margin: 0;">
                <br/>{{ changesResult }}
              </pre>
        </v-sheet>
      </v-card>
    </template>

    <template #dialogAction>
      <v-btn
          color="secondary"
          variant="outlined"
          @click="showChangesDialog=false"
      >
        {{ t("CLOSE") }}
      </v-btn>

    </template>
  </ReusableDialog>

</template>
