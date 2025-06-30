<script lang="ts" setup>
import {computed, defineAsyncComponent, reactive, ref} from "vue";
import {requiredRule} from "@/utils/rules.ts";
import type {OcservGroupCreateOcservGroupData} from "@/api";
import {useLocale} from "vuetify/framework";

const ReusableDialog = defineAsyncComponent(() => import("@/components/reusable/ReusableDialog.vue"))
const GroupForm = defineAsyncComponent(() => import('@/components/forms/Group.vue'));

const props = defineProps({
  modelValue: Boolean
})


const emit = defineEmits(["update:modelValue", "complete"])


const {t} = useLocale()
const valid = ref(true)
const validConfig = ref(true)
const createData = reactive<OcservGroupCreateOcservGroupData>({config: {}, name: ""})
const rules = {
  required: (v: string) => requiredRule(v, t),
}

const validCreateChecker = computed(() => {
  return !valid.value || !validConfig.value
})

const createGroup = () => {
  // TODO: get group api
  console.log("create group data", createData)
  emit("complete")
}

</script>

<template>
  <ReusableDialog
      v-model="props.modelValue"
      btnClose
      fullscreen
      persistent
      @update:modelValue="emit('update:modelValue', false)"
  >

    <template #dialogTitle>
      <span class="text-capitalize">{{ t("CREATE_OCSERV_GROUP") }}</span>
    </template>


    <template #dialogText>
      <v-form v-model="valid">

        <v-row align="center" justify="start">
          <v-col cols="12" md="12">

            <v-row align="center" dense justify="start">
              <v-col class="ma-0 pa-0" cols="12" md="11">
                <v-col class="ma-0 pa-0" cols="12" md="4">
                  <v-text-field
                      v-model="createData.name"
                      :label="t('GROUP_NAME')"
                      :rules="[rules.required]"
                      density="comfortable"
                      persistent-hint
                      type="text"
                      variant="underlined"
                  />
                </v-col>
              </v-col>

              <v-col cols="12" md="auto">
                <v-btn
                    :disabled="validCreateChecker"
                    color="primary"
                    @click="createGroup"
                >
                  {{ t("CREATE") }}
                </v-btn>
              </v-col>
            </v-row>

          </v-col>

          <v-col cols="12" md="12">
            <GroupForm
                v-model="createData.config"
                hideBtn
                @valid="(v: boolean) => validConfig = v"
            />
          </v-col>

        </v-row>
      </v-form>

    </template>

  </ReusableDialog>
</template>

<style scoped>

</style>