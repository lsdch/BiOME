<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="credentials.identifier"
      name="login"
      label="Account"
      placeholder="Login or email"
      variant="outlined"
      prepend-inner-icon="mdi-account"
    />
    <PasswordField
      v-model="credentials.password"
      label="Password"
      name="password"
      prepend-inner-icon="mdi-lock"
    />
    <div
      v-if="loginState.error?.status === StatusCodes.UNAUTHORIZED"
      class="text-red text-center mb-3"
    >
      Invalid credentials
    </div>
    <div v-else-if="loginState.error" class="text-red">An error occurred on the server</div>

    <v-btn
      block
      type="submit"
      size="large"
      rounded="sm"
      color="primary"
      text="Log in"
      class="mb-5"
      :loading="loginState.pending"
    />
  </v-form>
</template>

<script setup lang="ts">
import { UserCredentials } from '@/api'
import { useUserStore } from '@/stores/user'
import { useRouteQuery } from '@vueuse/router'
import { StatusCodes } from 'http-status-codes'
import { onBeforeMount, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import PasswordField from '../toolkit/forms/PasswordField.vue'
import { storeToRefs } from 'pinia'

const credentials: UserCredentials = reactive({
  identifier: '',
  password: ''
})

const userStore = useUserStore()
const { isAuthenticated, loginState } = storeToRefs(userStore)
const { login } = userStore

function submit() {
  login(credentials)
}

const router = useRouter()
watch(isAuthenticated, (ok) => {
  if (ok) redirect()
})

/**
 * Target page to redirect to after successful login,
 * or if user is already authenticated
 */
const target = useRouteQuery<string>('redirect', '/')

// If user is already authenticated, redirect to target
onBeforeMount(() => {
  if (isAuthenticated.value) redirect()
})

function redirect() {
  // Using replace to overwrite router history
  router.replace({ path: target.value })
  return
}
</script>

<style scoped></style>
