<template>
  <v-container>
    <!-- <HomeLinkTitle /> -->
    <v-row>
      <v-col cols="12" lg="8" offset-lg="2">
        <v-card
          variant="flat"
          :title="registrationDone ? 'New account registered' : 'Sign up for an account'"
        >
          <template v-slot:prepend>
            <v-icon
              :icon="registrationDone ? 'mdi-check-bold' : 'mdi-account-plus'"
              :color="registrationDone ? 'green' : 'primary'"
            />
          </template>
          <v-card-text>
            <div v-if="registrationDone">
              An email was sent to your address with a link to activate your account. Please check
              your inbox.
            </div>
            <SignUpForm v-else @created="registrationDone = true" />
            <div class="d-flex justify-center mt-3">
              <v-btn class="text-none" :to="{ name: 'login' }" variant="plain" :ripple="false">
                Back to login page
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import SignUpForm from '@/components/auth/SignUpForm.vue'
import { ref } from 'vue'

const registrationDone = ref(false)
</script>

<style scoped></style>
