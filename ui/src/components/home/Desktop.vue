<script lang="ts" setup>
import {Line} from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement
} from 'chart.js'
import type {ChartData} from 'chart.js'
import {computed, onMounted, ref} from "vue";
import {HomeApi, type ModelsDailyTraffic} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {useLocale} from "vuetify/framework";

const {t} = useLocale()
const generalInfo = ref("")
const currentStats = ref("")
const trafficData = ref<ModelsDailyTraffic[]>([])
// const trafficData = ref<ModelsDailyTraffic[]>([
//   {date: '2025-06-18', rx: 1.2, tx: 2.5},
//   {date: '2025-06-19', rx: 0.9, tx: 1.1},
//   // missing 2025-06-20
//   {date: '2025-06-21', rx: 0.7, tx: 0.8},
//   {date: '2025-06-22', rx: 1.0, tx: 1.3},
//   {date: '2025-06-23', rx: 0.5, tx: 0.6},
//   // missing 2025-06-24
//   {date: '2025-06-25', rx: 0.3, tx: 0.4},
//   {date: '2025-06-26', rx: 1.5, tx: 2.0},
//   {date: '2025-06-27', rx: 2.1, tx: 3.2},
//   {date: '2025-06-28', rx: 1.8, tx: 2.3},
// ])
const onlineUsers = ref<Array<string>>([])

const labels = computed(() => trafficData.value.map(item => item.date))
const rxValues = computed(() => trafficData.value.map(item => item.rx))
const txValues = computed(() => trafficData.value.map(item => item.tx))

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

onMounted(() => {
  const api = new HomeApi()
  api.homeGet(getAuthorization()).then((res) => {
    generalInfo.value = res.data.status.general_info.replaceAll('\n', '<br>').replaceAll('"', "")
    currentStats.value = res.data.status.current_stats.replaceAll('\n', '<br>').replaceAll('"', "")
    if (res.data.stats) {
      trafficData.value = res.data.stats
    }
    if (res.data.online_user) {
      onlineUsers.value = res.data.online_user
    }
  })
})


</script>

<template>
  <v-row align="start" class="ma-0 pa-0" justify="start">
    <v-col lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("GENERAL_INFO") }}
        </v-card-title>
        <v-card-text class="mt-5 text-info">
          <div class="text-justify ps-5 text-subtitle-2" v-html="generalInfo"></div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_STATS_PERIOD") }}
        </v-card-title>
        <v-card-text class="mt-5">
          <div class="text-justify text-subtitle-2 text-info" v-html="currentStats"></div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("RX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text class="pa-5">
          <Line
              :data="rxChartData"
              :height="220"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("TX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text class="pa-5">
          <Line
              :data="txChartData"
              :height="220"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("ONLINE_USERS") }}
        </v-card-title>
        <v-card-text class="my-1">
          <v-virtual-scroll
              v-if="onlineUsers.length"
              :items="onlineUsers"
              class="ma-2 mt-3"
              height="380"
          >
            <template v-slot:default="{ item, index }">
              <div :class="['pa-2', index % 2 === 0 ? 'bg-info' : '']">
                {{ index + 1 }} - {{ item }}
              </div>
            </template>
          </v-virtual-scroll>
          <div v-else class="mt-5 text-info">
            {{ t("T05") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-text>
          ban users
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-text>
          iroutes
        </v-card-text>
      </v-card>
    </v-col>

  </v-row>
</template>
