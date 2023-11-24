<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="1400"
          min-height="880"
        >
          <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
            <span class="mb-10">System Logs & Services</span>
            <div>
              <v-tabs
                v-model="subTabs"
                background-color="grey darken-1 accent-4"
                centered
                dark
                v-if="!dockerized"
              >
                <v-tabs-slider></v-tabs-slider>
                <v-tab
                  v-for="(tab, index) in subMenuTabs"
                  :key="index"
                  :href="'#tab-' + index"
                >
                  {{ tab.title }}
                </v-tab>
              </v-tabs>
              <v-tabs
                v-model="subTabs"
                background-color="grey darken-1 accent-4"
                centered
                dark
                v-else
              >
                <v-tabs-slider></v-tabs-slider>
                <v-tab
                  v-for="(tab, index) in subMenuTabsDocker"
                  :key="index"
                  :href="'#tab-' + index"
                >
                  {{ tab.title }}
                </v-tab>
              </v-tabs>
            </div>
          </v-card-subtitle>

          <v-card-text class="text-start">
            <v-tabs-items v-model="subTabs" v-if="!dockerized">
              <v-tab-item
                v-for="(tab, index) in subMenuTabs"
                :key="'sub' + index"
                :value="'tab-' + index"
              >
                <v-card flat>
                  <v-card-text>
                    <component
                      :is="tab.component"
                      v-if="subTabs == 'tab-' + index"
                    />
                  </v-card-text>
                </v-card>
              </v-tab-item>
            </v-tabs-items>
            <v-tabs-items v-model="subTabs" v-else>
              <v-tab-item
                v-for="(tab, index) in subMenuTabsDocker"
                :key="'sub' + index"
                :value="'tab-' + index"
              >
                <v-card flat>
                  <v-card-text>
                    <component
                      :is="tab.component"
                      v-if="subTabs == 'tab-' + index"
                    />
                  </v-card-text>
                </v-card>
              </v-tab-item>
            </v-tabs-items>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "System",
  data() {
    return {
      sentences: [] as string[],
      subTabs: "",
      dockerized: false,
      subMenuTabs: [
        {
          id: 1,
          title: "Action Logs",
          component: () => import("@/components/system/Action.vue"),
        },
        {
          id: 2,
          title: "Ocserv Systemd",
          component: () => import("@/components/system/OcservSyatemd.vue"),
        },
        {
          id: 3,
          title: "Ocserv Journal",
          component: () => import("@/components/system/Journal.vue"),
        },
      ],
      subMenuTabsDocker: [
        {
          id: 1,
          title: "Action Logs",
          component: () => import("@/components/system/Action.vue"),
        },
      ],
    };
  },

  created() {
    if (process.env.VUE_APP_DOCKERIZED == "true") {
      this.dockerized = true;
    }
  },
});
</script>