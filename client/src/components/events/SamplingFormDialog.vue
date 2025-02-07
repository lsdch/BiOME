<template>
  <CreateUpdateForm
    v-model="item"
    :initial
    :update-transformer
    :create
    :update
    @success="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-model="dialog"
        :title="mode == 'Create' ? 'New sampling' : 'Edit sampling'"
        :loading="loading.value"
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
                v-bind="field('target', 'kind')"
                :items="$SamplingTargetKind.enum"
                v-model="model.target.kind"
                hide-details
              />
            </v-col>
            <v-col>
              <TaxonPicker
                v-model="model.target.taxa"
                v-bind="field('target', 'taxa')"
                :disabled="model.target.kind !== 'Taxa'"
                label="Target taxa"
                item-value="name"
                :ranks="TaxonRank.ranksUpTo('Family')"
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
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import {
  $SamplingInputWithEvent,
  $SamplingTargetKind,
  $SamplingUpdate,
  EventInner,
  Sampling,
  SamplingInputWithEvent,
  SamplingUpdate
} from '@/api'
import { DateWithPrecision, TaxonRank } from '@/api/adapters'
import { createSamplingMutation, updateSamplingMutation } from '@/api/gen/@tanstack/vue-query.gen'
import HabitatPicker from '../habitat/HabitatPicker.vue'
import FixativePicker from '../samples/FixativePicker.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import AccessPointsPicker from './AccessPointsPicker.vue'
import HoursMinutesInput from './HoursMinutesInput.vue'
import SamplingMethodPicker from './SamplingMethodPicker.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Sampling>()

const props = defineProps<{ event: EventInner }>()

const initial: SamplingInputWithEvent = {
  event_id: props.event.id,
  target: { kind: 'Taxa' }
}

function updateTransformer({
  duration,
  access_points,
  fixatives,
  habitats,
  target,
  comments,
  methods
}: Sampling): SamplingUpdate {
  return {
    duration,
    access_points,
    comments,
    target: {
      kind: target.kind,
      taxa: target.taxa?.map(({ name }) => name)
    },
    habitats: habitats?.map(({ label }) => label),
    fixatives: fixatives?.map(({ code }) => code),
    methods: methods?.map(({ code }) => code)
  }
}

const create = {
  mutation: createSamplingMutation,
  schema: $SamplingInputWithEvent
}

const update = {
  mutation: updateSamplingMutation,
  schema: $SamplingUpdate,
  itemID: ({ id }: Sampling) => ({ id })
}
</script>

<style scoped lang="scss"></style>
