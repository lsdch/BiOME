<template>
  <SettingsForm :get="SettingsService.emailSettings" :update="SettingsService.updateEmailSettings">
    <template #prepend-toolbar="">
      <v-btn color="success" text="Test connection" variant="text" prepend-icon="mdi-network" />
    </template>
    <template #default="{ model }">
      <v-container>
        <v-row>
          <v-col cols="12" sm="8">
            <v-text-field v-model="model.host" label="SMTP Host" />
          </v-col>
          <v-col cols="12" sm="4">
            <VNumberInput
              v-model="model.port"
              label="SMTP Port"
              :min="1"
              control-variant="stacked"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field v-model="model.user" label="User" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <PasswordField v-model="model.password" label="Password" />
          </v-col>
        </v-row>
      </v-container>
    </template>
  </SettingsForm>
</template>

<script setup lang="ts">
import { SettingsService } from '@/api'
import SettingsForm from './SettingsForm.vue'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
</script>

<style scoped></style>
