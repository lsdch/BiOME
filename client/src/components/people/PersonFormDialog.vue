<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} person`"
        :loading="loading.value"
        @submit="submit"
      >
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
              <OrganisationPicker
                label="Organisations (optional)"
                v-model="model.organisations"
                item-color="primary"
                chips
                closable-chips
                multiple
                item-value="code"
                v-bind="field('organisations')"
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
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $PersonInput, $PersonUpdate, Person, PersonInput, PersonUpdate } from '@/api'
import { createPersonMutation, updatePersonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import OrganisationPicker from './OrganisationPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Person>()

const initial: PersonInput = {
  first_name: '',
  last_name: '',
  organisations: []
}

function updateTransformer({
  alias,
  comment,
  contact,
  first_name,
  last_name,
  organisations
}: Person): PersonUpdate {
  return {
    first_name,
    last_name,
    alias,
    comment,
    contact,
    organisations: organisations.map(({ code }) => code)
  }
}
const create = defineFormCreate(createPersonMutation(), {
  initial,
  schema: $PersonInput
})

const update = defineFormUpdate(updatePersonMutation(), {
  schema: $PersonUpdate,
  itemToModel: updateTransformer,
  requestData: ({ id }) => ({ path: { id } })
})
</script>

<style scoped></style>
