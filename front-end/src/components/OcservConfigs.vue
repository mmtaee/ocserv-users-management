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
    </v-row>
  </v-form>
</template>
<script lang="ts">
import Vue from "vue";
import { OcservConfigItems } from "@/utils/types";
import { number, ip } from "@/utils/rules";

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
    configs: object;
  } {
    return {
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
        {
          label: "IPV4 Network",
          model: "ipv4-network",
          type: "text",
          rules: [ip],
        },
        { label: "DNS-1", model: "dns1", type: "text", rules: [ip] },
        { label: "DNS-2", model: "dns2", type: "text", rules: [ip] },
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
      configs: {},
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
      console.log(newConfigs);
      if (this.vmodelEmit) this.$emit("input", newConfigs);
      else this.$emit("configs", newConfigs);
    },
  },

  watch: {
    initInput: {
      immediate: true,
      handler() {
        if (this.configs) this.configs = { ...this.initInput };
      },
    },
  },
});
</script>