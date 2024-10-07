<template>
  <SettingsForm @submit="submit" @reset="reload">
    <template #prepend-toolbar>
      <EmailSettingsTestConnection
        :settings="model"
        v-model:testing="status.testing"
        v-model:connectionOK="status.connectionOK"
      />
    </template>
    <template #default>
      <v-container>
        <v-row>
          <v-col cols="12" sm="8">
            <v-text-field v-model.trim="model.host" label="SMTP Host" v-bind="field('host')" />
          </v-col>
          <v-col cols="12" sm="4">
            <NumberInput
              v-model.number="model.port"
              label="SMTP Port"
              :min="1"
              v-bind="field('port')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field v-model.trim="model.user" label="User" v-bind="field('user')" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <PasswordField v-model="model.password" label="Password" v-bind="field('password')" />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </SettingsForm>
</template>

<script setup lang="ts">
import { $EmailSettingsInput, SettingsService } from '@/api'
import { errorFeedback } from '@/api/responses'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'
import SettingsForm from './SettingsForm.vue'
import { useFeedback } from '@/stores/feedback'

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
