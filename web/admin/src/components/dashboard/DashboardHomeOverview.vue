<script setup lang="ts">
import type { DeepReadonly } from "vue";
import { useI18n } from "vue-i18n";

import type { DashboardOverview } from "@/api/services/dashboard";
import TelegramStatus from "@/components/dashboard/TelegramStatus.vue";
import { telegramBotEnabled } from "@/config/features";
import { Alert, AlertDescription } from "@/components/ui/alert";

const props = defineProps<{
  overview: DeepReadonly<DashboardOverview> | null;
  loading: boolean;
  error: boolean;
}>();

const { t } = useI18n({ useScope: "global" });
</script>

<template>
  <section class="flex flex-col gap-6">
    <Alert v-if="error" variant="error">
      <AlertDescription>{{
        t("dashboard.homeOverviewError")
      }}</AlertDescription>
    </Alert>
    <TelegramStatus
      :available="telegramBotEnabled"
      :service="props.overview?.telegram_service ?? null"
      :loading="telegramBotEnabled && loading"
    />
  </section>
</template>
