<template>
  <v-card min-height="680" flat>
    <v-card-subtitle>
      <v-col md="12" class="ma-0 pa-0">
        <v-row align="center" justify="center">
          <v-col md="5">
            <v-divider style="border-top: 2px solid #1976d2" />
          </v-col>
          <v-col md="1"> Command Result</v-col>
          <v-col md="5">
            <v-divider style="border-top: 2px solid #1976d2" />
          </v-col>
        </v-row>
      </v-col>
    </v-card-subtitle>

    <v-card-text class="mx-15">
      <div v-if="show_status != null" class="text-start">
        <span v-html="show_status" />
      </div>
    </v-card-text>
  </v-card>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  name: "OcctlResult",
  props: {
    result: Object,
    command: String,
  },
  data(): {
    show_status: string | null;
  } {
    return {
      show_status: null,
    };
  },

  methods: {
    dataAnalizer() {
        this.show_status = null
      if (this.result.show_status !== undefined) {
        this.show_status = `<pre>${this.result.show_status}</pre>`;
      } else {
      }
    },
  },

  watch: {
    result: {
      immediate: false,
      handler() {
        console.log(this.result)
        if (this.result) {
          this.dataAnalizer();
        }
      },
    },
  },
});
</script>