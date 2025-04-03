<template>
  <SamplingMethodFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} sampling method`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </SamplingMethodFormDialog>
</template>

<script setup lang="ts">
import { $SamplingMethodInput, $SamplingMethodUpdate, SamplingMethod } from '@/api'
import {
  createSamplingMethodMutation,
  updateSamplingMethodMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { SamplingMethodModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import SamplingMethodFormDialog from './SamplingMethodFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<SamplingMethod>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createSamplingMethodMutation(), {
  initial: SamplingMethodModel.initialModel,
  schema: $SamplingMethodInput
})

const update = defineFormUpdate(updateSamplingMethodMutation(), {
  itemToModel: SamplingMethodModel.fromSamplingMethod,
  schema: $SamplingMethodUpdate,
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
      message: mode === 'Create' ? `Sampling method created` : `Sampling method updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
