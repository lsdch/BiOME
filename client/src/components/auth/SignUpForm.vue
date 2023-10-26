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
      :error-messages="vuelidateErrors(v$.login)"
    />
    <v-text-field
      v-model="state.email"
      prepend-inner-icon="mdi-at"
      name="email"
      type="email"
      label="Email address"
      @blur="v$.email.$touch"
      :error-messages="vuelidateErrors(v$.email)"
      class="mb-3"
    />
    <v-text-field
      v-model="state.password"
      name="password"
      label="Password"
      :type="show.pass1 ? 'text' : 'password'"
      required
      @blur="v$.password.$touch"
      :error-messages="vuelidateErrors(v$.password)"
      :append-inner-icon="show.pass1 ? 'mdi-eye' : 'mdi-eye-off'"
      @click:append-inner="show.pass1 = !show.pass1"
    >
      <template v-slot:loader>
        <PasswordStrengthMeter v-if="state.password.length > 0" :strength="passwordStrength" />
      </template>
    </v-text-field>
    <v-text-field
      v-model="state.password_confirmation"
      name="password-confirm"
      label="Password confirmation"
      :type="show.pass2 ? 'text' : 'password'"
      required
      @blur="v$.password_confirmation.$touch"
      :error-messages="vuelidateErrors(v$.password_confirmation)"
      :append-inner-icon="show.pass2 ? 'mdi-eye' : 'mdi-eye-off'"
      @click:append-inner="show.pass2 = !show.pass2"
    />

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
      :disabled="v$.$invalid || v$.password.$error || v$.password_confirmation.$error"
    />
  </v-form>
</template>

<script setup lang="ts">
import { ApiError, AuthService, UserInput } from '@/api'
import { reactive } from 'vue'

import { email, required, sameAs, minLength, helpers, maxLength } from '@vuelidate/validators'
import useVuelidate from '@vuelidate/core'
import { ValidationArgs } from '@vuelidate/core'
import PasswordStrengthMeter from './PasswordStrengthMeter.vue'
import { computed, Ref, ref } from 'vue'
import { zxcvbn } from '@zxcvbn-ts/core'
import { Validation } from '@vuelidate/core'
import { vuelidateErrors } from '@/api/validation'

const state: UserInput = reactive<UserInput>({
  email: '',
  email_public: false,
  login: '',
  password: '',
  password_confirmation: '',
  identity: {
    first_name: '',
    last_name: '',
    contact: undefined
  }
})

const show = ref({
  pass1: false,
  pass2: false
})

const passwordReference = computed(() => state.password)

const validators = {
  password: helpers.withMessage(
    'Password is too weak',
    () => zxcvbn(state.password, userInputs(state)).score >= 3
  ),
  pwdConfirm: helpers.withMessage('Passwords do not match', sameAs(passwordReference))
}

const rules: ValidationArgs<UserInput> = {
  email: { email, required },
  email_public: {},
  password: { required, strength: validators.password },
  password_confirmation: { sameAsPassword: validators.pwdConfirm, required },
  login: { minLength: minLength(5), maxLength: maxLength(16) },
  identity: {
    first_name: { minLength: minLength(2), required },
    last_name: { minLength: minLength(2), required },
    contact: {}
  }
}

const v$: Ref<Validation<ValidationArgs<UserInput>, UserInput>> = useVuelidate(rules, state)

function userInputs(state: UserInput) {
  return [state.email, state.login, state.identity.first_name, state.identity.last_name]
}

const passwordStrength = computed(() => {
  return zxcvbn(state.password, userInputs(state)).score
})

const userCreated = ref(false)

async function submit() {
  await AuthService.registerUser(state)
    .then(() => {
      userCreated.value = true
    })
    .catch((err: ApiError) => {
      console.log(err)
      // switch (err.status) {
      //   case 400:
      //   // do stuff
      //     break;

      //   // default:
      //   //   break;
    })
}
</script>

<style scoped></style>
