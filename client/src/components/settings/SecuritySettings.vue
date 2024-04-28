<template>
  <SettingsForm
    :get="SettingsService.securitySettings"
    :update="SettingsService.updateSecuritySettings"
    #="{ model }"
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
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="model.jwt_secret_key"
          class="flex-grow-1"
          label="Authentication token secret key"
          :min="32"
          prepend-inner-icon="mdi-key"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <VNumberInput
          v-model="model.auth_token_lifetime"
          class="flex-grow-0"
          label="User session lifetime (seconds)"
          :step="60"
          :min="60"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <VNumberInput
          v-model="model.account_token_lifetime"
          label="Account token lifetime (hours)"
          :step="1"
          :min="1"
        />
      </v-col>
    </v-row>
  </SettingsForm>
</template>

<script setup lang="ts">
import { SettingsService } from '@/api'
import SettingsForm from './SettingsForm.vue'
</script>

<style scoped></style>
