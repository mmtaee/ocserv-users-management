<template>
  <v-card flat>
    <v-card-title>
      <v-row>
        <v-col md="4">
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search Iroutes"
            single-line
            hide-details
          />
        </v-col>
      </v-row>
    </v-card-title>
    <v-data-table :headers="headers" :items="routes" :search="search" :hide-default-footer="routes.length < 5"></v-data-table>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
  name: "Iroutes",
  props: {
    routes: Array,
  },
  data() {
    return {
      search: "",
      headers: [
        {
          text: "ID",
          align: "start",
          filterable: true,
          value: "ID",
        },
        {
          text: "Username",
          align: "start",
          filterable: true,
          value: "Username",
        },
        {
          text: "vhost",
          align: "start",
          filterable: true,
          value: "vhost",
        },
        {
          text: "Device",
          align: "start",
          filterable: false,
          value: "Device",
        },
        {
          text: "IP",
          align: "start",
          filterable: true,
          value: "IP",
        },
        {
          text: "iRoutes",
          align: "start",
          filterable: false,
          value: "iRoutes",
        },
      ],
    };
  },
});
</script>