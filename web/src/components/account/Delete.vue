<script lang="ts" setup>
import {defineAsyncComponent} from "vue";
import type {ModelsUser} from "@/api";
import {useI18n} from "vue-i18n";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsUser
}>()

const emit = defineEmits(["update:modelValue", "done"])

const {t} = useI18n()

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      color="white"
      transition="dialog-top-transition"
      width="500"
  >
    <template #dialogTitle>
      <span class="text-capitalize">{{ t("DELETE_ADMIN_USER_TITLE") }} ({{ user.username }})</span>
    </template>

    <template #dialogText>
      {{ t("DELETE_ADMIN_USER_MESSAGE") }} <br/><br/>
    </template>

    <template #dialogAction>
      <v-btn
          color="secondary"
          variant="outlined"
          @click="emit('update:modelValue', false)"
      >
        {{ t("CANCEL") }}
      </v-btn>

      <v-btn
          color="error"
          variant="flat"
          @click="emit('done', user.uid)"
      >
        {{ t("DELETE") }}
      </v-btn>
    </template>
  </ReusableDialog>
</template>
