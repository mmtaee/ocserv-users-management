<script lang="ts" setup>
import {defineAsyncComponent, ref} from "vue";
import type {ModelsUser} from "@/api";
import {useLocale} from "vuetify/framework";
import {requiredRule} from "@/utils/rules.ts";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsUser
}>()

const emit = defineEmits(["update:modelValue", "save"])

const {t} = useLocale()

const rules = {
  required: (v: string) => requiredRule(v, t)
}
const password = ref("")
const valid = ref(true)
const showPassword = ref(false)
const passwordForm = ref()

const changePassword = () => {
  emit("save", props.user.uid, password.value)
  passwordForm.value?.reset()
}
</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      color="primary"
      transition="dialog-top-transition"
      width="500"
  >
    <template #dialogTitle>
      <v-icon class="mb-2">mdi-account-key</v-icon>
      {{ t("CHANGE_PASSWORD_TITLE") }}
    </template>

    <template #dialogText>
      <v-form ref="passwordForm" v-model="valid">
        <v-text-field
            v-model="password"
            :append-inner-icon="showPassword? 'mdi-eye-off' : 'mdi-eye'"
            :label="t('NEW_PASSWORD')"
            :rules="[rules.required]"
            :type="showPassword ? 'text': 'password'"
            density="comfortable"
            persistent-hint
            variant="underlined"
            @click:append-inner="showPassword = !showPassword"
        />
      </v-form>
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
          :disabled="!valid"
          color="primary"
          variant="flat"
          @click="changePassword"
      >
        {{ t("SAVE") }}
      </v-btn>
    </template>
  </ReusableDialog>
</template>
