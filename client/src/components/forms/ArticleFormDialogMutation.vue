<template>
  <ArticleFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} article`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </ArticleFormDialog>
</template>

<script setup lang="ts">
import { $PublicationInput, $PublicationUpdate, Publication, PublicationUpdate } from '@/api'
import {
  createPublicationMutation,
  updatePublicationMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/lib/mutations'
import { PublicationModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import ArticleFormDialog from './ArticleFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Publication>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createPublicationMutation(), {
  initial: PublicationModel.initialModel,
  schema: $PublicationInput,
  requestData: (model) => ({ body: PublicationModel.toRequestBody(model) })
})

const update = defineFormUpdate(updatePublicationMutation(), {
  itemToModel: PublicationModel.fromPublication,
  schema: $PublicationUpdate,
  requestData: ({ code }, model) => ({
    path: { code },
    body: PublicationModel.toRequestBody(model)
  })
})

const { feedback } = useFeedback()
const emit = defineEmits<{ success: [item: Publication] }>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    emit('success', item)
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Publication created` : `Publication updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
