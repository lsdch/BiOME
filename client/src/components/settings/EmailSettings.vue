<template>
  <CenteredSpinner v-if="isPending" text="Loading e-mail settings..." />
  <v-alert v-else-if="error" color="error" icon="mdi-alert">
    Failed to load e-mail settings
  </v-alert>
  <v-confirm-edit v-else v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @submit="submit(proxy.value).then(save)"
        @reset="cancel()"
        :loading="isUpdating"
      >
        <template #prepend>
          <EmailSettingsTestConnection
            :settings="proxy.value"
            v-model:testing="status.testing"
            v-model:connectionOK="status.connectionOK"
          />
        </template>
      </SettingsFormActions>
      <v-row>
        <v-col>
          <v-alert v-if="updateError" color="error" icon="mdi-alert">
            Failed to update settings
          </v-alert>
        </v-col>
      </v-row>
      <v-card title="Sender identity for automated e-mails" flat>
        <v-container>
          <v-row>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.trim="proxy.value.from_name"
                label="From identity"
                v-bind="field('from_name')"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.trim="proxy.value.from_address"
                label="From address"
                v-bind="field('from_address')"
              />
            </v-col>
          </v-row>
        </v-container>
      </v-card>
      <v-divider />
      <v-card title="SMTP server configuration" flat>
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
      </v-card>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $EmailSettingsInput, EmailSettingsInput } from '@/api'
import {
  emailSettingsOptions,
  updateEmailSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'
import SettingsFormActions from './SettingsFormActions.vue'

const status = ref<{
  testing: boolean
  connectionOK?: boolean
}>({
  testing: false,
  connectionOK: undefined
})

const { data: model, error, isPending, refetch } = useQuery(emailSettingsOptions())

const { field, dispatchErrors } = useSchema($EmailSettingsInput)
const { feedback } = useFeedback()

const {
  mutateAsync,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateEmailSettingsMutation(),
  onSuccess: () => {
    status.value.connectionOK = true
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors,
  onMutate() {
    status.value.testing = true
    status.value.connectionOK = undefined
  },
  onSettled() {
    status.value.testing = false
  }
})

async function submit(model: EmailSettingsInput) {
  await mutateAsync({ body: model })
}
</script>

<style scoped></style>
