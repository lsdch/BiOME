<template>
  <PersonFormDialog
    v-model="model"
    :mode
    :errors
    :loading="loading || activeMutation.isPending.value"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </PersonFormDialog>
</template>

<script setup lang="ts">
import { $PersonInput, $PersonUpdate, Person } from '@/api'
import { createPersonMutation, updatePersonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { PersonModel } from '@/models'
import PersonFormDialog from './PersonFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Person>()

defineProps<FormDialogProps>()

const create = defineFormCreate(createPersonMutation(), {
  initial: PersonModel.initialModel,
  schema: $PersonInput,
  requestData: (model) => ({
    body: PersonModel.toRequestBody(model)
  })
})

const update = defineFormUpdate(updatePersonMutation(), {
  schema: $PersonUpdate,
  itemToModel: PersonModel.fromPerson,
  requestData: ({ id }, model) => ({
    path: { id },
    body: PersonModel.toRequestBody(model)
  })
})

const { mode, model, errors, submit, activeMutation } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
  }
})
</script>

<style scoped></style>
