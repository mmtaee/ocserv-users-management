<script lang="ts" setup>
import {computed, defineAsyncComponent, reactive, ref} from "vue";
import {type SystemLoginData, SystemUsersApi} from "@/api";
import {useConfigStore} from "@/stores/config.ts";
import {useLocale} from "vuetify/framework";
import {requiredRule} from "@/utils/rules.ts";
import router from "@/plugins/router.ts";
import {useUserStore} from "@/stores/user.ts";

const Captcha = defineAsyncComponent(() => import('@/components/Captcha.vue'));

const {t} = useLocale()
const configStore = useConfigStore()

const valid = ref(true)
const resetCaptcha = ref(false)
const loading = ref(false)
const siteKey = configStore.config.googleCaptchaSiteKey

const data = reactive<SystemLoginData>({password: "", remember_me: false, token: "", username: ""})

const rules = {
  required: (v: string) => requiredRule(v, t),
}

const login = async () => {
  if (!checkValid) {
    return
  }
  loading.value = true
  const api = new SystemUsersApi()
  api.systemUsersLoginPost(
      {
        request: data,
      }
  ).then((res) => {
    const userStore = useUserStore()
    userStore.setUser(res.data.user)
    localStorage.setItem("token", res.data.token)
    router.push("/")
  }).finally(() => {
    resetCaptcha.value = true
    loading.value = false
  })
}

const checkValid = computed(() => {
  let allStringsFilled = data.username.trim() !== "" && data.password.trim() !== ""
  if (siteKey) {
    allStringsFilled = allStringsFilled && data.token?.trim() !== "";
  }
  return allStringsFilled;
});

</script>

<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col cols="12" lg="5" md="7" sm="9" xl="4" xs="10">
        <v-card>

          <v-card-title class="bg-info mb-2">
            {{ t("Login") }}
          </v-card-title>

          <v-card-text>
            <v-form v-model="valid">
              <v-row align="center" class="ma-0 pa-0" justify="center">

                <v-col cols="12" lg="9" md="8" sm="11" xl="8" xs="11">
                  <v-text-field
                      v-model="data.username"
                      :label="t('USERNAME')"
                      :rules="[rules.required]"
                      density="compact"
                      prepend-inner-icon="mdi-account"
                      variant="underlined"
                      @keyup.enter="login"
                  />
                </v-col>

                <v-col cols="12" lg="9" md="8" sm="11" xl="8" xs="11">
                  <v-text-field
                      v-model="data.password"
                      :label="t('PASSWORD')"
                      :rules="[rules.required]"
                      density="compact"
                      prepend-inner-icon="mdi-key"
                      type="password"
                      variant="underlined"
                      @keyup.enter="login"
                  />
                </v-col>

                <v-col cols="12" lg="9" md="8" sm="11" xl="8" xs="11">
                  <Captcha
                      v-if="siteKey !== ''"
                      v-model="data.token"
                      :reset="resetCaptcha"
                      :siteKey="siteKey"
                  />
                </v-col>

              </v-row>
            </v-form>
          </v-card-text>

          <v-card-actions>
            <v-spacer/>
            <v-btn
                :disabled="!checkValid"
                :loading="loading"
                color="primary"
                variant="outlined"
                @click="login"

            >
              {{ t("LOGIN") }}
            </v-btn>
          </v-card-actions>

        </v-card>
      </v-col>
    </v-row>
  </v-container>

</template>
