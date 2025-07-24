<script lang="ts" setup>
import {Line} from 'vue-chartjs'
import type {ChartData} from 'chart.js'
import {CategoryScale, Chart as ChartJS, Legend, LinearScale, LineElement, PointElement, Title, Tooltip} from 'chart.js'
import {computed, onMounted, ref} from "vue";
import {
  HomeApi,
  type HomeCurrentStats,
  type HomeGeneralInfo,
  type HomeGetHomeResponse,
  type ModelsDailyTraffic
} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {useLocale} from "vuetify/framework";
import {dummyBanIPs, dummyOnlineUsers} from "@/utils/dummy.ts";

const {t} = useLocale()
const trafficData = ref<ModelsDailyTraffic[]>([])

const homeData = ref<HomeGetHomeResponse>({
  ipbans: undefined,
  online_users_session: undefined,
  server_status: {},
  statistics: undefined
})

const labels = computed(() => trafficData.value.map(item => item.date))
const rxValues = computed(() => trafficData.value.map(item => item.rx))
const txValues = computed(() => trafficData.value.map(item => item.tx))

const generalInfo = ref<{ label: string; value: any }[]>([])
const currentStats = ref<{ label: string; value: any }[]>([])

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement)

const rxChartData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'RX (GB)',
      data: rxValues.value,
      borderColor: '#3b82f6',
      backgroundColor: 'rgba(59,130,246,0.2)',
      fill: true,
      tension: 0.3,
    }
  ]
} as ChartData<'line', number[], string>))

const txChartData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'TX (GB)',
      data: txValues.value,
      borderColor: '#10b981',
      backgroundColor: 'rgba(16,185,129,0.2)',
      fill: true,
      tension: 0.3,
    }
  ]
} as ChartData<'line', number[], string>))

const chartOptions = {
  responsive: true,
  plugins: {
    legend: {display: true},
    title: {display: false}
  },
  scales: {
    y: {beginAtZero: true}
  }
}

function formatWithDuration(dateKey: string, durationKey: string, source: Record<string, any>): string {
  const date = source[dateKey]?.trim()
  const duration = source[durationKey]?.trim()
  return date && duration ? `${date} (${duration})` : date ?? ''
}


const getGeneralInfo = (general: HomeGeneralInfo) => {
  return [
    {label: 'Server PID', value: general['Server PID']},
    {label: 'Sec-mod PID', value: general['Sec-mod PID']},
    {label: 'Sec-mod instance count', value: general['Sec-mod instance count']},
    {label: 'Up since', value: formatWithDuration('Up since', '_Up since', general)},
    {label: 'Active sessions', value: general['Active sessions']},
    {label: 'Total sessions', value: general['Total sessions']},
    {label: 'Total authentication failures', value: general['Total authentication failures']},
    {label: 'IPs in ban list', value: general['IPs in ban list']},
    {label: 'Median latency', value: general['Median latency']},
    {label: 'STDEV latency', value: general['STDEV latency']},
  ]
}

const getCurrentStats = (stats: HomeCurrentStats) => {
  return [
    {label: 'Last stats reset', value: formatWithDuration('Last stats reset', '_Last stats reset', stats)},
    {label: 'Sessions handled', value: stats['Sessions handled']},
    {label: 'Timed out sessions', value: stats['Timed out sessions']},
    {label: 'Timed out (idle) sessions', value: stats['Timed out (idle) sessions']},
    {label: 'Closed due to error sessions', value: stats['Closed due to error sessions']},
    {label: 'Authentication failures', value: stats['Authentication failures']},
    {label: 'Average auth time', value: stats['Average auth time']},
    {label: 'Max auth time', value: stats['Max auth time']},
    {label: 'Average session time', value: stats['Average session time']},
    {label: 'Max session time', value: stats['Max session time']},
  ]
}

onMounted(() => {
  const api = new HomeApi()
  api.homeGet(getAuthorization()).then((res) => {
    homeData.value = res.data
    generalInfo.value = getGeneralInfo(res.data.server_status.general_info || {})
    currentStats.value = getCurrentStats(res.data.server_status.current_stats || {})
    homeData.value.online_users_session = dummyOnlineUsers
    homeData.value.ipbans = dummyBanIPs
  })
})

