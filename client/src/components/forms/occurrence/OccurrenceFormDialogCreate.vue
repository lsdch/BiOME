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
import { ErrorModel, BaseOccurrenceSamplingOutline } from '@/api'
import {
  createOccurrenceMutation,
  samplingAddOccurrenceMutation,
  siteAddOccurrenceMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { hasID } from '@/lib/db'
import { IndexedValidationErrors } from '@/lib/mutations'
import { BiomatModel, OccurrenceModel, SamplingModel, SiteModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import OccurrenceFormDialog from './OccurrenceFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const model = ref<OccurrenceModel.OccurrenceModel>(OccurrenceModel.initialModel())

defineProps<FormDialogProps>()

const addFromSampling = useMutation(samplingAddOccurrenceMutation())
const addFromSite = useMutation(siteAddOccurrenceMutation())
const createFromScratch = useMutation(createOccurrenceMutation())

function getActiveMutation() {
  if (hasID(model.value?.sampling)) return addFromSampling
  else if (hasID(model.value?.site)) return addFromSite
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
  onSuccess: (data: BaseOccurrenceSamplingOutline) => {
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
  if (hasID(model.value?.sampling))
    return addFromSampling.mutate(
      {
        path: { number: model.value.sampling.number },
        body: BiomatModel.toRequestData(model.value.biomaterial!)
      },
      mutationCallbacks
    )
  else if (hasID(model.value?.site))
    return addFromSite.mutate(
      {
        path: { code: model.value.site.code },
        body: {
          sampling: SamplingModel.toRequestBody(model.value.sampling!),
          biomaterial: BiomatModel.toRequestData(model.value.biomaterial!)
        }
      },
      mutationCallbacks
    )
  else
    return createFromScratch.mutate(
      {
        body: {
          site: SiteModel.toRequestBody(model.value.site!),
          sampling: SamplingModel.toRequestBody(model.value.sampling!),
          bio_material: BiomatModel.toRequestData(model.value.biomaterial!)
        }
      },
      mutationCallbacks
    )
}
</script>

<style scoped lang="scss"></style>
