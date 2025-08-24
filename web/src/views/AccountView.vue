<script lang="ts" setup>
import {useI18n} from "vue-i18n";
import {defineAsyncComponent, onMounted, ref, watch} from "vue";
import {useRoute, useRouter} from "vue-router";

const Profile = defineAsyncComponent(() => import("@/components/account/Profile.vue"))

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const tab = ref('profile')

// const userStore = useUserStore()

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
          <Profile/>
        </v-card>
      </v-card>
    </v-col>
  </v-row>
</template>
<style scoped>

</style>