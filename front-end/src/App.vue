<template>
  <v-app style="background-color: #eee" @keyup.enter.prevent>
    <AppBar />
    <v-main>
      <router-view v-if="allowRouting" />
    </v-main>
    <SnackBar />
    <LoadingOverlay v-if="$store.state.loadingOverlay" />
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import { adminServiceApi } from "@/utils/services";

export default Vue.extend({
  name: "App",
  components: {
    AppBar: () => import("@/components/AppBar.vue"),
    SnackBar: () => import("@/components/SnackBar.vue"),
    LoadingOverlay: () => import("@/components/LoadingOverlay.vue")
  },
  data() {
    return {
      allowRouting: false,
    };
  },
  async mounted() {
    this.$store.commit("setLoadingOverlay", {
      active: true,
      text: "Loading ..."
    })
    await this.init();
    this.allowRouting = true;
  },

  methods: {
    async init() {
      let data = await adminServiceApi.config();
      if (!data.config) {
        this.$router.push({ name: "Config" });
      } else {
        if (!localStorage.getItem("token")) {
          this.$router.push({ name: "Login" });
        } else {
          this.$store.commit("setIsLogin", true);
          this.$router.push({ name: "Dashboard" });
        }
      }
    },
  },
});
</script>

