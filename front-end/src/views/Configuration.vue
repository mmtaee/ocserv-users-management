<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <ConfigsForm :initInput="initInput" editMode />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { adminServiceApi } from "@/utils/services";
import { AdminConfig } from "@/utils/types";

export default Vue.extend({
  name: "Configuration",
  components: {
    ConfigsForm: () => import("@/components/ConfigsForm.vue"),
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
      let data: AdminConfig = await adminServiceApi.get_configuration();
      this.initInput = data;
    },
  },
});
</script>