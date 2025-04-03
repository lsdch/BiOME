<template>
  <TaxonFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} taxon`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </TaxonFormDialog>
</template>

<script setup lang="ts">
import { $TaxonInput, $TaxonUpdate, Taxon } from '@/api'
import { createTaxonMutation, updateTaxonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { TaxonModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import TaxonFormDialog from './TaxonFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Taxon>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createTaxonMutation(), {
  initial: TaxonModel.initialModel,
  schema: $TaxonInput,
  requestData: (model) => ({ body: TaxonModel.toRequestBody(model) })
})

const update = defineFormUpdate(updateTaxonMutation(), {
  itemToModel: TaxonModel.fromTaxon,
  schema: $TaxonUpdate,
  requestData: ({ code }, model) => ({
    path: { code },
    body: TaxonModel.toRequestBody(model)
  })
})

const { feedback } = useFeedback()

const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Taxon created` : `Taxon updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
