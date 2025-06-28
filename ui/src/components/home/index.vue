<script lang="ts" setup>
import {Line} from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement
} from 'chart.js'
import type {ChartData} from 'chart.js'
import {computed, onMounted, ref} from "vue";
import {HomeApi, type ModelsDailyTraffic, type ModelsIPBan, type ModelsOnlineUserSession} from "@/api";
import {getAuthorization} from "@/utils/request.ts";
import {useLocale} from "vuetify/framework";
import {dummyBanIPs, dummyOnlineUsers, dummyTrafficData} from "@/components/home/dummy.ts";

const {t} = useLocale()
const generalInfo = ref("")
const currentStats = ref("")

const trafficData = ref<ModelsDailyTraffic[]>([])
const onlineUsers = ref<Array<ModelsOnlineUserSession>>([])
const banIPs = ref<Array<ModelsIPBan>>([])

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
    trafficData.value = dummyTrafficData

    if (res.data.online_users_session) {
      onlineUsers.value = res.data.online_users_session
    }
    onlineUsers.value = dummyOnlineUsers

    if (res.data.ipbans) {
      banIPs.value = res.data.ipbans
    }
    banIPs.value = dummyBanIPs
  })
})
</script>

<template>
  <v-row align="start" class="ma-0 pa-0" justify="center">
    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("GENERAL_INFO") }}
        </v-card-title>
        <v-card-text class="mt-5 text-info">
          <div class="text-justify ps-5 text-subtitle-2" v-html="generalInfo"></div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_STATS_PERIOD") }}
        </v-card-title>
        <v-card-text class="mt-5">
          <div class="text-justify text-subtitle-2 text-info" v-html="currentStats"></div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("RX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text class="pa-5">
          <Line
              :data="rxChartData"
              :height="280"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("TX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text class="pa-5">
          <Line
              :data="txChartData"
              :height="280"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("ONLINE_USERS") }} ({{ onlineUsers.length }} {{ t("USERS") }})
        </v-card-title>
        <v-card-text class="my-1">
          <v-virtual-scroll
              v-if="onlineUsers.length"
              :items="onlineUsers"
              class="ma-2 mt-3"
              height="320"
          >
            <template v-slot:default="{ item, index }">
              <v-card :class="['pa-2 my-1', index % 2 === 0 ? 'bg-info' : 'bg-odd']">
                <strong>{{ t("USERNAME") }}: </strong>{{ item.Username }}
                <br>
                <strong>{{ t("GROUP") }}</strong>: {{ item.Groupname === "(none)" ? "defaults" : item.Groupname }}
                <br>
                <strong>{{ t("CONNECTED_SINCE") }}</strong>: {{ item['_Connected at'] }}
                <br>
                <strong>RX {{ t("AVERAGES") }}</strong>: {{ item['Average RX'] }}
                <br>
                <strong>TX {{ t("AVERAGES") }}</strong>: {{ item['Average TX'] }}
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T05") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("BAN_IPS") }} ({{ banIPs.length }} {{ t("IPS") }})
        </v-card-title>

        <v-card-text class="my-1">
          <v-virtual-scroll
              v-if="banIPs.length"
              :items="banIPs"
              class="ma-2 mt-3"
              height="320"
              item-key="banIPs"
          >
            <template v-slot:default="{ item, index }">
              <v-card :class="['pa-2 my-1', index % 2 === 0 ? 'bg-info' : 'bg-odd']">
                <strong>IP:</strong> {{ item.IP }} ({{ item.Score }} {{ t("SCORES") }})
                <br>
                <strong> {{ t("SINCE") }}:</strong> {{ item.Since }} ({{ item._Since?.trim() }})
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T06") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="4" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("IROUTES") }}
        </v-card-title>
        <v-card-text>
          <div class="mt-5 text-info">
            this section is disable in ocserv version 1.2.4
          </div>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>
