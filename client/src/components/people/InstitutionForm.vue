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
          :error-messages="errorMsgs.name"
        />
      </v-row>
      <v-row class="mb-3">
        <v-text-field
          name="institution_shortname"
          label="Code or abbreviated name"
          id="institution-shortname"
          hint="A short label or code that identifies your lab."
          persistent-hint
          v-model="inst.code"
          :error-messages="errorMsgs.code"
        />
      </v-row>
      <v-row class="mb-3">
        <v-textarea
          variant="outlined"
          label="Description (optional)"
          v-model="inst.description"
          :error-messages="errorMsgs.description"
        />
      </v-row>
      <v-row>
        <v-spacer />
        <v-btn :loading="loading" color="primary" variant="plain" type="submit" text="Submit" />
      </v-row>
    </v-container>
  </v-form>
</template>

<script lang="ts">
const DEFAULT = { code: '', name: '' }
</script>

<script setup lang="ts">
import { Institution, InstitutionInput, PeopleService } from '@/api'
import { ref, watchEffect } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { VForm } from 'vuetify/components'
import { Emits, Props, useForm } from '../toolkit/form'

const innerForm = ref<ComponentExposed<typeof VForm> | null>(null)

const props = defineProps<Props<Institution>>()
const inst = defineModel<InstitutionInput>({ default: DEFAULT })
const emit = defineEmits<Emits<Institution>>()

watchEffect(() => {
  if (props.edit) inst.value = { ...props.edit }
  else inst.value = { ...DEFAULT }
})

function request() {
  if (props.edit) {
    return PeopleService.updateInstitution(props.edit.code, inst.value)
  } else {
    return PeopleService.createInstitution(inst.value)
  }
}

const { errorMsgs, submit, loading } = useForm<InstitutionInput, Institution>(props, emit, request)

function reset() {
  innerForm.value?.reset()
}

defineExpose({
  reset,
  submit
})
</script>

<style scoped></style>
