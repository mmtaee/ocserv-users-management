<script lang="ts" setup>
import {reactive, ref, watch} from "vue";
import type {ModelsOcservUserConfig} from "@/api";
import {useI18n} from "vue-i18n";
import {domainRule, ipOrRangeRule, ipRule, ipWithRangeRule} from "@/utils/rules.ts";


const props = withDefaults(defineProps<{
  modelValue: ModelsOcservUserConfig
  btnText?: string
  btnColor?: string
  hideBtn?: boolean
}>(), {
  btnText: 'Save',
  btnColor: 'primary',
  hideBtn: false
})

const emit = defineEmits(["update:modelValue", "save", "valid"])

const {t} = useI18n()
const valid = ref(true)
const cloneFormData = ref<ModelsOcservUserConfig>()

const rules = {
  ip: (v: string) => ipRule(v, t),
  ipOrRange: (v: string) => ipOrRangeRule(v, t),
  domain: (v: string) => domainRule(v, t),
  ipWithRange: (v: string) => ipWithRangeRule(v, t)
}

const fields = [
  // Network Configuration
  {
    key: 'nbns', label: 'NBNS', type: 'text', hint: "Net BIOS",
    example: "192.168.1.10",
    rules: [rules.ip]
  },
  {
    key: 'ipv4-network',
    label: 'IPv4 Network',
    type: 'text',
    hint: 'CIDR',
    example: "192.168.0.0/24",
    rules: [rules.ipWithRange]
  },
  {
    key: 'explicit-ipv4',
    label: 'Explicit IPv4',
    type: 'text',
    hint: t('SPECIFIC_IP_ADDRESS'),
    example: "192.168.1.5",
    rules: [rules.ip]
  },
  {
    key: 'iroute',
    label: 'Internal Route',
    type: 'text',
    hint: t('CUSTOM_INTERNAL_ROUTE'),
    example: " 10.0.0.0/8 ",
    rules: [rules.ipOrRange]
  },
  {
    key: 'restrict-user-to-ports', label: 'Restrict User To Ports', type: 'text', hint: t('ALLOWED_PORTS'),
    example: "80,443",
  },

  // Performance and Session Settings
  {key: 'idle-timeout', label: 'Idle Timeout', type: 'number', hint: t("INACTIVITY_TIMEOUT_S")},
  {key: 'mobile-idle-timeout', label: 'Mobile Idle Timeout', type: 'number', hint: t("MOBILE_INACTIVITY_TIMEOUT_S")},
  {key: 'session-timeout', label: 'Session Timeout', type: 'number', hint: t("MAX_SESSION_DURATION_S")},
  {
    key: 'rekey-time',
    label: 'Rekey Time', type: 'number',
    hint: t("TRIGGERS_KEY_RENEGOTIATION"),
    example: "86400 for 24 hours",
  },

  // Access and Feature Controls
  {
    key: 'restrict-user-to-routes',
    label: 'Restrict User To Routes',
    type: 'switch',
    hint: t("ALLOW_CLIENT_ACCESS_ONLY_TO_DEFINED_ROUTES")
  },
  {key: 'banner', label: 'Banner', type: 'text', hint: t("BANNER_HINT")},
]

const textFields = [
  // Routes
  {
    key: 'route',
    label: 'Route',
    type: 'text',
    example: "10.0.0.0/8",
    hint: t("ROUTES_ASSIGNED_TO_CLIENT"),
    rules: [rules.ipOrRange]
  },
  {
    key: 'no-route',
    label: 'No Route',
    type: 'text',
    hint: t("NON_VPN_NETWORKS"),
    example: "172.16.0.0/12",
    rules: [rules.ipOrRange]
  },
  {
    key: 'dns',
    label: 'DNS',
    type: 'text',
    hint: t("DNS_SERVERS_LIST"),
    example: "8.8.8.8/example.com",
    rules: [rules.ip]
  },
  {
    key: 'split-dns',
    label: 'Split DNS',
    type: 'text',
    hint: t("DNS_SPECIFIC_DOMAINS"),
    example: "example.com",
    rules: [rules.domain]
  }
]


const chipInputs = reactive<Record<string, string>>({
  dns: '',
  route: '',
  'no-route': '',
  'split-dns': '',
})

const addRoutes = (key: string) => {
  const typedKey = key as keyof ModelsOcservUserConfig;
  const input = chipInputs[typedKey];

  if (input) {
    if (!Array.isArray(props.modelValue[typedKey])) {
      props.modelValue[typedKey] = [] as any;
    }

    const arr = props.modelValue[typedKey] as string[];

    if (!arr.includes(input)) {
      arr.push(input);
      chipInputs[typedKey] = '';
    }
  }

}

