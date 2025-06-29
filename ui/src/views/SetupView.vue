<script lang="ts" setup>
import setupImage from "@/assets/setup2.png"
import {useLocale} from "vuetify/framework";
import {reactive, ref} from "vue";
import {SystemApi, type SystemSetupSystem} from "@/api";
import {requiredRule} from "@/utils/rules.ts";
import {useConfigStore} from "@/stores/config.ts";
import {useUserStore} from "@/stores/user.ts";
import router from "@/plugins/router.ts";

const {t} = useLocale()

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
  <div class="w-75">
    <v-row align="center" class="border" justify="center" style="border-radius: 10px">

      <v-col class="px-15" cols="12" md="7">
        <div class="text-h5 mb-2 mt-5">
          {{ t('T01') }} <span class="text-primary">Ocserv User Management Panel</span>.
        </div>

        <div class="text-subtitle-2 text-justify mb-2 text-justify">
          <v-icon class="mb-1 ma-0" color="primary">
            mdi-bullhorn
          </v-icon>
          <strong>{{ t("HINT") }}</strong>:
          {{ t('T02') }}
          <a
              class="text-info"
              href="https://www.google.com/recaptcha/admin/create"
              style="text-decoration: none"
          >
            {{ t("HERE") }}
          </a>.
        </div>

        <v-divider class="mt-8 mb-5"/>

        <v-form v-model="valid">
          <v-row align="start" justify="center" no-gutters>
            <v-col class="ma-0 pa-0 me-5" cols="12" lg="5" md="6" sm="12">
              <v-text-field
                  v-model="data.username"
                  :label="t('ADMIN_USERNAME')"
                  :rules="[rules.required]"
                  density="default"
                  prepend-inner-icon="mdi-account"
                  variant="underlined"
              />
            </v-col>

            <v-col class="ma-0 pa-0 me-5" cols="12" lg="5" md="6" sm="6">
              <v-text-field
                  v-model="data.password"
                  :label="t('ADMIN_PASSWORD')"
                  :rules="[rules.required]"
                  density="default"
                  prepend-inner-icon="mdi-key"
                  variant="underlined"
              />
            </v-col>

            <v-col class="ma-0 pa-0 me-5" cols="12" lg="5" md="6" sm="6">
              <v-text-field
                  v-model="data.google_captcha_site_key"
                  :label="t('GOOGLE_CAPTCHA_SITE_KEY')"
                  density="default"
                  prepend-inner-icon="mdi-shield-key-outline"
                  variant="underlined"
              />
            </v-col>

            <v-col class="ma-0 pa-0 me-5" cols="12" lg="5" md="6" sm="6">
              <v-text-field
                  v-model="data.google_captcha_secret_key"
                  :label="t('GOOGLE_CAPTCHA_SECRET_KEY')"
                  density="default"
                  prepend-inner-icon="mdi-shield-key-outline"
                  variant="underlined"
              />
            </v-col>

            <v-col class="text-end mt-5">
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

      <v-divider vertical/>

      <v-col
          class="align-center justify-center fill-height"
          cols="12"
          md="5"
          style="background-color: #0d47a1"
      >
        <v-img
            :src="setupImage"
            alt="ocserv logo"
            contain
            style="background-color: transparent"
        />
      </v-col>

    </v-row>
  </div>
</template>

