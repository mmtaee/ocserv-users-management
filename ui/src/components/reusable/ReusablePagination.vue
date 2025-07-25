<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import type {Meta} from "@/utils/interfaces.ts";
import {computed, ref, watch} from "vue";


const props = defineProps<{
  modelValue: Meta;
}>();

const emit = defineEmits(["update:modelValue"]);


const {t} = useLocale()
const totalRecords = ref(0)
const sort = ref(false)
const pages = computed(() => {
  const total = props.modelValue.total_records ?? 0;
  const size = props.modelValue.size || 1;
  return Math.ceil(total / size);
});

const refresh = () => {
  props.modelValue.sort = sort.value ? "DESC" : "ASC"
  emit("update:modelValue", props.modelValue)
}


watch(
    () => props.modelValue,
    (newVal) => {
      if (newVal) totalRecords.value = props.modelValue.total_records
    },
    {immediate: true, deep: true}
)

</script>

<template>

  <v-footer
      v-if="totalRecords >1"
      class="justify-space-between text-body-2 mt-4"
      color="surface-variant"
  >

    <v-row align="center" justify="center">

      <v-col class="ma-0" cols="12">
        <span>{{ t("TOTAL_RECORDS") }}: {{ totalRecords }}</span>
      </v-col>

      <v-divider opacity="10"/>

      <v-col cols="12" md="2">
        <v-select
            v-model="modelValue.size"
            :items="[5,10,25,50,100]"
            :label="t('ITEMS_PER_PAGE')"
            density="comfortable"
            hide-details
            variant="underlined"
            @update:modelValue='refresh'
        />

      </v-col>

      <v-col cols="12" md="8">
        <v-pagination
            v-model="modelValue.page"
            :length="pages"
            :total-visible="5"
            @update:modelValue='refresh'
        />
      </v-col>


      <v-col cols="12" md="2">
        <v-checkbox
            v-model="sort"
            :label="t('DESCENDING_SORT')"
            color="primary"
            density="compact"
            hide-details
            variant="underlined"
            @update:modelValue='refresh'
        />
      </v-col>

    </v-row>
  </v-footer>
</template>
