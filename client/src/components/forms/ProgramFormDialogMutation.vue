<template>
  <ProgramFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :title="title ?? `${mode} program`"
    :mode
    :errors
    @submit="submit()"
  />
</template>

<script setup lang="ts">
import { $ProgramInput, $ProgramUpdate, Program } from '@/api'
import { createProgramMutation, updateProgramMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { ProgramModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import ProgramFormDialog from './ProgramFormDialog.vue'
import { FormDialogProps } from '../toolkit/forms/FormDialog.vue'

const item = defineModel<Program>()
const dialog = defineModel<boolean>('dialog')

defineProps<FormDialogProps>()

const create = defineFormCreate(createProgramMutation(), {
  initial: ProgramModel.initialModel,
  schema: $ProgramInput,
  requestData(model) {
    return { body: ProgramModel.toRequestBody(model) }
  }
})

const update = defineFormUpdate(updateProgramMutation(), {
  schema: $ProgramUpdate,
  itemToModel: ProgramModel.fromProgram,
  requestData({ code }, model) {
    return { path: { code }, body: ProgramModel.toRequestBody(model) }
  }
})

const { feedback } = useFeedback()

const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Program created` : `Program updated`
    })
  }
})
</script>

<style scoped></style>
