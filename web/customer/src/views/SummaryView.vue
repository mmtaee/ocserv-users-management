<script setup lang="ts">
import { onMounted, shallowRef } from "vue";
import { useI18n } from "vue-i18n";
import type { SummaryResponse } from "@/api/generated";
import { normalizeApiError } from "@/api/http";
import { getSummary } from "@/api/services/customers";
import StatusMessage from "@/components/StatusMessage.vue";

const { t } = useI18n({ useScope: "global" });
const data = shallowRef<SummaryResponse | null>(null);
const loading = shallowRef(true);
const error = shallowRef("");
async function load(): Promise<void> {
  loading.value = true;
  error.value = "";
  try {
    data.value = await getSummary();
  } catch (cause) {
    error.value = normalizeApiError(cause).message;
  } finally {
    loading.value = false;
  }
}
onMounted(load);
</script>

<template>
  <section class="flex flex-col gap-4">
    <h1 class="text-2xl font-semibold">{{ t("summary.title") }}</h1>
    <StatusMessage :error="error" />
    <p v-if="loading">{{ t("common.loading") }}</p>
    <div v-else-if="data" class="grid gap-4 md:grid-cols-2">
      <article class="card grid grid-cols-2 gap-3">
        <span class="text-muted-foreground">{{ t("summary.username") }}</span
        ><strong>{{ data.ocserv_user?.username }}</strong>
        <span class="text-muted-foreground">{{ t("summary.owner") }}</span
        ><span>{{ data.ocserv_user?.owner ?? "—" }}</span>
        <span class="text-muted-foreground">{{ t("summary.expiry") }}</span
        ><span>{{ data.ocserv_user?.expire_at ?? "—" }}</span>
        <span class="text-muted-foreground">{{ t("summary.trafficPlan") }}</span
        ><span>{{ data.ocserv_user?.traffic_type ?? "—" }}</span>
        <span class="text-muted-foreground">{{ t("summary.allowance") }}</span
        ><span>{{ data.ocserv_user?.traffic_size ?? "—" }}</span>
        <span class="text-muted-foreground">{{
          data.ocserv_user?.is_locked
            ? t("summary.locked")
            : t("summary.active")
        }}</span>
      </article>
      <article class="card grid grid-cols-2 gap-3">
        <span class="text-muted-foreground">{{ t("summary.received") }}</span
        ><strong>{{ data.usage?.bandwidths?.rx ?? 0 }} GiB</strong>
        <span class="text-muted-foreground">{{ t("summary.transmitted") }}</span
        ><strong>{{ data.usage?.bandwidths?.tx ?? 0 }} GiB</strong>
      </article>
    </div>
  </section>
</template>
