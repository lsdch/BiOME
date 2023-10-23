<template>
  <v-form>
    <v-text-field
      v-model="state.identifier"
      name="login"
      label="Account"
      placeholder="Login or email"
      variant="outlined"
      prepend-inner-icon="mdi-account"
      :error-messages="v$.identifier.$errors.map((e) => String(e.$message))"
    />
    <v-text-field
      v-model="state.password"
      type="password"
      label="Password"
      name="password"
      variant="outlined"
      prepend-inner-icon="mdi-lock"
      :error-messages="v$.password.$errors.map((e) => String(e.$message))"
    />

    <v-btn
      block
      size="large"
      rounded="sm"
      color="primary"
      text="Log in"
      class="mb-5"
      @click="submit"
    />
  </v-form>
</template>

<script setup lang="ts">
import { required } from '@vuelidate/validators'
import useVuelidate from '@vuelidate/core'
import { reactive } from 'vue'
import { UserCredentials } from '@/api'

const state: UserCredentials = reactive({
  identifier: '',
  password: '',
  remember: true
})

const rules = {
  identifier: { required },
  password: { required },
  remember: { required }
}

const v$ = useVuelidate(rules, state)

function submit() {
  v$.value.$validate()
}
</script>

<style scoped></style>
