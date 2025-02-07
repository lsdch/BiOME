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

const create = {
  mutation: createSamplingMethodMutation,
  schema: $SamplingMethodInput
}

const update = {
  mutation: updateSamplingMethodMutation,
  schema: $SamplingMethodUpdate,
  itemID: ({ code }: SamplingMethod) => ({ code })
}
</script>

<style scoped lang="scss"></style>
