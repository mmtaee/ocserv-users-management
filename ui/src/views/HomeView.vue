<script lang="ts" setup>
import {defineAsyncComponent} from "vue";
import {useIsSmallDisplay} from "@/stores/display.ts";


const Desktop = defineAsyncComponent(() => import("@/components/home/Desktop.vue"));
const Mobile = defineAsyncComponent(() => import("@/components/home/Mobile.vue"));

const display = useIsSmallDisplay()

</script>

<template>
  <component :is="display.isSmallDisplay? Mobile: Desktop "/>
</template>

<style scoped>

</style>
