<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        title="Register program"
        v-model="dialog"
        v-bind="$attrs"
        @submit="submit"
        :loading="loading.value"
      >
        <v-container>
          <v-row>
            <v-col>
              <v-text-field label="Label" v-model="model.label" v-bind="field('label')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Code" v-model="model.code" v-bind="field('code')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="d-flex">
              <v-number-input
                label="Start year"
                v-model="model.start_year"
                v-bind="field('start_year')"
                rounded="e-0"
              />
              <v-number-input
                label="End year"
                v-model="model.end_year"
                v-bind="field('end_year')"
                rounded="s-0"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <OrganisationPicker
                label="Funding agencies"
                v-model="model.funding_agencies"
                chips
                closable-chips
                item-value="code"
                multiple
                clearable
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <PersonPicker
                label="Managers"
                v-model="model.managers"
                v-bind="field('managers')"
                item-value="alias"
                multiple
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                label="Description"
                v-model="model.description"
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
import { $ProgramInput, $ProgramUpdate, Program, ProgramInput, ProgramUpdate } from '@/api'
import { createProgramMutation, updateProgramMutation } from '@/api/gen/@tanstack/vue-query.gen'
import OrganisationPicker from '../people/OrganisationPicker.vue'
import PersonPicker from '../people/PersonPicker.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Program>()

const initial: ProgramInput = {
  label: '',
  code: '',
  managers: [],
  funding_agencies: []
}

function updateTransformer({
  id,
  funding_agencies,
  managers,
  meta,
  $schema,
  ...rest
}: Program): ProgramUpdate {
  return {
    ...rest,
    funding_agencies: funding_agencies.map(({ code }) => code),
    managers: managers.map(({ alias }) => alias)
  }
}

const create = defineFormCreate(createProgramMutation(), {
  initial,
  schema: $ProgramInput
})

const update = defineFormUpdate(updateProgramMutation(), {
  schema: $ProgramUpdate,
  itemToModel: updateTransformer,
  requestData({ code }) {
    return { path: { code } }
  }
})
</script>

<style scoped></style>
