<template>
  <v-container fluid fill-height>
    <v-row align="center" justify="center">
      <v-col class="d-flex justify-center" md="12" cols="12">
        <v-card
          class="text-center align-center justify-center"
          flat
          width="850"
        >
          <v-card-subtitle class="text-h5 grey darken-1 mb-8 white--text">
            Configuration
          </v-card-subtitle>

          <v-card-text>
            <v-form ref="configForm" v-model="formValid">
              <v-row align="center" justify="start">
                <v-col md="4" cols="12" class="ma-0 pa-1" v-if="!editMode">
                  <v-text-field
                    v-model="input.username"
                    label="Admin Username"
                    outlined
                    :rules="[rules.required]"
                    dense
                  />
                </v-col>
                <v-col md="4" cols="12" class="ma-0 pa-1" v-if="!editMode">
                  <v-text-field
                    v-model="input.password"
                    :type="passwordShow ? 'text' : 'password'"
                    :rules="[rules.required]"
                    label="Admin Password"
                    outlined
                    dense
                    :append-icon="
                      passwordShow ? 'mdi-eye-off-outline' : 'mdi-eye-outline'
                    "
                    @click:append="passwordShow = !passwordShow"
                    autocomplete="new-password"
                  />
                </v-col>
                <v-col :md="!editMode ? 3 : 12" cols="12" class="ma-0 pa-1">
                  <v-text-field
                    v-model="input.default_traffic"
                    label="Default Traffic(GB)"
                    outlined
                    :rules="[rules.required]"
                    dense
                  />
                </v-col>
                <v-col md="6" cols="12" class="ma-0 pa-1 my-0">
                  <v-textarea
                    v-model="input.captcha_site_key"
                    outlined
                    label="Google Captcha v2 site key"
                    dense
                    no-resize
                  />
                </v-col>
                <v-col md="6" cols="12" class="ma-0 pa-1 my-0">
                  <v-textarea
                    v-model="input.captcha_secret_key"
                    outlined
                    label="Google Captcha v2 secret key"
                    dense
                    no-resize
                  />
                </v-col>

                <v-col
                  md="12"
                  cols="12"
                  class="text-start ma-0 pa-1 primary--text text-h5"
                >
                  Default Group configs
                  <v-divider class="mb-5 mt-1" />
                </v-col>

                <v-col md="12" cols="12" class="ma-0 pa-1">
                  <OcservConfigs
                    v-model="input.default_configs"
                    label="Default Config keys"
                    valueLabel="Default Config Value"
                    :initInput="input.default_configs"
                    vmodelEmit
                    outlined
                    md="4"
                  />
                </v-col>
              </v-row>
              <v-row align="center" justify="start" v-if="editMode">
                <v-col md="3" cols="12" class="ma-0 pa-1">
                  <v-checkbox
                    class="ma-0 mt-1"
                    v-model="changePassword"
                    label="Change Password"
                  />
                </v-col>
                <v-col md="9" cols="12" class="ma-0 pa-1">
                  <v-divider />
                </v-col>

                <v-col md="6" cols="12" class="ma-0 pa-1" v-if="changePassword">
                  <v-text-field
                    v-model="input.password"
                    label="Old Password"
                    :append-icon="
                      passwordShow ? 'mdi-eye-off-outline' : 'mdi-eye-outline'
                    "
                    :type="passwordShow ? 'text' : 'password'"
                    :rules="[rules.required]"
                    @click:append="passwordShow = !passwordShow"
                    autocomplete="new-password"
                    outlined
                    dense
                  />
                </v-col>

                <v-col md="6" cols="12" class="ma-0 pa-1" v-if="changePassword">
                  <v-text-field
                    v-model="input.new_password"
                    label="New Password"
                    :append-icon="
                      newPasswordShow
                        ? 'mdi-eye-off-outline'
                        : 'mdi-eye-outline'
                    "
                    :type="newPasswordShow ? 'text' : 'password'"
                    :rules="[rules.required]"
                    @click:append="newPasswordShow = !newPasswordShow"
                    autocomplete="new-password"
                    outlined
                    dense
                  />
                </v-col>
              </v-row>
            </v-form>
          </v-card-text>
          <v-card-actions class="ma-0 py-0">
            <v-btn
              color="primary"
              class="mb-5"
              outlined
              block
              @click="save"
              :disabled="!formValid"
              :loading="loading"
            >
              Apply
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { required } from "@/utils/rules";
import { AdminConfig, Config } from "@/utils/types";
import { adminServiceApi } from "@/utils/services";

export default Vue.extend({
  name: "ConfigsForm",

  components: {
    OcservConfigs: () => import("./OcservConfigs.vue"),
  },

  props: {
    editMode: {
      type: Boolean,
      default: () => false,
    },
    initInput: {
      type: Object || null,
      default: null,
    },
  },

  data(): {
    input: AdminConfig;
    rules: object;
    passwordShow: boolean;
    formValid: boolean;
    loading: boolean;
    changePassword: boolean;
    newPasswordShow: boolean;
  } {
    return {
      input: {
        username: null,
        password: null,
        new_password: null,
        captcha_site_key: null,
        captcha_secret_key: null,
        default_configs: null,
        default_traffic: 10,
      },
      rules: { required: required },
      passwordShow: false,
      formValid: true,
      loading: false,
      changePassword: false,
      newPasswordShow: false,
    };
  },

  methods: {
    async save() {
      this.loading = true;
      let data: null | Config = this.editMode
        ? await adminServiceApi.patch_configuration(this.input)
        : await adminServiceApi.create_configs(this.input);
      let status: number = adminServiceApi.status();
      if (status == 201 && data !== null) {
        if (data.captcha_site_key) {
          this.$store.commit("setSiteKey", data.captcha_site_key);
        }
        if (data.token) {
          localStorage.setItem("token", data.token);
          this.$store.commit("setIsLogin", true);
          this.$router.push({ name: "Dashboard" });
        }
      }
      if (status == 202) {
        this.$store.commit("setSnackBar", {
          text: "Update successfully",
          color: "Success",
        });
      }
      this.loading = false;
    },
  },

  watch: {
    initInput: {
      immediate: true,
      handler() {
        if (this.initInput) this.input = { ...this.initInput };
      },
    },

    changePassword: {
      immediate: false,
      handler() {
        if (this.editMode && !this.changePassword) {
          this.input.password = null;
          this.input.new_password = null;
        }
      },
    },
  },
});
</script>
