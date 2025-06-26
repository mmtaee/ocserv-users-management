<script lang="ts" setup>
import {useTheme} from 'vuetify'
import {useIsSmallDevice} from '@/composables/useDevice';
import {defineAsyncComponent} from "vue";

const Snackbar = defineAsyncComponent(() => import("@/components/reusable/ReusableSnackbar.vue"))
const AppBar = defineAsyncComponent(() => import("@/components/AppBar.vue"))

useIsSmallDevice()
const theme = useTheme()

</script>

<template>
  <v-app :theme="theme.global.name.value">

    <AppBar/>

    <v-main scrollable>
      <v-card
          border="none"
          color="background"
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


    <Snackbar/>

  </v-app>
</template>



