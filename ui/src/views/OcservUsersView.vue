<script lang="ts" setup>

import {defineAsyncComponent, reactive, ref} from "vue";
import {useLocale} from "vuetify/framework";
import type {ModelsOcservUser} from "@/api";


const CreateOrEdit = defineAsyncComponent(() => import('@/components/ocserv_user/CreateOrUpdate.vue'));

const {t} = useLocale()

const createDialog = ref(true)
const users = reactive<ModelsOcservUser[]>([])


const completeCreate = (data: ModelsOcservUser) => {
  console.log(data)
}

</script>

<template>
  <v-row>
    <v-col>
      <v-card min-height="850">
        <v-toolbar :title="t('OCSERV_USERS')" color="primary"/>
        <v-card flat>
          <v-card-text>
            <v-row align="center" justify="center">
              <v-col cols="12" lg="10" md="10" sm="10">
                <span class="text-capitalize text-subtitle-1">{{ t("USERS") }}</span>
              </v-col>
              <v-col cols="12" lg="auto" md="2" sm="2">
                <v-btn
                    color="primary"
                    variant="outlined"
                    @click="createDialog = true"
                >
                  {{ t("CREATE") }}
                </v-btn>
              </v-col>

              <v-divider/>

              <v-col class="mt-3" cols="12" md="12">
                <v-row>
                  <v-col
                      v-for="(item, index) in users"
                      :key="`ocserv-users-${index}`"
                      cols="12"
                      lg="3"
                      md="4"
                      sm="6"
                  >
                    <v-card class="py-3" elevation="6">
                      <v-card-title class="text-subtitle-1">
                        <v-row align="center" justify="center">
                          <v-col cols="12" md="9" sm="8">{{ item.username }}</v-col>
                        </v-row>
                      </v-card-title>
                    </v-card>
                  </v-col>
                </v-row>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-card>
    </v-col>
  </v-row>

  <!-- Create Dialog -->
  <CreateOrEdit
      v-model="createDialog"
      @complete="completeCreate"
  />

</template>
