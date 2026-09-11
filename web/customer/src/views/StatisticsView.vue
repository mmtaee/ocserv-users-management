<script setup lang="ts">
import { onMounted, reactive, shallowRef } from "vue";
import { useI18n } from "vue-i18n";
import type { Bandwidth, DailyTraffic } from "@/api/generated";
import { normalizeApiError } from "@/api/http";
import { getBandwidth, getStats } from "@/api/services/customers";
import StatusMessage from "@/components/StatusMessage.vue";

const { t } = useI18n({ useScope: "global" });
const query = reactive({ date_start: "", date_end: "" });
const stats = shallowRef<DailyTraffic[]>([]);
const total = shallowRef<Bandwidth | null>(null);
const loading = shallowRef(false);
const error = shallowRef("");
async function load(): Promise<void> {
  loading.value = true;
  error.value = "";
  const params = {
    date_start: query.date_start || undefined,
    date_end: query.date_end || undefined,
  };
  try {
    [stats.value, total.value] = await Promise.all([
      getStats(params),
      getBandwidth(params),
    ]);
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
    <h1 class="text-2xl font-semibold">{{ t("statistics.title") }}</h1>
    <form class="card grid gap-3 sm:grid-cols-3" @submit.prevent="load">
      <label class="flex flex-col gap-1"
        ><span class="label">{{ t("activity.dateStart") }}</span
        ><input v-model="query.date_start" class="input" type="date"
      /></label>
      <label class="flex flex-col gap-1"
        ><span class="label">{{ t("activity.dateEnd") }}</span
        ><input v-model="query.date_end" class="input" type="date"
      /></label>
      <button class="button self-end" :disabled="loading">
        {{ t("common.apply") }}
      </button>
    </form>
    <StatusMessage :error="error" />
    <article v-if="total" class="card">
      <h2 class="font-semibold">{{ t("statistics.total") }}</h2>
      <p>
        {{ t("statistics.rx") }}: {{ total.rx }} · {{ t("statistics.tx") }}:
        {{ total.tx }}
      </p>
    </article>
    <p v-if="loading">{{ t("common.loading") }}</p>
    <p v-else-if="!stats.length" class="card">{{ t("common.empty") }}</p>
    <div v-else class="card overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t("statistics.date") }}</th>
            <th>{{ t("statistics.rx") }}</th>
            <th>{{ t("statistics.tx") }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in stats" :key="item.date">
            <td>{{ item.date }}</td>
            <td>{{ item.rx ?? 0 }}</td>
            <td>{{ item.tx ?? 0 }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
