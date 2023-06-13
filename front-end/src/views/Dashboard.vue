<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="1400"
          min-height="800"
        >
          <v-card-subtitle
            class="text-h5 grey darken-1 mb-8 white--text text-start"
          >
            Dashboard
          </v-card-subtitle>
          <v-card-text class="text-start">
            <v-tabs vertical>
              <v-tab v-for="tab in tabs" :key="tab.id">
                {{ tab.name }}
              </v-tab>
              <v-tab-item
                v-for="tab in tabs"
                :key="`item-${tab.id}`"
                class="ma-2"
              >
                <OnlineUsers
                  v-if="tab.key == 'online_users'"
                  :users="serverStats.online_users"
                />

                <v-card flat v-if="tab.key == 'show_status'">
                  <v-card-text>
                    <span v-html="serverStats.show_status" />
                  </v-card-text>
                </v-card>

                <Iroutes
                  v-if="tab.key == 'show_iroutes'"
                  :routes="iroutesToJSON(serverStats.show_iroutes)"
                />
              </v-tab-item>
            </v-tabs>
          </v-card-text>
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
  name: "Dashboard",

  components: {
    OnlineUsers: () => import("@/components/OnlineUsers.vue"),
    Iroutes: () => import("@/components/Iroutes.vue"),
  },

  data(): {
    serverStats: Dashboard;
    tabs: Array<object>;
  } {
    return {
      serverStats: {
        online_users: [],
        show_status: "",
        show_iroutes: "",
        show_ip_bans: "",
      },
      tabs: [
        { id: 1, name: "Online Users", key: "online_users" },
        { id: 2, name: "Show Iroutes", key: "show_iroutes" },
        { id: 3, name: "Show Status", key: "show_status" },

      ],
    };
  },

  async mounted() {
    let data = await adminServiceApi.dashboard();
    this.serverStats = data;
    this.serverStats.show_status = `<pre>${this.serverStats.show_status}</pre>`;
  },

  methods: {
    iroutesToJSON(data: string): Array<object> {
      data = "[" + data + "]";
      return JSON.parse(data);
    },
  },
});
</script>
