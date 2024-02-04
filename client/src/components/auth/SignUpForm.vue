<template>
  <v-form @submit.prevent="submit">
    <v-text-field
      v-model="state.login"
      prepend-inner-icon="mdi-account"
      name="login"
      label="Login"
      hint="Only used for authentication, not visible to other users"
      persistent-hint
      class="mb-3"
      required
      :error-messages="registerError('login')"
    />
    <v-text-field
      v-model="state.email"
      prepend-inner-icon="mdi-at"
      name="email"
      type="email"
      label="Email address"
      @blur="v$.email.$touch"
      :error-messages="registerError('email')"
      class="mb-3"
    />
    <PasswordFields v-model="state" :user-inputs="passwordSensitiveInputs" />

    <v-row>
      <v-col>
        <h2 class="text-subtitle-1">Personal informations</h2>
        <p>
          Your name will be publicly displayed only if you contribute to the database. This is
          required to ensure full traceability of the data.
        </p>
        <p>
          Otherwise it will be visible on your user profile, which is only accessible by yourself
          and administrators.
        </p>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="state.identity.first_name"
          name="first-name"
          label="First name"
          required
          @blur="v$.identity.first_name.$touch"
          :error-messages="vuelidateErrors(v$.identity.first_name)"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="state.identity.last_name"
          name="last-name"
          label="Last name"
          required
          @blur="v$.identity.last_name.$touch()"
          :error-messages="vuelidateErrors(v$.identity.last_name)"
        />
      </v-col>
    </v-row>
    <v-btn
      type="submit"
      color="primary"
      text="Create account"
      block
      size="large"
      rounded="md"
      :loading="loading"
      :disabled="v$.$invalid || v$.password.$error || v$.password_confirmation.$error"
    />
    <v-alert
      v-if="unhandledError"
      class="my-2"
      type="error"
      title="Error"
      text="Some unhandled error occurred on the server"
    />
  </v-form>
</template>

<script setup lang="ts">
import { ApiError, AuthService, InputValidationError, UserInput } from '@/api'
import { vuelidateErrors } from '@/api/validation'
import useVuelidate, { Validation, ValidationArgs } from '@vuelidate/core'
import { email, maxLength, minLength, required } from '@vuelidate/validators'
import { Ref, computed, ref } from 'vue'
import PasswordFields from './PasswordFields.vue'

const state: Ref<UserInput> = ref<UserInput>({
  email: '',
  login: '',
  password: '',
  password_confirmation: '',
  identity: {
    first_name: '',
    last_name: '',
    contact: undefined
  }
})

const rules: ValidationArgs<UserInput> = {
  email: { email, required },
  login: { minLength: minLength(5), maxLength: maxLength(16) },
  password: {},
  password_confirmation: {},
  identity: {
    first_name: { minLength: minLength(2), required },
    last_name: { minLength: minLength(2), required },
    contact: {}
  }
}

const v$: Ref<Validation<ValidationArgs<UserInput>, UserInput>> = useVuelidate(rules, state)

const passwordSensitiveInputs = computed(() => {
  return [
    state.value.email,
    state.value.login,
    state.value.identity.first_name,
    state.value.identity.last_name
  ]
})

const loading = ref(false)
const emits = defineEmits<{ (event: 'created'): void }>()
const errors: Ref<Record<string, InputValidationError[]>> = ref({})

function registerError(field: string) {
  return computed(() => {
    return errors.value?.[field]?.map((e) => e.message) ?? vuelidateErrors(v$.value[field])
  }).value
}

const unhandledError: Ref<undefined | ApiError> = ref(undefined)

async function submit() {
  loading.value = true
  await AuthService.registerUser(state.value)
    .then(() => {
      emits('created')
    })
    .catch((err: ApiError) => {
      switch (err.status) {
        case 400:
          errors.value = err.body
          break
        default:
          unhandledError.value = err
          break
      }
    })
    .finally(() => {
      loading.value = false
    })
}
</script>

<style scoped></style>
