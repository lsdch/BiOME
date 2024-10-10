<template>
  <SettingsForm @reset="reloadSettings" @submit="submit">
    <v-row>
      <v-col>
        <v-card
          class="flex-grow-1"
          :max-width="600"
          subtitle="Minimum password strength"
          variant="text"
        >
          <v-card-text>
            <v-slider
              v-model="model.min_password_strength"
              :step="1"
              :min="3"
              :max="5"
              :ticks="{ 3: 'Medium', 4: 'Strong', 5: 'Very strong' }"
              show-ticks="always"
              v-bind="field('min_password_strength')"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <NumberInput
          v-model="model.refresh_token_lifetime"
          label="User session lifetime (hours)"
          :step="1"
          v-bind="field('refresh_token_lifetime')"
        />
      </v-col>
    </v-row>
  </SettingsForm>
</template>

<script setup lang="ts">
import { $SecuritySettingsInput, SettingsService } from '@/api'
import { handleErrors } from '@/api/responses'
import { useFeedback } from '@/stores/feedback'
import { useClipboard } from '@vueuse/core'
import { ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import SettingsForm from './SettingsForm.vue'

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

<style scoped></style>
