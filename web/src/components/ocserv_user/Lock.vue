<script lang="ts" setup>
import {defineAsyncComponent} from "vue";
import type {ModelsOcservUser} from "@/api";
import {useI18n} from "vue-i18n";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))


const props = defineProps<{
  modelValue: boolean
  user: ModelsOcservUser
}>()

const emit = defineEmits(["update:modelValue", "done"])

const {t} = useI18n()

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      color="odd"
      transition="dialog-top-transition"
      width="500"
  >
    <template #dialogTitle>
      <v-icon class="mb-1" start>mdi-lock</v-icon>
      <span class="text-capitalize">{{ t("LOCK_OCSERV_USER_TITLE") }} ({{ user.username }})</span>
    </template>

    <template #dialogText>
      {{ t("LOCK_OCSERV_USER_MESSAGE") }} <br/><br/>
    </template>

    <template #dialogAction>
      <v-btn
          color="black"
          variant="outlined"
          @click="emit('update:modelValue', false)"
      >
        {{ t("CANCEL") }}
      </v-btn>

      <v-btn
          color="primary"
          variant="outlined"
          @click="emit('done', user.uid)"
      >
        {{ t("LOCK") }}
      </v-btn>
    </template>

  </ReusableDialog>

</template>
