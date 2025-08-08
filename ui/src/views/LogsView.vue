<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import {defineAsyncComponent, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from "vue-router";


const Ocserv = defineAsyncComponent(() => import("@/components/system_logs/Ocserv.vue"));
const AuditLogs = defineAsyncComponent(() => import("@/components/system_logs/AuditLogs.vue"));


const {t} = useLocale()
const route = useRoute()
const router = useRouter()

const tab = ref('audit_logs')

const createDialog = ref(false)

onMounted(() => {
  tab.value = route.query.tab?.toString() || 'audit_logs'
})
watch(tab, (newVal) => {
  router.replace({
    query: {
      ...route.query,
      tab: newVal
    }
  })
})
</script>

<template>
  <v-row>
    <v-col>
      <v-card min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title>
            {{ t('OCSERV_GROUPS') }}
          </v-toolbar-title>

          <template v-slot:append>
            <v-btn
                v-if="tab == 'others'"
                class="ma-5"
                color="primary"
                variant="elevated"
                @click="createDialog = true"
            >
              {{ t("CREATE") }}
            </v-btn>
          </template>
        </v-toolbar>

        <v-tabs
            v-model="tab"
            align-tabs="center"
            color="primary"
            direction="horizontal"
        >
          <v-tab :text="t('AUDIT_LOGS')" value="audit_logs"></v-tab>
          <v-tab :text="t('OCSERV')" value="defaults"></v-tab>
        </v-tabs>

        <v-tabs-window v-model="tab">

          <v-tabs-window-item value="audit_logs">
            <AuditLogs/>
          </v-tabs-window-item>

          <v-tabs-window-item value="ocserv">
            <Ocserv/>
          </v-tabs-window-item>

        </v-tabs-window>
      </v-card>
    </v-col>
  </v-row>

</template>

<style scoped>

</style>