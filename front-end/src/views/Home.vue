<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="850"
        >
          <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
            Home
          </v-card-subtitle>
          <v-card-text></v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { adminServiceApi } from "@/utils/services";
import { Dashboard } from "@/utils/types";
import Vue from "vue";

export default Vue.extend({
  name: "Home",
  data() {
    return {};
  },

  async mounted() {
    let data: Dashboard = await adminServiceApi.dashboard();
    console.log(data);
  },
});
</script>
