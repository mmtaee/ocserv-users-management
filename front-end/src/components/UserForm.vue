<template>
  <v-card :width="width" height="460">
    <v-card-title class="grey darken-1 mb-5 white--text text-start">
      {{ editMode ? "Update Ocserv User" : "Create Ocserv User" }}
      <v-spacer v-if="dialog" />
      <v-btn
        icon
        @click="$refs.userForm.reset(), $emit('dialog', false)"
        v-if="dialog"
      >
        <v-icon color="white">mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-form v-model="userFormValid" ref="userForm">
        <v-row align="center" justify="start">
          <v-col md="4">
            <v-select
              v-model="userInput.group"
              :items="groups"
              item-text="name"
              item-value="id"
              label="Group"
              :rules="[rules.required]"
              dense
              prepend-inner-icon="mdi-home-group"
            />
          </v-col>
          <v-col md="4">
            <v-text-field
              v-model="userInput.username"
              label="Username"
              :rules="[rules.required]"
              dense
              prepend-inner-icon="mdi-account-outline"
            />
          </v-col>
          <v-col md="4">
            <v-text-field
              v-model="userInput.password"
              label="Password"
              :rules="[rules.required]"
              dense
              autocomplete="new-password"
              prepend-inner-icon="mdi-account-key-outline"
            />
          </v-col>
          <v-col md="4">
            <v-dialog
              ref="dialog"
              v-model="dateModal"
              :return-value.sync="userInput.expire_date"
              persistent
              width="290px"
            >
              <template v-slot:activator="{ on, attrs }">
                <v-text-field
                  v-model="userInput.expire_date"
                  label="Expire Date"
                  prepend-inner-icon="mdi-calendar"
                  readonly
                  v-bind="attrs"
                  v-on="on"
                  :rules="[rules.required]"
                  dense
                />
              </template>
              <v-date-picker v-model="userInput.expire_date" scrollable>
                <v-spacer></v-spacer>
                <v-btn
                  text
                  color="primary"
                  @click="$refs.dialog.save(userInput.expire_date)"
                >
                  OK
                </v-btn>
              </v-date-picker>
            </v-dialog>
          </v-col>
          <v-col md="3">
            <v-select
              v-model="userInput.traffic"
              :items="trafficTypes"
              item-text="name"
              item-value="id"
              label="Traffic Type"
              :rules="[rules.required]"
              dense
              prepend-inner-icon="mdi-traffic-light-outline"
            />
          </v-col>
          <v-col md="2">
            <v-text-field
              v-model="userInput.default_traffic"
              :rules="userInput.traffic == 1 ? [] : [rules.required]"
              dense
              label="Default Traffic"
              :disabled="userInput.traffic == 1"
              prepend-inner-icon="mdi-numeric"
            />
          </v-col>
          <v-col md="auto">
            <v-checkbox
              v-model="userInput.active"
              label="Active"
              :disabled="activeDisabled"
            />
          </v-col>
          <v-col md="12">
            <v-textarea
              v-model="userInput.desc"
              label="Description"
              dense
              outlined
              rows="4"
              prepend-inner-icon="mdi-card-text-outline"
              hide-details
            />
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-row align="center" justify="center">
        <v-col md="auto" class="mb-2">
          <v-btn
            outlined
            color="primary"
            :disabled="!userFormValid"
            @click="save"
          >
            {{ editMode ? "Update" : "Create" }}
          </v-btn>
        </v-col>
      </v-row>
    </v-card-actions>
  </v-card>
</template>
<script lang="ts">
import Vue from "vue";
import { OcservUser, OcservGroup, URLParams } from "@/utils/types";
import { required } from "@/utils/rules";
import { ocservGroupApi, ocservUserApi } from "@/utils/services";

export default Vue.extend({
  name: "UserForm",

  props: {
    dialog: {
      type: Boolean,
      default: false,
    },
    width: {
      type: String,
      default: "auto",
    },
    editMode: {
      type: Boolean,
      default: false,
    },
    initInput: {
      type: Object,
      default: () => ({}),
    },
  },

  data(): {
    userInput: OcservUser;
    userFormValid: boolean;
    rules: object;
    groups: Array<OcservGroup | null>;
    dateModal: boolean;
    trafficTypes: Array<object>;
    activeDisabled: Boolean;
  } {
    return {
      userInput: {
        id: null,
        group: null,
        group_name: null,
        username: null,
        password: null,
        active: true,
        expire_date: new Date(Date.now()).toISOString().slice(0, 10),
        desc: "desc",
        traffic: 1,
        default_traffic: 20,
        online: false,
      },
      rules: { required: required },
      userFormValid: true,
      groups: [],
      dateModal: false,
      trafficTypes: [
        { name: "Free", id: 1 },
        { name: "Monthly", id: 2 },
        { name: "Totally", id: 3 },
      ],
      activeDisabled: false,
    };
  },

  async mounted() {
    let params: URLParams = {
      args : "defaults"
    }
    let data = await ocservGroupApi.groups(params);
    this.groups = data.result;
  },

  methods: {
    async save() {
      let data: OcservUser;
      let meitMethodName = "create";

      if (this.editMode) {
        let pk = this.userInput.id;
        meitMethodName = "update";
        data = await ocservUserApi.update_user(pk!, this.userInput);
      } else {
        data = await ocservUserApi.create_user(this.userInput);
      }
      if ([200, 201, 202].includes(ocservUserApi.status())) {
        this.$emit(meitMethodName, data);
        if (this.$refs.userForm) {
          (this.$refs.userForm as HTMLFormElement).reset();
        }
        this.$emit("dialog", false);
      }
    },
  },

  watch: {
    initInput: {
      immediate: true,
      handler() {
        if (Boolean(this.userInput))
          this.userInput = Object.assign({}, this.initInput);
      },
    },
    "userInput.expire_date": {
      immediate: true,
      handler() {
        let today = new Date(Date.now()).toISOString().slice(0, 10);
        if (today == this.userInput.expire_date) {
          this.activeDisabled = true;
          this.userInput.active = false;
        } else {
          this.activeDisabled = false;
        }
      },
    },
  },
});
</script>