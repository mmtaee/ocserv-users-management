<template>
  <v-form ref="refConfigForm">
    <v-row align="center" justify="start">
      <v-col
        v-for="(item, index) in OcservconfigItems"
        :key="index"
        cols="12"
        :md="md"
        class="my-0 py-0 mx-0 my-0"
      >
        <v-text-field
          v-if="item.type == 'text'"
          v-model="configs[item.model]"
          :label="item.label"
          :outlined="outlined"
          dense
          clearable
          :rules="item.rules || []"
          @change="emitter()"
          @click:clear="(configs[item.model] = null), emitter()"
        />
        <v-select
          v-if="item.type == 'select'"
          v-model="configs[item.model]"
          :items="item.items"
          item-text="text"
          item-value="value"
          :label="item.label"
          :outlined="outlined"
          @change="emitter()"
          @click:clear="(configs[item.model] = null), emitter()"
          dense
          clearable
        />
      </v-col>
      <v-col md="6" class="my-0 py-0">
        <v-text-field
          v-model="route"
          label="routes"
          :outlined="outlined"
          dense
          clearable
          class="ma-0 pa-0 mb-1"
          hide-details
          append-icon="mdi-content-save-outline"
          @click:append="updateRoutes(route), (route = null)"
          @keyup.enter="updateRoutes(route), (route = null)"
        />
      </v-col>
      <v-col md="6" class="my-0 py-0">
        <v-text-field
          v-model="noRoute"
          label="no routes"
          :outlined="outlined"
          dense
          clearable
          class="ma-0 pa-0 mb-1"
          hide-details
          append-icon="mdi-content-save-outline"
          @click:append="updateNo_routes(noRoute), (noRoute = null)"
          @keyup.enter="updateNo_routes(noRoute), (noRoute = null)"
        />
      </v-col>

      <v-col md="6" class="mb-5 text-start">
        <v-card outlined class="ma-0 pa-0 overflow-auto" height="200">
          <v-chip
            v-for="(item, index) in configs['routes']"
            :key="index"
            class="ma-2"
            close
            color="primary"
            outlined
            @click:close="removeRoute(item)"
          >
            {{ item }}
          </v-chip>
        </v-card>
      </v-col>
      <v-col md="6" class="mb-5 text-start">
        <v-card outlined class="ma-0 pa-0 overflow-auto" height="200">
          <v-chip
            v-for="(item, index) in configs['no_routes']"
            :key="index"
            class="ma-2"
            close
            color="primary"
            outlined
            @click:close="removeNoRoute(item)"
          >
            {{ item }}
          </v-chip>
        </v-card>
      </v-col>
    </v-row>
  </v-form>
</template>
<script lang="ts">
import Vue from "vue";
import { OcservConfigItems } from "@/utils/types";
import { number, ip, ipOrRange } from "@/utils/rules";

interface Configs {
  [key: string]: any;
}

