<template>
  <v-card flat>
    <v-card-title>
      <v-row v-if="users.length > 5">
        <v-col md="4">
          <v-text-field
            v-model="search"
            append-icon="mdi-magnify"
            label="Search Online User"
            single-line
            hide-details
          />
        </v-col>
      </v-row>
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="users"
      :search="search"
      :hide-default-footer="users.length < 5"
    >
      <template v-slot:[`item.host`]="{ item }">
        <span class="primary--text">Hostname:</span> {{ item.hostname }}
        <br />
        <span class="primary--text">Device:</span> {{ item.device }}
        <br />
      </template>
      <template v-slot:[`item.remote_ip`]="{ item }">
        <span class="primary--text">Remote IP:</span> {{ item.remote_ip }}
        <br />
        <span class="primary--text">User Agent:</span> {{ item.user_agent }}
      </template>
      <template v-slot:[`item.since`]="{ item }">
        <span class="primary--text">Since:</span> {{ item.since }}
        <br />
        <span class="primary--text">Connected At:</span> {{ item.connected_at }}
      </template>
      <template v-slot:[`item.averages`]="{ item }">
        <span class="primary--text">Average RX:</span> {{ item.average_rx }}
        <br />
        <span class="primary--text">Average TX:</span> {{ item.average_tx }}
      </template>
    </v-data-table>
  </v-card>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "OnlineUsers",
  props: {
    users: Array,
  },
  data() {
    return {
      search: "",
      headers: [
        {
          text: "Username",
          align: "start",
          filterable: true,
          value: "username",
        },
        {
          text: "Host",
          align: "start",
          filterable: false,
          value: "host",
        },
        {
          text: "Remote",
          align: "start",
          filterable: true,
          value: "remote_ip",
        },
        {
          text: "Since",
          align: "start",
          filterable: false,
          value: "since",
        },
        {
          text: "Averages",
          align: "start",
          filterable: false,
          value: "averages",
        },
      ],
    };
  },
});
</script>