<script lang="ts" setup>
import {Bar, Line} from "vue-chartjs";
import {computed, type PropType, reactive, ref, watch} from "vue";
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
  Tooltip
} from "chart.js";
import {useLocale} from "vuetify/framework";
import {formatDate} from "@/utils/convertors.ts";

interface Traffic {
  date?: string
  rx?: number
  tx?: number
}

interface Date {
  date_start: string
  date_end: string,
}

const props = defineProps({
  traffic: {
    type: Array as PropType<Traffic[]>,
    required: true,
  },
  immediate: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(["search"])

ChartJS.register(
    Title, Tooltip, Legend,
    LineElement, PointElement,
    BarElement,
    CategoryScale, LinearScale
)


const {t} = useLocale()
const dateStartMenu = ref(false)
const dateEndMenu = ref(false)
const width = window.innerWidth

const date = reactive<Date>({
  date_start: "",
  date_end: "",
})

// X-axis labels
const labels = computed(() => props.traffic.map(d => d.date))

const lineData = computed(() => ({
  labels: labels.value,
  datasets: [
    {
      label: 'Rx (GiB)',
      data: props.traffic.map(d => d.rx),
      borderColor: '#2C7BE5',
      backgroundColor: 'rgba(44, 123, 229, 0.2)',
      tension: 0.3,
      fill: false,
    },
    {
      label: 'Tx (GiB)',
      data: props.traffic.map(d => d.tx),
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
      label: 'Rx (GiB)',
      data: props.traffic.map(d => d.rx),
      backgroundColor: 'rgba(44, 123, 229, 0.7)',
    },
    {
      label: 'Tx (GiB)',
      data: props.traffic.map(d => d.tx),
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

const search = () => {
  emit("search", date.date_start, date.date_end)
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

const dynamicHeight = computed(() => {
  if (width < 600) {
    // Mobile
    return 220
  } else if (width >= 600 && width < 960) {
    // Tablet
    return 190
  } else if (width >= 960 && width <= 1280) {
    // Desktop
    return 300
  } else {
    // Large Desktop
    return 200
  }
})


const chartCardHeight = computed(() => {
  if (width < 600) return 300
  return 650
})

const initData = () => {
  let today = new Date()
  let lastMonth = new Date(today)
  lastMonth.setMonth(today.getMonth() - 1)
  date.date_start = formatDate(lastMonth.toString())
  date.date_end = formatDate(today.toString())
}

watch(
    () => props.immediate,
    (newVal) => {
      if (newVal) {
        if (!date.date_start) initData();
        search()
      }
    },
    {immediate: true}
)


</script>

<template>
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

  <v-divider class="mt-3 mb-5"/>

  <v-row align="center" justify="center">
    <v-col class="ma-0" lg="6" md="6" sm="12" xs="12">
      <v-card :height="chartCardHeight" class="my-0" elevation="4">
        <v-card-title class="bg-primary">
          RX
        </v-card-title>
        <v-card-text>
          <Line
              :data="lineData"
              :height="dynamicHeight"
              :options="chartOptions"
          />
        </v-card-text>
      </v-card>
    </v-col>


    <v-col class="ma-0" lg="6" md="6" sm="12" xs="12">
      <v-card :height="chartCardHeight" class="my-0" elevation="4">
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

<style scoped>

</style>