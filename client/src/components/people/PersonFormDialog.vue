<template>
  <FormDialog v-model="dialog" title="Create person" :loading @submit="submit">
    <v-form @submit.prevent="submit">
      <v-container fluid>
        <PersonFormFields v-model="model" :bindings="nameBindings" />
        <v-row>
          <v-col>
            <v-text-field
              label="Contact (optional)"
              v-model="model.contact"
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
              v-model="model.comment"
              variant="outlined"
              label="Comments (optional)"
              v-bind="field('comment')"
            />
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </FormDialog>
</template>

<script lang="ts">
const DEFAULT: PersonInput = {
  first_name: '',
  last_name: '',
  institutions: []
}
</script>

<script setup lang="ts">
import { $PersonInput, PeopleService, Person, PersonInput } from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import { ref } from 'vue'
import { VForm } from 'vuetify/components'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import InstitutionPicker from './InstitutionPicker.vue'
import PersonFormFields from './PersonFormFields.vue'
import { useToggle } from '@vueuse/core'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Person>>()
const emit = defineEmits<FormEmits<Person>>()
const { model } = useForm(props, {
  initial: DEFAULT,
  transformers: {
    institutions: (v) => v.map(({ code }) => code) ?? []
  }
})

const { field, errorHandler } = useSchema($PersonInput)

const nameBindings = ref({
  firstName: field('first_name'),
  lastName: field('last_name')
})

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  const data = model.value
  const request = props.edit
    ? PeopleService.updatePerson({ path: { id: props.edit.id }, body: data })
    : PeopleService.createPerson({ body: data })

  return await request
    .then(errorHandler)
    .then((person) => emit('success', person))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped></style>
