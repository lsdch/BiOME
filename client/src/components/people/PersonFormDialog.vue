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
              :item-props="({ code, name }: Institution) => ({ title: code, subtitle: name })"
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
import { handleErrors } from '@/api/responses'
import { FormEmits, FormProps, useForm } from '@/components/toolkit/forms/form'
import { ref } from 'vue'
import { VForm } from 'vuetify/components'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import InstitutionKindChip from './InstitutionKindChip.vue'
import PersonFormFields from './PersonFormFields.vue'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Person>>()
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
const institutions = ref<Institution[]>(
  await PeopleService.listInstitutions().then(
    handleErrors((err) => console.error('Failed to fetch institutions: ', err))
  )
)

async function submit() {
  const data = model.value
  const request = props.edit
    ? PeopleService.updatePerson({ path: { id: props.edit.id }, body: data })
    : PeopleService.createPerson({ body: data })

  return await request.then(errorHandler).then((person) => emit('success', person))
}
</script>

<style scoped></style>
