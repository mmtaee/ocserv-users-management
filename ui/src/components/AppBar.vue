<script lang="ts" setup>
import logoUrl from "@/assets/ocserv.png"
import {useLocale, useTheme} from "vuetify/framework";
import {computed, defineAsyncComponent, onBeforeMount, ref} from "vue";
import {useUserStore} from "@/stores/user.ts";
import router from "@/plugins/router.ts";
import {useDisplay} from "vuetify";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"));

const {t} = useLocale()
const theme = useTheme()
const display = useDisplay()
const logoutDialog = ref(false)
const userStore = useUserStore()
const isSmallDisplay = computed(() => display.mdAndDown.value)
const drawer = ref(!isSmallDisplay.value)

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
  router.push('/login')
  logoutDialog.value = false
}


const items = [
  {
    title: 'Foo',
    value: 'foo',
  },
  {
    title: 'Bar',
    value: 'bar',
  },
  {
    title: 'Fizz',
    value: 'fizz',
  },
  {
    title: 'Buzz',
    value: 'buzz',
  },
]
</script>

<template>
  <v-navigation-drawer
      v-if="userStore.user?.username"
      v-model="drawer"
      :location="isSmallDisplay ? 'bottom' : undefined"
      :temporary="isSmallDisplay"
  >
    <v-list
        :items="items"
    ></v-list>
  </v-navigation-drawer>

  <v-app-bar :elevation="12" density="comfortable">

    <template v-slot:prepend>
      <v-app-bar-nav-icon v-if="userStore.user?.username" variant="text" @click.stop="drawer = !drawer"/>
      <v-img :src="logoUrl" alt="ocserv logo" width="45"/>
    </template>

    <template v-slot:title>
      <span class="text-subtitle-1">Ocserv User Management</span>
    </template>

    <template v-slot:append>
      <v-btn density="comfortable" icon @click="toggleTheme">
        <v-icon>mdi-theme-light-dark</v-icon>
      </v-btn>
      <v-btn v-if="userStore.user?.username" density="comfortable" icon @click="logoutDialog = true">
        <v-icon>mdi-logout</v-icon>
      </v-btn>
    </template>

  </v-app-bar>

  <ReusableDialog
      v-if="userStore.user?.username"
      v-model="logoutDialog"
      color="error"
  >
    <template #dialogTitle>
      <v-icon>mdi-logout</v-icon>
      {{ t("LOGOUT_TITLE") }}
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

</template>
