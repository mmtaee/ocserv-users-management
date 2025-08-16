<script lang="ts" setup>
import {onMounted, onUnmounted, ref} from 'vue'

import {useI18n} from "vue-i18n";

const {t} = useI18n()
const logs = ref<string[]>([])


let eventSource: EventSource | null = null

const SSE_URL = import.meta.env.VITE_LOG_SOCKET_URL || 'http://localhost:8082/logs'
const connected = ref(false)
const containerHeight = ref(window.innerHeight)
const maxLogs = ref(0)

const updateMaxLogs = () => {
  containerHeight.value = window.innerHeight - 325
  maxLogs.value = Math.floor(containerHeight.value / 25)
  if (logs.value.length > maxLogs.value) {
    logs.value.splice(0, logs.value.length - maxLogs.value)
  }
}

const addLog = async (newLog: string) => {
  if (logs.value.length >= maxLogs.value) {
    logs.value.splice(0, 2)
  }
  logs.value.push(newLog)
}

const connect = () => {
  if (connected.value) return

  addLog(t("START_CONNECTING") + "...")
  eventSource = new EventSource(SSE_URL)

  eventSource.onmessage = (event) => {
    addLog(event.data)
  }

  eventSource.onerror = (error) => {
    console.error('EventSource error:', error)
    disconnect()
  }

  connected.value = true
  addLog(t("SERVER_CONNECTED"))
}

const disconnect = () => {
  addLog(t("SERVER_DISCONNECTING") + "...")
  eventSource?.close()
  eventSource = null
  connected.value = false
  addLog(t("SERVER_DISCONNECT"))
}

onMounted(() => {
  updateMaxLogs()
  window.addEventListener('resize', updateMaxLogs)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateMaxLogs)
  disconnect()
})
</script>

<template>
  <v-card class="pa-4" flat>
    <v-card-title class="text-capitalize mt-5">
      <v-icon class="me-2" start>mdi-account-outline</v-icon>
      {{ t("OCSERV_SERVER_LOGS") }}
      <v-btn
          v-if="!connected"
          color="primary"
          density="compact"
          variant="plain"
          @click="connect"
      >
        {{ t("CONNECT") }}
      </v-btn>
      <v-btn
          v-if="connected"
          color="error"
          density="compact"
          variant="plain"
          @click="disconnect"
      >
        {{ t("DISCONNECT") }}
      </v-btn>
    </v-card-title>

    <v-card-text>
      <div
          id="log-container"
          :style="{
            height: containerHeight + 'px',
            background: '#222',
            color: '#0f0',
            fontFamily: 'monospace',
            whiteSpace: 'pre-line',
            borderRadius: '4px',
          }"
          class="pa-5"
      >
        <div
            v-for="(log, i) in logs"
            :id="`log-${i}`"
            :key="i"
            title="log"
        >
          {{ log }}
        </div>
      </div>
    </v-card-text>
  </v-card>
</template>
