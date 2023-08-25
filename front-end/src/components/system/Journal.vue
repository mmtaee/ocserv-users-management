<template>
  <div>
    <v-row align="center" justify="start" class="mb-5 mx-1">
      <v-col md="2" align-self="start" class="ma-0 pa-0">
        <v-switch v-model="autoRefresh" inset class="ms-3" color="primary">
          <template v-slot:label>
            <span :class="autoRefresh ? 'primary--text' : 'error--text'">
              Auto Refresh {{ autoRefresh ? `(${intervalTime})` : "(Off)" }}
            </span>
          </template>
        </v-switch>
      </v-col>

      <v-col md="1" align-self="start" class="ma-0 pa-0">
        <v-text-field
          v-model.number="intervalTimeClone"
          label="Refresh Time"
          hint="Refresh Time"
          persistent-hint
          single-line
          :rules="rules.required"
          @change="intervalTime=intervalTimeClone"
        />
      </v-col>

      <v-spacer />

      <v-col md="4" align-self="start" class="ma-0 pa-0">
        <v-row align="center" justify="end">
          <v-col md="auto">
            <v-text-field
              v-model="lines"
              label="Journalctl lines"
              :hint="`Journalctl -n ${lines} ocserv.service`"
              single-line
              persistent-hint
              :rules="rules.required"
            />
          </v-col>

          <v-col md="auto" class="mt-3">
            <v-btn
              small
              outlined
              color="primary"
              @click="getJournal(false, true)"
              :disabled="lines ? false : true"
            >
              Refresh
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <div class="black white--text py-5">
      <v-virtual-scroll height="640" item-height="20" :items="journalData">
        <template v-slot:default="{ index, item }">
          <span class="pa-8">{{ index + 1 }}- {{ item }}</span>
        </template>
      </v-virtual-scroll>
    </div>
  </div>
</template>

<script lang="ts">
import { required } from "@/utils/rules";
import { systemServiceApi } from "@/utils/services";
import Vue from "vue";
export default Vue.extend({
  name: "Journal",
  data(): {
    lines: number;
    cloneLines: number;
    journalData: string[];
    autoRefresh: boolean;
    rules: object;
    intervalTime: number;
    intervalTimeDefault: number;
    intervalId: null | number;
    intervalTimeClone: number;
  } {
    return {
      lines: 100,
      cloneLines: 0,
      journalData: [],
      autoRefresh: true,
      rules: {
        required: [required],
      },
      intervalTimeDefault: 60,
      intervalTime: 0,
      intervalId: null,
      intervalTimeClone: 0,
    };
  },

  mounted() {
    this.intervalTimeClone = this.intervalTime
    this.cloneLines = this.lines
    this.getJournal();
  },

  methods: {
    async getJournal(overlay: boolean = true, reset: boolean = false) {
      let data = await systemServiceApi.journal(this.lines, overlay);
      if (reset) {
        this.cloneLines = this.lines;
        this.journalData = data.logs;
      }
      this.makeSetLogs(data.logs);
    },

    makeSetLogs(logs: string[]) {
      const newElements = logs.filter(
        (item) => !this.journalData.includes(item)
      );
      this.journalData.push(...newElements);
      if (this.journalData.length > this.cloneLines) {
        let extraLines = this.journalData.length - this.cloneLines;
        this.journalData.splice(0, extraLines);
      }
    },
  },

  watch: {
    autoRefresh: {
      immediate: true,
      handler() {
        this.intervalTime = this.intervalTimeDefault;
        if (this.autoRefresh) {
          this.intervalId = setInterval(() => {
            if (this.intervalTime > 0) {
              this.intervalTime--;
            } else {
              this.getJournal(false);
              this.intervalTime = this.intervalTimeDefault;
            }
          }, 1000);
        } else {
          clearInterval(this.intervalId!);
        }
      },
    },
  },

  beforeDestroy() {
    clearInterval(this.intervalId!);
  },
});
</script>