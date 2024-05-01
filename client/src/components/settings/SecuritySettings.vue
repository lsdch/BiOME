<template>
  <SettingsForm
    :get="SettingsService.securitySettings"
    :update="SettingsService.updateSecuritySettings"
    :schema="$SecuritySettings"
    #="{ model, schema }"
  >
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
              v-bind="schema('min_password_strength')"
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
          v-bind="schema('jwt_secret_key')"
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
          v-bind="schema('auth_token_lifetime')"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <NumberInput
          v-model="model.account_token_lifetime"
          label="Account token lifetime (hours)"
          :step="1"
          v-bind="schema('account_token_lifetime')"
        />
      </v-col>
    </v-row>
  </SettingsForm>
</template>

<script setup lang="ts">
import { $SecuritySettings, SettingsService } from '@/api'
import { ref } from 'vue'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import SettingsForm from './SettingsForm.vue'

const copySecretKeyinitial = {
  icon: 'mdi-content-copy',
  color: 'primary'
}
const copySecretKey = ref(copySecretKeyinitial)

async function toClipboard(text: string) {
  await navigator.clipboard.writeText(text)
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
