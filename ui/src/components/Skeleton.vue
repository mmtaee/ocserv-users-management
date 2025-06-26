<script lang="ts" setup>
import logoUrl from "@/assets/ocserv.png"
import {useLocale, useTheme} from "vuetify/framework";
import {computed, onBeforeMount, ref} from "vue";
import {useUserStore} from "@/stores/user.ts";
import router from "@/plugins/router.ts";
import {useDisplay} from "vuetify";
import avatarUrl from "@/assets/torvalds.jpg"

const {t} = useLocale()
const theme = useTheme()
const display = useDisplay()
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

// const logout = () => {
//   const userStore = useUserStore()
//   userStore.clearUser()
//   localStorage.removeItem('token')
//   router.push('/login')
// }

</script>

<template>
  <v-navigation-drawer
      v-if="userStore.user?.username"
      v-model="drawer"
      :location="isSmallDisplay ? 'bottom' : undefined"
      :temporary="isSmallDisplay"
  >
    <v-list>
      <v-list-item :prepend-avatar="avatarUrl">
        <v-list-item-title>
          <span class="text-capitalize">
            {{ userStore.user.username }}
          </span>
        </v-list-item-title>

        <v-list-item-subtitle>
          <span>
            {{ userStore.user.isAdmin ? t('ADMIN') : t('STAFF') }}
          </span>
        </v-list-item-subtitle>
      </v-list-item>
    </v-list>

    <v-divider/>

    <v-list></v-list>

    <template #append>
      <v-divider class="mb-2"/>
      <div style="text-align: center; font-size: 0.9rem; color: #555; margin-bottom: 10px">
        <div>Built with ❤️ in 2025</div>
        <div>
          Need help? Contact
          <a
              href="https://github.com/mmtaee/ocserv-users-management/issues"
              style="color: #007BFF; text-decoration: none;"
              target="_blank"
          >
            Github
          </a>
        </div>
      </div>
    </template>

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

      <v-icon class="me-3" @click="router.push('/account')">
        mdi-account-cog-outline
      </v-icon>

    </template>

  </v-app-bar>

</template>
