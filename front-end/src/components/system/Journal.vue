<template>
  <div>
    <v-row align="center" justify="start" class="mb-5">
      <v-col md="2" align-self="start" class="ma-0 pa-0">
        <v-switch v-model="autoRefresh" inset class="ms-3" color="primary">
          <template v-slot:label>
            <span :class="autoRefresh ? 'primary--text' : 'error--text'">
              Auto Refresh {{ autoRefresh ? `(${intervalTime})` : "(Off)" }}
            </span>
          </template>
        </v-switch>
      </v-col>

      <v-spacer />

      <v-col md="4" align-self="start" class="ma-0 pa-0">
        <v-row align="center" justify="end">
          <v-col md="8">
            <v-text-field
              v-model="lines"
              label="Journalctl lines"
              hint="Journalctl -n ${lines} ocserv.service, blank=20"
              single-line
              persistent-hint
              :rules="rules.required"
            />
          </v-col>

          <v-col md="auto">
            <v-btn
              small
              outlined
              color="primary"
              @click="getJournal()"
              :disabled="lines ? false : true"
            >
              Get
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <div v-if="journalData.length">
      <div v-for="(line, index) in journalData" :key="index">
        {{ index + 1 }}- {{ line }}
      </div>
    </div>

    <v-dialog v-model="dialogJournal" width="450">
      <v-card>
        <v-card-title class="grey white--text">
          Ocserv Service Journalctl
          <v-spacer v-if="dialogJournal" />
          <v-btn icon @click="dialogJournal = false">
            <v-icon color="white">mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-row align="center" justify="center">
            <v-col md="9">
              <v-text-field
                v-model.number="lines"
                label="Number of lines"
                single-line
                hint="journalctl -n ${lines} ocserv.service. blank=20"
                persistent-hint
              />
            </v-col>
            <v-col md="auto">
              <v-btn outlined color="primary" @click="getJournal">Get</v-btn>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts">
import { required } from "@/utils/rules";
import { systemServiceApi } from "@/utils/services";
import Vue from "vue";
export default Vue.extend({
  name: "Journal",
  data(): {
    lines: number | null;
    journalData: string[];
    dialogJournal: boolean;
    autoRefresh: boolean;
    rules: object;
    intervalTime: number;
    intervalTimeDefault: number;
    intervalId: null | number;
  } {
    return {
      lines: 20,
      journalData: [],
      dialogJournal: false,
      autoRefresh: true,
      rules: {
        required: [required],
      },
      intervalTimeDefault: 15,
      intervalTime: 0,
      intervalId: null,
    };
  },

  mounted() {
    this.getJournal();
  },

  methods: {
    async getJournal() {
      let data = await systemServiceApi.journal(this.lines || 20);
      this.journalData = data.logs;
      this.dialogJournal = false;
    },
    makeSetLogs() {},
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
              this.getJournal();
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