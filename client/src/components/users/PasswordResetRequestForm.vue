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
        v-model="email"
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
        :disabled="!isValid"
      />
    </template>
  </v-form>
</template>

<script setup lang="ts">
import { $PasswordResetRequestInputBody, PasswordResetRequestInputBody } from '@/api'
import { ref } from 'vue'
import { useSchemaBinding } from '@/composables/schema'
import { useRouter } from 'vue-router'
import { requestPasswordResetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation } from '@tanstack/vue-query'

const router = useRouter()

const email = ref<string>('')
const { schema } = useSchemaBinding($PasswordResetRequestInputBody)

const requestSent = ref(false)

const { mutateAsync } = useMutation(requestPasswordResetMutation())
async function submit() {
  await mutateAsync({
    body: { email: email.value, client_path: router.resolve('password-reset').path }
  }).then(() => {
    requestSent.value = true
  })
}
</script>

<style scoped></style>
