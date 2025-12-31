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
import {
  $SamplingInput,
  $SamplingUpdate,
  CreateSamplingAtSiteData,
  Sampling,
  SiteItem
} from '@/api'
import {
  createSamplingAtSiteMutation,
  updateSamplingMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { defineFormCreate, defineFormUpdate, RequestData, useMutationForm } from '@/lib/mutations'
import { SamplingModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import SamplingFormDialog from '@/components/forms/SamplingFormDialog.vue'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { Options } from '@/api/gen/client'

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
  requestData: (model): Options<CreateSamplingAtSiteData> => ({
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
const emit = defineEmits<{ success: [item: Sampling] }>()
const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    emit('success', item)
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Sampling created` : `Sampling updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
