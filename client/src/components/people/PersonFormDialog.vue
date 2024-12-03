<template>
  <FormDialog v-model="dialog" :title="`${mode} person`" :loading @submit="submit">
    <v-container fluid>
      <v-row>
        <v-col cols="12" sm="6">
          <v-text-field
            name="first_name"
            label="First name(s)"
            v-model.trim="model.first_name"
            v-bind="field('first_name')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <v-text-field
            name="last_name"
            label="Last name"
            v-model.trim="model.last_name"
            v-bind="field('last_name')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field
            label="Contact (optional)"
            v-model.trim="model.contact"
            prepend-inner-icon="mdi-at"
            v-bind="field('contact')"
            hint="An e-mail address to contact this person"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <InstitutionPicker
            label="Institutions (optional)"
            v-model="model.institutions"
            item-color="primary"
            chips
            closable-chips
            multiple
            item-value="code"
            v-bind="field('institutions')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea
            v-model.trim="model.comment"
            variant="outlined"
            label="Comments (optional)"
            v-bind="field('comment')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import {
  $PersonInput,
  $PersonUpdate,
  PeopleService,
  Person,
  PersonInput,
  PersonUpdate
} from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import { reactiveComputed, useToggle } from '@vueuse/core'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import InstitutionPicker from './InstitutionPicker.vue'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Person>>()
const emit = defineEmits<FormEmits<Person>>()
const initial: PersonInput = {
  first_name: '',
  last_name: '',
  institutions: []
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({
    alias,
    comment,
    contact,
    first_name,
    last_name,
    institutions
  }): PersonUpdate {
    return {
      first_name,
      last_name,
      alias,
      comment,
      contact,
      institutions: institutions.map(({ code }) => code)
    }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $PersonInput : $PersonUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: PeopleService.createPerson,
    edit: ({ id }, model) => PeopleService.updatePerson({ path: { id }, body: model })
  })
    .then(errorHandler)
    .then((person) => emit('success', person))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped></style>
