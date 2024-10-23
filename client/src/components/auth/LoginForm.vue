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
    <div v-if="error?.status === 401" class="text-red text-center mb-3">Invalid credentials</div>
    <div v-else-if="error" class="text-red">An error occurred on the server</div>

    <v-btn
      block
      type="submit"
      size="large"
      rounded="sm"
      color="primary"
      text="Log in"
      class="mb-5"
      :loading="loading"
    />
  </v-form>
</template>

<script setup lang="ts">
import { ErrorModel, UserCredentials } from '@/api'
import { useUserStore } from '@/stores/user'
import { useRouteQuery } from '@vueuse/router'
import { onBeforeMount, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PasswordField from '../toolkit/ui/PasswordField.vue'
import { useToggle } from '@vueuse/core'

const credentials: UserCredentials = reactive({
  identifier: '',
  password: ''
})

const error = ref<ErrorModel>()

const [loading, toggleLoading] = useToggle(false)

const router = useRouter()
const target = useRouteQuery<string>('redirect', '/')
const { isAuthenticated, login } = useUserStore()

onBeforeMount(() => {
  if (isAuthenticated) redirect()
})

function redirect() {
  // Using replace to overwrite router history
  router.replace({ path: target.value })
  return
}

async function submit() {
  toggleLoading(true)
  await login(credentials)
    .then((err) => {
      error.value = err
      if (err === undefined) redirect()
    })
    .finally(() => toggleLoading(false))
}
</script>

<style scoped></style>
