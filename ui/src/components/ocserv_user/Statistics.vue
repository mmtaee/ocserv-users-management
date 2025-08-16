<script lang="ts" setup>
import {defineAsyncComponent, ref} from "vue";
import {type ModelsOcservUser, OcservUsersApi} from "@/api";
import {useI18n} from "vue-i18n";
import {getAuthorization} from "@/utils/request.ts";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))
const ReusableStatistics = defineAsyncComponent(() => import("@/components/reusable/ReusableStatistics.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsOcservUser
}>()

const emit = defineEmits(["update:modelValue"])

const {t} = useI18n()
const traffic = ref<any[]>([])

const search = (dateStart: string, dateEnd: string) => {
  const api = new OcservUsersApi()
  api.ocservUsersUidStatisticsGet({
    ...getAuthorization(),
    uid: props.user.uid,
    dateStart,
    dateEnd,
  }).then((res) => {
    traffic.value = res.data || []
  })
}

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      btnClose
      color="primary"
      fullscreen
      transition="dialog-bottom-transition"
      width="500"
      @update:modelValue="val => $emit('update:modelValue', val)"
  >
    <template #dialogTitle>
      <v-icon class="mb-1">mdi-delete</v-icon>
      <span class="text-capitalize">{{ t("OCSERV_USER_STATISTICS_TITLE") }} ({{ user.username }})</span>
    </template>

    <template #dialogText>
      <ReusableStatistics :traffic="traffic" immediate @search="search"/>
    </template>
  </ReusableDialog>
</template>

