<template>
  <CenteredSpinner v-if="isPending" text="Loading e-mail settings..." />
  <v-alert v-else-if="error" color="error" icon="mdi-alert">
    Failed to load e-mail settings
  </v-alert>
  <v-form v-else-if="model">
    <template #="{ isValid }">
      <v-confirm-edit v-model="model">
        <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
          <v-row>
            <v-col>
              <v-alert v-if="updateError" color="error" icon="mdi-alert">
                Failed to update settings
              </v-alert>
            </v-col>
          </v-row>
          <v-card title="E-mail service settings" prepend-icon="mdi-email-fast" flat>
            <v-divider></v-divider>
            <v-container fluid>
              <v-list-subheader>Sender identity for automated e-mails</v-list-subheader>
              <v-row>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model.trim="proxy.value.from_name"
                    label="From identity"
                    v-bind="schema('from_name')"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model.trim="proxy.value.from_address"
                    label="From address"
                    v-bind="schema('from_address')"
                  />
                </v-col>
              </v-row>
              <v-divider></v-divider>
              <v-list-subheader>SMTP server configuration</v-list-subheader>
              <v-row>
                <v-col class="d-flex">
                  <v-text-field
                    class="flex-grow-1"
                    v-model.trim="proxy.value.host"
                    label="SMTP Host"
                    v-bind="schema('host')"
                    rounded="e-0"
                  />
                  <v-number-input
                    :min-width="120"
                    rounded="s-0"
                    class="flex-grow-0"
                    v-model.number="proxy.value.port"
                    label="Port"
                    :min="1"
                    v-bind="schema('port')"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-text-field
                    v-model.trim="proxy.value.user"
                    label="User"
                    v-bind="schema('user')"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <PasswordField
                    v-model="proxy.value.password"
                    label="Password"
                    v-bind="schema('password')"
                  />
                </v-col>
              </v-row>
            </v-container>
            <v-divider></v-divider>
            <v-card-actions class="d-flex justify-space-between">
              <EmailSettingsTestConnection
                :disabled="!isValid"
                :settings="proxy.value"
                v-model:testing="status.testing"
                v-model:connectionOK="status.connectionOK"
              />
              <div>
                <v-btn color="" @click="cancel()" :disabled="isPristine || isUpdating">
                  Reset changes
                </v-btn>
                <v-btn
                  text="Save"
                  @click="
                    mutateAsync({ body: proxy.value }, { onSuccess: save, onError: dispatchErrors })
                  "
                />
              </div>
            </v-card-actions>
          </v-card>
        </template>
      </v-confirm-edit>
    </template>
  </v-form>
</template>

<script setup lang="ts">
import { $EmailSettingsInput, EmailSettingsInput } from '@/api'
import {
  emailSettingsOptions,
  updateEmailSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import PasswordField from '@/components/toolkit/forms/PasswordField.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useSchema } from '@/composables/schema'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'

const status = ref<{
  testing: boolean
  connectionOK?: boolean
}>({
  testing: false,
  connectionOK: undefined
})

const { data: model, error, isPending, refetch } = useQuery(emailSettingsOptions())

const {
  bind: { schema },
  dispatchErrors
} = useSchema($EmailSettingsInput)
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
