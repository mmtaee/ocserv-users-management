<script lang="ts" setup>

const props = defineProps({
  modelValue: Boolean,
  persistent: {type: Boolean, default: false},
  width: {type: String, default: "auto"},
  divider: {type: Boolean, default: false},
  hide_action: {type: Boolean, default: false},
  color: {type: String, default: "primary"}
})

const emit = defineEmits(['update:modelValue'])


</script>

<template>
  <v-dialog
      v-model="props.modelValue"
      :persistent="persistent"
      :width="width"
      transition="dialog-top-transition"
      @update:modelValue="emit('update:modelValue', false)"
  >
    <v-card>
      <v-card-title :class="`bg-${color}`">
        <slot name="dialogTitle"/>
      </v-card-title>

      <v-card-text class="text-subtitle-1 text-capitalize">
        <slot name="dialogText"/>
      </v-card-text>

      <v-divider v-if="divider" class="mb-3"/>

      <v-card-actions v-if="!hide_action" class="justify-end me-2 mb-2">
        <slot name="dialogAction"/>
      </v-card-actions>

    </v-card>
  </v-dialog>
</template>