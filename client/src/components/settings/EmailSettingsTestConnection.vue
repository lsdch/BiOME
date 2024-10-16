<template>
  <v-btn
    :text="testing ? 'Dialing...' : btnProps.text"
    variant="text"
    size="small"
    :readonly="connectionOK !== undefined"
    :color="btnProps.color"
    :prepend-icon="btnProps.prependIcon"
    @click="testConnection(settings)"
  />
  <v-progress-circular v-if="testing" indeterminate />
</template>

<script setup lang="ts">
import { EmailSettingsInput, SettingsService } from '@/api'
import { computed, ref, watch } from 'vue'

const props = defineProps<{ settings: EmailSettingsInput }>()

const testing = defineModel<boolean>('testing', { default: false })
const connectionOK = defineModel<boolean | undefined>('connectionOK', { default: undefined })

const abortController = ref(new AbortController())

watch(
  () => props.settings,
  () => {
    abortController.value.abort()
    testing.value = false
    connectionOK.value = undefined
  },
  { deep: true }
)

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

async function testConnection(settings: EmailSettingsInput) {
  testing.value = true
  const { data: ok, error } = await SettingsService.testSmtp({
    body: settings,
    signal: abortController.value.signal
  }).finally(() => (testing.value = false))
  connectionOK.value = !error && ok
}
</script>

<style scoped></style>
