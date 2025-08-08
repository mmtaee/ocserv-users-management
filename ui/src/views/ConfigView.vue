<script lang="ts" setup>
import {onMounted, ref} from "vue";
import {SystemApi, type SystemPatchSystemUpdateData} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {useLocale} from "vuetify/framework";

const {t} = useLocale()
const api = new SystemApi()
const showConfigEdit = ref(false)
const configData = ref<SystemPatchSystemUpdateData>({
  google_captcha_site_key: "",
  google_captcha_secret_key: ""
})
const loading = ref(false)

const updateSetting = () => {
  loading.value = true
  api.systemPatch({
    ...getAuthorization(),
    request: configData.value,
  }).then((res) => {
    console.log(res.data)
    showConfigEdit.value = false
  }).finally(() => {
    loading.value = false
  })
}


onMounted(() => {
  api.systemGet({
    ...getAuthorization()
  }).then((res) => {
    configData.value.google_captcha_site_key = res.data.google_captcha_site_key || ""
    configData.value.google_captcha_secret_key = res.data.google_captcha_secret_key || ""
  })
})


</script>

<template>

  <v-row align="start" justify="center">
    <v-col cols="12">
      <v-card class="p" min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title class="text-capitalize">
            {{ t('CONFIG') }}
          </v-toolbar-title>
        </v-toolbar>

        <v-card-title class="text-capitalize mt-5 mx-5">
          <v-icon class="me-2" start>mdi-account-outline</v-icon>
          {{ t("SYSTEM_CONFIGS") }}
          <v-btn
              v-if="!showConfigEdit"
              color="primary"
              density="compact"
              variant="plain"
              @click="showConfigEdit=true"
          >
            {{ t("EDIT_CONFIG") }}
          </v-btn>
        </v-card-title>

        <v-card-text class="mx-5">
          <v-list max-width="800">
            <v-list-item :class="!showConfigEdit ? 'mb-5 mt-10': 'mt-5'">
              <template v-if="!showConfigEdit" #prepend>
                <v-icon size="large">mdi-shield-key-outline</v-icon>
              </template>

              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">
                {{ t("GOOGLE_SITE_KEY") }}
              </v-list-item-title>

              <v-list-item-subtitle class="text-subtitle-1">
                <span v-if="!showConfigEdit" class="mb-5">
                  {{ configData.google_captcha_site_key }}
                </span>
                <v-text-field
                    v-else
                    v-model="configData.google_captcha_site_key"
                    :prepend-inner-icon="showConfigEdit ? 'mdi-shield-key-outline': ''"
                >
                </v-text-field>
              </v-list-item-subtitle>
            </v-list-item>


            <v-list-item>
              <template v-if="!showConfigEdit" #prepend>
                <v-icon size="large">mdi-shield-key-outline</v-icon>
              </template>

              <v-list-item-title class="text-subtitle-2 text-capitalize mb-2">
                {{ t("GOOGLE_SECRET_KEY") }}
              </v-list-item-title>

              <v-list-item-subtitle class="text-subtitle-1">

                <div v-if="!showConfigEdit" class="mb-5">
                  {{ configData.google_captcha_secret_key }}
                </div>

                <v-text-field
                    v-else
                    v-model="configData.google_captcha_secret_key"
                    :prepend-inner-icon="showConfigEdit ? 'mdi-shield-key-outline': ''"
                >
                </v-text-field>

              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions v-if="showConfigEdit" class="w-50 ms-8">
          <v-spacer/>

          <v-btn
              color="secondary"
              variant="outlined"
              @click="showConfigEdit = false"
          >
            {{ t("CANCEL") }}
          </v-btn>

          <v-btn
              :loading="loading"
              class="mx-2 me-5"
              color="primary"
              variant="flat"
              @click="updateSetting"
          >
            {{ t("SAVE") }}
          </v-btn>
        </v-card-actions>

      </v-card>
    </v-col>
  </v-row>

</template>

<style scoped>

</style>