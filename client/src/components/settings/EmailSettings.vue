<template>
  <v-confirm-edit v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @submit="save(), submit()"
        @reset="reload().then(cancel)"
      >
        <template #prepend>
          <EmailSettingsTestConnection
            :settings="proxy.value"
            v-model:testing="status.testing"
            v-model:connectionOK="status.connectionOK"
          />
        </template>
      </SettingsFormActions>
      <v-container>
        <v-row>
          <v-col class="d-flex">
            <v-text-field
              class="flex-grow-1"
              v-model.trim="proxy.value.host"
              label="SMTP Host"
              v-bind="field('host')"
              rounded="e-0"
            />
            <v-number-input
              :min-width="100"
              rounded="s-0"
              class="flex-grow-0"
              v-model.number="proxy.value.port"
              label="SMTP Port"
              :min="1"
              v-bind="field('port')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field v-model.trim="proxy.value.user" label="User" v-bind="field('user')" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <PasswordField
              v-model="proxy.value.password"
              label="Password"
              v-bind="field('password')"
            />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $EmailSettingsInput, SettingsService } from '@/api'
import { errorFeedback } from '@/api/responses'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { useFeedback } from '@/stores/feedback'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'
import SettingsFormActions from './SettingsFormActions.vue'

const status = ref<{
  testing: boolean
  connectionOK?: boolean
}>({
  testing: false,
  connectionOK: undefined
})

async function fetch() {
  return SettingsService.emailSettings().then(errorFeedback('Failed to retrieve email settings'))
}

const model = ref(await fetch())

async function reload() {
  model.value = await fetch()
}

const { field, errorHandler } = useSchema($EmailSettingsInput)
const { feedback } = useFeedback()

async function submit() {
  status.value.testing = true
  await SettingsService.updateEmailSettings({ body: model.value })
    .then(errorHandler)
    .then(() => {
      status.value.connectionOK = true
      feedback({ message: 'Settings updated', type: 'success' })
    })
    .catch(() => {
      status.value.connectionOK = false
    })
  status.value.testing = false
}
</script>

<style scoped></style>
