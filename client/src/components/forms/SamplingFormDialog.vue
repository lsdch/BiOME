<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} sampling`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container>
      <v-row v-if="hasID(event)">
        <v-col>
          <v-card variant="tonal">
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
            v-bind="schema('target', 'kind')"
            :items="$SamplingTargetKind.enum"
            v-model="model.target.kind"
            hide-details
          >
          </v-select>
        </v-col>
        <v-col>
          <TaxonPicker
            v-model="model.target.taxa"
            v-bind="
              addRules(
                schema('target', 'taxa'),
                (v: TaxonWithParentRef[] | undefined) =>
                  model.target.kind !== 'Taxa' || (v?.length ?? 0) > 0 || 'Target taxa are required'
              )
            "
            :disabled="model.target.kind !== 'Taxa'"
            label="Target taxa"
            item-value="name"
            return-object
            :ranks="TaxonRank.ranksUpTo('Family')"
            multiple
            chips
            closable-chips
            clearable
            :class="{ required: model.target.kind === 'Taxa' }"
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" sm="6">
          <FixativePicker
            label="Fixatives"
            v-model="model.fixatives"
            v-bind="schema('fixatives')"
            multiple
            return-object
            chips
            closable-chips
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6">
          <SamplingMethodPicker
            label="Sampling methods"
            v-model="model.methods"
            v-bind="schema('methods')"
            multiple
            return-object
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
            v-bind="schema('duration')"
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
            v-bind="schema('habitats')"
          />
        </v-col>
        <v-col cols="12" sm="6">
          <AccessPointsPicker
            v-model="model.access_points"
            v-bind="schema('access_points')"
            label="Access points"
            hint="Pick existing terms already in use, or enter new terms"
            persistent-hint
            clearable
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-textarea v-model="model.comments" v-bind="schema('comments')" label="Comments" />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { $SamplingInput, $SamplingTargetKind, $SamplingUpdate, EventInner } from '@/api'
import { DateWithPrecision, TaxonRank, TaxonWithParentRef } from '@/api/adapters'
import { hasID } from '@/functions/db'
import { FormProps } from '@/functions/mutations'
import { SamplingModel } from '@/models'
import { EventModel } from '@/models/event'
import HabitatPicker from '@/components/occurrence/habitat/HabitatPicker.vue'
import FixativePicker from '@/components/occurrence/FixativePicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import AccessPointsPicker from '@/components/events/AccessPointsPicker.vue'
import HoursMinutesInput from '@/components/events/HoursMinutesInput.vue'
import SamplingMethodPicker from '@/components/events/SamplingMethodPicker.vue'
import { addRules, useSchema } from '@/composables/schema'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<SamplingModel.SamplingFormModel>({
  default: SamplingModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<
  {
    event: EventInner | EventModel
  } & FormProps &
    FormDialogProps
>()

const emit = defineEmits<{
  submit: [model: SamplingModel.SamplingFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $SamplingInput : $SamplingUpdate))
</script>

<style scoped lang="scss"></style>
