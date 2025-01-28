<template>
  <CenteredSpinner v-if="isPending" text="Loading security settings..." />
  <v-alert v-else-if="error" color="error" icon="mdi-alert">
    Failed to load security settings
  </v-alert>
  <v-confirm-edit v-else v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @reset="cancel()"
        @submit="submit(proxy.value).then(() => save())"
        :loading="isUpdating"
      />
      <v-container>
        <v-row>
          <v-col>
            <v-alert v-if="updateError" color="error" icon="mdi-alert">
              Failed to update settings
            </v-alert>
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-list class="settings-list">
              <v-list-item title="Minimum password strength">
                <v-slider
                  class="mx-4"
                  v-model="proxy.value.min_password_strength"
                  :step="1"
                  :min="3"
                  :max="5"
                  :thumb-size="15"
                  :ticks="{ 3: 'Medium', 4: 'Strong', 5: 'Very strong' }"
                  show-ticks="always"
                  v-bind="field('min_password_strength')"
                />
              </v-list-item>
              <v-divider />
              <v-list-item title="User session lifetime">
                <DaysHoursInput
                  v-model="proxy.value.refresh_token_lifetime"
                  v-bind="field('refresh_token_lifetime')"
                />
              </v-list-item>
              <v-divider />
              <v-list-item title="Registration settings">
                <VNumberInput
                  v-model="proxy.value.invitation_token_lifetime"
                  label="User invitation lifetime (days)"
                  :step="1"
                  v-bind="field('invitation_token_lifetime')"
                  control-variant="stacked"
                />
              </v-list-item>
            </v-list>
          </v-col>
        </v-row>
      </v-container>
    </template>
  </v-confirm-edit>
</template>

<script setup lang="ts">
import { $SecuritySettingsInput, SecuritySettingsInput } from '@/api'
import {
  securitySettingsOptions,
  updateSecuritySettingsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { useFeedback } from '@/stores/feedback'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { useSchema } from '../toolkit/forms/schema'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import DaysHoursInput from './DaysHoursInput.vue'
import SettingsFormActions from './SettingsFormActions.vue'

const { feedback } = useFeedback()

const { field, dispatchErrors } = useSchema($SecuritySettingsInput)

const { data: model, refetch, error, isPending } = useQuery(securitySettingsOptions())

const {
  mutateAsync,
  error: updateError,
  isPending: isUpdating
} = useMutation({
  ...updateSecuritySettingsMutation(),
  onSuccess: () => {
    feedback({ message: 'Updated settings', type: 'success' })
  },
  onError: dispatchErrors
})

async function submit(model: SecuritySettingsInput) {
  await mutateAsync({ body: model })
}
</script>

<style lang="scss">
.settings-list {
  .v-list-item-title {
    margin-top: 10px;
    margin-bottom: 20px;
  }
}
.v-slider-track__tick-label {
  font-size: small;
}
</style>
