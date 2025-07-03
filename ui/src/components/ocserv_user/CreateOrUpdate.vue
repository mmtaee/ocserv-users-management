<script lang="ts" setup>
import {computed, defineAsyncComponent, onBeforeMount, reactive, ref, watch} from "vue";
import {type ModelsOcservGroup, type ModelsOcservUser, ModelsOcservUserTrafficTypeEnum, OcservGroupsApi} from "@/api";
import {useLocale} from "vuetify/framework";
import {requiredRule} from "@/utils/rules.ts";
import {getAuthorization} from "@/utils/request.ts";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  initValue?: ModelsOcservGroup
}>()

const emit = defineEmits(["update:modelValue", "complete"])

const {t} = useLocale()
const valid = ref(true)
const editMode = ref(false)
const validConfig = ref(true)
const groups = ref<string[]>([])
const trafficTypes = ref([
  {
    label: t('FREE'),
    value: ModelsOcservUserTrafficTypeEnum.FREE,
  },
  {
    label: t('MONTHLY_TRANSMIT'),
    value: ModelsOcservUserTrafficTypeEnum.MONTHLY_TRANSMIT,
  },
  {
    label: t('MONTHLY_RECEIVE'),
    value: ModelsOcservUserTrafficTypeEnum.MONTHLY_RECEIVE,
  },
  {
    label: t('TOTALLY_RECEIVE'),
    value: ModelsOcservUserTrafficTypeEnum.TOTALLY_RECEIVE,
  },
  {
    label: t('TOTALLY_TRANSMIT'),
    value: ModelsOcservUserTrafficTypeEnum.TOTALLY_TRANSMIT,
  }
])


const data = reactive<ModelsOcservUser>({
  created_at: "",
  group: "defaults",
  is_locked: false,
  is_online: false,
  password: "",
  rx: 0,
  traffic_size: 10,
  traffic_type: ModelsOcservUserTrafficTypeEnum.TOTALLY_TRANSMIT,
  tx: 0,
  uid: "",
  username: ""
})

const rules = {
  required: (v: string) => requiredRule(v, t),
}

const save = () => {
  emit("complete", data)
}

const getRules = computed(() => {
  if (editMode.value) {
    return []
  }
  return [rules.required]
})

const clearData = () => {
  for (const key in data) {
    delete data[key as keyof typeof data]
  }
}


watch(
    () => props.initValue,
    (newVal) => {
      if (newVal) {
        clearData()
        let cloneVal = JSON.parse(JSON.stringify(newVal))
        Object.assign(data, cloneVal.config)
        editMode.value = true
      }
    },
    {immediate: true, deep: true}
)

watch(
    () => props.modelValue,
    () => {
      if (!props.initValue) {
        clearData()
      }
    }
)

onBeforeMount(
    () => {
      const api = new OcservGroupsApi()
      api.ocservGroupsLookupGet({
        ...getAuthorization()
      }).then((res) => {
        groups.value = res.data
        console.log(" res.data: ", res.data)
      })
    }
)

</script>

<template>

  <ReusableDialog
      v-model="props.modelValue"
      btnClose
      fullscreen
      persistent
      @update:modelValue="emit('update:modelValue', false);Object.assign(data, {}) "
  >

    <template #dialogTitle>
      <span class="text-capitalize">{{ editMode ? t("UPDATE_OCSERV_USER") : t("CREATE_OCSERV_USER") }}</span>
    </template>


    <template #dialogText>
      <v-row align="center" class="mb-5" dense justify="center">
        <v-col cols="12" md="10">
          <h3 class="text-capitalize">{{ editMode ? t("UPDATE") : t("NAME") }}</h3>
        </v-col>

        <v-col cols="12" md="1">
          <v-btn
              :disabled="!valid"
              color="primary"
              variant="outlined"
              @click="save"
          >
            {{ editMode ? t("UPDATE") : t("CREATE") }}
          </v-btn>
        </v-col>

        <v-divider/>
      </v-row>

      <v-form v-model="valid">


        <v-row align="center" dense justify="start">

          <v-col cols="12" lg="2" md="4">
            <v-text-field
                v-model="data.username"
                :label="t('USERNAME')"
                :readonly="editMode"
                :rules="[rules.required]"
                density="comfortable"
                variant="underlined"
            />
          </v-col>

          <v-col cols="12" lg="2" md="4">
            <v-select
                v-model="data.group"
                :items="groups"
                :label="t('GROUP')"
                :rules="[rules.required]"
                density="comfortable"
                variant="underlined"
            />
          </v-col>


          <v-col cols="12" lg="2" md="4">
            <v-text-field
                v-model="data.password"
                :label="t('PASSWORD')"
                :rules="[rules.required]"
                density="comfortable"
                variant="underlined"
            />
          </v-col>

          <v-col cols="12" lg="2" md="4">
            <v-select
                v-model="data.traffic_type"
                :items="trafficTypes"
                :label="t('TRAFFIC_TYPE')"
                :rules="[rules.required]"
                density="comfortable"
                item-title="label"
                item-value="value"
                variant="underlined"
                @update:modelValue="(v)=>v ==ModelsOcservUserTrafficTypeEnum.FREE ? data.traffic_size = 0 : false"
            />
          </v-col>

          <v-col cols="12" lg="2" md="4">
            <v-number-input
                v-model="data.traffic_size"
                :items="trafficTypes"
                :label="t('TRAFFIC_SIZE')"
                :readonly="data.traffic_type == ModelsOcservUserTrafficTypeEnum.FREE"
                :rules="ModelsOcservUserTrafficTypeEnum.FREE ? []: [rules.required]"
                control-variant="hidden"
                density="comfortable"
                suffix="GB"
                variant="underlined"
            />
          </v-col>

        </v-row>
      </v-form>

    </template>

  </ReusableDialog>

</template>
