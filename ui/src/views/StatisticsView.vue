<script lang="ts" setup>
import {defineAsyncComponent, ref} from "vue";
import {useLocale} from "vuetify/framework";
import {OcservStatisticsApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";

const ReusableStatistics = defineAsyncComponent(() => import("@/components/reusable/ReusableStatistics.vue"))

const {t} = useLocale()
const traffic = ref<any[]>([])

const search = (dateStart: string, dateEnd: string) => {
  console.log(dateStart, dateEnd)
  const api = new OcservStatisticsApi()
  api.ocservUsersStatisticsGet({
    ...getAuthorization(),
    dateStart,
    dateEnd,
  }).then((res) => {
    traffic.value = res.data
  })
}

</script>

<template>
  <v-row align="start" justify="center">
    <v-col>
      <v-card min-height="850">
        <v-toolbar color="secondary">
          <v-toolbar-title>
            {{ t('STATISTICS') }}
          </v-toolbar-title>
        </v-toolbar>
        <v-card flat>
          <v-card-text>
            <ReusableStatistics :traffic="traffic" immediate @search="search"/>
          </v-card-text>
        </v-card>
      </v-card>
    </v-col>
  </v-row>
</template>
