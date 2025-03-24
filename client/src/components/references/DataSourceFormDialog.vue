<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} data source`"
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
              <FTextField
                label="Code"
                v-model.trim="model.code"
                v-bind="field('code')"
                class="input-font-monospace"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="URL"
                v-model.trim="model.url"
                v-bind="field('url')"
                class="input-font-monospace"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-text-field
                label="Link template"
                v-model.trim="model.link_template"
                v-bind="field('link_template')"
                class="input-font-monospace"
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
import {
  $DataSourceInput,
  $DataSourceUpdate,
  DataSource,
  DataSourceInput,
  DataSourceUpdate
} from '@/api'
import {
  createDataSourceMutation,
  updateDataSourceMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FTextField from '../toolkit/forms/FTextField'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<DataSource>()

const initial: DataSourceInput = {
  code: '',
  label: ''
}

function updateTransformer({
  code,
  label,
  description,
  link_template
}: DataSource): DataSourceUpdate {
  return { code, label, description, link_template }
}

const create = defineFormCreate(createDataSourceMutation(), {
  initial,
  schema: $DataSourceInput
})

const update = defineFormUpdate(updateDataSourceMutation(), {
  schema: $DataSourceUpdate,
  itemToModel: updateTransformer,
  requestData: ({ code }) => ({ path: { code } })
})
</script>

<style scoped lang="scss"></style>
