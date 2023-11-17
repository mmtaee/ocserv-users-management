<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="1400"
        >
          <v-card-subtitle
            class="text-h5 grey darken-1 mb-8 white--text text-start"
          >
            Ocserv Users
          </v-card-subtitle>

          <v-card-text>
            <v-row align="center" justify="start" class="my-3 ms-2">
              <v-col md="3" align-self="start">
                <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search Ocserv User"
                  single-line
                  hide-details
                  @keyup="searchInit"
                  @click:clear="(search = null), (page = 1), init()"
                  clearable
                />
              </v-col>
              <v-spacer />
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
              <v-col md="auto" class="me-5">
                <v-btn @click="init" outlined>
                  refresh
                  <v-icon right>mdi-refresh</v-icon>
                </v-btn>
              </v-col>
            </v-row>

            <v-data-table
              :headers="headers"
              :items="users"
              :items-per-page="100"
              :hide-default-footer="true"
              no-data-text="No users"
              disable-pagination
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

                <v-btn
                  color="error"
                  dark
                  icon
                  @click="(dialogDisconnect = true), (disconnectUserObj = item)"
                  v-if="item.online"
                >
                  <v-icon color="error"> mdi-lan-disconnect </v-icon>
                </v-btn>

                <v-btn
                  color="error"
                  icon
                  @click="(dialogDelete = true), (deleteUserObj = item)"
                >
                  <v-icon color="error"> mdi-delete </v-icon>
                </v-btn>
              </template>

              <template v-slot:[`item.desc`]="{ item }">
                <v-tooltip bottom v-if="item.desc">
                  <template v-slot:activator="{ on, attrs }">
                    <v-btn color="error" v-bind="attrs" v-on="on" icon>
                      <v-icon color="primary" dark> mdi-email-outline </v-icon>
                    </v-btn>
                  </template>
                  <span v-html="item.desc" />
                </v-tooltip>
                <span v-else> No description </span>
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
                <span class="primary--text">Create Date: </span>
                {{ item.create }}
                <br />
                <span class="primary--text">Expire Date: </span>
                {{ item.expire_date || '- - - - - - - - - -' }}
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

            <v-col cols="auto" class="ma-0 pa-0 text-start" v-if="pages > 1">
              <Pagination
                v-if="Boolean(users.length)"
                :pages="pages"
                :page="page"
                :count="total_count"
                :perPage="item_per_page"
                @changePage="change"
              />
            </v-col>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-dialog
      v-model="dialogDisconnect"
      max-width="450"
      v-if="dialogDisconnect"
    >
      <v-card>
        <v-card-title class="text-h5">
          Disconnect User ({{ disconnectUserObj.username }})
        </v-card-title>
        <v-card-text>
          Are you sure to want to disconnect user
          <b>({{ disconnectUserObj.username }})?</b>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="dialogDisconnect = false">
            Cancel
          </v-btn>
          <v-btn color="error" text @click="disconnectUser(disconnectUserObj)">
            Disconnect
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="dialogDelete" max-width="450" v-if="dialogDelete">
      <v-card>
        <v-card-title class="text-h5">
          Delete User ({{ deleteUserObj.username }})
        </v-card-title>
        <v-card-text>
          Are you sure to want to delete user
          <b>({{ deleteUserObj.username }})?</b>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="dialogDelete = false">
            Cancel
          </v-btn>
          <v-btn color="error" text @click="deleteUser(deleteUserObj)">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

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
import { OcservUser, URLParams } from "@/utils/types";

export default Vue.extend({
  name: "Users",
  components: {
    UserForm: () => import("@/components/UserForm.vue"),
    Pagination: () => import("@/components/Pagination.vue"),
  },
  data(): {
    users: Array<OcservUser | null>;
    headers: Array<object>;
    page: number;
    pages: number;
    item_per_page: number;
    total_count: number;
    search: string;
    traffics: object;
    userFormDialog: boolean;
    initInput: OcservUser | null;
    editMode: boolean;
    dialogDelete: boolean;
    dialogDisconnect: boolean;
    disconnectUserObj: OcservUser | null;
    deleteUserObj: OcservUser | null;
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
      item_per_page: 30,
      total_count: 0,
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
      disconnectUserObj: null,
      deleteUserObj: null,
    };
  },

  async mounted() {
    await this.init();
  },

  methods: {
    async init() {
      let params: URLParams = {
        page: this.page,
        item_per_page: this.item_per_page,
      };
      if (this.search) {
        params["username"] = this.search;
      }
      let data = await ocservUserApi.users(params);
      this.users = data.result;
      this.page = data.page || 1;
      this.pages = data.pages || 1;
      this.total_count = data.total_count || 0;
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
        this.disconnectUserObj = null;
        this.dialogDisconnect = false;
      }
    },
    async deleteUser(user: OcservUser) {
      await ocservUserApi.delete_user(user.id!);
      if (ocservUserApi.status() == 204) {
        let index = this.users.findIndex((item) => item?.id == user.id);
        this.users.splice(index, 1);
        this.dialogDelete = false;
        this.deleteUserObj = null;
      }
    },
    change(page: number, item_per_page: number) {
      this.page = page;
      this.item_per_page = item_per_page;
      this.init();
    },

    searchInit() {
      if (this.search.length > 2) {
        this.init();
      }
    },
  },
});
</script>