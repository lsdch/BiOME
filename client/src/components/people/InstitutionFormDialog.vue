<template>
  <CreateUpdateForm
    v-model="item"
    :initial
    :update-transformer
    :create
    :update
    @success="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} institution`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field
                id="institution-input"
                v-model="model.name"
                name="institution"
                label="Institution name"
                persistent-hint
                required
                v-bind="field('name')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field
                id="institution-shortname"
                v-model="model.code"
                name="institution_shortname"
                label="Code or abbreviated name"
                persistent-hint
                v-bind="field('code')"
              />
            </v-col>
            <v-col cols="12" sm="6">
              <InstitutionKindPicker
                v-model="model.kind"
                v-bind="field('kind')"
                label="Kind"
                variant="outlined"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                v-model="model.description"
                variant="outlined"
                label="Description (optional)"
                v-bind="field('description')"
              />
            </v-col>
          </v-row>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import {
  $InstitutionInput,
  $InstitutionUpdate,
  Institution,
  InstitutionInput,
  InstitutionUpdate
} from '@/api'
import {
  createInstitutionMutation,
  updateInstitutionMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import InstitutionKindPicker from './InstitutionKindPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Institution>()

const initial: InstitutionInput = { code: '', name: '', kind: 'Lab' }

function updateTransformer({ code, name, kind, description }: Institution): InstitutionUpdate {
  return { code, name, kind, description }
}

const create = {
  mutation: createInstitutionMutation,
  schema: $InstitutionInput
}

const update = {
  mutation: updateInstitutionMutation,
  schema: $InstitutionUpdate,
  itemID: ({ code }: Institution) => ({ code })
}
</script>

<style scoped></style>
