<template>
  <FixativeFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} fixative`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </FixativeFormDialog>
</template>

<script setup lang="ts">
import { $FixativeInput, $FixativeUpdate, Fixative } from '@/api'
import { createFixativeMutation, updateFixativeMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/lib/mutations'
import { FixativeModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import FixativeFormDialog from './FixativeFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Fixative>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createFixativeMutation(), {
  initial: FixativeModel.initialModel,
  schema: $FixativeInput,
  requestData: (model) => ({ body: model })
})

const update = defineFormUpdate(updateFixativeMutation(), {
  itemToModel: FixativeModel.fromFixative,
  schema: $FixativeUpdate,
  requestData: ({ code }, model) => ({
    path: { code },
    body: model
  })
})

const { feedback } = useFeedback()
const emit = defineEmits<{ success: [item: Fixative] }>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    emit('success', item)
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Fixative created` : `Fixative updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
