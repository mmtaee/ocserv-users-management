<script lang="ts" setup>
import {useLocale} from "vuetify/framework";
import {defineAsyncComponent, onMounted, reactive} from "vue";
import type {ModelsOcservGroupConfig} from "@/api";
import {dummyGroupConfig} from "@/utils/dummy.ts";

const GroupForm = defineAsyncComponent(() => import('@/components/ocserv_group/Form.vue'));

const {t} = useLocale()
const data = reactive<ModelsOcservGroupConfig>({} as ModelsOcservGroupConfig)

const updateDefaultGroup = () => {
  // TODO: call update api
  console.log("form update: ", data)
}

onMounted(
    () => {
      //   TODO: call get default api
      Object.assign(data, dummyGroupConfig)
    }
)

</script>

<template>
  <v-card flat>
    <v-card-text>
      <GroupForm
          v-model="data"
          :btnText="t('SAVE')"
          @save="updateDefaultGroup"
      />
    </v-card-text>
  </v-card>
</template>

