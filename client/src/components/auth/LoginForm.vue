<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="state.identifier"
      name="login"
      label="Account"
      placeholder="Login or email"
      variant="outlined"
      prepend-inner-icon="mdi-account"
    />
    <v-text-field
      v-model="state.password"
      type="password"
      label="Password"
      name="password"
      variant="outlined"
      prepend-inner-icon="mdi-lock"
    />
    <div v-if="feedback" class="mb-3 w-100 text-center">
      <v-alert v-if="feedback == 'ConfirmationSent'" type="info" variant="outlined">
        A new activation link was sent to your email address.
      </v-alert>
      <v-alert v-else-if="feedback == 'ServerError'" type="error" variant="outlined">
        An error occurred on the server.
      </v-alert>
      <span v-else-if="feedback?.reason == 'InvalidCredentials'" class="text-red">
        Invalid credentials
      </span>
      <v-alert v-else-if="feedback?.reason == 'Inactive'" type="warning" variant="outlined">
        Your email was not confirmed yet, please check your inbox<br />
        or <a href="#" @click="resendConfirmation">request another confirmation link</a> if needed.
      </v-alert>
    </div>
    <v-btn
      block
      type="submit"
      size="large"
      rounded="sm"
      color="primary"
      text="Log in"
      class="mb-5"
      @click="submit"
      :loading="loading"
    />
  </v-form>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { ApiError, AuthService, LoginFailedError, UserCredentials } from '@/api'
import { ref } from 'vue'
import { Ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

const state: UserCredentials = reactive({
  identifier: '',
  password: '',
  remember: true
})

const loading = ref(false)

const feedback: Ref<undefined | LoginFailedError | 'ConfirmationSent' | 'ServerError'> =
  ref(undefined)

defineEmits<{ (event: 'sendConfirmation'): void }>()

const { getUser } = useUserStore()
const router = useRouter()

async function submit() {
  loading.value = true
  await AuthService.login(state)
    .then(() => {
      feedback.value = undefined
      getUser()
      router.push({ name: 'home' })
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          feedback.value = err.body
          break
        default:
          feedback.value = 'ServerError'
          break
      }
      return undefined
    })
    .finally(() => {
      loading.value = false
    })
}

function resendConfirmation() {
  AuthService.resendConfirmationEmail(state)
    .then(() => {
      feedback.value = 'ConfirmationSent'
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          feedback.value = err.body
          break

        default:
          feedback.value = 'ServerError'
          break
      }
    })
}
</script>

<style scoped></style>
