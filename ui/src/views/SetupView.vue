<script lang="ts" setup>
import logoUrl from "@/assets/ocserv.png"
import {useLocale} from "vuetify/framework";
import {reactive, ref} from "vue";
import {SystemApi, type SystemSetupSystem} from "@/api";
import {requiredRule} from "@/utils/rules.ts";
import {useConfigStore} from "@/stores/config.ts";
import {useUserStore} from "@/stores/user.ts";
import router from "@/plugins/router.ts";

const {t} = useLocale()

const showSetup = ref(false)
const valid = ref(true)
const loading = ref(false)
const data = reactive<SystemSetupSystem>({
  google_captcha_secret_key: "",
  google_captcha_site_key: "",
  password: "",
  username: ""
})

const rules = {
  required: (v: string) => requiredRule(v, t),
}

const setup = () => {
  loading.value = true
  const api = new SystemApi()
  api.systemSetupPost({
    request: data
  }).then((res) => {
    const configStore = useConfigStore()
    const userStore = useUserStore()

    configStore.setConfig(res.data.system.google_captcha_site_key)
    userStore.setUser({
      uid: res.data.user.uid,
      username: res.data.user.username,
      isAdmin: res.data.user.is_admin,
    })
    localStorage.setItem("token", res.data.token)
    router.push({name: 'HomePage'})
  }).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <v-container>
    <v-card>
      <v-row align="center" justify="center">

        <v-col class="mb-10" cols="12" md="5">
          <v-img :src="logoUrl" alt="ocserv logo"/>
        </v-col>

        <v-col v-if="!showSetup" cols="12" md="4">
          <div class="text-h6 mb-8">
            {{ t('T01') }} <span class="text-primary">Ocserv User Management Panel</span>.
          </div>

          <div class="text-subtitle-1 text-justify">
            {{ t('T02') }}
          </div>
          <div class="text-end mt-5">
            <v-btn color="primary" variant="outlined" @click="showSetup=true">Start</v-btn>
          </div>
        </v-col>

        <v-col v-else cols="12" md="4">
          <v-form v-model="valid">
            <v-row align="start" justify="start" no-gutters>
              <v-col md="12">
                <v-text-field
                    v-model="data.username"
                    :label="t('USERNAME')"
                    :rules="[rules.required]"
                    density="comfortable"
                    prepend-inner-icon="mdi-account"
                    variant="underlined"
                />
              </v-col>
              <v-col md="12">
                <v-text-field
                    v-model="data.password"
                    :label="t('PASSWORD')"
                    :rules="[rules.required]"
                    density="comfortable"
                    prepend-inner-icon="mdi-key"
                    variant="underlined"
                />
              </v-col>
              <v-col md="12">
                <v-text-field
                    v-model="data.google_captcha_site_key"
                    :label="t('GOOGLE_CAPTCHA_SITE_KEY')"
                    density="comfortable"
                    prepend-inner-icon="mdi-shield-key-outline"
                    variant="underlined"
                />
              </v-col>
              <v-col md="12">
                <v-text-field
                    v-model="data.google_captcha_secret_key"
                    :label="t('GOOGLE_CAPTCHA_SECRET_KEY')"
                    density="comfortable"
                    prepend-inner-icon="mdi-shield-key-outline"
                    variant="underlined"
                />
              </v-col>

              <v-col class="text-end mt-5" md="12">
                <v-btn
                    :disabled="!valid"
                    :loading="loading"
                    color="primary"
                    variant="outlined"
                    @click="setup"
                >
                  {{ t("SAVE") }} & {{ t("FINISH") }}
                </v-btn>
              </v-col>

            </v-row>

          </v-form>
        </v-col>

      </v-row>
    </v-card>
  </v-container>

</template>

