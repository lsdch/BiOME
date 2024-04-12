<template>
  <v-row>
    <v-col cols="12" sm="6">
      <v-text-field
        name="first_name"
        label="First name(s)"
        v-model.trim="person.first_name"
        :error-messages="errorMsgs?.first_name"
        :rules="inlineRules([required])"
      />
    </v-col>
    <v-col cols="12" sm="6">
      <v-text-field
        name="last_name"
        label="Last name"
        v-model.trim="person.last_name"
        :error-messages="errorMsgs?.last_name"
        :rules="inlineRules([required])"
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { PersonInput } from '@/api'
import { required } from '@vuelidate/validators'
import { ErrorMsgs, inlineRules } from '../toolkit/form'

type PartialPersonInput = Pick<PersonInput, 'first_name' | 'last_name'>

const person = defineModel<PartialPersonInput>({ required: true })

defineProps<{
  errorMsgs?: ErrorMsgs<PartialPersonInput>
}>()
</script>

<style scoped></style>
