<template>
  <v-container class="fill-height">
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
        <v-alert variant="text" type="error"> Password reset token is invalid or expired. </v-alert>
      </v-card-text>
      <v-card-text v-else-if="status === Status.ServerError">
        <v-alert variant="text" type="error">
          Failed to validate the password reset token. This is likely due to a server error.
        </v-alert>
      </v-card-text>
      <v-card-text v-else-if="status === Status.PasswordResetDone">
        <v-alert variant="text" type="success">
          Your password was successfully reset! You are now being redirected to the login page...
        </v-alert>
      </v-card-text>
      <v-card-text v-else> Waiting for token validation... </v-card-text>
      <div class="d-flex justify-center mt-3">
        <v-btn class="text-none" :to="{ name: 'login' }" variant="plain" :ripple="false">
          Back to login page
        </v-btn>
      </div>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { PasswordResetInputBody } from '@/api'
import {
  resetPasswordMutation,
  verifyPasswordResetTokenMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import PasswordFields from '@/components/users/PasswordFields.vue'
import { useMutation } from '@tanstack/vue-query'
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

const state = ref<{
  password: PasswordResetInputBody['new_password']
  password_confirmation: PasswordResetInputBody['new_password']
}>({
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

const { mutateAsync: verifyPasswordResetToken } = useMutation(verifyPasswordResetTokenMutation())
async function validateToken(token: string) {
  await verifyPasswordResetToken(
    { body: { token } },
    {
      onError: (error) => {
        switch (error?.status) {
          case 500:
            status.value = Status.ServerError
            break
          default:
            status.value = Status.InvalidToken
        }
      },
      onSuccess: () => {
        status.value = Status.TokenOK
      }
    }
  )
}

const { mutateAsync: resetPassword } = useMutation(resetPasswordMutation())
async function submit() {
  loading.value = true
  await resetPassword(
    { body: { token: token.value ?? '', new_password: state.value.password } },
    {
      onError: (error) => {
        switch (error?.status) {
          case 500:
            status.value = Status.ServerError
            break
          default:
            status.value = Status.InvalidToken
        }
      },
      onSuccess: () => {
        status.value = Status.PasswordResetDone
        setTimeout(() => {
          router.replace({ name: 'login' })
        }, 3000)
      },
      onSettled: () => {
        loading.value = false
      }
    }
  )
}
</script>

<style scoped></style>
