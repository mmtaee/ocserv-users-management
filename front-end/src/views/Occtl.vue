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
          <v-card-subtitle
            class="text-h5 grey darken-1 mb-8 white--text text-start"
          >
            Occtl Command Runner
          </v-card-subtitle>
          <v-card-text>
            <v-row align="center" justify="start" class="mx-15">
              <v-col md="4" class="my-0 py-0 mx-3">
                <v-select
                  v-model="command"
                  label="Command"
                  :items="commands"
                  item-text="text"
                  item-value="command"
                  :rules="[rules.required]"
                  persistent-hint
                  return-object
                  :hint="command ? command.help : ''"
                  @change="CommandVars"
                />
              </v-col>

              <v-col md="3">
                <v-text-field
                  v-model="args"
                  :label="argLable"
                  :rules="argRequired ? [rules.required] : []"
                  :disabled="argDisable"
                />
              </v-col>

              <v-col md="auto">
                <v-btn
                  color="primary"
                  outlined
                  :disabled="disabledComput"
                  @click="run"
                >
                  Run
                </v-btn>
              </v-col>
            </v-row>
           
            <v-col md="12">
              <Result
                :result="result"
                :command="command ? command.command : ''"
              />
            </v-col>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { required } from "@/utils/rules";
import { occtlServiceApi } from "@/utils/services";

interface OcservCommands {
  text: string;
  command: string;
  help: string;
  needArg: boolean;
  label: string | null;
}

export default Vue.extend({
  name: "Occtl",

  components: {
    Result: () => import("@/components/occtl/Result.vue"),
  },

  data(): {
    commands: Array<OcservCommands>;
    command: OcservCommands | null;
    rules: object;
    args: string | null;
    argLable: string | null;
    argRequired: boolean;
    argDisable: boolean;
    result: any;
  } {
    return {
      commands: [
        {
          text: "Show Ip Bans",
          command: "show_ip_bans",
          help: "Prints the banned IP addresses",
          needArg: false,
          label: null,
        },
        {
          text: "Show Ip Bans Points",
          command: "show_ip_ban_points",
          help: "Prints all the known IP addresses which have points",
          needArg: false,
          label: null,
        },
        {
          text: "Unban IP",
          command: "unban_ip",
          help: "Unban the specified IP",
          needArg: true,
          label: "Banned IP",
        },
        {
          text: "Reload Configs",
          command: "reload_configs",
          help: "Reloads the server configuration(throttle 1/per minutes)",
          needArg: false,
          label: null,
        },
        {
          text: "Show Status",
          command: "show_status",
          help: "Prints the status and statistics of the server",
          needArg: false,
          label: null,
        },
        {
          text: "Show User",
          command: "show_user",
          help: "Prints information on the specified user",
          needArg: true,
          label: "Username",
        },
        {
          text: "Show Users",
          command: "show_users",
          help: "Prints the connected users",
          needArg: false,
          label: null,
        },
        {
          text: "Show Iroutes",
          command: "show_iroutes",
          help: "Prints the routes provided by users of the server",
          needArg: false,
          label: null,
        },
        {
          text: "Show All Sessions",
          command: "show_sessions_all",
          help: "Prints all the session IDs",
          needArg: false,
          label: null,
        },
        {
          text: "Show Valid Sessions",
          command: "show_sessions_valid",
          help: "Prints all the valid for reconnection sessions",
          needArg: false,
          label: null,
        },
        {
          text: "Disconnect User",
          command: "disconnect_user",
          help: "Disconnect the specified user",
          needArg: true,
          label: "Username",
        },
        {
          text: "Disconnect ID",
          command: "disconnect_id",
          help: "Disconnect the specified ID",
          needArg: true,
          label: "User ID",
        },
      ],
      command: null,
      args: null,
      argLable: null,
      argRequired: false,
      argDisable: true,
      rules: { required: required },
      result: null,
    };
  },

  computed: {
    disabledComput() {
      if (this.command == null) {
        return true;
      } else {
        if (this.command.needArg && !this.args) {
          return true;
        }
        return false;
      }
    },
  },

  methods: {
    async run() {
      if (this.command?.command == "reload_configs") {
        await occtlServiceApi.reload();
        let status = occtlServiceApi.status();
        if (status == 202) {
          console.log("reload occerv");
          // TODO: snackbar 202
        }
        if (status == 429) {
          // TODO : snackbar 429
          console.log("status 429");
        }
      } else {
        this.result = await occtlServiceApi.config(
          this.command?.command!,
          this.args!
        );
      }
    },

    iroutesToJSON(data: string): Array<object> {
      let result = [];
      if (data.length < 2) {
        data = "[" + data + "]";
        result = JSON.parse(data);
      }
      return result;
    },

    CommandVars(command: OcservCommands) {
      if (command && command.needArg) {
        this.argRequired = true;
        this.argLable = command.label;
        this.argDisable = false;
      } else {
        this.argRequired = false;
        this.argLable = null;
        this.argDisable = true;
      }
    },
  },
});
</script>