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
import {dummyBanIPs, dummyOnlineUsers, dummyTrafficData} from "@/utils/dummy.ts";
import {tr} from "vuetify/locale";

type StatusResult = {
  cleanedText: string;
  status: string | null;
};

type TrafficStats = {
  rx: string | null;
  tx: string | null;
  cleanedText: string;
};

const {t} = useLocale()
const generalInfo = ref("")
const generalStatus = ref("")
const currentStats = ref("")
const currentTX = ref("")
const currentRX = ref("")

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

const statusWrapper = (text: string): StatusResult => {
  const statusMatch = text.match(/^Status:\s*(.+)$/m);
  const status = statusMatch ? statusMatch[1] : null;

  const cleanedText = text.replace(/^Status:.*\n?/m, '');
  return {cleanedText, status};
};


const trafficStatsWrapper = (text: string): TrafficStats => {
  const rxMatch = text.match(/^RX:\s*(.+)$/m);
  const txMatch = text.match(/^TX:\s*(.+)$/m);

  const cleanedText = text
      .replace(/^RX:.*\n?/m, '')
      .replace(/^TX:.*\n?/m, '');

  return {
    rx: rxMatch ? rxMatch[1] : null,
    tx: txMatch ? txMatch[1] : null,
    cleanedText,
  };
};

onMounted(() => {
  const api = new HomeApi()
  api.homeGet(getAuthorization()).then((res) => {
    const infoResult = statusWrapper(res.data.status.general_info)
    generalInfo.value = infoResult.cleanedText.replaceAll('\n', '<br>').replaceAll('"', "")
    generalStatus.value = infoResult.status || ""

    const trafficInfo = trafficStatsWrapper(res.data.status.current_stats)
    currentStats.value = trafficInfo.cleanedText.replaceAll('\n', '<br>').replaceAll('"', "")
    currentRX.value = trafficInfo.rx
    currentTX.value = trafficInfo.tx

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
  <v-row align="start" justify="center">

    <v-col class="ma-0" lg="3" md="6" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_SERVER_INFO") }}
        </v-card-title>

        <v-card-text class="mt-3">
          <v-btn :color="generalStatus === 'online'?'green': 'red'" block class="my-3" readonly>
            ({{ generalStatus }})
            <v-icon v-if="generalStatus === 'online'" end>mdi-check</v-icon>
            <v-icon v-else end>mdi-alert-circle</v-icon>
          </v-btn>

          <div class="text-justify text-subtitle-2 text-info ps-1" v-html="generalInfo"></div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="6" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="400">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_STATS_PERIOD") }}
        </v-card-title>
        <v-card-text class="mt-3">
          <v-btn block class="my-3" color="info" readonly>
            RX: {{ currentRX }}
            <v-icon class="ms-1 me-1"> mdi-download-outline</v-icon>
            |
            TX: {{ currentTX }}
            <v-icon class="ms-1 me-1">mdi-upload-outline</v-icon>
          </v-btn>

          <div class="text-justify text-subtitle-2 text-info ps-1" v-html="currentStats"></div>
        </v-card-text>
      </v-card>
    </v-col>


    <v-col class="ma-0" lg="3" md="6" sm="12" xs="12">
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
                {{ item.Username }}
                [{{ t("SINCE") }} {{ item['_Connected at'] }}]
                <br>
                RX: {{ item['Average RX'] }} <strong>|</strong>
                TX: {{ item['Average TX'] }} [{{ t("AVERAGES") }}]
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T05") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="3" md="6" sm="12" xs="12">
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
                <strong>IP:</strong> {{ item.IP }} [{{ item.Score }} {{ t("SCORES") }}]
                <br>
                <strong> {{ t("SINCE") }}:</strong> {{ item.Since }} [{{ item._Since?.trim() }}]
              </v-card>
            </template>
          </v-virtual-scroll>

          <div v-else class="mt-5 text-info">
            {{ t("T06") }}
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="6" md="6" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="430">
        <v-card-title class="bg-primary">
          {{ t("RX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text>
          <Line
              :data="rxChartData"
              :height="140"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="6" md="6" sm="12" xs="12">
      <v-card class="my-0" elevation="4" height="430">
        <v-card-title class="bg-primary">
          {{ t("TX_STATISTICS") }} (10 {{ t("DAYS") }})
        </v-card-title>
        <v-card-text>
          <Line
              :data="txChartData"
              :height="140"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>


  </v-row>
</template>