</script>

<template>
  <v-row align="start" justify="center">

    <v-col class="ma-0" lg="4" md="6" sm="12" xl="3" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_SERVER_INFO") }}
        </v-card-title>

        <v-card-text class="mt-3">
          <v-btn
              :color="homeData.server_status.general_info?.Status === 'online'?'green': 'red'"
              block
              class="my-3"
              readonly
          >
            {{ t("STATUS") }}: {{ homeData.server_status.general_info?.Status }}
            <v-icon v-if="homeData.server_status.general_info?.Status === 'online'" end>mdi-check</v-icon>
            <v-icon v-else end>mdi-alert-circle</v-icon>
          </v-btn>

          <v-card class="my-2" style="max-height: 280px; overflow-y: auto;">
            <v-data-table
                :items="generalInfo"
                density="compact"
                hide-default-footer
                hide-default-header
                striped="odd"
            />
          </v-card>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12" xl="3" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_STATS_PERIOD") }}
        </v-card-title>
        <v-card-text class="mt-3">
          <v-btn block class="my-3" color="info" readonly>
            RX: {{ homeData.server_status.current_stats?.RX }}
            <v-icon class="ms-1 me-1"> mdi-download-outline</v-icon>
            |
            TX: {{ homeData.server_status.current_stats?.TX }}
            <v-icon class="ms-1 me-1">mdi-upload-outline</v-icon>
          </v-btn>

          <v-card class="my-2" style="max-height: 280px; overflow-y: auto;">
            <v-data-table
                :items="currentStats"
                density="compact"
                hide-default-footer
                hide-default-header
            />
          </v-card>
        </v-card-text>
      </v-card>
    </v-col>


    <v-col class="ma-0" lg="4" md="6" sm="12" xl="3" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("ONLINE_USERS") }} ({{ homeData.online_users_session?.length }} {{ t("USERS") }})
        </v-card-title>
        <v-card-text class="my-1">
          <v-virtual-scroll
              v-if="homeData.online_users_session?.length"
              :items="homeData.online_users_session"
              class="ma-2 mt-3"
              height="320"
          >
            <template v-slot:default="{ item, index }">
              <v-card :class="['pa-2 my-1', index % 2 === 0 ? 'bg-odd' : 'bg-even']">
                {{ item.Username }}
                <span class="text-grey">[{{ t("SINCE") }} {{ item['_Connected at'] }}]</span>
                <br>
                <span class="text-grey">RX:</span> {{ item['Average RX'] }}
                <strong class="text-grey"> | </strong>
                <span class="text-grey">TX:</span> {{ item['Average TX'] }}
                <span class="text-grey">
                  [{{ t("AVERAGES") }}]
                </span>
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T05") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12" xl="3" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("BAN_IPS") }} ({{ homeData.ipbans?.length }} {{ t("IPS") }})
        </v-card-title>

        <v-card-text class="my-1">
          <v-virtual-scroll
              v-if="homeData.ipbans?.length"
              :items="homeData.ipbans"
              class="ma-2 mt-3"
              height="320"
              item-key="banIPs"
          >
            <template v-slot:default="{ item, index }">
              <v-card :class="['pa-2 my-1', index % 2 === 0 ? 'bg-odd' : 'bg-even']">
                <strong class="text-grey">IP:</strong>
                {{ item.IP }}
                <span class="text-grey">[{{ item.Score }} {{ t("SCORES") }}]</span>
                <br>
                <strong class="text-grey"> {{ t("SINCE") }}:</strong> {{ item.Since }}
                <span class="text-grey">[{{ item._Since?.trim() }}]</span>
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T06") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="6" md="6" sm="12" xl="6" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("RX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text>
          <Line
              :data="rxChartData"
              :height="120"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="6" md="6" sm="12" xl="6" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("TX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text>
          <Line
              :data="txChartData"
              :height="120"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>


  </v-row>
</template>

<style>
tbody tr:nth-of-type(odd) {
  background-color: rgb(var(--v-theme-odd, 230, 240, 255)); /* fallback: #E6F0FF */
}

tbody tr:nth-of-type(even) {
  background-color: rgb(var(--v-theme-even, 248, 249, 250)); /* fallback: #F8F9FA */
}
</style>

