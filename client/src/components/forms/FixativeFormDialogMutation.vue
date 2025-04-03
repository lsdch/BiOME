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
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { FixativeModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import FixativeFormDialog from './FixativeFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Fixative>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createFixativeMutation(), {
  initial: FixativeModel.initialModel,
  schema: $FixativeInput
})

const update = defineFormUpdate(updateFixativeMutation(), {
  itemToModel: FixativeModel.fromFixative,
  schema: $FixativeUpdate,
  requestData: ({ code }, model) => ({
    path: { code }
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
      message: mode === 'Create' ? `Fixative created` : `Fixative updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