export default Vue.extend({
  name: "OcservConfigs",
  props: {
    label: {
      type: String,
      default: () => "keys of group config",
    },
    valueLabel: {
      type: String,
      default: () => "Value of group config",
    },
    vmodelEmit: {
      type: Boolean,
      default: () => false,
    },
    md: {
      type: String,
      default: () => "4",
    },
    outlined: {
      type: Boolean,
      default: false,
    },
    initInput: {
      type: Object,
      default: () => ({}),
    },
  },
  data(): {
    OcservconfigItems: OcservConfigItems[];
    configs: Configs;
    route: string | null;
    noRoute: string | null;
    test: object;
  } {
    return {
      configs: {
        routes: [],
        no_routes: [],
      },
      route: null,
      noRoute: null,
      OcservconfigItems: [
        {
          label: "RX Data(bytes/sec)",
          model: "rx-data-per-sec",
          type: "text",
          rules: [number],
        },
        {
          label: "TX Data(bytes/sec)",
          model: "tx-data-per-sec",
          type: "text",
          rules: [number],
        },
        {
          label: "Max Session Per Client",
          model: "max-same-clients",
          type: "text",
          rules: [number],
        },
        { label: "DNS-1", model: "dns1", type: "text", rules: [ip] },
        { label: "DNS-2", model: "dns2", type: "text", rules: [ip] },
        { label: "DNS-3", model: "dns3", type: "text", rules: [ip] },
        { label: "DNS-4", model: "dns3", type: "text", rules: [ip] },
        { label: "DNS-5", model: "dns3", type: "text", rules: [ip] },
        { label: "DNS-6", model: "dns3", type: "text", rules: [ip] },
        {
          label: "IPV4 Network",
          model: "ipv4-network",
          type: "text",
          rules: [ipOrRange],
        },
        {
          label: "No UDP",
          model: "no-udp",
          type: "select",
          items: [
            { text: "True", value: "true" },
            { text: "False", value: "false" },
          ],
        },
        {
          label: "Keepalive(Seconds)",
          model: "keepalive",
          type: "text",
          rules: [number],
        },
        { label: "DPD(Seconds)", model: "dpd", type: "text", rules: [number] },
        {
          label: "Mobile DPD(Seconds)",
          model: "mobile-dpd",
          type: "text",
          rules: [number],
        },
        {
          label: "Tunnel All DNS",
          model: "tunnel-all-dns",
          type: "select",
          items: [
            { text: "True", value: "true" },
            { text: "False", value: "false" },
          ],
        },
        {
          label: "Restrict User To Routes",
          model: "restrict-user-to-routes",
          type: "select",
          items: [
            { text: "True", value: "true" },
            { text: "False", value: "false" },
          ],
        },
        {
          label: "Stats Report Time(Seconds)",
          model: "stats-report-time",
          type: "text",
          rules: [number],
        },
        { label: "MTU", model: "mtu", type: "text", rules: [number] },
        {
          label: "IDLE Timeout(Seconds)",
          model: "idle-timeout",
          type: "text",
          rules: [number],
        },
        {
          label: "Mobile IDLE Timeout(Seconds)",
          model: "mobile-idle-timeout",
          type: "text",
          rules: [number],
        },
        {
          label: "Session Timeout(Seconds)",
          model: "session-timeout",
          type: "text",
          rules: [number],
        },
      ],
      test: {
        xc: 1,
        xy: 2,
      },
    };
  },

  methods: {
    emitter() {
      const newConfigs: { [key: string]: any } = {};
      Object.entries(this.configs).forEach(([key, val]) => {
        if (Boolean(val)) {
          newConfigs[key] = val;
        }
      });
      if (this.vmodelEmit) this.$emit("input", newConfigs);
      else this.$emit("configs", newConfigs);
    },

    updateRoutes(route: string | null) {
      if (
        Boolean(route) &&
        !this.configs.routes.find((elm: string) => elm == route)
      ) {
        this.configs.routes.push(route);
        this.emitter();
      }
    },

    removeRoute(route: string) {
      let index = this.configs.routes.findIndex((elm: string) => elm == route);
      this.configs.routes.splice(index, 1);
      this.emitter();
    },

    updateNo_routes(noRoute: string | null) {
      if (
        Boolean(noRoute) &&
        !this.configs.no_routes.find((elm: string) => elm == noRoute)
      ) {
        this.configs.no_routes.push(noRoute);
        this.emitter();
      }
    },

    removeNoRoute(noRoute: string) {
      let index = this.configs.no_routes.findIndex(
        (elm: string) => elm == noRoute
      );
      this.configs.no_routes.splice(index, 1);
      this.emitter();
    },
  },

  watch: {
    initInput: {
      immediate: true,
      handler() {
        if (Boolean(this.initInput) && Object.keys(this.initInput).length) this.configs = { ...this.initInput };       
      },
    },
  },
});
</script>