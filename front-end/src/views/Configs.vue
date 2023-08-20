<template>
  <ConfigsForm v-if="firstInit" />
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "Configs",

  components: {
    ConfigsForm: () => import("../components/ConfigsForm.vue"),
  },

  data(): {
    firstInit: Boolean;
  } {
    return {
      firstInit: false,
    };
  },

  mounted() {
    if (localStorage.getItem("token")) {
      this.$router.push({ name: "Dashboard" });
    } else {
      this.firstInit = true
    }
  },
});
</script>
