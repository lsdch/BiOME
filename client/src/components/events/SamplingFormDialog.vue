<template>
  <FormDialog
    v-model="dialog"
    :title="mode == 'Create' ? 'New sampling' : 'Edit sampling'"
    :loading
    @submit="submit"
  >
    <v-container>
      <v-row>
        <v-col>
          <v-card>
            <template #title>
              {{ event.site.name }}
            </template>
            <template #subtitle>
              {{ DateWithPrecision.format(event.performed_on) }}
            </template>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="3">
          <v-select
            label="Sampling target"
            v-bind="field('target_kind')"
            :items="$SamplingTargetKind.enum"
            v-model="model.target_kind"
            hide-details
          />
        </v-col>
        <v-col>
          <TaxonPicker
            v-model="model.target_taxa"
            v-bind="field('target_taxa')"
            :disabled="model.target_kind !== 'Taxa'"
            label="Target taxa"
            item-value="name"
            :ranks="ranksUpTo('Family')"
            multiple
            chips
            closable-chips
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" sm="6">
          <FixativePicker
            label="Fixatives"
            v-model="model.fixatives"
            v-bind="field('fixatives')"
            multiple
            item-value="code"
            chips
            closable-chips
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6">
          <SamplingMethodPicker
            label="Sampling methods"
            v-model="model.methods"
            v-bind="field('methods')"
            multiple
            item-value="code"
            chips
            closable-chips
            clearable
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" md="6">
          <HoursMinutesInput
            label="Duration"
            class="mt-2"
            v-model="model.duration"
            v-bind="field('duration')"
            clearable
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6">
          <HabitatPicker
            v-model="model.habitats"
            label="Habitat tags"
            item-value="label"
            v-bind="field('habitats')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <AccessPointsPicker
            v-model="model.access_points"
            v-bind="field('access_points')"
            label="Access points"
            hint="Pick existing terms already in use, or enter new terms"
            persistent-hint
            clearable
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea v-model="model.comments" v-bind="field('comments')" label="Comments" />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import {
  $SamplingInput,
  $SamplingTargetKind,
  $SamplingUpdate,
  Event,
  Sampling,
  SamplingInput,
  SamplingService,
  SamplingUpdate
} from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import { reactiveComputed, useToggle } from '@vueuse/core'
import HabitatPicker from '../habitat/HabitatPicker.vue'
import FixativePicker from '../samples/FixativePicker.vue'
import { ranksUpTo } from '../taxonomy/rank'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import { FormProps, useForm, useSchema } from '../toolkit/forms/form'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import AccessPointsPicker from './AccessPointsPicker.vue'
import HoursMinutesInput from './HoursMinutesInput.vue'
import SamplingMethodPicker from './SamplingMethodPicker.vue'

const dialog = defineModel<boolean>()
const [loading, toggleLoading] = useToggle(false)

const emit = defineEmits<{
  created: [sampling: Sampling]
  updated: [sampling: Sampling]
}>()

const props = defineProps<
  FormProps<Sampling> & {
    event: Event
  }
>()

const initial: SamplingInput = {
  event_id: props.event.id,
  target_kind: 'Taxa',
  target_taxa: []
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({
    duration,
    access_points,
    fixatives,
    habitats,
    target,
    comments,
    methods
  }): SamplingUpdate {
    return {
      duration,
      access_points,
      comments,
      target_kind: target.kind,
      target_taxa: target.target_taxa?.map(({ name }) => name),
      habitats: habitats.map(({ label }) => label),
      fixatives: fixatives?.map(({ code }) => code),
      methods: methods?.map(({ code }) => code)
    }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $SamplingInput : $SamplingUpdate)
)

async function submit() {
  toggleLoading(true)
  return makeRequest({
    create: SamplingService.createSampling,
    edit: ({ id }, model) => SamplingService.updateSampling({ path: { id }, body: model })
  })
    .then(errorHandler)
    .then((sampling) => {
      if (mode.value == 'Create') emit('created', sampling)
      else emit('updated', sampling)
    })
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
