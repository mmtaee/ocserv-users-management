<script lang="ts" setup>
import {defineAsyncComponent} from "vue";
import type {ModelsOcservUser} from "@/api";
import {useLocale} from "vuetify/framework";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsOcservUser
}>()

const emit = defineEmits(["update:modelValue", "done"])

const {t} = useLocale()

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      color="error"
      transition="dialog-top-transition"
      width="500"
  >
    <template #dialogTitle>
      <v-icon class="mb-1">mdi-delete</v-icon>
      <span class="text-capitalize">{{ t("DELETE_OCSERV_USER_TITLE") }} ({{ user.username }})</span>
    </template>

    <template #dialogText>
      {{ t("DELETE_OCSERV_USER_MESSAGE") }} <br/><br/>
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
          color="error"
          variant="outlined"
          @click="emit('done', user.uid)"
      >
        {{ t("DELETE") }}
      </v-btn>
    </template>
  </ReusableDialog>
</template>
