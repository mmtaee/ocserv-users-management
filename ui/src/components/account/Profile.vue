<script lang="ts" setup>
import {useUserStore} from "@/stores/user.ts";
import {useLocale} from "vuetify/framework";
import {reactive, ref} from "vue";
import {requiredRule} from "@/utils/rules.ts";
import {formatDateTimeWithRelative} from "@/utils/convertors.ts";
import {type SystemChangeUserPasswordBySelf, SystemUsersApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";

const userStore = useUserStore()

const {t} = useLocale()
const validPassword = ref(true)
const passwordForm = ref()
const showChangePassword = ref(false)
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const data = reactive<SystemChangeUserPasswordBySelf>({
  old_password: "",
  new_password: "",
})

const rules = {
  required: (v: string) => requiredRule(v, t)
}

const changePassword = () => {
  const api = new SystemUsersApi()
  api.systemUsersPasswordPost({
    ...getAuthorization(),
    request: data,
  }).then(() => {
    showChangePassword.value = false
  })
}

const resetForm = () => {
  passwordForm.value?.reset()
}

</script>

<template>
  <v-card class="mx-auto pa-4" flat>
    <v-card-title class="text-capitalize">
      <v-icon class="me-2" start>mdi-account-outline</v-icon>
      {{ t("USER_PROFILE") }}
    </v-card-title>

    <v-card-text>

      <v-row align="start" justify="center">
        <v-col>
          <v-list>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">mdi-account-circle</v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">{{ t("USERNAME") }}</v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">
                {{ userStore.user?.username || 'guest' }}
              </v-list-item-subtitle>
            </v-list-item>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">mdi-fingerprint</v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">UID</v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">{{ userStore.user?.uid || '—' }}</v-list-item-subtitle>
            </v-list-item>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">
                  {{ userStore.user?.is_admin ? 'mdi-shield-account' : 'mdi-account-outline' }}
                </v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">{{ t("ROLE") }}</v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">
                {{ userStore.user?.is_admin ? t('ADMIN') : t('USER') }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-col>

        <v-divider opacity="1" vertical/>

        <v-col>
          <v-list>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">mdi-calendar</v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize  mb-2">{{ t('ACCOUNT_CREATED') }}
              </v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">
                {{ formatDateTimeWithRelative(userStore.user?.created_at, 'NOT_AVAILABLE') }}
              </v-list-item-subtitle>
            </v-list-item>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">mdi-clock-outline</v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">{{ t("LAST_LOGIN") }}</v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">
                {{ formatDateTimeWithRelative(userStore.user?.last_login, t('NOT_AVAILABLE')) }}
              </v-list-item-subtitle>
            </v-list-item>
            <v-list-item class="mb-3">
              <template #prepend>
                <v-icon size="large">mdi-update</v-icon>
              </template>
              <v-list-item-title class="text-subtitle-2 text-capitalize  mb-2">
                {{ t("LAST_UPDATED") }}
              </v-list-item-title>
              <v-list-item-subtitle class="text-subtitle-1">
                {{ formatDateTimeWithRelative(userStore.user?.updated_at, 'NOT_AVAILABLE') }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <v-divider class="mx-10" opacity="3"/>

  <v-card class="mx-auto pa-4 max-w-md" flat>
    <v-card-title class="text-capitalize">
      <v-icon class="me-2" start>mdi-shield-key-outline</v-icon>
      {{ t("PASSWORD") }}

      <v-btn
          v-if="!showChangePassword"
          color="primary"
          density="compact"
          variant="plain"
          @click="showChangePassword = true"
      >
        {{ t("CHANGE") }}
      </v-btn>
    </v-card-title>

    <v-card-text v-if="showChangePassword">
      <v-form ref="passwordForm" v-model="validPassword">
        <v-row align="center" class="mx-2" justify="start">
          <v-col md="3">
            <v-text-field
                v-model="data.old_password"
                :append-inner-icon="showOldPassword? 'mdi-eye-off' : 'mdi-eye'"
                :label="t('OLD_PASSWORD')"
                :rules="[rules.required]"
                :type="showOldPassword ? 'text':'password'"
                density="comfortable"
                prepend-inner-icon="mdi-key"
                variant="underlined"
                @click:append-inner="showOldPassword = !showOldPassword"
            />
          </v-col>

          <v-col md="3">
            <v-text-field
                v-model="data.new_password"
                :append-inner-icon="showNewPassword? 'mdi-eye-off' : 'mdi-eye'"
                :label="t('NEW_PASSWORD')"
                :rules="[rules.required]"
                :type="showNewPassword ? 'text':'password'"
                density="comfortable"
                prepend-inner-icon="mdi-key"
                variant="underlined"
                @click:append-inner="showNewPassword = !showNewPassword"
            />
          </v-col>

          <v-col md="3">
            <v-btn
                color="secondary"
                variant="outlined"
                @click="resetForm, showChangePassword=false"
            >
              {{ t("CANCEL") }}
            </v-btn>
            <v-btn
                :disabled="!data.old_password || !data.new_password"
                class="mx-2"
                color="primary"
                variant="flat"
                @click="changePassword"
            >
              {{ t("CHANGE_PASSWORD") }}
            </v-btn>
          </v-col>
        </v-row>

      </v-form>
    </v-card-text>
  </v-card>
</template>
