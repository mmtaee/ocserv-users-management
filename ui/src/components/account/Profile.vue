<script lang="ts" setup>
import {useUserStore} from "@/stores/user.ts";
import {useLocale} from "vuetify/framework";
import {ref} from "vue";
import {requiredRule} from "@/utils/rules.ts";

const userStore = useUserStore()

const {t} = useLocale()
const validPassword = ref(true)
const data = ref({
  old_password: "",
  new_password: "",
})

const rules = {
  required: (v: string) => requiredRule(v, t)
}

const changePassword = () => {
//   TODO: complete this
}


const passwordForm = ref()

const resetForm = () => {
  passwordForm.value?.reset()
}

</script>

<template>
  <v-card class="mx-auto my-6 pa-4 max-w-md" flat>
    <v-card-title>
      <v-icon class="me-2" start>mdi-account</v-icon>
      User Profile
    </v-card-title>

    <v-card-text>
      <v-list density="compact">
        <v-list-item>
          <template #prepend>
            <v-icon>mdi-account-circle</v-icon>
          </template>
          <v-list-item-title>Username</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.username || 'guest' }}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-fingerprint</v-icon>
          </template>
          <v-list-item-title>UID</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.uid || '—' }}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon>{{ userStore.user?.is_admin ? 'mdi-shield-account' : 'mdi-account-outline' }}</v-icon>
          </template>
          <v-list-item-title>Role</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.is_admin ? 'Admin' : 'User' }}</v-list-item-subtitle>
        </v-list-item>

        <v-divider class="my-3" opacity="4"/>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-clock-outline</v-icon>
          </template>
          <v-list-item-title>Last Login</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.last_login || 'Not available' }}</v-list-item-subtitle>
        </v-list-item>


        <v-list-item>
          <template #prepend>
            <v-icon>mdi-calendar</v-icon>
          </template>
          <v-list-item-title>Account Created</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.created_at || 'Not available' }}</v-list-item-subtitle>
        </v-list-item>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-update</v-icon>
          </template>
          <v-list-item-title>Last Updated</v-list-item-title>
          <v-list-item-subtitle>{{ userStore.user?.updated_at || 'Not available' }}</v-list-item-subtitle>
        </v-list-item>

        <v-divider class="my-3" opacity="4"/>

        <v-list-item>
          <template #prepend>
            <v-icon>mdi-password</v-icon>
          </template>
          <v-list-item-title>Change Password</v-list-item-title>
          <v-card>
            <v-card-text>
              <v-form ref="passwordForm" v-model="validPassword">
                <v-row align="start" justify="end">
                  <v-col cols="12" md="6">
                    <v-text-field
                        v-model="data.old_password"
                        :label="t('OLD_PASSWORD')"
                        :rules="[rules.required]"
                        density="comfortable"
                        variant="underlined"
                    />
                  </v-col>

                  <v-col cols="12" md="6">
                    <v-text-field
                        v-model="data.new_password"
                        :label="t('NEW_PASSWORD')"
                        :rules="[rules.required]"
                        density="comfortable"
                        variant="underlined"
                    />
                  </v-col>
                </v-row>

              </v-form>

            </v-card-text>
            <v-card-actions>
              <v-spacer/>
              <v-btn
                  color="secondary"
                  variant="outlined"
                  @click="resetForm"
              >
                {{ t("CANCEL") }}
              </v-btn>
              <v-btn
                  :disabled="!data.old_password || !data.new_password"
                  class="mx-2"
                  color="primary"
                  variant="outlined"
                  @click="changePassword"
              >
                {{ t("CHANGE_PASSWORD") }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<style scoped>

</style>