<script lang="ts" setup>

import {useLocale} from "vuetify/framework";

const {t} = useLocale()

import {Line} from 'vue-chartjs'
import {
  Chart as ChartJS,
  Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement, type TooltipItem
} from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement)

const chartData = {
  labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', "1", "2"],
  datasets: [
    {
      label: 'Value',
      data: [10, 20, 15, 100, 25, 0, 200],
      borderColor: 'blue',
      tension: 0.4,
    },
  ],
}

const chartRXOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      title: {
        display: true,
        text: 'RX',
      },
    },
    x: {
      title: {
        display: true,
        text: 'Day',
      },
    },
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      callbacks: {
        label: (context: TooltipItem<'line'>) => `RX: ${context.parsed.y}`,
      },
    },
  },
}
const chartTXOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      title: {
        display: true,
        text: 'TX',
      },
    },
    x: {
      title: {
        display: true,
        text: 'Day',
      },
    },
  },
  plugins: {
    legend: {
      display: false,
    },
    tooltip: {
      callbacks: {
        label: (context: TooltipItem<'line'>) => `TX: ${context.parsed.y}`,
      },
    },
  },
}

const onlineUsers = [
  "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash",
  "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash",
  "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash",
  "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash", "masoud", "kiarash",
]

</script>

<template>
  <v-row align="start" class="ma-0 pa-0" justify="start">
    <v-col lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("GENERAL_INFO") }}
        </v-card-title>
        <v-card-text class="mt-5">
          <div class="text-justify ps-5 text-subtitle-1">
            Status: online<br>
            Server PID: 37<br>
            Sec-mod PID: 45<br>
            Sec-mod instance count: 3<br>
            Up since: 2025-06-27 10:37 ( 4h:26m)<br>
            Active sessions: 0<br>
            Total sessions: 0<br>
            Total authentication failures: 0<br>
            IPs in ban list: 0<br>
            Median latency: <1ms<br>
            STDEV latency: <1ms<br>
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("CURRENT_STATS_PERIOD") }}
        </v-card-title>
        <v-card-text class="mt-5">
          <div class="text-justify ps-5 text-subtitle-1">
            Last stats reset: 2025-06-27 10:37 ( 4h:26m)<br>
            Sessions handled: 0<br>
            Timed out sessions: 0<br>
            Timed out (idle) sessions: 0<br>
            Closed due to error sessions: 0<br>
            Authentication failures: 0<br>
            Average auth time: 0s<br>
            Max auth time: 0s<br>
            Average session time: 0s<br>
            Max session time: 0s<br>
            RX: 0 bytes<br>
            TX: 0 bytes<br>
          </div>
        </v-card-text>

      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("RX_STATISTICS") }}
        </v-card-title>
        <v-card-text class="pa-5">
          <Line :data="chartData" :height="380"
                :options="chartRXOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>

    <v-col class="ma-0" lg="4" md="6" sm="12">
      <v-card class="my-0" elevation="4" height="450">
        <v-card-title class="bg-primary">
          {{ t("TX_STATISTICS") }}
        </v-card-title>
        <v-card-text class="pa-5">
          <Line :data="chartData" :height="380"
                :options="chartTXOptions"
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

<style scoped>

</style>