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
            Ocserv Groups
          </v-card-subtitle>
          <v-card-text>
            <v-row align="center" justify="start" class="my-3 ms-2">
              <v-col md="3">
                <v-text-field
                  v-model="search"
                  append-icon="mdi-magnify"
                  label="Search Ocserv Group"
                  single-line
                  hide-details
                  @keyup="searchInit"
                  @click:clear="(search = null), (page = 1), init()"
                  clearable
                />
              </v-col>

              <v-col md="3" align-self="start">
                <v-checkbox v-model="ascending" label="Ascending Sort" />
              </v-col>

              <v-spacer />
              <v-col md="auto">
                <v-btn color="primary" outlined @click="groupFormDialog = true">
                  <v-icon left>mdi-home-group-plus</v-icon>
                  Create New Group
                </v-btn>
              </v-col>
              <v-col md="auto" class="me-5">
                <v-btn @click="init" outlined>
                  <v-icon left>mdi-refresh</v-icon>
                  refresh
                </v-btn>
              </v-col>
            </v-row>

            <v-data-table
              :headers="headers"
              :items="groups"
              :hide-default-footer="groups.length < 5"
            >
              <template v-slot:[`item.edit`]="{ item }">
                <v-icon
                  color="primary"
                  @click="
                    (initInput = { ...item }),
                      (groupFormDialog = true),
                      (editMode = true)
                  "
                >
                  mdi-home-edit
                </v-icon>
                <v-icon
                  color="error"
                  right
                  dark
                  @click="(dialogDelete = true), (deleteGroupObj = item)"
                >
                  mdi-delete
                </v-icon>
              </template>

              <template v-slot:[`item.desc`]="{ item }">
                <v-tooltip bottom>
                  <template v-slot:activator="{ on, attrs }">
                    <v-icon
                      color="primary"
                      dark
                      v-bind="attrs"
                      v-on="on"
                      style="cursor: context-menu"
                    >
                      mdi-email-outline
                    </v-icon>
                  </template>
                  <span v-html="item.desc" />
                </v-tooltip>
              </template>

              <template v-slot:[`item.configs`]="{ item }">
                <v-row align="start" justify="start">
                  <v-col
                    class="mx-3 ma-0 pa-0"
                    md="auto"
                    v-for="(val, conf, index) in item.configs"
                    :key="index"
                  >
                    <span class="primary--text">{{ conf }}=</span>{{ val }}
                  </v-col>
                </v-row>
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

    <v-dialog v-model="dialogDelete" max-width="450" v-if="dialogDelete">
      <v-card>
        <v-card-title class="text-h5">
          Delete Group ({{ deleteGroupObj.name }})
        </v-card-title>
        <v-card-text>
          Are you sure to want to delete Group
          <b>({{ deleteGroupObj.name }})?</b>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text @click="dialogDelete = false">
            Cancel
          </v-btn>
          <v-btn color="error" text @click="deleteGroup(deleteGroupObj)">
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="groupFormDialog" width="850">
      <GroupForm
        v-if="groupFormDialog"
        :editMode="editMode"
        @create="createGroup"
        @update="updateGroup"
        @dialog="
          (groupFormDialog = false), (initInput = null), (editMode = false)
        "
        :initInput="initInput || {}"
        dialog
      />
    </v-dialog>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { ocservGroupApi } from "@/utils/services";
import { OcservGroup, URLParams } from "@/utils/types";

export default Vue.extend({
  name: "Groups",
  components: {
    GroupForm: () => import("@/components/GroupForm.vue"),
    Pagination: () => import("@/components/Pagination.vue"),
  },
  data(): {
    groups: Array<OcservGroup | null>;
    headers: Array<object>;
    page: number;
    pages: number;
    item_per_page: number;
    total_count: number;
    search: string;
    traffics: object;
    groupFormDialog: boolean;
    initInput: OcservGroup | null;
    editMode: boolean;
    dialogDelete: boolean;
    deleteGroupObj: OcservGroup | null;
    ascending: boolean;
  } {
    return {
      groups: [],
      headers: [
        {
          text: "Name",
          align: "start",
          filterable: true,
          value: "name",
        },

        {
          text: "Configs",
          align: "start",
          filterable: true,
          value: "configs",
        },
        {
          text: "Description",
          align: "center",
          filterable: true,
          value: "desc",
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
      groupFormDialog: false,
      initInput: null,
      editMode: false,
      dialogDelete: false,
      deleteGroupObj: null,
      ascending: true,
    };
  },

  async mounted() {
    await this.init();
  },

  watch: {
    ascending: {
      immediate: false,
      handler() {
        this.init();
      },
    },
  },

  methods: {
    async init() {
      let params: URLParams = {
        page: this.page,
        item_per_page: this.item_per_page,
        ascending: this.ascending,
      };
      if (this.search) {
        params["name"] = this.search;
      }
      let data = await ocservGroupApi.groups(params);
      this.groups = data.result;
      this.page = data.page || 1;
      this.pages = data.pages || 1;
      this.total_count = data.total_count || 0;
    },
    createGroup(group: OcservGroup) {
      this.groups.unshift(group);
    },
    updateGroup(group: OcservGroup) {
      let index = this.groups.findIndex((item) => item?.id == group.id);
      this.groups.splice(index, 1, group);
    },
    async deleteGroup(group: OcservGroup) {
      await ocservGroupApi.delete_group(group.id!);
      if (ocservGroupApi.status() == 204) {
        let index = this.groups.findIndex((item) => item?.id == group.id);
        this.groups.splice(index, 1);
        this.dialogDelete = false;
        this.deleteGroupObj = null;
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