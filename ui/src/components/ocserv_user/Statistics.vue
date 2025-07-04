<script lang="ts" setup>
import {defineAsyncComponent, reactive, ref} from "vue";
import type {ModelsOcservUser} from "@/api";
import {useLocale} from "vuetify/framework";
import {formatDate} from "@/utils/convertors.ts";

type Date = {
  dateStart: String | undefined;
  dateEnd: String | undefined;
}


const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))

const props = defineProps<{
  modelValue: boolean
  user: ModelsOcservUser
}>()

const emit = defineEmits(["update:modelValue", "done"])

const {t} = useLocale()
const dateStartMenu = ref(false)
const dateEndMenu = ref(false)
const today = new Date()

const lastMonth = new Date(today)
lastMonth.setMonth(today.getMonth() - 1)

const date = reactive<Date>({
  dateStart: formatDate(lastMonth.toString()),
  dateEnd: formatDate(today.toString()),
})


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
      <v-row>
        <v-col cols="12" lg="2" md="4">
          <v-menu
              v-model="dateStartMenu"
              :close-on-content-click="false"
              transition="scale-transition"
          >
            <template #activator="{ props }">
              <v-text-field
                  v-model="date.dateStart"
                  :label="t('DATE_START')"
                  clearable
                  density="compact"
                  readonly
                  v-bind="props"
                  variant="underlined"
              />
            </template>
            <v-date-picker
                v-model="date.dateStart"
                :header="t('DATE_START')"
                elevation="24"
                title=""
                @update:model-value="(val:any)=>{date.dateStart = formatDate(val); dateStartMenu = false}"
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
                  v-model="date.dateEnd"
                  :label="t('EXPIRE_AT')"
                  clearable
                  density="compact"
                  readonly
                  v-bind="props"
                  variant="underlined"
              />
            </template>
            <v-date-picker
                v-model="date.dateEnd"
                :header="t('EXPIRE_AT')"
                elevation="24"
                title=""
                @update:model-value="(val:any)=>{date.dateEnd=formatDate(val); dateEndMenu = false}"
            />
          </v-menu>
        </v-col>
      </v-row>


    </template>
  </ReusableDialog>
</template>

