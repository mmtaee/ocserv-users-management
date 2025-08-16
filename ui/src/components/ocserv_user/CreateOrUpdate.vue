<script lang="ts" setup>
import {computed, defineAsyncComponent, reactive, ref, watch} from "vue";
import {type ModelsOcservUser, type ModelsOcservUserConfig, ModelsOcservUserTrafficTypeEnum,} from "@/api";
import {useI18n} from "vue-i18n";
import {requiredRule} from "@/utils/rules.ts";
import {formatDate} from "@/utils/convertors.ts";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))
const UserForm = defineAsyncComponent(() => import("@/components/ocserv_user/ConfigForm.vue"))


const props = defineProps<{
  modelValue: boolean
  initValue?: ModelsOcservUser
  groups: string[]
}>()

const emit = defineEmits(["update:modelValue", "complete"])

const {t} = useI18n()
const valid = ref(true)
const editMode = ref(false)
const validConfig = ref(true)
const dateMenu = ref(false)
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
  username: "",
  expire_at: ""
})

const config = reactive<ModelsOcservUserConfig>({})

const rules = {
  required: (v: string) => requiredRule(v, t),
}

const save = () => {
  emit("complete", data, config)
}

const clearData = () => {
  for (const key in data) {
    delete data[key as keyof typeof data]
  }
}

const validChecker = computed(() => {
  if (editMode.value) {
    return !validConfig.value
  }
  return !valid.value || !validConfig.value
})


watch(
    () => props.initValue,
    (newVal) => {
      if (newVal) {
        clearData()
        let cloneVal = JSON.parse(JSON.stringify(newVal))
        Object.assign(data, cloneVal)
        Object.assign(config, data.config)
        if (data.expire_at) {
          data.expire_at = formatDate(data.expire_at)
        }
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
      <v-form v-model="valid">
        <v-row align="center" justify="start">
          <v-col cols="12">
            <v-row align="center" justify="start">
              <v-col cols="12" md="11">
                <h3 class="text-capitalize">{{ t("USER_SETUP") }}</h3>
              </v-col>

              <v-col cols="12" md="auto">
                <v-btn
                    :disabled="validChecker"
                    color="primary"
                    variant="outlined"
                    @click="save"
                >
                  {{ editMode ? t("UPDATE") : t("CREATE") }}
                </v-btn>
              </v-col>
            </v-row>

            <v-divider/>
          </v-col>

          <v-col class="px-4" cols="12" md="12">
            <v-row align="center" justify="start">

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
                    :rules="data.traffic_type == ModelsOcservUserTrafficTypeEnum.FREE ? []: [rules.required]"
                    control-variant="hidden"
                    density="comfortable"
                    suffix="GB"
                    variant="underlined"
                />
              </v-col>
              <v-col cols="12" lg="2" md="4">
                <v-menu
                    v-model="dateMenu"
                    :close-on-content-click="false"
                    transition="scale-transition"
                >
                  <template #activator="{ props }">
                    <v-text-field
                        v-model="data.expire_at"
                        :hint="t('SET_EXPIRE_HINT')"
                        :label="t('EXPIRE_AT')"
                        class="mt-7"
                        clearable
                        density="compact"
                        persistent-hint
                        readonly
                        v-bind="props"
                        variant="underlined"
                    />
                  </template>
                  <v-date-picker
                      v-model="data.expire_at"
                      :header="t('EXPIRE_AT')"
                      elevation="24"
                      title=""
                      @update:model-value="data.expire_at=formatDate(data.expire_at); dateMenu = false"
                  />
                </v-menu>
              </v-col>
            </v-row>
          </v-col>

          <v-col cols="12" md="12">
            <UserForm
                v-if="props.modelValue"
                v-model="config"
                hideBtn
                @valid="(v: boolean) => validConfig = v"
            />
          </v-col>

        </v-row>
      </v-form>
    </template>
  </ReusableDialog>

</template>
