<template>
  <div class="d-flex">
    <TaxonPicker
      v-model="model.taxon_id"
      :min-width="500"
      item-value="id"
      :ranks="TaxonRank.ranksUpTo('FAMILY')"
      v-bind="schema('taxon_id')"
    >
      <template #append-inner>
        <v-btn
          :active="model.confer"
          @click="model.confer = !model.confer"
          color=""
          active-color="primary"
          :variant="model.confer ? 'outlined' : 'plain'"
          rounded="sm"
          text="cf."
          size="small"
          class="font-monospace"
          v-tooltip="{
            location: 'top',
            text: 'When active, indicates a tentative identification'
          }"
        />
      </template>
    </TaxonPicker>
    <v-text-field
      v-model="model.addendum"
      label="Addendum"
      class="ml-2 flex-grow-0"
      :min-width="250"
      v-bind="schema('addendum')"
      placeholder="e.g. form A, group B..."
    />
  </div>
  <v-combobox
    v-model.trim="model.identified_by"
    label="Curator(s)"
    chips
    closable-chips
    v-bind="schema('identified_by')"
  />
  <DateWithPrecisionField v-model="model.identified_on" v-bind="schema('identified_on')" />
</template>

<script setup lang="ts">
import { $IdentificationInput, IdentificationInput, Taxon, TaxonRank } from '@/api'
import DateWithPrecisionField from '@/components/toolkit/forms/DateWithPrecisionField.vue'
import { useSchema } from '@/composables/schema'
import TaxonPicker from '@/features/taxonomy/components/TaxonPicker.vue'
import { FormProps } from '@/lib/mutations'
import { DateWithPrecisionModel } from '@/models/date_with_precision'
import { reactiveComputed } from '@vueuse/core'

export type IdentificationModel = {
  identified_on: DateWithPrecisionModel
  identified_by?: string[]
  confer?: boolean
  addendum?: string
  taxon?: Taxon
}

const model = defineModel<Partial<IdentificationInput>>({ required: true })

const { mode = 'Create' } = defineProps<FormProps>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema($IdentificationInput))
</script>

<style scoped lang="scss"></style>
