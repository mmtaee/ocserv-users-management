<script lang="ts" setup>
import logoUrl from "@/assets/ocserv.png"
import {useTheme} from "vuetify/framework";
import {defineAsyncComponent, onBeforeMount, ref} from "vue";
import {useUserStore} from "@/stores/user.ts";
import router from "@/plugins/router.ts";
import {useIsSmallDisplay} from "@/stores/display.ts";
import {useI18n} from "vue-i18n";
import {languages} from '@/plugins/i18n'

const SideBar = defineAsyncComponent(() => import("@/components/SideBar.vue"));
const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"));

const {t, locale} = useI18n()
const theme = useTheme()
const userStore = useUserStore()
const logoutDialog = ref(false)
const smallDisplay = useIsSmallDisplay()
const drawer = ref(!smallDisplay.isSmallDisplay)
const currentLang = ref(locale.value)
const languageDialog = ref(false)


function changeLanguage() {
  localStorage.setItem("language", currentLang.value)
  locale.value = currentLang.value
  window.location.reload()
  languageDialog.value = false
}

onBeforeMount(() => {
  theme.global.name.value = localStorage.getItem('theme') === 'dark' ? 'dark' : 'light';
});

function toggleTheme() {
  theme.global.name.value = theme.global.current.value.dark ? 'light' : 'dark'
  localStorage.setItem('theme', theme.global.name.value)
}

const logout = () => {
  const userStore = useUserStore()
  userStore.clearUser()
  localStorage.removeItem('token')
  logoutDialog.value = false
  router.push('/login')
}


</script>

<template>
  <SideBar v-model="drawer"/>

  <v-app-bar :elevation="4" density="comfortable">
    <template v-slot:prepend>
      <v-app-bar-nav-icon v-if="userStore.user?.username" variant="text" @click.stop="drawer = !drawer"/>
      <v-img :src="logoUrl" alt="ocserv logo" width="45"/>
    </template>

    <template v-slot:title>
      <span class="text-subtitle-1">Ocserv User Management</span>
    </template>

    <template v-slot:append>
      <v-icon class="me-3" @click="languageDialog=true">
        mdi-translate
      </v-icon>

      <v-icon @click="toggleTheme">mdi-theme-light-dark</v-icon>

      <v-icon v-if="userStore.user?.username" class="me-5 mx-3" color="error" @click="logoutDialog = true">
        mdi-logout
      </v-icon>
    </template>

  </v-app-bar>

  <ReusableDialog
      v-if="userStore.user?.username"
      v-model="logoutDialog"
      color="error"
      transition="dialog-top-transition"
  >
    <template #dialogTitle>
      <v-icon>mdi-logout</v-icon>
      <span class="text-capitalize">{{ t("LOGOUT_TITLE") }}</span>
    </template>

    <template #dialogText>
      {{ t("LOGOUT_MESSAGE") }} <br/><br/>
      <span class="text-subtitle-2">{{ t("LOGOUT_MESSAGE_SUB") }}</span> <br/>
      <span class="text-subtitle-2">{{ t("LOGOUT_MESSAGE_SUB_2") }}</span>
    </template>

    <template #dialogAction>
      <v-btn
          color="black"
          variant="outlined"
          @click="logoutDialog = false"
      >
        {{ t("CANCEL") }}
      </v-btn>

      <v-btn
          color="error"
          variant="outlined"
          @click="logout"
      >
        {{ t("LOGOUT") }}
      </v-btn>
    </template>
  </ReusableDialog>

  <ReusableDialog
      v-model="languageDialog"
      color="white"
      transition="dialog-top-transition"
  >
    <template #dialogText>
      <div class="pa-0" dir="ltr">
        <v-radio-group
            v-model="currentLang"
            class="mt-3"
            density="comfortable"
            dir="ltr"
            hide-details
            inline
            @change="changeLanguage"
        >
          <v-radio
              v-for="item in languages"
              :key="item.code"
              :label="item.label"
              :value="item.code"
              class="mx-1"
              density="comfortable"
              dir="ltr"
              inline
          >
          </v-radio>
        </v-radio-group>
      </div>
    </template>

  </ReusableDialog>

</template>
