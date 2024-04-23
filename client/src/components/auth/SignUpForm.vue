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
    <PasswordFields v-model="state" :user-inputs="state" />

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
    <PersonFormFields v-model="state.identity" />
    <v-row>
      <v-col cols="12" sm="6">
        <v-text-field
          v-model="state.identity.contact"
          name="contact"
          label="Contact (optional)"
          :disabled="exposeEmail"
          @blur="v$.identity.contact.$touch"
          :error-messages="registerError('identity.contact')"
        />
      </v-col>
      <v-col cols="12" sm="6">
        <v-switch
          label="Expose my account e-mail as contact address"
          v-model="exposeEmail"
          color="primary"
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
      text="An unexpected error occurred on the server"
    />
  </v-form>
</template>

<script setup lang="ts">
import { ApiError, AccountService, InputValidationError, UserInput } from '@/api'
import { vuelidateErrors } from '@/api/validation'
import useVuelidate, { Validation, ValidationArgs } from '@vuelidate/core'
import { email, maxLength, minLength, required } from '@vuelidate/validators'
import { Ref, computed, ref, watchEffect } from 'vue'
import PasswordFields from './PasswordFields.vue'
import PersonFormFields from '../people/PersonFormFields.vue'
import _ from 'lodash'

const exposeEmail = ref(false)

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

watchEffect(() => {
  if (exposeEmail.value) {
    state.value.identity.contact = state.value.email
  } else {
    state.value.identity.contact = undefined
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
    contact: { email }
  }
}

const v$: Ref<Validation<ValidationArgs<UserInput>, UserInput>> = useVuelidate(rules, state)

const loading = ref(false)
const emits = defineEmits<{ (event: 'created'): void }>()
const errors: Ref<Record<string, InputValidationError[]>> = ref({})

function registerError(field: string) {
  return computed(() => {
    return (
      _.get(errors.value, field)?.map((e) => e.message) ?? vuelidateErrors(_.get(v$.value, field))
    )
  }).value
}

const unhandledError: Ref<undefined | ApiError> = ref(undefined)

async function submit() {
  loading.value = true
  await AccountService.registerUser(state.value)
    .then(() => emits('created'))
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
