<template>
  <div>
    <v-text-field
      v-model="state.password"
      v-bind="$attrs"
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
      v-bind="$attrs"
      name="password-confirm"
      label="Password confirmation"
      :type="show.pass2 ? 'text' : 'password'"
      required
      @blur="v$.password_confirmation.$touch"
      :error-messages="vuelidateErrors(v$.password_confirmation)"
      :append-inner-icon="show.pass2 ? 'mdi-eye' : 'mdi-eye-off'"
      @click:append-inner="show.pass2 = !show.pass2"
    />
  </div>
</template>

<script setup lang="ts">
import { PasswordInput } from '@/api'
import { vuelidateErrors } from '@/api/validation'
import useVuelidate, { Validation, ValidationArgs } from '@vuelidate/core'
import { helpers, required, sameAs } from '@vuelidate/validators'
import { zxcvbn } from '@zxcvbn-ts/core'
import { Ref, computed, defineModel, ref } from 'vue'
import PasswordStrengthMeter from './PasswordStrengthMeter.vue'

const props = withDefaults(defineProps<{}>(), { userInputs: () => [] })

const state = defineModel<PasswordInput>({
  default: {
    password: '',
    password_confirmation: ''
  }
})

const passwordReference = computed(() => state.value.password)

const validators = {
  password: helpers.withMessage(
    'Password is too weak',
    () => zxcvbn(state.value.password, props.userInputs).score >= 3
  ),
  pwdConfirm: helpers.withMessage('Passwords do not match', sameAs(passwordReference))
}

const rules: ValidationArgs<PasswordInput> = {
  password: { required, strength: validators.password },
  password_confirmation: { sameAsPassword: validators.pwdConfirm, required }
}

const show = ref({
  pass1: false,
  pass2: false
})

const v$: Ref<Validation<ValidationArgs<PasswordInput>, PasswordInput>> = useVuelidate(rules, state)

const passwordStrength = computed(() => {
  return zxcvbn(state.value.password, props.userInputs).score
})
</script>

<style scoped></style>
