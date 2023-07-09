<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="1400"
          min-height="450"
          max-height="800"
        >
          <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
            Statistics
          </v-card-subtitle>
          <v-card-text>
            <v-row align="center" justify="center">
              <v-col cols="12" md="auto" class="ma-0 pa-0">
                <v-card
                  width="auto"
                  class="text-h7 mb-1 pa-5 px-15 text-start"
                  elevation="2"
                >
                  Year: {{ total.year }} <br />
                  Total RX: {{ total.total_rx }} GB <br />
                  Total TX: {{ total.total_tx }} GB
                </v-card>
              </v-col>
            </v-row>
            <v-row align="center" justify="center">
              <v-col md="12" class="ma-0 pa-0">
                <v-card
                  width="auto"
                  class="text-h7 mb-1 pa-5 px-15 text-start"
                  flat
                  elevation="1"
                >
                  <VueApexCharts
                    v-if="chartOptions.xaxis.categories.length"
                    :options="chartOptions"
                    :series="series"
                    height="500"
                  />
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { statsServiceApi } from "@/utils/services";
import Vue from "vue";

export default Vue.extend({
  name: "Statistics",
  components: {
    VueApexCharts: () => import("vue-apexcharts"),
  },
  data() {
    return {
      chartOptions: {
        chart: {
          type: "line",
          zoom: {
            enabled: false,
          },
          toolbar: {
            show: false,
          },
          stroke: {
            curve: "smooth",
          },
          animations: {
            enabled: true,
            easing: "easeinout",
            speed: 800,
            animateGradually: {
              enabled: true,
              delay: 150,
            },
            dynamicAnimation: {
              enabled: true,
              speed: 350,
            },
          },
        },
        title: {
          text: "GB",
          style: {
            color: undefined,
            fontSize: "12px",
            fontFamily: "Helvetica, Arial, sans-serif",
            fontWeight: 600,
            cssClass: "apexcharts-xaxis-title",
          },
        },
        xaxis: {
          categories: [] as string[],
          tickPlacement: "between",
          position: "bottom",
        },
      },
      series: [
        {
          name: "rx",
          data: [] as number[],
        },
        {
          name: "tx",
          data: [] as number[],
        },
      ],
      total: {},
    };
  },
  async mounted() {
    let data = await statsServiceApi.get_stats();
    data.months.forEach((i) => {
      this.series[0].data.push(data.result[i].total_rx);
      this.series[1].data.push(data.result[i].total_tx);
    });
    this.chartOptions.xaxis.categories = data.months.map(
      (number) => `Month(${number})`
    );
    this.total = data.total;
  },
});
</script>