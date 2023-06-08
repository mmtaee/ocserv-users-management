<template>
  <v-app-bar color="grey" class="mx-4 rounded" elevate-on-scroll app dark>
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
import httpRequest from "@/plugins/axios";
import { AxiosResponse } from "axios";
import Vue from "vue";

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
          title: "Home",
          icon: "mdi-home",
          to: "/",
        },
        {
          title: "Users",
          icon: "mdi-account-group-outline",
          to: "/users",
        },
        {
          title: "Groups",
          icon: "mdi-router-network",
          to: "/groups",
        },
        {
          title: "Occtl",
          icon: "mdi-bash",
          to: "/occtl",
        },
        {
          title: "Logs",
          icon: "mdi-math-log",
          to: "/logs",
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
        let res: AxiosResponse = await httpRequest("delete", {
          urlName: "admin",
          urlPath: "logout",
        });
        if (res.status == 204) {
          this.$store.commit("setIsLogin", false);
          localStorage.removeItem("token");
          this.$router.push({ name: "Login" });
        }
      }
    },
  },
});
</script>
