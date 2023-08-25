<template>
  <v-card>
    <v-card-subtitle>
      <v-row align="center" justify="end">
        <v-col md="auto" class="pa-0 ma-0">
          <v-btn @click="init" outlined>
            refresh
            <v-icon right>mdi-refresh</v-icon>
          </v-btn>
        </v-col>
        <v-col md="auto">
          <v-btn
            outlined
            color="error"
            @click="clearLogs"
            :disabled="checkLogsIsCleared"
          >
            Clear Logs
            <v-icon right>mdi-delete</v-icon>
          </v-btn>
        </v-col>
      </v-row>
    </v-card-subtitle>
    <v-divider />
    <v-card-text>
      <div v-for="(line, index) in logs" :key="index" class="my-2">
        <span :style="{ color: rowStyles(line) }">
          <span v-if="logs.length > 1">{{ index + 1 }}-</span> {{ line }}
        </span>
      </div>
    </v-card-text>
  </v-card>
</template>
<script lang="ts">
import { systemServiceApi } from "@/utils/services";
import Vue from "vue";
export default Vue.extend({
  name: "Action",
  data(): {
    logs: string[];
  } {
    return {
      logs: [],
    };
  },

  async mounted() {
    await this.init();
  },

  computed: {
    checkLogsIsCleared() {
      if (this.logs.length && this.logs[0].startsWith("##")) return true;
      return false;
    },
  },

  methods: {
    async init() {
      let data = await systemServiceApi.get_action_logs();
      this.logs = data.logs;
    },
    async clearLogs() {
      await systemServiceApi.clear_action_logs();
      this.logs = [];
    },
    rowStyles(line: string) {
      if (line.startsWith("[Warning]")) return "orange";
      if (line.startsWith("[Critical]")) return "red";
      else return "darkblue";
    },
  },
});
</script>