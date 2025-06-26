<script lang="ts" setup>
import {useTheme} from 'vuetify'

const theme = useTheme()
</script>

<template>
  <v-app :theme="theme.global.name.value">
    <v-main scrollable>
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
              <RouterView/>
            </v-col>
          </v-row>
        </v-container>
      </v-card>
    </v-main>
  </v-app>
</template>


<style scoped>

</style>

