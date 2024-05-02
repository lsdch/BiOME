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
        <v-btn color="primary" type="submit" variant="plain"> Set new password </v-btn>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { AccountService, UpdatePasswordInput } from '@/api'
import PasswordFields from '@/components/auth/PasswordFields.vue'
import PasswordField from '@/components/toolkit/ui/PasswordField.vue'
import { useFeedback } from '@/stores/feedback'
import { ref } from 'vue'

const model = ref<UpdatePasswordInput>({
  password: '',
  new_password: {
    password: '',
    password_confirmation: ''
  }
})

const { feedback } = useFeedback()

async function submit() {
  AccountService.updatePassword({ requestBody: model.value })
    .then(() => {
      feedback({ type: 'success', message: 'Password updated' })
    })
    .catch(() => {
      feedback({ type: 'error', message: 'Password update failed' })
    })
}
</script>

<style scoped></style>
