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
            <v-row align="start" justify="start" class="my-3">
              <v-col md="auto">
                <v-btn
                  color="primary"
                  outlined
                  @click="(initInput = null), (userFormDialog = true)"
                >
                  <v-icon left>mdi-account-plus-outline</v-icon>
                  Create New User
                </v-btn>
              </v-col>
              <v-col md="auto">
                <v-btn icon @click="init">
                  <v-icon>mdi-refresh</v-icon>
                </v-btn>
              </v-col>
            </v-row>
            <v-row
              v-if="users.length > 5"
              align="start"
              justify="start"
              class="my-3"
            >
              <v-col md="4">
                <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search Ocserv User"
                  single-line
                  hide-details
                />
              </v-col>
            </v-row>
            <v-data-table
              :headers="headers"
              :items="users"
              :search="search"
              :hide-default-footer="users.length < 5"
            >
              <template v-slot:[`item.status`]="{ item }">
                <span class="error--text" v-if="!item.online">Offline</span>
                <span class="success--text" v-else>Online</span>
              </template>

              <template v-slot:[`item.edit`]="{ item }">
                <v-btn icon>
                  <v-icon
                    color="primary"
                    @click="
                      (initInput = { ...item }),
                        (userFormDialog = true),
                        (editMode = true)
                    "
                  >
                    mdi-account-edit-outline
                  </v-icon>
                </v-btn>

                <v-dialog
                  v-model="dialogDisconnect"
                  max-width="450"
                  v-if="item.online"
                >
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn color="error" dark v-bind="attrs" v-on="on" icon>
                      <v-icon color="error"> mdi-lan-disconnect </v-icon>
                    </v-btn>
                  </template>
                  <v-card>
                    <v-card-title class="text-h5">
                      Disconnect User ({{ item.username }})
                    </v-card-title>
                    <v-card-text>
                      Are you sure to want to disconnect user
                      <b>({{ item.username }})?</b>
                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn color="primary" text @click="dialogDelete = false">
                        Cancel
                      </v-btn>
                      <v-btn color="error" text @click="disconnectUser(item)">
                        Disconnect
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-dialog>

                <v-dialog v-model="dialogDelete" max-width="450">
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn color="error" v-bind="attrs" v-on="on" icon>
                      <v-icon color="error"> mdi-delete </v-icon>
                    </v-btn>
                  </template>
                  <v-card>
                    <v-card-title class="text-h5">
                      Delete User ({{ item.username }})
                    </v-card-title>
                    <v-card-text>
                      Are you sure to want to delete user
                      <b>({{ item.username }})?</b>
                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn color="primary" text @click="dialogDelete = false">
                        Cancel
                      </v-btn>
                      <v-btn color="error" text @click="deleteUser(item)">
                        Delete
                      </v-btn>
                    </v-card-actions>
                  </v-card>
                </v-dialog>
              </template>
              <template v-slot:[`item.desc`]="{ item }">
                <v-tooltip bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn color="error" v-bind="attrs" v-on="on" icon>
                      <v-icon color="primary" dark> mdi-email-outline </v-icon>
                    </v-btn>
                  </template>
                  <span v-html="item.desc" />
                </v-tooltip>
              </template>
              <template v-slot:[`item.group`]="{ item }">
                {{ item.group_name }}
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
    <v-dialog v-model="userFormDialog" width="850">
      <UserForm
        v-if="userFormDialog"
        dialog
        :editMode="editMode"
        @create="createUser"
        @update="updateUser"
        @dialog="
          (userFormDialog = false), (initInput = null), (editMode = false)
        "
        :initInput="initInput || {}"
      />
    </v-dialog>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { ocservUserApi } from "@/utils/services";
import { OcservUser } from "@/utils/types";

export default Vue.extend({
  name: "Users",
  components: {
    UserForm: () => import("@/components/UserForm.vue"),
  },
  data(): {
    users: Array<OcservUser | null>;
    headers: Array<object>;
    page: number;
    pages: number;
    search: string;
    traffics: object;
    userFormDialog: boolean;
    initInput: OcservUser | null;
    editMode: boolean;
    dialogDelete: boolean;
    dialogDisconnect: boolean;
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
          align: "start",
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
          align: "start",
          filterable: true,
          value: "default_traffic",
          sortable: false,
        },
        {
          text: "Avrages",
          align: "start",
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
          text: "Status",
          align: "center",
          filterable: false,
          value: "status",
          sortable: false,
        },
        {
          text: "Edit",
          align: "center",
          filterable: false,
          value: "edit",
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
      userFormDialog: false,
      initInput: null,
      editMode: false,
      dialogDelete: false,
      dialogDisconnect: false,
    };
  },

  async mounted() {
    await this.init();
  },

  methods: {
    async init() {
      let data = await ocservUserApi.users();
      this.users = data.result;
      this.page = data.page;
      this.pages = data.pages;
    },
    createUser(user: OcservUser) {
      this.users.unshift(user);
    },
    updateUser(user: OcservUser) {
      let index = this.users.findIndex((item) => item?.id == user.id);
      this.users.splice(index, 1, user);
    },
    async disconnectUser(user: OcservUser) {
      await ocservUserApi.disconnect_user(user.id!);
      if (ocservUserApi.status() == 202) {
        user.online = false;
      }
    },
    async deleteUser(user: OcservUser) {
      await ocservUserApi.delete_user(user.id!);
      if (ocservUserApi.status() == 204) {
        let index = this.users.findIndex((item) => item?.id == user.id);
        this.users.splice(index, 1);
        this.dialogDelete = false;
      }
    },
  },
});
</script>