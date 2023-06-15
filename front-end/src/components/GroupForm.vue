<template>
  <v-card :width="width" height="480">
    <v-card-title class="grey darken-1 mb-5 white--text text-start">
      {{ editMode ? "Update Ocserv Group" : "Create Ocserv Group" }}
      <v-spacer v-if="dialog" />
      <v-btn
        icon
        @click="$refs.groupForm.reset(), $emit('dialog', false)"
        v-if="dialog"
      >
        <v-icon color="white">mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text>
      <v-form v-model="groupFormValid" ref="groupForm">
        <v-row align="center" justify="start">
          <v-col md="4">
            <v-text-field
              v-model="groupInput.name"
              label="Group Name"
              :rules="[rules.required]"
              dense
              prepend-inner-icon="mdi-home-group"
            />
          </v-col>
          <v-col md="12">
            <v-col md="12" cols="12" class="ma-0 pa-1">
              <OcservConfigs
                v-model="groupInput.configs"
                :initInput="initInput.configs"
                label="Config keys"
                valueLabel="Config Value"
                vmodelEmit
                innerIcon
                md="4"
              />
            </v-col>
          </v-col>
          <v-col md="12">
            <v-textarea
              v-model="groupInput.desc"
              label="Description"
              dense
              outlined
              rows="4"
              prepend-inner-icon="mdi-card-text-outline"
              hide-details
            />
          </v-col>
        </v-row>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-row align="center" justify="center">
        <v-col md="auto" class="mb-2">
          <v-btn
            outlined
            color="primary"
            :disabled="!groupFormValid"
            @click="save"
          >
            {{ editMode ? "Update" : "Create" }}
          </v-btn>
        </v-col>
      </v-row>
    </v-card-actions>
  </v-card>
</template>
<script lang="ts">
import Vue from "vue";
import { OcservGroup } from "@/utils/types";
import { required } from "@/utils/rules";
import { ocservGroupApi } from "@/utils/services";

export default Vue.extend({
  name: "GroupForm",

  components: {
    OcservConfigs: () => import("@/components/OcservConfigs.vue"),
  },

  props: {
    dialog: {
      type: Boolean,
      default: false,
    },
    width: {
      type: String,
      default: "auto",
    },
    editMode: {
      type: Boolean,
      default: false,
    },
    initInput: {
      type: Object,
      default: () => ({}),
    },
  },

  data(): {
    groupInput: OcservGroup;
    groupFormValid: boolean;
    rules: object;
    dateModal: boolean;
  } {
    return {
      groupInput: {
        id: null,
        name: null,
        desc: "desc",
        configs: {},
      },
      rules: { required: required },
      groupFormValid: true,
      dateModal: false,
    };
  },

  methods: {
    async save() {
      let data: OcservGroup;
      let meitMethodName = "create";
      if (this.editMode) {
        let pk = this.groupInput.id;
        meitMethodName = "update";
        data = await ocservGroupApi.update_group(pk!, this.groupInput);
      } else {
        data = await ocservGroupApi.create_group(this.groupInput);
      }
      if ([200, 201, 202].includes(ocservGroupApi.status())) {
        this.$emit(meitMethodName, data);
        if (this.$refs.groupForm) {
          (this.$refs.groupForm as HTMLFormElement).reset();
        }
        this.$emit("dialog", false);
      }
    },
  },

  watch: {
    initInput: {
      immediate: true,
      handler() {
        if (this.groupInput) this.groupInput = { ...this.initInput };
      },
    },
  },
});
</script>