<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} gene`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field label="Label" v-model.trim="model.label" v-bind="field('label')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="Code"
                v-model.trim="model.code"
                v-bind="field('code')"
                class="input-font-monospace"
              />
            </v-col>
            <v-col>
              <v-switch
                label="MOTU delimiter"
                v-model="model.is_MOTU_delimiter"
                color="primary"
                v-bind="field('is_MOTU_delimiter')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                label="Description"
                v-model.trim="model.description"
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
import { $GeneInput, $GeneUpdate, Gene, GeneInput, GeneUpdate } from '@/api'
import { createGeneMutation, updateGeneMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Gene>()

const initial: GeneInput = {
  code: '',
  label: '',
  is_MOTU_delimiter: false
}

function updateTransformer({ code, label, description, is_MOTU_delimiter }: Gene): GeneUpdate {
  return { code, label, description, is_MOTU_delimiter }
}

const create = defineFormCreate(createGeneMutation(), {
  initial,
  schema: $GeneInput
})

const update = defineFormUpdate(updateGeneMutation(), {
  schema: $GeneUpdate,
  itemToModel: updateTransformer,
  requestData: ({ code }) => ({ path: { code } })
})
</script>

<style scoped lang="scss"></style>
