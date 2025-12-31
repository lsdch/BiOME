<template>
  <CollectionFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} Collection`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
    v-bind="$attrs"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </CollectionFormDialog>
</template>

<script setup lang="ts">
import { $CollectionInput, $CollectionUpdate, Collection } from '@/api'
import {
  createCollectionMutation,
  updateCollectionMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/lib/mutations'
import { CollectionModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import CollectionFormDialog from './CollectionFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Collection>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createCollectionMutation(), {
  initial: CollectionModel.initialModel,
  schema: $CollectionInput
  // requestData: (model): RequestData<CollectionInput> => ({ body: model })
})

const update = defineFormUpdate(updateCollectionMutation(), {
  itemToModel: CollectionModel.fromCollection,
  schema: $CollectionUpdate,
  requestData: ({ code }, model) => ({
    path: { code }
    // body: model
  })
})

const { feedback } = useFeedback()
const emit = defineEmits<{ success: [item: Collection] }>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    emit('success', item)
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Collection created` : `Collection updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
