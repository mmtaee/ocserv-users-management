<script lang="ts" setup>
import {useI18n} from "vue-i18n";
import router from "@/plugins/router.ts";


const {t} = useI18n()

</script>

<template>
  <v-card color="warning">
    <v-card-text class="text-white">
      {{ t("T03") }}
    </v-card-text>
    <v-card-actions>
      <v-spacer/>
      <v-btn
          class="mt-5"
          color="white"
          variant="outlined"
          @click="router.push('/')"
      >
        {{ t("GO_BACK_HOME") }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<style scoped>

</style>