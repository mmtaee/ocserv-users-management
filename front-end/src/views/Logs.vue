<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="1400"
          min-height="800"
        >
          <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
            Logs
          </v-card-subtitle>
          <v-card-text class="text-start">
            <div class="black white--text pa-10" style="min-height: 700px">
              <div>
                <span>
                  Server Status
                  {{ status ? "Coneccted" : "Disconnected" }}
                </span>
                <br />
                <span id="ws_conn" v-show="status">
                  Receiving Logs
                  <span id="loading"></span>
                </span>
              </div>
              <div ref="sentence" class="my-10"></div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Typed, { TypedOptions } from "typed.js";

export default Vue.extend({
  name: "Log",
  data() {
    return {
      sentences: [] as string[],
      index: 1,
      socket: undefined as WebSocket | undefined,
      status: false,
    };
  },
  async mounted() {
    this.socket = new WebSocket(
      `ws://127.0.0.1:8080/?user=${localStorage.getItem(
        "user"
      )}&token=${localStorage.getItem("token")}`
    );
    this.socket.onopen = (event: Event) => {
      this.status = true;
      this.loadingDots();
    };
    this.socket.onmessage = (event: MessageEvent) => {
      this.sentences.push(event.data);
    };
    this.socket.onclose = (event: Event) => {
      this.sentences = [];
      console.log("closed");
    };
  },
  methods: {
    async sleep(ms: number) {
      return new Promise((resolve) => setTimeout(resolve, ms));
    },
    loadingDots() {
      let options: TypedOptions = {
        strings: [". . . . . . . ."],
        typeSpeed: 90,
        loop: true,
        showCursor: false,
        smartBackspace: true,
      };
      new Typed("#loading", options);
    },
    async typdCreator() {
      let elm = this.$refs.sentence as HTMLElement;
      let text = this.sentences.shift();
      if (text) {
        let options: TypedOptions = {
          strings: [text],
          typeSpeed: 20,
          loop: false,
          showCursor: false,
          smartBackspace: false,
        };
        let spanElement = document.createElement("span");
        spanElement.id = `text-${this.index}`;
        this.index += 1;
        if (this.index > 22) {
          let ind = this.index - 22;
          document.getElementById(`text-${ind}`)?.remove();
          document.getElementById(`br-${ind}`)?.remove();
        }
        elm.appendChild(spanElement);
        let brElement = document.createElement("br");
        brElement.id = `br-${this.index}`;
        elm.appendChild(brElement);
        new Typed(`#${spanElement.id}`, options);
      }
    },
  },

  watch: {
    sentences: {
      immediate: false,
      deep: true,
      async handler() {
        await this.typdCreator();
      },
    },
  },

  beforeDestroy() {
    this.socket?.close();
  },
});
</script>