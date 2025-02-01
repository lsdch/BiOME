<template>
  <v-form @submit.prevent="submit" validate-on="blur lazy">
    <v-row>
      <v-col>
        <PasswordField v-model="model.password" label="Your current password" />
      </v-col>
    </v-row>
    <PasswordFields v-model="model.new_password" />
    <v-row>
      <v-col>
        <v-btn
          color="primary"
          type="submit"
          variant="plain"
          :loading="isPending"
          text="Set new password"
        />
      </v-col>
    </v-row>
    <v-alert v-if="error" type="error">
      {{
        error.status == StatusCodes.UNPROCESSABLE_ENTITY ? 'Invalid password' : 'An error occurred'
      }}
    </v-alert>
  </v-form>
</template>

<script setup lang="ts">
import { UpdatePasswordInput } from '@/api'
import { updatePasswordMutation } from '@/api/gen/@tanstack/vue-query.gen'
import PasswordFields from '@/components/auth/PasswordFields.vue'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
import { StatusCodes } from 'http-status-codes'
import { ref } from 'vue'

const model = ref<UpdatePasswordInput>({
  password: '',
  new_password: {
    password: '',
    password_confirmation: ''
  }
})

const { feedback } = useFeedback()

const { mutate, isPending, error } = useMutation({
  ...updatePasswordMutation(),
  onSuccess: () => {
    feedback({ type: 'success', message: 'Password updated' })
  },
  onError: (error) => {
    console.error(error)
    feedback({ type: 'error', message: 'Password update failed' })
  }
})

function submit() {
  mutate({ body: model.value })
}
</script>

<style scoped></style>
