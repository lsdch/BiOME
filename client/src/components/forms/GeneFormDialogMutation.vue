<template>
  <GeneFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} gene`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </GeneFormDialog>
</template>

<script setup lang="ts">
import { $GeneInput, $GeneUpdate, Gene } from '@/api'
import { createGeneMutation, updateGeneMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { GeneModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import GeneFormDialog from './GeneFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Gene>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createGeneMutation(), {
  initial: GeneModel.initialModel,
  schema: $GeneInput
})

const update = defineFormUpdate(updateGeneMutation(), {
  itemToModel: GeneModel.fromGene,
  schema: $GeneUpdate,
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
      message: mode === 'Create' ? `Gene created` : `Gene updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
