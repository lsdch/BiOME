<template>
  <v-btn
    :text="testing ? 'Testing connection...' : btnProps.text"
    variant="text"
    :readonly="connectionOK !== undefined"
    :color="btnProps.color"
    :prepend-icon="btnProps.prependIcon"
    @click="testConnection(settings)"
  />
  <v-progress-circular v-if="testing" indeterminate />
</template>

<script setup lang="ts">
import { EmailSettingsInput, SettingsService } from '@/api'
import { computed, watch } from 'vue'
import { ref } from 'vue'

const props = defineProps<{ settings: EmailSettingsInput }>()

const testing = defineModel<boolean>('testing', { default: false })
const connectionOK = defineModel<boolean | undefined>('connectionOK', { default: undefined })

const response = ref<ReturnType<typeof SettingsService.testSmtp> | undefined>(undefined)

watch(
  () => props.settings,
  () => {
    response.value?.cancel()
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
  response.value = SettingsService.testSmtp({ requestBody: settings })
  await response.value
    .then((ok) => (connectionOK.value = ok))
    .catch(() => {
      console.info('CAUGHT ERROR')
      connectionOK.value = false
    })
  testing.value = false
}
</script>

<style scoped></style>
