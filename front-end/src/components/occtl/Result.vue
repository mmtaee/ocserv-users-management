<template>
  <v-card min-height="680" flat>
    <v-card-subtitle>
      <v-col md="12" class="ma-0 pa-0">
        <v-row align="center" justify="center">
          <v-col md="5">
            <v-divider style="border-top: 2px solid #1976d2" />
          </v-col>
          <v-col md="1"> Command Result</v-col>
          <v-col md="5">
            <v-divider style="border-top: 2px solid #1976d2" />
          </v-col>
        </v-row>
      </v-col>
    </v-card-subtitle>

    <v-card-text class="mx-15">
      <div v-if="show_status" class="text-start">
        <span v-html="show_status" />
      </div>
      <OnlineUsers :users="show_users" v-if="show_users" />
      <Iroutes :routes="StringToJson(show_iroutes)" v-if="show_iroutes" />
      <IPBans :ips="show_ip_bans" v-if="show_ip_bans" />
    </v-card-text>
  </v-card>
</template>
<script lang="ts">
import Vue from "vue";
import { StringToJson } from "@/utils/methods";

export default Vue.extend({
  name: "OcctlResult",
  components: {
    OnlineUsers: () => import("./OnlineUsers.vue"),
    Iroutes: () => import("./Iroutes.vue"),
    IPBans: () => import("./IPBans.vue"),
  },
  props: {
    result: Object,
    command: String,
  },
  data(): {
    StringToJson: Function;
    show_status: string | null;
    show_users: Array<object> | null;
    show_iroutes: string | null;
    show_ip_bans: Array<object> | null;
  } {
    return {
      StringToJson: StringToJson,
      show_status: null,
      show_users: null,
      show_iroutes: null,
      show_ip_bans: null,
    };
  },

  methods: {
    dataAnalizer() {
      this.show_status = null;
      this.show_users = null;
      this.show_iroutes = null;
      this.show_ip_bans = null;
      if (this.result.show_status !== undefined) {
        this.show_status = `<pre>${this.result.show_status}</pre>`;
      }
      if (
        this.result.show_users != undefined ||
        this.result.show_user != undefined
      ) {
        this.show_users =
          this.result.show_users !== undefined
            ? this.result.show_users
            : this.result.show_user;
      }
      if (this.result.show_iroutes) {
        this.show_iroutes = this.result.show_iroutes;
      }
      if (
        this.result.show_ip_ban_points !== undefined ||
        this.result.show_ip_bans !== undefined
      ) {
        this.show_ip_bans =
          this.result.show_ip_bans !== undefined
            ? this.result.show_ip_bans
            : this.result.show_ip_ban_points;
      }
      if (
        this.result.disconnect_user ||
        this.result.unban_ip ||
        this.result.disconnect_id
      ) {
        // TODO: snackbar result
      }
    },
  },

  watch: {
    result: {
      immediate: false,
      handler() {
        if (this.result) {
          this.dataAnalizer();
        }
      },
    },
  },
});
</script>