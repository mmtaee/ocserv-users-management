<template>
  <v-app-bar
    color="grey"
    class="mx-4 rounded"
    elevate-on-scroll
    app
    dark
    absolute
  >
    <v-img :src="logo" max-width="120" />
    <span>Ocserv Panel</span>

    <v-tabs v-model="tab" centered v-if="$store.state.isLogin">
      <v-tab v-for="(tab, index) in menuTabs" :key="index" :to="tab.to">
        <v-icon left>{{ tab.icon }}</v-icon>
        {{ tab.title }}
      </v-tab>
    </v-tabs>

    <div v-if="$store.state.isLogin">
      <v-menu
        v-model="menu"
        :close-on-content-click="false"
        :nudge-width="200"
        offset-x
      >
        <template v-slot:activator="{ on, attrs }">
          <v-btn color="white" dark v-bind="attrs" v-on="on" small outlined>
            Profile
          </v-btn>
        </template>

        <v-card>
          <v-list>
            <v-list-item>
              <v-list-item-content>
                <v-list-item-title class="text-capitalize">
                  {{ $store.getters.getUserstate.username }}
                </v-list-item-title>
                <v-list-item-subtitle class="muted--text">
                  {{
                    $store.getters.getUserstate.is_admin
                      ? "Admin User"
                      : "Staff User"
                  }}
                </v-list-item-subtitle>
              </v-list-item-content>
            </v-list-item>
          </v-list>

          <v-divider></v-divider>

          <v-list>
            <v-list-item
              v-for="(item, index) in menuItemsComput"
              :key="index"
              link
              @click="
                Boolean(item.dialog) ? (dialogs[item.dialog] = true) : logout(),
                  (menu = false)
              "
            >
              <v-list-item-title :class="`${item.color}--text`">
                {{ item.btn }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card>
      </v-menu>
    </div>

    <v-dialog
      v-model="dialogs.changePassword"
      v-if="dialogs.changePassword"
      width="450"
      persistent
      hide-overlay
    >
      <v-card>
        <v-card-title class="grey darken-1 mb-5 white--text text-start">
          Change Password
          <v-spacer />
          <v-btn
            icon
            @click="(dialogs.changePassword = false), $refs.validForm.reset()"
            v-if="dialogs.changePassword"
          >
            <v-icon color="white">mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-form v-model="validForm" ref="validForm">
            <v-text-field
              v-model="changePasswordModel.oldPassword"
              label="Current Password"
              :rules="[rules.required]"
              :type="changePasswordModel.oldPasswordshow ? 'text' : 'password'"
              :append-icon="
                changePasswordModel.oldPasswordshow
                  ? 'mdi-eye-off-outline'
                  : 'mdi-eye-outline'
              "
              @click:append="
                changePasswordModel.oldPasswordshow =
                  !changePasswordModel.oldPasswordshow
              "
              autocomplete="new-password"
            />
            <v-text-field
              v-model="changePasswordModel.password"
              label="New Password"
              :rules="[rules.required]"
              :type="changePasswordModel.passwordShow ? 'text' : 'password'"
              :append-icon="
                changePasswordModel.passwordShow
                  ? 'mdi-eye-off-outline'
                  : 'mdi-eye-outline'
              "
              @click:append="
                changePasswordModel.passwordShow =
                  !changePasswordModel.passwordShow
              "
              autocomplete="new-password"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            outlined
            color="primary"
            small
            class="mb-2"
            :disabled="!validForm"
            @click="changePassword"
            :loading="changePasswordModel.loading"
          >
            Change Password
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="dialogs.staffUsers"
      v-if="dialogs.staffUsers"
      width="850"
      persistent
      hide-overlay
    >
      <v-card>
        <v-card-title class="grey darken-1 mb-5 white--text text-start">
          Staff Users
          <v-spacer />
          <v-btn
            icon
            @click="dialogs.staffUsers = false"
            v-if="dialogs.staffUsers"
          >
            <v-icon color="white">mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text>
          <v-data-table
            :headers="staff.staffListHeaders"
            :items="staff.staffListItems"
            hide-default-footer
            :search="staff.search"
          >
            <template v-slot:[`item.is_admin`]="{ item }">
              <span :class="item.is_admin ? 'primary--text' : 'muted--text'">
                {{ item.is_admin ? "Admin" : "Staff" }}
              </span>
            </template>
            <template v-slot:[`item.action`]="{ item }">
              <v-btn v-if="!item.is_admin" icon>
                <v-icon color="error" small @click="deleteStaff(item.id)">
                  mdi-delete-outline
                </v-icon>
              </v-btn>
              <v-icon v-else small> mdi-minus-circle-outline </v-icon>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="dialogs.createStaff"
      v-if="dialogs.createStaff"
      width="400"
      persistent
      hide-overlay
    >
      <v-card>
        <v-card-title class="grey darken-1 mb-5 white--text text-start">
          Create Staff User
          <v-spacer />
          <v-btn
            icon
            @click="(dialogs.createStaff = false), $refs.validForm.reset()"
          >
            <v-icon color="white">mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text>
          <v-form v-model="validForm" ref="validForm">
            <v-text-field
              v-model="staff.user.username"
              label="Username"
              :rules="[rules.required]"
            />
            <v-text-field
              v-model="staff.user.password"
              label="Password"
              :rules="[rules.required]"
              :type="staff.createPasswordShow ? 'text' : 'password'"
              :append-icon="
                staff.createPasswordShow
                  ? 'mdi-eye-off-outline'
                  : 'mdi-eye-outline'
              "
              @click:append="
                staff.createPasswordShow = !staff.createPasswordShow
              "
              autocomplete="new-password"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            outlined
            color="primary"
            small
            class="mb-2"
            :disabled="!validForm"
            @click="createStaff"
            :loading="staff.createLoading"
          >
            Create Staff
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-app-bar>
</template>


<script lang="ts">
import Vue from "vue";
import { adminServiceApi } from "@/utils/services";
import { required } from "@/utils/rules";
import { User } from "@/utils/types";

declare interface MenuItems {
  btn: string;
  dialog: string;
  color: string;
}

declare interface ChangePasswordModel {
  password: string | null;
  passwordShow: boolean;
  oldPassword: string | null;
  oldPasswordshow: boolean;
  loading: boolean;
}

declare interface Dialogs {
  [key: string]: boolean;
}

declare interface Staff {
  createPasswordShow: boolean;
  createLoading: boolean;
  user: object;
  search: string;
  staffListHeaders: Array<object>;
  staffListItems: Array<User>;
}

export default Vue.extend({
  name: "AppBar",
  data(): {
    staff: Staff;
    dialogs: Dialogs;
    validForm: boolean;
    rules: object;
    changePasswordModel: ChangePasswordModel;
    menu: boolean;
    logo: string;
    menuTabs: Array<{
      title: string;
      icon: string;
      to: string;
    }>;
    tab: number;
  } {
    return {
      staff: {
        createPasswordShow: false,
        createLoading: false,
        user: {
          username: null,
          password: null,
        },
        search: "",
        staffListHeaders: [
          {
            text: "ID",
            align: "start",
            filterable: true,
            value: "id",
          },
          {
            text: "Username",
            align: "start",
            filterable: true,
            value: "username",
          },
          {
            text: "Role",
            align: "center",
            filterable: true,
            value: "is_admin",
          },
          {
            text: "action",
            align: "center",
            filterable: false,
            value: "action",
          },
        ],
        staffListItems: [],
      },
      dialogs: {
        changePassword: false,
        staffUsers: false,
        createStaff: false,
      },
      validForm: true,
      rules: { required: required },
      changePasswordModel: {
        password: null,
        passwordShow: false,
        oldPassword: null,
        oldPasswordshow: false,
        loading: false,
      },
      menu: false,
      logo: require("@/assets/oc_logo.png"),
      menuTabs: [
        {
          title: "Dashboard",
          icon: "mdi-monitor-dashboard",
          to: "/",
        },
        {
          title: "Groups",
          icon: "mdi-home-group",
          to: "/groups",
        },
        {
          title: "Users",
          icon: "mdi-account-group-outline",
          to: "/users",
        },
        {
          title: "Occtl",
          icon: "mdi-bash",
          to: "/occtl",
        },
        {
          title: "User Statistics",
          icon: "mdi-chart-bar",
          to: "/stats",
        },
        {
          title: "System",
          icon: "mdi-linux",
          to: "/system",
        },
        {
          title: "Configuration",
          icon: "mdi-cog-outline",
          to: "/configuration",
        },
      ],
      tab: 0,
    };
  },
  methods: {
    async logout() {
      if (localStorage.getItem("token")) {
        await adminServiceApi.logout();
        let status: number = adminServiceApi.status();
        if (status == 204) {
          this.$store.commit("setIsLogin", false);
          localStorage.removeItem("token");
          this.$store.commit("setUser", {
            username: null,
            is_admin: false,
          });
          this.$router.push({ name: "Login" });
        }
      }
    },

    async changePassword() {
      this.changePasswordModel.loading = true;
      let data: object = {
        old_password: this.changePasswordModel.oldPassword,
        password: this.changePasswordModel.password,
      };
      await adminServiceApi.change_password(data);
      let status: number = adminServiceApi.status();
      if (status == 202) {
        this.dialogs.changePassword = false;
        this.$store.commit("setSnackBar", {
          text: "password changed successfully",
          color: "success",
        });
        if (this.$refs.validForm) {
          (this.$refs.validForm as HTMLFormElement).reset();
        }
      }
      this.changePasswordModel.loading = false;
    },

    async createStaff() {
      this.staff.createLoading = true;
      let staff = Object.assign({}, this.staff.user);
      let data = await adminServiceApi.create_staff(staff);
      let status: number = adminServiceApi.status();
      if (status == 202) {
        this.dialogs.createStaff = false;
        this.$store.commit("setSnackBar", {
          text: "Staff created successfully",
          color: "success",
        });
      }
      this.staff.createLoading = false;
    },

    async staffUsers() {
      let data = await adminServiceApi.get_staff();
      let status: number = adminServiceApi.status();
      if (status == 200) {
        this.staff.staffListItems = JSON.parse(JSON.stringify(data));
      }
    },

    async deleteStaff(id: number) {
      await adminServiceApi.delete_staff(id);
      if (adminServiceApi.status() == 204) {
        let itemIndex = this.staff.staffListItems.findIndex(
          (elm) => elm.id == id
        );
        this.staff.staffListItems.splice(itemIndex, 1);
        this.$store.commit("setSnackBar", {
          text: "staff deleted successfully",
          color: "success",
        });
      }
    },
  },

  computed: {
    menuItemsComput(): Array<MenuItems> {
      let dafaultMenuItems: Array<MenuItems> = [
        { btn: "Change Password", dialog: "changePassword", color: "primary" },
        { btn: "Logout", dialog: "", color: "red" },
      ];
      if (this.$store.getters.getUserstate.is_admin) {
        dafaultMenuItems.unshift({
          btn: "Staff List",
          dialog: "staffUsers",
          color: "primary",
        });
        dafaultMenuItems.unshift({
          btn: "Create Staff",
          dialog: "createStaff",
          color: "primary",
        });
      }
      return dafaultMenuItems;
    },
  },

  watch: {
    "dialogs.staffUsers": {
      immediate: false,
      deep: true,
      handler() {
        if (this.dialogs.staffUsers) {
          this.staffUsers();
        }
      },
    },
  },
});
</script>