const removeRoute = (key: string, value: string) => {
  const typedKey = key as keyof ModelsOcservUserConfig
  const arr = props.modelValue[typedKey] as string[]

  let index = arr.findIndex(i => i == value)
  if (index !== -1) {
    arr.splice(index, 1)
  }
}

// const checkValid = computed(() => {
//   const isChanged = JSON.stringify(props.modelValue) !== JSON.stringify(cloneFormData.value)
//   return valid.value && isChanged
// })


watch(
    () => props.modelValue,
    (newVal) => {
      cloneFormData.value = {...newVal}
    },
    {immediate: true, once: true}
)

watch(
    () => valid.value,
    (newVal) => {
      emit('valid', newVal)
    },
)

</script>

<template>
  <v-form v-model="valid">
    <v-row>
      <v-col cols="12">
        <v-row align="center" justify="start">
          <v-col cols="12" md="11">
            <h3 class="text-capitalize">{{ t("NETWORK_CONFIGURATION") }}</h3>
          </v-col>
        </v-row>

        <v-divider/>
      </v-col>

      <template v-for="field in fields.filter(f => f.type === 'text')" :key="field.key">
        <v-col
            cols="12"
            lg="2" md="4"
        >
          <v-text-field
              v-model="props.modelValue[field.key as keyof ModelsOcservUserConfig]"
              :hint="field.hint"
              :label="field.label"
              :rules="field.rules"
              density="comfortable"
              persistent-hint
              type="text"
              variant="underlined"
          />
        </v-col>
      </template>

      <!-- Number Fields Section -->
      <v-col class="mt-6" cols="12">
        <h3 class="text-capitalize">{{ t("T07") }}</h3>
        <v-divider/>
      </v-col>

      <template v-for="field in fields.filter(f => f.type === 'number')" :key="field.key">
        <v-col
            cols="12"
            lg="3"
            md="4"
            xl="2"
        >
          <v-number-input
              v-model="props.modelValue[field.key as keyof ModelsOcservUserConfig] as number"
              :hint="field.hint"
              :label="field.label"
              control-variant="hidden"
              density="comfortable"
              persistent-hint
              variant="underlined"
          />
        </v-col>
      </template>

      <!-- Switch Fields Section -->
      <v-col class="mt-6" cols="12">
        <h3 class="text-capitalize">{{ t("T08") }}</h3>
        <v-divider/>
      </v-col>

      <template v-for="field in fields.filter(f => f.type === 'switch')" :key="field.key">
        <v-col
            cols="12"
            md="3"
        >
          <v-row align="center" justify="center">
            <v-col cols="6" md="12">
              <v-switch
                  v-model="props.modelValue[field.key as keyof ModelsOcservUserConfig]"
                  :hint="field.hint"
                  :label="field.label"
                  class="ms-1"
                  color="primary"
                  density="compact"
                  persistent-hint
              />
            </v-col>
          </v-row>

        </v-col>
      </template>
    </v-row>
  </v-form>

  <v-row align="center" justify="center">
    <!-- New TextFields with chips section -->
    <v-col class="mt-10" cols="12">
      <h3 class="text-capitalize">{{ t("ROUTES") }}</h3>
      <v-divider/>
    </v-col>

    <template v-for="field in textFields" :key="field.key">
      <v-col lg="3" md="6" sm="12">
        <v-card min-height="300">

          <v-toolbar class="text-subtitle-1 px-3" color="primary" density="compact">
            {{ field.label }}
          </v-toolbar>

          <v-card-title>
            <v-row align="start" justify="start">
              <v-col cols="12" md="9">
                <v-text-field
                    v-model="chipInputs[field.key]"
                    :hint="field.hint"
                    :placeholder="field.example"
                    :rules="field.rules"
                    density="comfortable"
                    persistent-hint
                    variant="underlined"
                    @keydown.enter="addRoutes(field.key)"
                />
              </v-col>

              <v-col cols="12" md="2">
                <v-btn
                    class="mt-5"
                    color="success"
                    density="compact"
                    variant="outlined"
                    @click="addRoutes(field.key)"
                >
                  {{ t("ADD") }}
                </v-btn>
              </v-col>
            </v-row>
          </v-card-title>

          <v-card-text>
            <v-chip
                v-for="chip in props.modelValue[field.key as keyof ModelsOcservUserConfig]"
                :key="`${field.key}-${chip}`"
                class="me-2"
                color="primary"
            >
              {{ chip }}
              <v-icon color="error" end @click="removeRoute(field.key, chip)">mdi-delete</v-icon>
            </v-chip>
          </v-card-text>

        </v-card>
      </v-col>
    </template>
  </v-row>

</template>
