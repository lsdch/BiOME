<template>
  <FormDialog v-model="dialog" title="Create person" :loading="loading" @submit="submit">
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
            <v-autocomplete
              label="Institutions (optional)"
              v-model="model.institutions"
              :items="institutions"
              item-color="primary"
              chips
              closable-chips
              multiple
              :item-props="({ code, name }) => ({ title: code, subtitle: name })"
              item-value="code"
              v-bind="field('institutions')"
              prepend-inner-icon="mdi-domain"
            >
              <template v-slot:chip="{ item, props }">
                <InstitutionKindChip :kind="item.raw.kind" v-bind="props" size="x-small">
                  {{ item.raw.code }}
                </InstitutionKindChip>
              </template>
            </v-autocomplete>
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
import { $PersonInput, Institution, PeopleService, Person, PersonInput } from '@/api'
import { ref } from 'vue'
import { VForm } from 'vuetify/components'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { FormEmits, FormProps, useForm } from '../toolkit/forms/form'
import InstitutionKindChip from './InstitutionKindChip.vue'
import PersonFormFields from './PersonFormFields.vue'

const props = defineProps<FormProps<Person>>()
const dialog = defineModel<boolean>()
const emit = defineEmits<FormEmits<Person>>()
const { loading, field, errorHandler, model } = useForm(props, $PersonInput, {
  initial: DEFAULT,
  transformers: {
    institutions: (v) => v.map(({ code }) => code) ?? []
  }
})
const nameBindings = ref({
  firstName: field('first_name'),
  lastName: field('last_name')
})

/**
 * List of known institutions from the DB
 */
const institutions = ref<Institution[]>(await PeopleService.listInstitutions())

async function submit() {
  const data = model.value
  const request = props.edit
    ? PeopleService.updatePerson({ id: props.edit.id, requestBody: data })
    : PeopleService.createPerson({ requestBody: data })

  return await request.then((person) => emit('success', person)).catch(errorHandler)
}
</script>

<style scoped></style>
