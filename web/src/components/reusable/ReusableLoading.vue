<script lang="ts" setup>
import {useLoadingStore} from "@/stores/loading.js";

const loader = useLoadingStore();

</script>

<template>
  <v-overlay v-model="loader.isLoading" class="d-flex align-center justify-center" persistent>
    <span class="loader">
      <span class="blitz"></span>
      <span class="blitz"></span>
    </span>
  </v-overlay>
</template>

<style scoped>
.loader {
  width: 180px;
  height: 140px;
  display: block;
  margin: 0 auto 20px;
  background-image: radial-gradient(circle 25px at 25px 25px, #b0b0b0 100%, transparent 0),
  radial-gradient(circle 50px at 50px 50px, #d0d0d0 100%, transparent 0),
  radial-gradient(circle 25px at 25px 25px, #a0a0a0 100%, transparent 0),
  radial-gradient(circle 15px at 15px 15px, #8c8c8c 100%, transparent 0);
  background-size: 50px 50px, 100px 75px, 50px 50px, 30px 32px;
  background-repeat: no-repeat;
  background-position: 0px 30px, 30px 0px, 113px 29px, 147px 50px;
  position: relative;
  box-sizing: border-box;
}

/* ⚡ Main lightning bolt */
.loader::after,
.loader::before {
  content: '';
  position: absolute;
  width: 20px;
  height: 40px;
  background: #ffeb3b;
  clip-path: polygon(40% 0%, 60% 0%, 50% 40%, 70% 40%, 30% 100%, 45% 60%, 25% 60%);
  animation: flash 0.6s infinite;
}

/* ⚡ First bolt */
.loader::after {
  left: 70px;
  top: 90px;
  transform: rotate(-10deg);
  animation-delay: 0s;
}

/* ⚡ Second bolt */
.loader::before {
  left: 120px;
  top: 85px;
  transform: rotate(15deg);
  animation-delay: 0.3s;
}

/* ⚡ Additional blitz using span */
.loader span.blitz {
  position: absolute;
  width: 16px;
  height: 32px;
  background: #ffeb3b;
  clip-path: polygon(40% 0%, 60% 0%, 50% 40%, 70% 40%, 30% 100%, 45% 60%, 25% 60%);
  animation: flash 0.6s infinite;
  opacity: 0.8;
}

/* ⚡ Third bolt via internal span */
.loader span.blitz:nth-child(1) {
  left: 40px;
  top: 95px;
  transform: rotate(-20deg);
  animation-delay: 0.15s;
}

/* ⚡ Fourth bolt */
.loader span.blitz:nth-child(2) {
  left: 100px;
  top: 105px;
  transform: rotate(5deg);
  animation-delay: 0.45s;
}

@keyframes flash {
  0%, 100% {
    opacity: 0.2;
  }
  50% {
    opacity: 1;
  }
}

.loader::after,
.loader::before {
  width: 20px;
  height: 60px; /* increased height from 40px to 60px */
  background: #ffeb3b;
  clip-path: polygon(40% 0%, 60% 0%, 50% 40%, 70% 40%, 30% 100%, 45% 60%, 25% 60%);
  animation: flash 0.6s infinite;
}

/* inside .loader span.blitz */
.loader span.blitz {
  width: 18px; /* slightly wider for balance */
  height: 54px; /* increased from 32px to 54px */
  background: #ffeb3b;
  clip-path: polygon(40% 0%, 60% 0%, 50% 40%, 70% 40%, 30% 100%, 45% 60%, 25% 60%);
  animation: flash 0.6s infinite;
  opacity: 0.8;
}


</style>
