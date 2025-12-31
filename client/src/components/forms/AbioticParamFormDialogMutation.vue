<template>
  <AbioticParamFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} name`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </AbioticParamFormDialog>
</template>

<script setup lang="ts">
import { $AbioticParameterInput, $AbioticParameterUpdate, AbioticParameter } from '@/api'
import {
  createAbioticParameterMutation,
  updateAbioticParameterMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/lib/mutations'
import { AbioticParamModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import AbioticParamFormDialog from './AbioticParamFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<AbioticParameter>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createAbioticParameterMutation(), {
  initial: AbioticParamModel.initialModel,
  schema: $AbioticParameterInput
})

const update = defineFormUpdate(updateAbioticParameterMutation(), {
  itemToModel: AbioticParamModel.fromAbioticParam,
  schema: $AbioticParameterUpdate,
  requestData: ({ code }, model) => ({
    path: { code }
  })
})

const { feedback } = useFeedback()
const emit = defineEmits<{ success: [item: AbioticParameter] }>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    emit('success', item)
    feedback({
      type: 'success',
      message: mode === 'Create' ? `AbioticParameter created` : `AbioticParameter updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
