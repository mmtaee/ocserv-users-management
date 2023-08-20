<template>
  <v-app-bar color="grey" class="mx-4 rounded" elevate-on-scroll app dark absolute>
    <v-img :src="logo" max-width="120" />
    <span>Ocserv Panel</span>

    <v-tabs v-model="tab" centered v-if="$store.state.isLogin">
      <v-tab v-for="(tab, index) in menuTabs" :key="index" :to="tab.to">
        <v-icon left>{{ tab.icon }}</v-icon>
        {{ tab.title }}
      </v-tab>
    </v-tabs>

    <div v-if="$store.state.isLogin">
      <v-btn color="warning" small class="mx-4" @click="logout">
        <v-icon> mdi-logout </v-icon>
        Logout
      </v-btn>
    </div>
  </v-app-bar>
</template>

<script lang="ts">
import Vue from "vue";
import { adminServiceApi } from "@/utils/services";

export default Vue.extend({
  name: "AppBar",
  data(): {
    logo: string;
    menuTabs: Array<{
      title: string;
      icon: string;
      to: string;
    }>;
    tab: number;
  } {
    return {
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
        // {
        //   title: "User Statistics",
        //   icon: "mdi-chart-bar",
        //   to: "/stats",
        // },
        // {
        //   title: "Logs",
        //   icon: "mdi-math-log",
        //   to: "/logs",
        // },
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
          localStorage.removeItem("user");
          this.$router.push({ name: "Login" });
        }
      }
    },
  },
});
</script>
