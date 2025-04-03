<template>
  <div v-if="requestSent">
    An e-mail was sent to the address you provided, if it was found in our database.
  </div>
  <v-form v-else validate-on="input" @submit.prevent="submit">
    <template #default="{ isValid }">
      <p>Please input your account email address below.</p>
      <p class="mb-5">An email with a link to reset your password will be sent to this address.</p>
      <v-text-field
        name="email"
        label="Email"
        type="email"
        class="mb-3"
        v-model="state.email"
        variant="outlined"
        prepend-inner-icon="mdi-at"
        v-bind="schema('email')"
      />
      <v-btn
        block
        size="large"
        rounded="sm"
        color="primary"
        text="Request password reset"
        class="mb-5"
        type="submit"
        :disabled="!isValid.value"
      />
    </template>
  </v-form>
</template>

<script setup lang="ts">
import { $PasswordResetRequest, AccountService, PasswordResetRequest } from '@/api'
import { ref } from 'vue'
import { useSchema } from '../../composables/schema'
import { useRouter } from 'vue-router'

const router = useRouter()

const state = ref<PasswordResetRequest>({
  email: '',
  handler: window.location.origin + router.resolve('password-reset').fullPath
})
const {
  bind: { schema }
} = useSchema($PasswordResetRequest)

const requestSent = ref(false)

async function submit() {
  await AccountService.requestPasswordReset({ body: state.value }).then(() => {
    requestSent.value = true
  })
}
</script>

<style scoped></style>
