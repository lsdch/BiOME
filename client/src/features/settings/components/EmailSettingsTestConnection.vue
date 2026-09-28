<template>
  <v-btn
    :text="testing ? 'Dialing...' : btnProps.text"
    variant="text"
    size="small"
    :readonly="connectionOK !== undefined"
    :color="btnProps.color"
    :prepend-icon="btnProps.prependIcon"
    @click="testConnection()"
  />
  <v-progress-circular v-if="testing" indeterminate />
</template>

<script setup lang="ts">
import { SettingsService } from '@/api'
import { computed } from 'vue'

const testing = defineModel<boolean>('testing', { default: false })
const connectionOK = defineModel<boolean | undefined>('connectionOK', { default: undefined })

const btnProps = computed(() => {
  switch (connectionOK.value) {
    case true:
      return {
        text: 'Connection OK',
        color: 'success',
        prependIcon: 'mdi-check-network'
      }
    case false:
      return {
        text: 'Connection failed',
        color: 'warning',
        prependIcon: 'mdi-close-network'
      }
    default:
      return {
        text: 'Test connection',
        color: 'primary',
        prependIcon: 'mdi-network'
      }
  }
})

async function testConnection() {
  testing.value = true
  const { data: ok, error } = await SettingsService.testSmtpConnection().finally(
    () => (testing.value = false)
  )
  if (error) {
    console.error('Error testing SMTP connection', error)
  }
  connectionOK.value = !error && ok
}
</script>

<style scoped></style>
