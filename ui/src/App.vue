<script lang="ts" setup>
import {useTheme} from 'vuetify'
import {defineAsyncComponent, onBeforeUnmount, onMounted} from "vue";
import {useIsSmallDisplay} from "@/stores/display.ts";

const Loading = defineAsyncComponent(() => import("@/components/reusable/ReusableLoading.vue"))
const Snackbar = defineAsyncComponent(() => import("@/components/reusable/ReusableSnackbar.vue"))
const Skeleton = defineAsyncComponent(() => import("@/components/Skeleton.vue"))

const theme = useTheme()

const smallDisplay = useIsSmallDisplay()

onMounted(() => {
  checkIsMobile()
  window.addEventListener('resize', checkIsMobile)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkIsMobile)
})

const checkIsMobile = () => {
  smallDisplay.setIsSmall(
      /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent)
      ||
      window.innerWidth <= 890
  )
}

</script>

<template>
  <v-app :theme="theme.global.name.value">
    <Skeleton/>

    <v-main scrollable>
      <v-container class="fill-height d-flex justify-center align-center">
        <v-row align="center" class="ma-0 pa-0" justify="center">
          <v-col align-self="center" class="d-flex justify-center" cols="12">
            <RouterView/>
          </v-col>
        </v-row>
      </v-container>
    </v-main>

    <Snackbar/>

    <Loading/>
  </v-app>
</template>



