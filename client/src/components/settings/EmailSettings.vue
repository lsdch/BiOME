<template>
  <SettingsForm
    :get="SettingsService.emailSettings"
    :update="SettingsService.updateEmailSettings"
    :schema="$EmailSettings"
  >
    <template #prepend-toolbar="">
      <v-btn color="success" text="Test connection" variant="text" prepend-icon="mdi-network" />
    </template>
    <template #default="{ model, schema }">
      <v-container>
        <v-row>
          <v-col cols="12" sm="8">
            <v-text-field v-model="model.host" label="SMTP Host" v-bind="schema('host')" />
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
            <v-text-field v-model="model.user" label="User" v-bind="schema('user')" />
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
import NumberInput from '../toolkit/ui/NumberInput.vue'
import SettingsForm from './SettingsForm.vue'
</script>

<style scoped></style>
