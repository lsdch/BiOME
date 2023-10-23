<template>
  <v-container>
    <v-row>
      <v-col lg="6" offset-lg="3">
        <v-card>
          <v-card-title primary-title> Sign up </v-card-title>
          <v-card-text>
            <v-form>
              <v-text-field
                v-model="state.email"
                prepend-inner-icon="mdi-at"
                name="email"
                type="email"
                label="Email address"
              />
              <v-text-field
                v-model="state.password"
                name="password"
                type="password"
                label="Password"
              />
              <v-text-field
                v-model="state.password_confirmation"
                name="password-confirm"
                type="password"
                label="Password confirmation"
                required
              />
              <v-row>
                <v-col cols="12" sm="6">
                  <v-text-field name="first-name" label="First name" />
                </v-col>
                <v-col cols="12" sm="6">
                  <v-text-field name="last-name" label="Last name" />
                </v-col>
              </v-row>
              <v-btn color="primary" text="Create account" block size="large" rounded="md" />
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { UserInput } from '@/api'
import { reactive } from 'vue'
import { email, required, sameAs, minLength } from '@vuelidate/validators'

const state: UserInput = reactive({
  email: '',
  password: '',
  password_confirmation: '',
  identity: {
    first_name: '',
    last_name: '',
    contact: ''
  }
})
</script>

<style scoped>
div[aria-required='true'].v-input .v-label::after {
  content: ' *';
  color: red;
}
</style>
