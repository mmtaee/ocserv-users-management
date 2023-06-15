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
          :outlined="outlined"
          dense
          clearable
          @keyup="emitter()"
          :prepend-inner-icon="innerIcon ? item.icon : ''"
        />
      </v-col>
    </v-row>
  </v-form>
</template>
<script lang="ts">
import Vue from "vue";
import { OcservConfigItems } from "@/utils/types";

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
    innerIcon: {
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
        { label: "TCP Port", model: "tcp-port", icon: "mdi-ethernet-cable" },
        { label: "UDP Port", model: "udp-port", icon: "mdi-ethernet-cable" },
        {
          label: "Max Session Per Client",
          model: "max-same-clients",
          icon: "mdi-account-network-outline",
        },
        {
          label: "IPV4 Network",
          model: "ipv4-network",
          icon: "mdi-ip-network-outline",
        },
        { label: "DNS-1", model: "dns1", icon: "mdi-dns-outline" },
        { label: "DNS-2", model: "dns2", icon: "mdi-dns-outline" },
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