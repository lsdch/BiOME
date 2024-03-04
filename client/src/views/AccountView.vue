<template>
  <v-container>
    <h1>Account settings</h1>
    <v-divider class="my-5" />
    <v-row>
      <v-col>
        <v-text-field label="E-mail" prepend-inner-icon="mdi-at" />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-expansion-panels>
          <v-expansion-panel>
            <v-expansion-panel-title>
              <v-icon>mdi-lock</v-icon>
              <h3>Password</h3>
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-form>
                <v-text-field label="Your current password" density="compact" />
                <PasswordFields />
                <v-btn color="primary" type="submit" variant="plain"> Set new password </v-btn>
              </v-form>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import PasswordFields from '@/components/auth/PasswordFields.vue'
</script>

<style scoped></style>
