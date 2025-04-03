<template>
  <DataSourceFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} data source`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </DataSourceFormDialog>
</template>

<script setup lang="ts">
import { $DataSourceInput, $DataSourceUpdate, DataSource } from '@/api'
import {
  createDataSourceMutation,
  updateDataSourceMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { DataSourceModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import DataSourceFormDialog from './DataSourceFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<DataSource>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createDataSourceMutation(), {
  initial: DataSourceModel.initialModel,
  schema: $DataSourceInput
})

const update = defineFormUpdate(updateDataSourceMutation(), {
  itemToModel: DataSourceModel.fromDataSource,
  schema: $DataSourceUpdate,
  requestData: ({ code }, model) => ({ path: { code } })
})

const { feedback } = useFeedback()

const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `DataSource created` : `DataSource updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
