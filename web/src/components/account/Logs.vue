<script lang="ts" setup>
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import {type AuditLogAuditLog, LogsUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import type {Meta} from "@/utils/interfaces.ts";
import {formatDateTimeWithRelative} from "@/utils/convertors.ts";
import {useI18n} from "vue-i18n";

const ReusablePagination = defineAsyncComponent(() => import("@/components/reusable/ReusablePagination.vue"))
const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const {t} = useI18n()
const userLogs = reactive<AuditLogAuditLog[]>([])
const meta = reactive<Meta>({
  page: 1,
  size: 25,
  sort: "DESC",
  total_records: 0
})

const changesResult = ref<any>({})
const showChangesDialog = ref(false)

const getLogs = () => {
  const api = new LogsUsersApi()
  api.logsUsersGet({
    ...getAuthorization(),
    ...meta
  }).then((res) => {
    userLogs.splice(0, userLogs.length, ...(res.data.result ?? []))
    Object.assign(meta, res.data.meta)
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

onMounted(() => {
  getLogs()
})
</script>

<template>
  <v-card class="mx-auto pa-4" flat>
    <v-card-text>
      <v-data-iterator :items="userLogs" :items-per-page="meta.size">
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
