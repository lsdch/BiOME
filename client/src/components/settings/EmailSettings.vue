<template>
  <SettingsForm
    :get="SettingsService.emailSettings"
    :update="updateHandler"
    :schema="$EmailSettings"
  >
    <template #prepend-toolbar="{ model }">
      <EmailSettingsTestConnection
        :settings="model"
        v-model:testing="status.testing"
        v-model:connectionOK="status.connectionOK"
      />
    </template>
    <template #default="{ model, schema }">
      <v-container>
        <v-row>
          <v-col cols="12" sm="8">
            <v-text-field v-model.trim="model.host" label="SMTP Host" v-bind="schema('host')" />
          </v-col>
          <v-col cols="12" sm="4">
            <NumberInput
              v-model.number="model.port"
              label="SMTP Port"
              :min="1"
              v-bind="schema('port')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field v-model.trim="model.user" label="User" v-bind="schema('user')" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <PasswordField v-model="model.password" label="Password" v-bind="schema('password')" />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </SettingsForm>
</template>

<script setup lang="ts">
import { $EmailSettings, SettingsService } from '@/api'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { ref } from 'vue'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'
import SettingsForm from './SettingsForm.vue'

const status = ref<{
  testing: boolean
  connectionOK?: boolean
}>({
  testing: false,
  connectionOK: undefined
})

async function updateHandler(data: Parameters<typeof SettingsService.updateEmailSettings>[0]) {
  status.value.testing = true
  const req = SettingsService.updateEmailSettings(data)
  await req
    .then(() => {
      status.value.connectionOK = true
      return req
    })
    .catch(() => {
      status.value.connectionOK = false
      return req
    })
  status.value.testing = false
  return req
}
</script>

<style scoped></style>
