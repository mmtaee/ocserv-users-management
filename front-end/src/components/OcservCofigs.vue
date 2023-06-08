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
          v-model="configs[item.model]"
          :label="item.label"
          outlined
          dense
          clearable
          @keyup="emitter()"
        />
      </v-col>
    </v-row>
  </v-form>
</template>
<script lang="ts">
import { OcservConfigItems } from "@/utils/types";
import Vue from "vue";

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
  },
  data(): {
    OcservconfigItems: OcservConfigItems[];
    configs: object;
  } {
    return {
      OcservconfigItems: [
        { label: "TCP Port", model: "tcp-port" },
        { label: "UDP Port", model: "udp-port" },
        { label: "Max Session Per Client", model: "max-same-clients" },
        { label: "IPV4 Network", model: "ipv4-network" },
        { label: "DNS-1", model: "dns1" },
        { label: "DNS-2", model: "dns2" },
      ],
      configs: {},
    };
  },
  methods: {
    emitter() {
      if (this.vmodelEmit) this.$emit("input", this.configs);
      else this.$emit("configs", this.configs);
    },
  },
});
</script>