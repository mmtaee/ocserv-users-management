<script lang="ts" setup>
import {useI18n} from "vue-i18n";
import {defineAsyncComponent, onMounted, reactive} from "vue";
import {type ModelsOcservGroupConfig, OcservGroupsApi} from "@/api";
import {getAuthorization} from "@/utils/request.ts";

const GroupForm = defineAsyncComponent(() => import('@/components/ocserv_group/ConfigForm.vue'));

const {t} = useI18n()
const data = reactive<ModelsOcservGroupConfig>({} as ModelsOcservGroupConfig)

const api = new OcservGroupsApi()

const updateDefaultGroup = () => {
  api.ocservGroupsDefaultsPatch({
    ...getAuthorization(),
    request: {
      config: data
    }
  }).then((res) => {
    // TODO: update group not completed

    console.log(res.data)
  })
}

const getDefaultGroup = () => {
  api.ocservGroupsDefaultsGet({
    ...getAuthorization()
  }).then((res) => {
    Object.assign(data, res.data)
  })
}


onMounted(
    () => {
      getDefaultGroup()
    }
)

</script>

<template>
  <v-card flat>
    <v-card-text>
      <GroupForm
          v-model="data"
          :btnText="t('UPDATE')"
          @save="updateDefaultGroup"
      />
    </v-card-text>
  </v-card>
</template>

