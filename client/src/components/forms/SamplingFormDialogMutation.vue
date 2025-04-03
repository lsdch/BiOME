<template>
  <SamplingFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :title="title ?? `${mode} sampling`"
    :mode
    :event
    :errors
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData" />
    </template>
  </SamplingFormDialog>
</template>

<script setup lang="ts">
import { $SamplingInput, $SamplingUpdate, EventInner, Sampling } from '@/api'
import {
  createSamplingAtEventMutation,
  updateSamplingMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { SamplingModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import SamplingFormDialog from '../forms/SamplingFormDialog.vue'
import { FormDialogProps } from '../toolkit/forms/FormDialog.vue'

const item = defineModel<Sampling>()
const dialog = defineModel<boolean>('dialog')

const { event } = defineProps<
  {
    event: EventInner
  } & FormDialogProps
>()

const create = defineFormCreate(createSamplingAtEventMutation(), {
  initial: SamplingModel.initialModel,
  schema: $SamplingInput,
  requestData: (model) => ({
    path: { id: event.id },
    body: SamplingModel.toRequestBody(model)
  })
})

const update = defineFormUpdate(updateSamplingMutation(), {
  schema: $SamplingUpdate,
  itemToModel: SamplingModel.fromSampling,
  requestData: ({ id }, model) => ({
    path: { id },
    body: SamplingModel.toRequestBody(model)
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
      message: mode === 'Create' ? `Sampling created` : `Sampling updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
