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
        <v-text-field
          v-model="model.jwt_secret_key"
          class="flex-grow-1"
          label="Authentication token secret key"
          prepend-inner-icon="mdi-key"
          readonly
          v-bind="field('jwt_secret_key')"
        >
          <template #append-inner>
            <v-btn
              variant="plain"
              v-bind="copySecretKey"
              @click="toClipboard(model.jwt_secret_key)"
            />
            <v-btn variant="plain" icon="mdi-refresh" />
          </template>
        </v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6"> </v-col>
      <v-col cols="12" sm="6">
        <NumberInput
          v-model="model.auth_token_lifetime"
          class="flex-grow-0"
          label="User session lifetime (minutes)"
          :step="15"
          v-bind="field('auth_token_lifetime')"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <NumberInput
          v-model="model.account_token_lifetime"
          label="Account token lifetime (hours)"
          :step="1"
          v-bind="field('account_token_lifetime')"
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

const copySecretKeyinitial = {
  icon: 'mdi-content-copy',
  color: 'primary'
}
const copySecretKey = ref(copySecretKeyinitial)
const { copy } = useClipboard()
async function toClipboard(text: string) {
  await copy(text)
  copySecretKey.value = {
    icon: 'mdi-check-circle',
    color: 'success'
  }
  setTimeout(() => {
    copySecretKey.value = copySecretKeyinitial
  }, 2000)
}
</script>

<style scoped></style>
