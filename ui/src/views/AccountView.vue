<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import {defineAsyncComponent, onMounted, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";
import {useUserStore} from "@/stores/user.ts";

const Profile = defineAsyncComponent(() => import("@/components/account/Profile.vue"))
const Admin = defineAsyncComponent(() => import("@/components/account/Admin.vue"))


const {t} = useLocale()
const route = useRoute()
const router = useRouter()

const tab = ref('profile')

const userStore = useUserStore()

onMounted(() => {
  tab.value = route.query.tab?.toString() || 'profile'
})

watch(tab, (newVal) => {
  router.replace({
    query: {
      ...route.query,
      tab: newVal
    }
  })
})

</script>

<template>
  <v-row align="start" justify="center">
    <v-col>
      <v-card min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title class="text-capitalize">
            {{ t('ACCOUNT') }}
          </v-toolbar-title>
        </v-toolbar>
        <v-card flat>
          <v-tabs
              v-model="tab"
              align-tabs="center"
              color="primary"
              direction="horizontal"
          >
            <v-tab :text="t('PROFILE')" value="profile"/>
            <v-tab v-if="userStore.isAdmin" :text="t('ADMINS')" value="admins"/>
            <v-tab :text="t('LOGS')" value="logs"/>
          </v-tabs>

          <v-tabs-window v-model="tab">

            <v-tabs-window-item value="profile">
              <Profile/>
            </v-tabs-window-item>

            <v-tabs-window-item value="logs">
            </v-tabs-window-item>

            <v-tabs-window-item value="admins">
              <Admin/>
            </v-tabs-window-item>

          </v-tabs-window>
        </v-card>
      </v-card>
    </v-col>
  </v-row>
</template>
<style scoped>

</style>