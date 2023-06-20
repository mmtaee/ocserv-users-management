<template>
  <v-card flat max-width="1000">
    <v-card-title>
      <v-row v-if="ips.length > 5">
        <v-col md="4">
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search IP bans"
            single-line
            hide-details
          />
        </v-col>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="ips"
      :search="search"
      :hide-default-footer="ips.length < 5"
    >
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "IPBans",
  props: {
    ips: Array,
  },
  data(): {
    search: string;
    headers: Array<object>;
  } {
    return {
      search: "",
      headers: [
        {
          text: "IP",
          align: "start",
          filterable: true,
          value: "ip",
        },
        {
          text: "Since",
          align: "start",
          filterable: true,
          value: "since",
        },

        {
          text: "Score",
          align: "start",
          filterable: true,
          value: "score",
        },
      ],
    };
  },
});
</script>