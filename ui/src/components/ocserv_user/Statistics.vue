<script lang="ts" setup>
import {computed, defineAsyncComponent, onMounted, reactive, ref, watch} from "vue";
import {type ModelsDailyTraffic, type ModelsOcservUser, OcservUsersApi, type OcservUserStatisticsData} from "@/api";
import {useLocale} from "vuetify/framework";
import {formatDate} from "@/utils/convertors.ts";
import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  type ChartData,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from 'chart.js'
import {Bar, Line} from 'vue-chartjs'
import {getAuthorization} from "@/utils/request.ts";


const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsOcservUser
}>()

const emit = defineEmits(["update:modelValue", "done"])


ChartJS.register(
    Title, Tooltip, Legend,
    LineElement, PointElement,
    BarElement,
    CategoryScale, LinearScale
)

const {t} = useLocale()
const dateStartMenu = ref(false)
const dateEndMenu = ref(false)

const date = reactive<OcservUserStatisticsData>({
  date_start: "",
  date_end: "",
})

const traffic = ref<ModelsDailyTraffic[]>([])

// X-axis labels
const labels = computed(() => traffic.value.map(d => d.date))

// Line Chart – Rx only
const rxLineData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'Rx (MB)',
      data: traffic.value.map(d => d.rx),
      borderColor: '#2C7BE5',
      backgroundColor: 'rgba(44, 123, 229, 0.2)',
      tension: 0.3,
      fill: false,
    }
  ]
}) as ChartData<'line', number[], string>)

// Line Chart – Tx only
const txLineData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'Tx (MB)',
      data: traffic.value.map(d => d.tx),
      borderColor: '#6C757D',
      backgroundColor: 'rgba(108, 117, 125, 0.2)',
      tension: 0.3,
      fill: false,
    }
  ]
}) as ChartData<'line', number[], string>)


// Bar Chart – Rx and Tx
const barData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'Rx (MB)',
      data: traffic.value.map(d => d.rx),
      backgroundColor: 'rgba(44, 123, 229, 0.7)',
    },
    {
      label: 'Tx (MB)',
      data: traffic.value.map(d => d.tx),
      backgroundColor: 'rgba(108, 117, 125, 0.7)',
    }
  ]
}) as ChartData<'bar', number[], string>)

const chartOptions = {
  responsive: true,
  plugins: {
    legend: {
      position: 'top' as const
    },
    title: {
      display: true
    }
  }
}

const dynamicHeight = computed(() => {
  return window.innerWidth < 450 ? 190 : 160;
});

const search = () => {
  const api = new OcservUsersApi()
  api.ocservUsersUidStatisticsGet({
    ...getAuthorization(),
    uid: props.user.uid,
    request: date,
  }).then((res) => {
    traffic.value = res.data || []
  })
}

const searchBtnDisabled = computed(() => {
  const start = date.date_start
  const end = date.date_end

  if (!start && !end) return true

  if (start && end) {
    const startDate = new Date(start)
    const endDate = new Date(end)
    return startDate > endDate
  }

  return false
})

onMounted(
    () => {
      let today = new Date()
      let lastMonth = new Date(today)
      lastMonth.setMonth(today.getMonth() - 1)
      date.date_start = formatDate(lastMonth.toString())
      date.date_end = formatDate(today.toString())
    }
)

watch(
    () => props.user.uid,
    () => {
      search()
    },
    {deep: true}
)

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      btnClose
      color="primary"
      fullscreen
      transition="dialog-bottom-transition"
      width="500"
      @update:modelValue="val => $emit('update:modelValue', val)"
  >
    <template #dialogTitle>
      <v-icon class="mb-1">mdi-delete</v-icon>
      <span class="text-capitalize">{{ t("OCSERV_USER_STATISTICS_TITLE") }} ({{ user.username }})</span>
    </template>

    <template #dialogText>
      <v-row align="start" class="mt-3" justify="center">
        <v-col cols="12" lg="2" md="4">
          <v-menu
              v-model="dateStartMenu"
              :close-on-content-click="false"
              transition="scale-transition"
          >
            <template #activator="{ props }">
              <v-text-field
                  v-model="date.date_start"
                  :label="t('DATE_START')"
                  clearable
                  density="compact"
                  readonly
                  v-bind="props"
                  variant="underlined"
              />
            </template>
            <v-date-picker
                v-model="date.date_start"
                :header="t('DATE_START')"
                elevation="24"
                title=""
                @update:model-value="(val:any)=>{date.date_start = formatDate(val); dateStartMenu = false}"
            />
          </v-menu>
        </v-col>

        <v-col cols="12" lg="2" md="4">
          <v-menu
              v-model="dateEndMenu"
              :close-on-content-click="false"
              transition="scale-transition"
          >
            <template #activator="{ props }">
              <v-text-field
                  v-model="date.date_end"
                  :label="t('DATE_END')"
                  clearable
                  density="compact"
                  readonly
                  v-bind="props"
                  variant="underlined"
              />
            </template>
            <v-date-picker
                v-model="date.date_end"
                :header="t('DATE_END')"
                elevation="24"
                title=""
                @update:model-value="(val:any)=>{date.date_end=formatDate(val); dateEndMenu = false}"
            />
          </v-menu>
        </v-col>

        <v-col cols="12" lg="2" md="4">
          <v-btn
              :disabled="searchBtnDisabled"
              color="primary"
              variant="elevated"
              @click="search"
          >
            <v-icon start>mdi-magnify</v-icon>
            {{ t("SEARCH") }}
          </v-btn>
        </v-col>

      </v-row>

      <v-divider/>


      <v-row align="center" justify="center">
        <v-col class="ma-0" lg="4" md="6" sm="12" xs="12">
          <v-card class="my-0" elevation="4" height="430">
            <v-card-title class="bg-primary">
              RX
            </v-card-title>
            <v-card-text>
              <Line
                  :data="rxLineData"
                  :height="dynamicHeight"
                  :options="chartOptions"
              />
            </v-card-text>
          </v-card>
        </v-col>

        <v-col class="ma-0" lg="4" md="6" sm="12" xs="12">
          <v-card class="my-0" elevation="4" height="430">
            <v-card-title class="bg-primary">
              TX
            </v-card-title>
            <v-card-text>
              <Line
                  :data="txLineData"
                  :height="dynamicHeight"
                  :options="chartOptions"
              />
            </v-card-text>
          </v-card>
        </v-col>


        <v-col class="ma-0" lg="4" md="6" sm="12" xs="12">
          <v-card class="my-0" elevation="4" height="430">
            <v-card-title class="bg-primary">
              RX/TX
            </v-card-title>
            <v-card-text>
              <Bar
                  :data="barData"
                  :height="dynamicHeight"
                  :options="chartOptions"
              />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

    </template>
  </ReusableDialog>
</template>

