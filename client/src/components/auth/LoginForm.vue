<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="formData.identifier"
      name="login"
      label="Account"
      placeholder="Login or email"
      variant="outlined"
      prepend-inner-icon="mdi-account"
    />
    <PasswordField
      v-model="formData.password"
      label="Password"
      name="password"
      prepend-inner-icon="mdi-lock"
    />
    <!-- <div v-if="errors" class="mb-3 w-100 text-center">
      <v-alert v-if="errors == 'ConfirmationSent'" type="info" variant="outlined">
        A new activation link was sent to your email address.
      </v-alert>
      <v-alert v-else-if="errors == 'ServerError'" type="error" variant="outlined">
        An error occurred on the server.
      </v-alert>
      <span v-else-if="errors?.reason == 'InvalidCredentials'" class="text-red">
        Invalid credentials
      </span>
      <v-alert v-else-if="errors?.reason == 'Inactive'" type="warning" variant="outlined">
        Your email was not confirmed yet, please check your inbox<br />
        or <a href="#" @click="resendConfirmation">request another confirmation link</a> if needed.
      </v-alert>
    </div> -->
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
import { UserCredentials } from '@/api'
import { useUserStore } from '@/stores/user'
import { useRouteQuery } from '@vueuse/router'
import { onBeforeMount, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PasswordField from '../toolkit/ui/PasswordField.vue'

const formData: UserCredentials = reactive({
  identifier: '',
  password: ''
})

const loading = ref(false)

defineEmits<{ (event: 'sendConfirmation'): void }>()

const router = useRouter()
const { getUser, isAuthenticated, login, error } = useUserStore()
const target = useRouteQuery<string | 'home'>('redirect', '/')

onBeforeMount(() => {
  if (isAuthenticated) redirect()
})

function redirect() {
  // Using replace to overwrite router history
  router.replace({ path: target.value })
  return
}

async function submit() {
  loading.value = true
  await login(formData)
    .then(() => {
      if (error !== undefined) {
        return
      }
      redirect()
    })
    .finally(() => {
      loading.value = false
    })
}

// function resendConfirmation() {
//   AccountService.resendEmailConfirmation(state)
//     .then(() => {
//       errors.value = 'ConfirmationSent'
//     })
//     .catch((err: ApiError) => {
//       switch (err.status) {
//         case 400:
//           errors.value = err.body
//           break

//         default:
//           errors.value = 'ServerError'
//           break
//       }
//     })
// }
</script>

<style scoped></style>
