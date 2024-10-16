<template>
  <v-confirm-edit v-model="model">
    <template #default="{ isPristine, save, cancel, model: proxy, actions: _ }">
      <SettingsFormActions
        :model-value="!isPristine"
        @reset="reloadSettings().then(cancel)"
        @submit="save(), submit()"
      />
      <v-container>
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
import { $SecuritySettingsInput, SettingsService } from '@/api'
import { handleErrors } from '@/api/responses'
import { useFeedback } from '@/stores/feedback'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import DaysHoursInput from './DaysHoursInput.vue'
import SettingsFormActions from './SettingsFormActions.vue'

const { feedback } = useFeedback()

async function fetch() {
  return SettingsService.securitySettings().then(
    handleErrors((err) => {
      feedback({ message: 'Failed to retrieve security settings', type: 'error' })
      console.error(err)
    })
  )
}

async function reloadSettings() {
  model.value = await fetch()
}

async function submit() {
  return SettingsService.updateSecuritySettings({ body: model.value })
    .then(errorHandler)
    .then(() => feedback({ message: 'Settings updated', type: 'success' }))
}

const model = ref(await fetch())

const { field, errorHandler } = useSchema($SecuritySettingsInput)
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
