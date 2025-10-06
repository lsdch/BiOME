<template>
  <SamplingFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :title="title ?? `${mode} sampling`"
    :mode
    :site
    :errors
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData" />
    </template>
  </SamplingFormDialog>
</template>

<script setup lang="ts">
import { $SamplingInput, $SamplingUpdate, Sampling, SiteItem } from '@/api'
import {
  createSamplingAtSiteMutation,
  updateSamplingMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { SamplingModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import SamplingFormDialog from '../forms/SamplingFormDialog.vue'
import { FormDialogProps } from '../toolkit/forms/FormDialog.vue'

const item = defineModel<Sampling>()
const dialog = defineModel<boolean>('dialog')

const { site } = defineProps<
  {
    site: SiteItem
  } & FormDialogProps
>()

const create = defineFormCreate(createSamplingAtSiteMutation(), {
  initial: SamplingModel.initialModel,
  schema: $SamplingInput,
  requestData: (model) => ({
    path: { code: site.code },
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
