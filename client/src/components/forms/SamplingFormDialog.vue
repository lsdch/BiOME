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
      <v-row>
        <v-col>
          <v-card variant="tonal">
            <template #title>
              {{ site.name }}
            </template>
            <v-card-text>
              <v-row>
                <v-col>
                  <!-- No schema binding, component enforces constraints on its own -->
                  <DateWithPrecisionField v-model="model.performed_on" />
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <v-combobox
                    label="Performed by"
                    v-model.trim="model.performed_by"
                    multiple
                    chips
                    closable-chips
                    v-bind="schema('performed_by')"
                  />
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <TaxonPicker
            v-model="model.target_taxa"
            v-bind="schema('target_taxa')"
            label="Target taxa"
            item-value="name"
            return-object
            multiple
            chips
            closable-chips
            clearable
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
import { $SamplingInput, $SamplingUpdate } from '@/api'
import { SiteItem } from '@/api/adapters'
import DateWithPrecisionField from '@/components/toolkit/forms/DateWithPrecisionField.vue'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import HoursMinutesInput from '@/components/toolkit/forms/HoursMinutesInput.vue'
import { useSchema } from '@/composables/schema'
import AccessPointsPicker from '@/features/occurrences/components/sampling/AccessPointsPicker.vue'
import FixativePicker from '@/features/registries/components/FixativePicker.vue'
import HabitatPicker from '@/features/registries/components/HabitatPicker.vue'
import SamplingMethodPicker from '@/features/registries/components/SamplingMethodPicker.vue'
import TaxonPicker from '@/features/taxonomy/components/TaxonPicker.vue'
import { FormProps } from '@/lib/mutations'
import { SamplingModel } from '@/models'
import { SiteFormModel } from '@/models/site'
import { reactiveComputed } from '@vueuse/core'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<SamplingModel.SamplingFormModel>({
  default: SamplingModel.initialModel
})

const { mode = 'Create', ...props } = defineProps<
  {
    site: SiteItem | SiteFormModel
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
