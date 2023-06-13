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
            Ocserv Users
          </v-card-subtitle>
          <v-card-text>
            <v-btn>Create New User</v-btn>

            <v-data-table
              :headers="headers"
              :items="users"
              :search="search"
              :hide-default-footer="users.length < 5"
            >
              <template v-slot:[`item.action`]>
                <v-icon color="primary">mdi-account-edit-outline</v-icon>
              </template>

              <template v-slot:[`item.desc`]>
                <v-icon color="primary">mdi-email-edit-outline</v-icon>
              </template>

              <template v-slot:[`item.avrages`]="{ item }">
                <span class="primary--text">RX:</span>
                {{ item.rx }}
                <br />
                <span class="primary--text">TX:</span>
                {{ item.tx }}
                <br />
              </template>
              <template v-slot:[`item.default_traffic`]="{ item }">
                <span class="primary--text">Default Traffic:</span>
                {{ item.default_traffic }}
                <br />
                <span class="primary--text">Traffic Type:</span>
                {{ traffics[item.traffic] }}
                <br />
              </template>
              <template v-slot:[`item.expire_date`]="{ item }">
                <span class="primary--text">Create Date:</span>
                {{ item.create }}
                <br />
                <span class="primary--text">Expire Date:</span>
                {{ item.expire_date }}
                <br />
              </template>
              <template v-slot:[`item.username`]="{ item }">
                <span class="primary--text">Username:</span> {{ item.username }}
                <br />
                <span class="primary--text">Password:</span> {{ item.password }}
                <br />
              </template>
              <template v-slot:[`item.active`]="{ item }">
                <v-icon v-if="item.active" color="success">
                  mdi-account-check-outline
                </v-icon>
                <v-icon v-else color="error">
                  mdi-account-cancel-outline
                </v-icon>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { ocservUserApi } from "@/utils/services";
import { OcservUser } from "@/utils/types";

export default Vue.extend({
  name: "Users",
  data(): {
    users: Array<OcservUser | null>;
    headers: Array<object>;
    page: number;
    pages: number;
    search: string;
    traffics: object;
  } {
    return {
      users: [],
      headers: [
        {
          text: "User",
          align: "start",
          filterable: true,
          value: "username",
        },
        {
          text: "Group",
          align: "start",
          filterable: true,
          value: "group",
        },
        {
          text: "Active",
          align: "center",
          filterable: true,
          value: "active",
        },
        {
          text: "Dates",
          align: "start",
          filterable: true,
          value: "expire_date",
          sortable: false,
        },
        {
          text: "Traffic Details",
          align: "center",
          filterable: true,
          value: "default_traffic",
          sortable: false,
        },
        {
          text: "Avrages",
          align: "center",
          filterable: false,
          value: "avrages",
          sortable: false,
        },
        {
          text: "Description",
          align: "center",
          filterable: false,
          value: "desc",
          sortable: false,
        },

        {
          text: "Action",
          align: "center",
          filterable: false,
          value: "action",
          sortable: false,
        },
      ],
      page: 1,
      pages: 1,
      search: "",
      traffics: {
        1: "Free",
        2: "Monthly",
        3: "Totally",
      },
    };
  },

  async mounted() {
    let data = await ocservUserApi.users();
    this.users = data.result;
    this.page = data.page;
    this.pages = data.pages;
  },
});
</script>