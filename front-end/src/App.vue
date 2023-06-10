<template>
  <v-app style="background-color: #eee">
    <AppBar />
    <v-main>
      <router-view v-if="allowRouting" />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import { adminServiceApi } from "@/utils/services";

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
      let data = await adminServiceApi.config();
      let status: number = adminServiceApi.status();
      if (status == 401) {
        this.$store.commit("setIsLogin", false);
        localStorage.removeItem("token");
        this.$router.push({ name: "Login" });
      } else {
        if (!data.config) {
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

function AdminServiceApi() {
  throw new Error("Function not implemented.");
}
