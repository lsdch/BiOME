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
          <v-card
            title="E-mail service settings"
            -
            :subtitle="`Last updated: ${DateTime.fromJSDate(model.last_updated).toLocaleString(DateTime.DATETIME_MED)}`"
            prepend-icon="mdi-email-fast"
            flat
          >
            <v-divider></v-divider>
            <v-container fluid>
              <v-row>
                <v-col>
                  <v-card variant="outlined" class="mb-3">
                    <v-list>
                      <ListItemInput
                        title="E-mail service enabled"
                        subtitle="Enable or disable e-mail sending for the entire system."
                      >
                        <v-switch
                          color="primary"
                          base-color="error"
                          @update:model-value="(v) => toggleMailing(v)"
                          :model-value="proxy.value.enabled"
                        />
                      </ListItemInput>
                    </v-list>
                  </v-card>
                </v-col>
              </v-row>
              <v-list-subheader>Sender identity for automated e-mails</v-list-subheader>
              <v-row>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model.trim="proxy.value.mail_from_name"
                    label="From identity"
                    v-bind="schema('mail_from_name')"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model.trim="proxy.value.mail_from_address"
                    label="From address"
                    v-bind="schema('mail_from_address')"
                  />
                </v-col>
              </v-row>
              <v-divider></v-divider>
              <v-list-subheader>SMTP server configuration</v-list-subheader>
              <v-row>
                <v-col class="d-flex">
                  <v-text-field
                    class="flex-grow-1"
                    v-model.trim="proxy.value.smtp_host"
                    label="SMTP Host"
                    v-bind="schema('smtp_host')"
                    rounded="e-0"
                  />
                  <v-number-input
                    :min-width="120"
                    rounded="s-0"
                    class="flex-grow-0"
                    v-model.number="proxy.value.smtp_port"
                    label="Port"
                    :min="1"
                    v-bind="schema('smtp_port')"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-text-field
                    v-model.trim="proxy.value.smtp_user"
                    label="User"
                    v-bind="schema('smtp_user')"
                  />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <PasswordField
                    :model-value="proxy.value.smtp_password ?? ''"
                    @update:model-value="proxy.value.smtp_password = $event"
                    label="Password"
                    v-bind="schema('smtp_password')"
                  />
                </v-col>
              </v-row>
            </v-container>
            <v-divider></v-divider>
            <v-card-actions class="d-flex justify-space-between">
              <EmailSettingsTestConnection
                :disabled="!isPristine || isUpdating"
                v-model:testing="status.testing"
                v-model:connectionOK="status.connectionOK"
              />
              <div>
                <v-btn color="" @click="cancel()" :disabled="isPristine || isUpdating">
                  Reset changes
                </v-btn>
                <v-btn
                  text="Save"
                  @click="updateSettings({ body: updateParams(proxy.value) }, { onSuccess: save })"
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
import { $UpsertMailingParams, Mailing, UpsertMailingParams } from '@/api'
import {
  getEmailSettingsOptions,
  toggleMailingMutation,
  updateEmailSettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import PasswordField from '@/components/toolkit/forms/PasswordField.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import { useSchemaBinding } from '@/composables/schema'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import EmailSettingsTestConnection from './EmailSettingsTestConnection.vue'
import { DateTime } from 'luxon'

const status = ref<{
  testing: boolean
  connectionOK?: boolean
}>({
  testing: false,
  connectionOK: undefined
})

const { data: model, error, isPending, refetch } = useQuery(getEmailSettingsOptions())
const { mutateAsync: _toggleMailing } = useMutation(toggleMailingMutation())
async function toggleMailing(enabled: boolean | null) {
  if (enabled === null) return
  await _toggleMailing(
    { body: enabled },
    {
      onSuccess: () => {
        feedback({ message: `E-mail sending ${enabled ? 'enabled' : 'disabled'}`, type: 'success' })
        refetch()
      }
    }
  )
}

const { schema } = useSchemaBinding($UpsertMailingParams)
const { feedback } = useFeedback()

const {
  mutateAsync: updateSettings,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateEmailSettingsMutation(),
  onSuccess: async () => {
    await refetch()
    status.value.connectionOK = true
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onMutate() {
    status.value.testing = true
    status.value.connectionOK = undefined
  },
  onSettled() {
    status.value.testing = false
  }
})

function updateParams({ enabled, last_updated, ...rest }: Mailing): UpsertMailingParams {
  return rest
}
</script>

<style scoped></style>
