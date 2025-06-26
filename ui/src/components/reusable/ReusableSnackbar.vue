<template>
  <div class="snackbar-stack">
    <v-snackbar
        v-for="(snack, index) in snackbars"
        :key="snack.id"
        v-model="visible[snack.id!]"
        :color="snack.color"
        :style="getStyle(index)"
        :timeout="snack.timeout"
        @update:model-value="val => { if (!val) remove(snack.id!) }"
    >
      {{ snack.message }}
    </v-snackbar>
  </div>
</template>

<script lang="ts" setup>
import {type CSSProperties, ref, watch} from 'vue'
import {storeToRefs} from 'pinia'
import {useSnackbarStore} from '@/stores/snackbar.ts'

const snackbarStore = useSnackbarStore()
const {snackbars} = storeToRefs(snackbarStore)
const {remove} = snackbarStore

const visible = ref<Record<number, boolean>>({})

watch(snackbars, (newList) => {
  newList.forEach((snack) => {
    if (snack?.id && !(snack.id in visible.value)) {
      visible.value[snack.id] = true
    }
  })
}, {immediate: true, deep: true})

const getStyle = (index: number): CSSProperties => {
  return {
    bottom: `${16 + index * 60}px`,
    right: '16px',
    position: 'absolute',
  }
}
</script>

<style scoped>
.snackbar-stack {
  position: fixed;
  right: 0;
  bottom: 0;
  z-index: 9999;
  pointer-events: none;
}
</style>
