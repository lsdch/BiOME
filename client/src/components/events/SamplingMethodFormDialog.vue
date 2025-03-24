<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} sampling method`"
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
  $SamplingMethodInput,
  $SamplingMethodUpdate,
  SamplingMethod,
  SamplingMethodInput,
  SamplingMethodUpdate
} from '@/api'
import {
  createSamplingMethodMutation,
  updateSamplingMethodMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<SamplingMethod>()

const initial: SamplingMethodInput = {
  code: '',
  label: ''
}

function updateTransformer({ code, label, description }: SamplingMethod): SamplingMethodUpdate {
  return { code, label, description }
}

const create = defineFormCreate(createSamplingMethodMutation(), {
  initial,
  schema: $SamplingMethodInput
})

const update = defineFormUpdate(updateSamplingMethodMutation(), {
  schema: $SamplingMethodUpdate,
  itemToModel: updateTransformer,
  requestData: ({ code }) => ({ path: { code } })
})
</script>

<style scoped lang="scss"></style>
