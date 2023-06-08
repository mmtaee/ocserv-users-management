<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          style="background-color: #eee"
          width="450"
        >
          <v-card-text>
            <v-card elevation="3" outlined rounded>
              <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
                Login
              </v-card-subtitle>

              <v-card-text>
                <v-form ref="configForm" v-model="formValid">
                  <v-row align="center" justify="start">
                    <v-col md="12" cols="12" class="px-5">
                      <v-text-field
                        v-model="input.username"
                        label="Username"
                        outlined
                        :rules="[rules.required]"
                        dense
                      />
                    </v-col>
                    <v-col md="12" cols="12" class="px-5">
                      <v-text-field
                        v-model="input.password"
                        :type="passwordShow ? 'text' : 'password'"
                        :rules="[rules.required]"
                        label="Password"
                        outlined
                        dense
                        :append-icon="
                          passwordShow
                            ? 'mdi-eye-off-outline'
                            : 'mdi-eye-outline'
                        "
                        @click:append="passwordShow = !passwordShow"
                        autocomplete="new-password"
                      />
                    </v-col>
                  </v-row>
                </v-form>
              </v-card-text>
              <v-card-actions class="px-5 ma-0 py-0 mb-4">
                <v-btn
                  color="secondary"
                  outlined
                  @click="login"
                  block
                  :disabled="!formValid"
                  :loading="loading"
                >
                  Apply
                </v-btn>
              </v-card-actions>
            </v-card>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import httpRequest from "@/plugins/axios";
import { required } from "@/utils/rules";
import { AminLogin } from "@/utils/types";
import { AxiosResponse } from "axios";

export default Vue.extend({
  name: "Login",

  data(): {
    input: AminLogin;
    rules: object;
    passwordShow: boolean;
    formValid: boolean;
    loading: boolean
  } {
    return {
      input: {
        username: null,
        password: null,
      },
      rules: { required: required },
      passwordShow: false,
      formValid: true,
      loading: false,
    };
  },

  methods: {
    async login() {
      this.loading = true
      let res: AxiosResponse = await httpRequest(
        "post",
        {
          urlName: "admin",
          urlPath: "login",
        },
        this.input
      );
      localStorage.setItem("token", res.data.token);
      this.$store.commit("setIsLogin", true);
      this.$router.push({ name: "Home" });
      this.loading = false
    },
  },
});
</script>
