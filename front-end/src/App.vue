<template>
  <v-app style="background-color: #eee">
    <AppBar />
    <v-main>
      <router-view v-if="allowRouting" />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { AxiosResponse } from "axios";
import Vue from "vue";
import httpRequest from "./plugins/axios";

export default Vue.extend({
  name: "App",
  components: {
    AppBar: () => import("@/components/AppBar.vue"),
  },
  data() {
    return {
      allowRouting: false,
    };
  },
  async mounted() {
    await this.init();
    this.allowRouting = true;
  },

  methods: {
    async init() {
      interface ConfigResponse {
        config: Boolean;
        captcha_site_key: String | null;
      }
      let res: ConfigResponse | AxiosResponse = await httpRequest("get", {
        urlName: "admin",
        urlPath: "config",
      });
      if ((res as AxiosResponse).status === 401) {
        this.$store.commit("setIsLogin", false);
        localStorage.removeItem("token");
        this.$router.push({ name: "Login" });
      } else {
        this.$store.commit("setSiteKey", res.data.captcha_site_key);
        if (!res.data.config) {
          this.$router.push({ name: "Config" });
        } else {
          if (!localStorage.getItem("token")) {
            this.$router.push({ name: "Login" });
          } else {
            this.$store.commit("setIsLogin", true);
            this.$router.push({ name: "Home" });
          }
        }
      }
    },
  },
});
</script>
