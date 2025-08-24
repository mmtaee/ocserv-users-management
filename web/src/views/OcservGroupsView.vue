<script lang="ts" setup>
import {useI18n} from "vue-i18n";
import {defineAsyncComponent, onMounted, ref, watch} from 'vue'
import {useRoute, useRouter} from "vue-router";

const Defaults = defineAsyncComponent(() => import('@/components/ocserv_group/Defaults.vue'));
const Others = defineAsyncComponent(() => import('@/components/ocserv_group/Others.vue'));


const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const tab = ref('others')

const createDialog = ref(false)

onMounted(() => {
  tab.value = route.query.tab?.toString() || 'others'
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
        <v-toolbar color="secondary">
          <v-toolbar-title>
            {{ t('OCSERV_GROUPS') }}
          </v-toolbar-title>

          <template v-slot:append>
            <v-btn
                v-if="tab == 'others'"
                class="ma-5"
                color="primary"
                variant="elevated"
                @click="createDialog = true"
            >
              {{ t("CREATE") }}
            </v-btn>
          </template>
        </v-toolbar>

        <v-tabs
            v-model="tab"
            align-tabs="center"
            color="primary"
            direction="horizontal"
        >
          <v-tab :text="t('OTHERS')" value="others"></v-tab>
          <v-tab :text="t('DEFAULTS')" value="defaults"></v-tab>
        </v-tabs>

        <v-tabs-window v-model="tab">

          <v-tabs-window-item value="defaults">
            <Defaults/>
          </v-tabs-window-item>

          <v-tabs-window-item value="others">
            <Others v-model="createDialog"/>
          </v-tabs-window-item>

        </v-tabs-window>
      </v-card>
    </v-col>
  </v-row>

</template>
