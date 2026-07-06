<template>
  <!-- @update:model-value="console.log" -->
  <OccurrenceFormDialog
    v-model="model"
    v-model:dialog="dialog"
    mode="Create"
    :errors
    :title="`Register occurrence`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </OccurrenceFormDialog>
</template>

<script setup lang="ts">
import { BaseOccurrence, ErrorModel } from '@/api'
import {
  createOccurrenceAtSamplingMutation,
  createOccurrenceMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { hasID } from '@/lib/db'
import { IndexedValidationErrors } from '@/lib/mutations'
import { OccurrenceModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import OccurrenceFormDialog from './OccurrenceFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const model = ref<OccurrenceModel.OccurrenceModel>(OccurrenceModel.initialModel())

defineProps<FormDialogProps>()

const addFromSampling = useMutation(createOccurrenceAtSamplingMutation())
const createFromScratch = useMutation(createOccurrenceMutation())

function getActiveMutation() {
  if (hasID(model.value?.sampling)) return addFromSampling
  else return createFromScratch
}

const activeMutation = computed(getActiveMutation)

const errors = computed<IndexedValidationErrors | undefined>(() => {
  return activeMutation.value.error.value?.errors?.reduce<IndexedValidationErrors>(
    (acc, { location, message }) => {
      if (location?.startsWith('body.')) {
        const loc = location.replace('body.', '')
        acc[loc] = (acc[loc] ?? []).concat(message ?? 'Invalid value')
      } else {
        acc.rest.push(message ?? 'Invalid value')
      }
      return acc
    },
    { rest: [] }
  )
})

const { feedback } = useFeedback()

const mutationCallbacks = {
  onSuccess: (data: BaseOccurrence) => {
    feedback({
      type: 'success',
      message: `Occurrence ${data.code} created`
    })
    dialog.value = false
  },
  onError: (error: ErrorModel) => {
    console.error('Error submitting form:', error)
  }
}

function submit() {
  if (typeof model.value?.sampling === 'string')
    return addFromSampling.mutate(
      {
        path: { id: model.value.sampling },
        body: model.value
      },
      mutationCallbacks
    )
  else
    return createFromScratch.mutate(
      {
        body: model.value
      },
      mutationCallbacks
    )
}
</script>

<style scoped lang="scss"></style>
