<script lang="ts" setup>
import {defineAsyncComponent} from "vue";

const Home = defineAsyncComponent(() => import("@/components/home/index.vue"));
</script>

<template>
  <Home/>
</template>

