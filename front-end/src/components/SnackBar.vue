<template>
  <v-snackbar
    class="mt-15 ms-2 text-capitalize"
    v-if="snackbar"
    v-model="snackbar"
    timeout="3000"
    :color="color"
    top
    left
    multi-line
    absolute
    transition="fade-transition"
  >
    {{ text }}
  </v-snackbar>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
  name: "SnackBar",
  data() {
    return {
      snackbar: false,
      text: null,
      color: "success",
    };
  },

  watch: {
    "$store.state.snacbar.text": {
      immediate: false,
      handler() {
        if (this.$store.state.snacbar.text) {
          this.text = this.$store.state.snacbar.text;
          this.color = this.$store.state.snacbar.color;
          this.$store.commit("setSnackBar", {});
          this.snackbar = true;
        }
      },
    },
  },
});
</script>