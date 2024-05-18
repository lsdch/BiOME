<template>
  <v-row>
    <v-col cols="12" sm="6">
      <v-text-field
        name="first_name"
        label="First name(s)"
        v-model.trim="person.first_name"
        v-bind="bindings?.firstName"
      />
    </v-col>
    <v-col cols="12" sm="6">
      <v-text-field
        name="last_name"
        label="Last name"
        v-model.trim="person.last_name"
        v-bind="bindings?.lastName"
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { PersonInput } from '@/api'
import { type FieldBinding } from '../toolkit/forms/form'

type PartialPersonInput = Pick<PersonInput, 'first_name' | 'last_name'>

const person = defineModel<PartialPersonInput>({ required: true })

defineProps<{
  bindings?: {
    firstName: FieldBinding
    lastName: FieldBinding
  }
}>()
</script>

<style scoped></style>
