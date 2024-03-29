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

                    <!-- captcha field -->
                    <v-col
                      cols="12"
                      md="12"
                      class="pa-0 ma-0 mb-2"
                      style="height: 85px; border: 2px solid rgba(0, 0, 0, 0)"
                      v-if="$store.getters.getGoogleSiteKey"
                    >
                      <center>
                        <Captcha
                          ref="recaptcha"
                          @token="(token) => (input.token = token)"
                          @validForm="
                            (valid) => (!valid ? (input.token = null) : false)
                          "
                        />
                      </center>
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
                  :disabled="validateFormComput"
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
import { required } from "@/utils/rules";
import { AminLogin, User } from "@/utils/types";
import { adminServiceApi } from "@/utils/services";

export default Vue.extend({
  name: "Login",

  components: {
    Captcha: () => import("@/components/Captcha.vue"),
  },

  data(): {
    input: AminLogin;
    rules: object;
    passwordShow: boolean;
    formValid: boolean;
    loading: boolean;
  } {
    return {
      input: {
        username: null,
        password: null,
        token: null,
      },
      rules: { required: required },
      passwordShow: false,
      formValid: true,
      loading: false,
    };
  },

  computed: {
    validateFormComput() {
      if (!this.formValid) return true;
      if (
        Boolean(this.$store.getters.getGoogleSiteKey) &&
        !Boolean(this.input.token)
      ) {
        return true;
      }
      return false;
    },
  },

  methods: {
    async login() {
      this.loading = true;
      let data: {
        token: string;
        user: User;
      } = await adminServiceApi.login(this.input);
      let status: number = adminServiceApi.status();
      if (status == 200) {
        localStorage.setItem("token", data.token);
        this.$store.commit("setUser", data.user);
        this.$store.commit("setIsLogin", true);
        this.$router.push({ name: "Dashboard" });
      }
      this.loading = false;
    },
  },
});
</script>
