<script lang="ts" setup>
</script>

<template>
  <v-card
      border="none"
      elevation="40"
      height="100%"
      max-width="100%"
      variant="outlined"
  >
    <v-container class="fill-height d-flex justify-center align-center">
      <v-row align="center" justify="center">
        <v-col align-self="center" class="d-flex justify-center" cols="12" md="12">
          <slot/>
        </v-col>
      </v-row>
    </v-container>
  </v-card>
</template>