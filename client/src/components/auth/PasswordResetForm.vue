<template>
  <v-form validate-on="blur">
    <p>Please input your account email address below.</p>
    <p class="mb-5">An email with a link to reset your password will be sent to this address.</p>
    <v-text-field
      name="email"
      label="Email"
      type="email"
      v-model="state.email"
      variant="outlined"
      prepend-inner-icon="mdi-at"
      :error-messages="v$.$errors.map((e) => String(e.$message))"
    />
    <v-btn
      block
      size="large"
      rounded="sm"
      color="primary"
      text="Request password reset"
      class="mb-5"
      @click="submit()"
    />
  </v-form>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { required, email } from '@vuelidate/validators'
import useVuelidate from '@vuelidate/core'

const state = reactive({
  email: ''
})

const rules = {
  email: { required, email }
}

const v$ = useVuelidate(rules, state)

function submit() {
  v$.value.$validate()
}
</script>

<style scoped></style>
