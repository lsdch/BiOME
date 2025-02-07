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
        :title="`${mode} fixative`"
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
import { $FixativeInput, $FixativeUpdate, Fixative, FixativeInput, FixativeUpdate } from '@/api'
import { createFixativeMutation, updateFixativeMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'

const dialog = defineModel<boolean>()
const item = defineModel<Fixative>()

const initial: FixativeInput = {
  code: '',
  label: ''
}

function updateTransformer({ id, $schema, meta, ...rest }: Fixative): FixativeUpdate {
  return rest
}

const create = {
  mutation: createFixativeMutation,
  schema: $FixativeInput
}

const update = {
  mutation: updateFixativeMutation,
  schema: $FixativeUpdate,
  itemID: ({ code }: Fixative) => ({ code })
}
</script>

<style scoped lang="scss"></style>
