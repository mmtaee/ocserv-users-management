<script lang="ts" setup>
import {computed, defineAsyncComponent, reactive, ref, watch} from "vue";
import {requiredRule} from "@/utils/rules.ts";
import type {ModelsOcservGroup, ModelsOcservGroupConfig} from "@/api";
import {useLocale} from "vuetify/framework";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))
const GroupForm = defineAsyncComponent(() => import('@/components/ocserv_group/ConfigForm.vue'));


const props = defineProps<{
  modelValue: boolean
  initValue?: ModelsOcservGroup
}>()


const emit = defineEmits(["update:modelValue", "complete"])

const {t} = useLocale()
const valid = ref(true)
const editMode = ref(false)
const validConfig = ref(true)
const name = ref("")
const data = reactive<ModelsOcservGroupConfig>({})
const rules = {
  required: (v: string) => requiredRule(v, t),
}

const validChecker = computed(() => {
  if (editMode.value) {
    return !validConfig.value
  }
  return !valid.value || !validConfig.value
})

const save = () => {
  emit("complete", data, name)
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
        name.value = newVal.name
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
      <span class="text-capitalize">{{ editMode ? t("UPDATE_OCSERV_GROUP") : t("CREATE_OCSERV_GROUP") }}</span>
    </template>


    <template #dialogText>
      <v-form v-model="valid">
        <v-row align="center" justify="start">

          <v-col class="px-4" cols="12" md="12">
            <v-row align="center" dense justify="start">

              <v-col class="ma-0 pa-0" cols="12" md="11">
                <h3 class="text-capitalize">{{ editMode ? t("UPDATE") : t("NAME") }}</h3>
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

              <v-divider/>

              <v-col class="ma-0 pa-0 mt-5" cols="12" md="3">
                <v-text-field
                    v-model="name"
                    :label="t('GROUP_NAME')"
                    :readonly="editMode"
                    :rules="getRules"
                    density="comfortable"
                    persistent-hint
                    type="text"
                    variant="underlined"
                />
              </v-col>
            </v-row>
          </v-col>

          <v-col cols="12" md="12">
            <GroupForm
                v-if="props.modelValue"
                v-model="data"
                hideBtn
                @valid="(v: boolean) => validConfig = v"
            />
          </v-col>

        </v-row>
      </v-form>

    </template>

  </ReusableDialog>
</template>

