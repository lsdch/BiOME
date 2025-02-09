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
        :title="`${mode} organisation`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field
                id="organisation-input"
                v-model="model.name"
                name="organisation"
                label="Organisation name"
                persistent-hint
                required
                v-bind="field('name')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="6">
              <v-text-field
                id="organisation-shortname"
                v-model="model.code"
                name="organisation_shortname"
                label="Code or abbreviated name"
                persistent-hint
                v-bind="field('code')"
              />
            </v-col>
            <v-col cols="12" sm="6">
              <OrgKindPicker
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
  $OrganisationInput,
  $OrganisationUpdate,
  Organisation,
  OrganisationInput,
  OrganisationUpdate
} from '@/api'
import {
  createOrganisationMutation,
  updateOrganisationMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import OrgKindPicker from './OrgKindPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Organisation>()

const initial: OrganisationInput = { code: '', name: '', kind: 'Lab' }

function updateTransformer({ code, name, kind, description }: Organisation): OrganisationUpdate {
  return { code, name, kind, description }
}

const create = {
  mutation: createOrganisationMutation,
  schema: $OrganisationInput
}

const update = {
  mutation: updateOrganisationMutation,
  schema: $OrganisationUpdate,
  itemID: ({ code }: Organisation) => ({ code })
}
</script>

<style scoped></style>
