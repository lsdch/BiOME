<template>
  <div v-if="requestAccepted">
    An email was sent to <code>{{ state.email }}</code> with a link to reset your password.
  </div>
  <v-form v-else validate-on="blur" @submit.prevent="submit">
    <p>Please input your account email address below.</p>
    <p class="mb-5">An email with a link to reset your password will be sent to this address.</p>
    <v-text-field
      name="email"
      label="Email"
      type="email"
      v-model="state.email"
      variant="outlined"
      prepend-inner-icon="mdi-at"
      :error-messages="error"
    />
    <v-btn
      block
      size="large"
      rounded="sm"
      color="primary"
      text="Request password reset"
      class="mb-5"
      type="submit"
    />
  </v-form>
</template>

<script setup lang="ts">
import { AccountService, ApiError } from '@/api'
import useVuelidate from '@vuelidate/core'
import { email, required } from '@vuelidate/validators'
import { Ref, ref } from 'vue'

const state = ref({
  email: ''
})

const rules = {
  email: { required, email }
}

const requestAccepted = ref(false)
const error: Ref<string | undefined> = ref(undefined)

const v$ = useVuelidate(rules, state)

async function submit() {
  v$.value.$validate()
  await AccountService.requestPasswordReset({ requestBody: state.value })
    .then(() => {
      requestAccepted.value = true
    })
    .catch((err: ApiError) => {
      console.log(err)
      switch (err.status) {
        case 400:
          error.value = 'Invalid email address'
          break
        case 404:
          error.value = 'No account matches the provided email address'
          break
        default:
          break
      }
    })
}
</script>

<style scoped></style>
