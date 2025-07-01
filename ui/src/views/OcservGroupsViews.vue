<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import {defineAsyncComponent, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from "vue-router";

const Defaults = defineAsyncComponent(() => import('@/components/ocserv_group/Defaults.vue'));
const Others = defineAsyncComponent(() => import('@/components/ocserv_group/Others.vue'));

const {t} = useLocale()
const route = useRoute()
const router = useRouter()

const tab = ref('defaults')

onMounted(() => {
  tab.value = route.query.tab?.toString() || 'defaults'
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
  <v-row>
    <v-col>
      <v-card min-height="850">
        <v-toolbar :title="t('OCSERV_GROUPS')" color="primary"/>

        <v-tabs
            v-model="tab"
            align-tabs="center"
            direction="horizontal"
        >
          <v-tab :text="t('DEFAULTS')" value="defaults"></v-tab>
          <v-tab :text="t('OTHERS')" value="others"></v-tab>
        </v-tabs>

        <v-tabs-window v-model="tab">

          <v-tabs-window-item value="defaults">
            <Defaults/>
          </v-tabs-window-item>

          <v-tabs-window-item value="others">
            <Others/>
          </v-tabs-window-item>

        </v-tabs-window>
      </v-card>
    </v-col>
  </v-row>

</template>
