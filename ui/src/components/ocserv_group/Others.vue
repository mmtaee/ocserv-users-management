<script lang="ts" setup>
import {dummyGroupList} from "@/utils/dummy.ts";
import {defineAsyncComponent, onMounted, reactive, ref} from "vue";
import type {ModelsOcservGroup} from "@/api";
import {useLocale} from "vuetify/framework";

const Create = defineAsyncComponent(() => import('@/components/ocserv_group/Create.vue'));

const {t} = useLocale()


const otherGroups = reactive<ModelsOcservGroup[]>([])
const createDialog = ref(false)
const editDialog = ref(false)
const deleteDialog = ref(false)


const getGroups = async () => {
  console.log("getGroups")
}

const complete = () => {
  createDialog.value = false
  console.log("complete")
  getGroups()
}

onMounted(() => {
  // Fetch and assign group list
  // TODO: call group api list
  getGroups()
  Object.assign(otherGroups, dummyGroupList)
})

</script>

<template>
  <v-card flat>
    <v-card-text>
      <v-row align="center" dense justify="start">
        <v-col cols="12" md="12">
          <v-row>
            <v-col cols="12" md="11">
              <span class="text-capitalize text-subtitle-1">{{ t("OTHER") }} {{ t("GROUPS") }}</span>
              <v-tooltip location="top">
                <template #activator="{ props }">
                  <v-icon
                      class="ms-2"
                      color="primary"
                      icon="mdi-plus-circle-outline"
                      size="x-large"
                      v-bind="props"
                      @click="createDialog = true"
                  />
                </template>
                <span>{{ t("CREATE") }}</span>
              </v-tooltip>

            </v-col>
          </v-row>
        </v-col>

        <v-divider/>
        <v-col class="mt-3" cols="12" md="12">
          <v-row>
            <v-col
                v-for="(item, index) in otherGroups"
                :key="`other-groups-${index}`"
                cols="12"
                md="2"
            >
              <v-card elevation="6">
                <v-card-title class="text-subtitle-1 bg-primary-darken-1">
                  <v-row align="start" justify="start">

                    <v-col cols="12" md="9"> {{ item.name }}</v-col>

                    <v-col class="bg-white" cols="12" md="3">
                      <v-icon color="odd" size="small">mdi-pencil</v-icon>
                      <v-icon color="error" size="small">mdi-delete</v-icon>
                    </v-col>
                  </v-row>

                </v-card-title>
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>

  <Create
      v-if="createDialog"
      v-model="createDialog"
      @complete="complete"
  />
</template>
