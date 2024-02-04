<template>
  <v-form ref="innerForm" @submit.prevent="submit">
    <v-container fluid>
      <v-row class="mb-3">
        <v-text-field
          name="institution"
          label="Institution name"
          id="institution-input"
          hint="The name of the structure for which this instance is deployed, e.g. the name of your lab."
          persistent-hint
          v-model="inst.name"
          required
          :error-messages="errors?.name?.map(({ message }) => message)"
        />
      </v-row>
      <v-row class="mb-3">
        <v-text-field
          name="institution_shortname"
          label="Acronym or abbreviated name"
          id="institution-shortname"
          hint="A short label or acronym that identifies your lab."
          persistent-hint
          v-model="inst.acronym"
          :error-messages="errors?.acronym?.map(({ message }) => message)"
        />
      </v-row>
      <v-row class="mb-3">
        <v-textarea
          variant="outlined"
          label="Description (optional)"
          v-model="inst.description"
          :error-messages="errors?.description?.map(({ message }) => message)"
        />
      </v-row>
      <v-row>
        <v-spacer />
        <v-btn :loading="loading" color="primary" variant="plain" type="submit" text="Submit" />
      </v-row>
    </v-container>
  </v-form>
</template>

<script setup lang="ts">
import { Institution, InstitutionInput, PeopleService } from '@/api'
import { Ref, ref } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { VForm } from 'vuetify/components'
import { Emits, Props, useForm } from '../toolkit/form'

const innerForm = ref<ComponentExposed<typeof VForm> | null>(null)

const props = defineProps<Props<Institution>>()

const inst: Ref<InstitutionInput> = ref(
  props.edit
    ? { ...props.edit }
    : {
        name: '',
        acronym: '',
        description: ''
      }
)

function request() {
  if (props.edit) {
    return PeopleService.updateInstitution({ ...props.edit, ...inst.value })
  } else {
    return PeopleService.createInstitution(inst.value)
  }
}

const emit = defineEmits<Emits<Institution>>()

const { errors, submit, loading } = useForm<InstitutionInput, Institution>(props, emit, request)

function reset() {
  innerForm.value?.reset()
}

defineExpose({
  reset,
  submit
})
</script>

<style scoped></style>
