<template>
  <v-row>
    <v-col>
      <PasswordField
        class="password-strength"
        :model-value="model.password"
        @update:model-value="(password) => (model = { ...model, password })"
        v-bind="$attrs"
        name="password"
        label="New password"
        :rules="[() => passwordStrength >= MIN_PWD_STRENGTH || 'Password is too weak']"
      >
        <template v-slot:loader>
          <PasswordStrengthMeter v-if="model.password.length > 0" :strength="passwordStrength" />
        </template>
      </PasswordField>
    </v-col>
  </v-row>
  <v-row>
    <v-col>
      <PasswordField
        :model-value="model.password_confirmation"
        @update:model-value="
          (password_confirmation) => (model = { ...model, password_confirmation })
        "
        v-bind="$attrs"
        name="password-confirm"
        label="Password confirmation"
        :rules="[(v: string) => v === model.password || 'Inputs do not match']"
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { PasswordInput, PersonIdentity } from '@/api'
import { zxcvbn } from '@zxcvbn-ts/core'
import { computed, defineModel } from 'vue'
import PasswordField from '../toolkit/ui/PasswordField.vue'
import PasswordStrengthMeter from './PasswordStrengthMeter.vue'

const MIN_PWD_STRENGTH = 3

const props = defineProps<{
  email?: string
  login?: string
  identity?: PersonIdentity
}>()

const model = defineModel<PasswordInput>({
  default: {
    password: '',
    password_confirmation: ''
  }
})

const passwordSensitiveInputs = computed(() => {
  return [props.email, props.login, props.identity?.first_name, props.identity?.last_name].filter(
    (v): v is string => v !== undefined
  )
})

const passwordStrength = computed(() => {
  return zxcvbn(model.value.password, passwordSensitiveInputs.value).score
})
</script>

<style lang="less">
.password-strength .v-field__loader {
  top: calc(100% - 12px);
}
</style>
