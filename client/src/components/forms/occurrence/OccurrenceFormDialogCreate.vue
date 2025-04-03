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
import { BioMaterialWithDetails, ErrorModel } from '@/api'
import {
  createExternalBioMatMutation,
  eventAddExternalOccurrenceMutation,
  samplingAddExternalOccurrenceMutation,
  siteAddExternalOccurrenceMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { hasID } from '@/functions/db'
import { IndexedValidationErrors } from '@/functions/mutations'
import { BiomatModel, EventModel, OccurrenceModel, SamplingModel, SiteModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
import { reactiveComputed } from '@vueuse/core'
import { computed, ref, unref } from 'vue'
import OccurrenceFormDialog from './OccurrenceFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const model = ref<OccurrenceModel.OccurrenceModel>(OccurrenceModel.initialModel())

defineProps<FormDialogProps>()

const addFromSampling = useMutation(samplingAddExternalOccurrenceMutation())
const addFromEvent = useMutation(eventAddExternalOccurrenceMutation())
const addFromSite = useMutation(siteAddExternalOccurrenceMutation())
const createFromScratchExternal = useMutation(createExternalBioMatMutation())

function getActiveMutation() {
  if (hasID(model.value?.sampling)) return addFromSampling
  else if (hasID(model.value?.event)) return addFromEvent
  else if (hasID(model.value?.site)) return addFromSite
  else return createFromScratchExternal
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
  onSuccess: (data: BioMaterialWithDetails) => {
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
  if (model.value?.biomaterial.category !== 'External') {
    console.error('Internal pipeline is not implemented yet')
    console.log(model.value)
    return
  }
  if (hasID(model.value?.sampling))
    return addFromSampling.mutate(
      {
        path: { id: model.value.sampling.id },
        body: BiomatModel.toRequestData(model.value.biomaterial.external!)
      },
      mutationCallbacks
    )
  else if (hasID(model.value?.event))
    return addFromEvent.mutate(
      {
        path: { id: model.value.event.id },
        body: {
          sampling: SamplingModel.toRequestBody(model.value.sampling!),
          biomaterial: BiomatModel.toRequestData(model.value.biomaterial.external!)
        }
      },
      mutationCallbacks
    )
  else if (hasID(model.value?.site))
    return addFromSite.mutate(
      {
        path: { code: model.value.site.code },
        body: {
          event: EventModel.toRequestData(model.value.event!),
          sampling: SamplingModel.toRequestBody(model.value.sampling!),
          biomaterial: BiomatModel.toRequestData(model.value.biomaterial.external!)
        }
      },
      mutationCallbacks
    )
  else
    return createFromScratchExternal.mutate(
      {
        body: {
          site: SiteModel.toRequestBody(model.value.site!),
          event: EventModel.toRequestData(model.value.event!),
          sampling: SamplingModel.toRequestBody(model.value.sampling!),
          bio_material: BiomatModel.toRequestData(model.value.biomaterial.external!)
        }
      },
      mutationCallbacks
    )
}
</script>

<style scoped lang="scss"></style>
