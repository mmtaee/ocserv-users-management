<script lang="ts" setup>
import avatarUrl from "@/assets/torvalds.jpg";
import {useLocale} from "vuetify/framework";
import {useUserStore} from "@/stores/user.ts";
import {useIsSmallDisplay} from "@/stores/display.ts";
import {useServerStore} from "@/stores/server.ts";

const props = defineProps({
  modelValue: Boolean
})

const emit = defineEmits(["update:modelValue"])

const userStore = useUserStore()
const serverStore = useServerStore()

const {t} = useLocale()

const smallDisplay = useIsSmallDisplay()

const items = [
  {
    "id": 0,
    "value": "home",
    "icon": "mdi-home",
    "title": t("HOME"),
    "to": "/"
  },
  {
    "id": 1,
    "value": "ocserv-group",
    "icon": "mdi-router-network",
    "title": t("OCSERV_GROUPS"),
    "to": "/ocserv-groups"
  },
  {
    "id": 2,
    "value": "ocserv-user",
    "icon": "mdi-account-network",
    "title": t("OCSERV_USERS"),
    "to": "/ocserv-users"
  },
  {
    "id": 3,
    "value": "occtl",
    "icon": "mdi-console",
    "title": "Occtl",
    "to": "/occtl"
  },
  {
    "id": 4,
    "value": "stats",
    "icon": "mdi-chart-bar-stacked",
    "title": t("STATISTICS"),
    "to": "/statistics"
  },
  {
    "id": 5,
    "value": "logs",
    "icon": "mdi-server-network",
    "title": t("LOGS"),
    "to": "/logs"
  },
]

</script>

<template>
  <v-navigation-drawer
      v-if="userStore.user?.username"
      v-model="props.modelValue"
      :location="smallDisplay.isSmallDisplay ? 'bottom' : undefined"
  >
    <v-list>
      <v-list-item :prepend-avatar="avatarUrl">
        <v-list-item-title>
          <v-row>
            <v-col>
              <span class="text-capitalize">
                {{ userStore.user.username }}
                (<span>{{ userStore.isAdmin ? t('ADMIN') : t('STAFF') }}</span>)
             </span>
            </v-col>
            <v-col v-if="smallDisplay.isSmallDisplay" class="text-end">
              <v-icon @click="emit('update:modelValue')">mdi-close</v-icon>
            </v-col>
          </v-row>
        </v-list-item-title>
      </v-list-item>
    </v-list>

    <v-divider/>

    <v-list>
      <v-list-item
          v-for="(item, i) in items"
          :key="`${item.value}-${i}`"
          :prepend-icon="item.icon"
          :title="item.title"
          :to="item.to"
          :value="item.value"
          color="primary"
      />
      <v-divider class="my-3" opacity="1"/>
      <v-list-item
          color="primary"
          prepend-icon="mdi-account-cog"
          to="/account"
          value="account"
      >
        <v-list-item-title class="text-capitalize">{{ t('ACCOUNT') }}</v-list-item-title>
      </v-list-item>

      <v-list-item
          v-if="userStore.isAdmin"
          color="primary"
          prepend-icon="mdi-cog"
          to="/config"
          value="config"
      >
        <v-list-item-title class="text-capitalize">{{ t('CONFIG') }}</v-list-item-title>
      </v-list-item>
    </v-list>

    <template #append>
      <div v-if="!smallDisplay.isSmallDisplay">
        <v-divider class="mb-2"/>
        <div style="text-align: center; font-size: 0.9rem; color: #555; margin-bottom: 10px">
          <div class="me-3 text-primary">Ocserv Version: {{ serverStore.versionInfo }}</div>
          <div>
            Need help? Contact
            <a
                href="https://github.com/mmtaee/ocserv-users-management/issues"
                style="color: #007BFF; text-decoration: none;"
                target="_blank"
            >
              <v-icon>
                mdi-github
              </v-icon>
            </a>
          </div>
        </div>
      </div>
    </template>

  </v-navigation-drawer>
</template>
