<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <CofigsForm :initInput="initInput" editMode />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import httpRequest from "@/plugins/axios";
import { AdminConfig } from "@/utils/types";
import { AxiosResponse } from "axios";
import Vue from "vue";

export default Vue.extend({
  name: "Configuration",
  components: {
    CofigsForm: () => import("@/components/CofigsForm.vue"),
  },
  data(): {
    initInput: AdminConfig | null;
  } {
    return {
      initInput: null,
    };
  },

  async mounted() {
    await this.getInit();
  },

  methods: {
    async getInit() {
      let res: AxiosResponse = await httpRequest("get", {
        urlName: "admin",
        urlPath: "configuration",
      });
      if (res.data) this.initInput = { ...res.data };
    },
  },
});
</script>