<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} abiotic parameter`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field label="Label" v-model="model.label" v-bind="field('label')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field label="Code" v-model="model.code" v-bind="field('code')" />
            </v-col>
            <v-col>
              <v-text-field label="Unit" v-model="model.unit" v-bind="field('unit')" />
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
import {
  $AbioticParameterInput,
  $AbioticParameterUpdate,
  AbioticParameter,
  AbioticParameterInput,
  AbioticParameterUpdate
} from '@/api'
import {
  createAbioticParameterMutation,
  updateAbioticParameterMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<AbioticParameter>()

const initial: AbioticParameterInput = {
  label: '',
  code: '',
  unit: ''
}

function updateTransformer({
  id,
  $schema,
  meta,
  ...rest
}: AbioticParameter): AbioticParameterUpdate {
  return rest
}

const create = defineFormCreate(createAbioticParameterMutation(), {
  schema: $AbioticParameterInput,
  initial
})

const update = defineFormUpdate(updateAbioticParameterMutation(), {
  schema: $AbioticParameterUpdate,
  itemToModel: updateTransformer,
  requestData({ code }) {
    return {
      path: { code }
    }
  }
})
</script>

<style scoped lang="scss"></style>
