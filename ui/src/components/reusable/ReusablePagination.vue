<script lang="ts" setup>
import {ref, watch} from "vue";
import {useLocale} from "vuetify/framework";

const props = defineProps({
  totalRecords: {
    type: Number,
    default: 0
  },
  pageCount: Number,
  page: Number,
  pageSize: {
    type: Number,
    default: 5
  },
})

const emit = defineEmits(["reFetch"])

const {t} = useLocale()
const itemPerPage = ref(5)
const pageNumber = ref(1)
const desc = ref(false)

watch(
    () => props.page,
    (newVal) => {
      pageNumber.value = newVal || 1
    }
)

watch(
    () => props.pageSize,
    (newVal) => {
      itemPerPage.value = newVal || 5
    },
    {immediate: true}
)

</script>

<template>
  <v-divider class="mt-5"/>
  <div v-if="totalRecords>1" class="pt-2">
    <v-row align="center" class="ma-0 pa-0" justify="center">

      <v-col class="ma-0 pa-0 mx-5 me-5" cols="12" lg="2" md="2" sm="2">
        <v-select
            v-model="itemPerPage"
            :items="[5,10,25,50,100]"
            :label="t('ITEMS_PER_PAGE')"
            density="comfortable"
            hide-details
            variant="underlined"
            @update:modelValue='emit("reFetch",pageNumber,  itemPerPage, desc)'
        />
      </v-col>


      <v-col class="ma-0 pa-0 mx-5" cols="12" lg="auto" md="auto" sm="auto">
        <v-pagination
            v-model="pageNumber"
            :length="pageCount"
            :total-visible="5"
            @update:modelValue='emit("reFetch",  pageNumber,  itemPerPage, desc)'
        />
      </v-col>

      <v-spacer/>

      <v-col class="ma-0 pa-0 mx-5 me-5" cols="12" lg="auto" md="auto" sm="auto">
        <v-checkbox
            v-model="desc"
            :label="t('DESCENDING_SORT')"
            color="primary"
            density="compact"
            hide-details
            variant="underlined"
            @update:modelValue='emit("reFetch",  pageNumber,  itemPerPage, desc)'
        />
      </v-col>

    </v-row>
  </div>
</template>
