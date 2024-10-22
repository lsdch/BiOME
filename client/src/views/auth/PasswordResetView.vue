<template>
  <v-container class="fill-height">
    <v-row :align="smAndDown ? 'baseline' : 'center'">
      <v-col lg="6" offset-lg="3">
        <v-card
          :variant="smAndDown ? 'flat' : 'elevated'"
          :disabled="status === Status.ValidationPending"
          :loading="status === Status.ValidationPending ? 'primary' : false"
          :title="
            status == Status.TokenOK ? 'Please set up a new password for your account' : undefined
          "
        >
          <v-card-text v-if="status === Status.TokenOK">
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
import { AccountService, PasswordInput } from '@/api'
import PasswordFields from '@/components/auth/PasswordFields.vue'
import { useRouteQuery } from '@vueuse/router'
import { onBeforeMount, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()

const router = useRouter()
const token = useRouteQuery<string | undefined>('token', undefined, {
  transform(v) {
    if (Array.isArray(v)) return undefined
    return v
  }
})

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
  if (token.value === undefined) router.replace({ name: 'home' })
  else validateToken(token.value)
})

async function validateToken(token: string) {
  AccountService.validatePasswordToken({ query: { token } }).then(({ error }) => {
    if (error?.status === 500) {
      status.value = Status.ServerError
    } else if (error) {
      status.value = Status.InvalidToken
    } else {
      status.value = Status.TokenOK
    }
  })
}

async function submit() {
  loading.value = true
  AccountService.resetPassword({ query: { token: token.value as string }, body: state.value })
    .then(({ error }) => {
      if (error?.status === 500) {
        status.value = Status.ServerError
      } else if (error) {
        status.value = Status.InvalidToken
      } else {
        status.value = Status.PasswordResetDone
      }
      setTimeout(() => {
        router.replace({ name: 'login' })
      }, 3000)
    })
    .finally(() => (loading.value = false))
}
</script>

<style scoped></style>
