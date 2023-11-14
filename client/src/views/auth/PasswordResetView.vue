<template>
  <v-container class="fill-height">
    <v-row :align="smAndDown ? 'baseline' : 'center'">
      <v-col lg="6" offset-lg="3">
        <HomeLinkTitle />
        <v-card
          :variant="smAndDown ? 'flat' : 'elevated'"
          :disabled="status === Status.ValidationPending"
          :loading="status === Status.ValidationPending ? 'primary' : false"
        >
          <v-card-text v-if="status === Status.TokenOK">
            <p class="mb-3">Please set up a new password for your account</p>
            <v-form @submit.prevent="submit">
              <PasswordFields v-model="state" :user-inputs="[]" />
              <v-btn block type="submit" text="Confirm" :rounded="false" :loading="loading" />
            </v-form>
          </v-card-text>
          <v-card-text v-else-if="status === Status.InvalidToken">
            <v-alert variant="text" type="error">
              Password reset token is invalid or expired.
            </v-alert>
          </v-card-text>
          <v-card-text v-else-if="status === Status.ServerError">
            <v-alert variant="text" type="error">
              Failed to validate the password reset token. This is likely due to a server error.
            </v-alert>
          </v-card-text>
          <v-card-text v-else-if="status === Status.PasswordResetDone">
            <v-alert variant="text" type="success">
              Your password was successfully reset! You are now being redirected to the login
              page...
            </v-alert>
          </v-card-text>
          <v-card-text v-else> Waiting for token validation... </v-card-text>
          <div class="d-flex justify-center mt-3">
            <v-btn class="text-none" :to="{ name: 'login' }" variant="plain" :ripple="false">
              Back to login page
            </v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ApiError, AuthService, PasswordInput } from '@/api'
import PasswordFields from '@/components/auth/PasswordFields.vue'
import HomeLinkTitle from '@/components/navigation/HomeLinkTitle.vue'
import { onBeforeMount } from 'vue'
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()

const route = useRoute()
const router = useRouter()
const token = route.params.token.toString()

const state = ref<PasswordInput>({
  password: '',
  password_confirmation: ''
})

const loading = ref(false)

enum Status {
  ValidationPending,
  TokenOK,
  InvalidToken,
  ServerError,
  PasswordResetDone
}
const status = ref(Status.ValidationPending)

onBeforeMount(() => {
  validateToken(token)
})

async function validateToken(token: string) {
  AuthService.validatePasswordToken(token)
    .then(() => {
      status.value = Status.TokenOK
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          status.value = Status.InvalidToken
          break

        default:
          status.value = Status.ServerError
          break
      }
    })
}

async function submit() {
  loading.value = true
  AuthService.resetPassword(token, state.value)
    .then(() => {
      status.value = Status.PasswordResetDone
      setTimeout(() => {
        router.push({ name: 'login' })
      }, 3000)
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          status.value = Status.InvalidToken
          break

        default:
          status.value = Status.ServerError
          break
      }
    })
}
</script>

<style scoped></style>
