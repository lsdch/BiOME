<template>
  <TaxonPicker v-model="model.taxon" return-object v-bind="schema('taxon')" />
  <PersonPicker
    v-model="model.identified_by"
    label="Curator"
    return-object
    v-bind="schema('identified_by')"
  />
  <DateWithPrecisionField v-model="model.identified_on" v-bind="schema('identified_on')" />
</template>

<script setup lang="ts">
import {
  $IdentificationInput,
  $IdentificationUpdate,
  DateWithPrecisionInput,
  PersonUser,
  Taxon
} from '@/api'
import DateWithPrecisionField from '@/components/toolkit/forms/DateWithPrecisionField.vue'
import PersonPicker from '@/components/people/PersonPicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import { useSchema } from '@/composables/schema'
import { reactiveComputed } from '@vueuse/core'
import { FormProps } from '@/functions/mutations'

export type IdentificationModel = {
  identified_on: DateWithPrecisionInput
  identified_by?: PersonUser
  taxon?: Taxon
}

const model = defineModel<IdentificationModel>({ required: true })

const { mode = 'Create' } = defineProps<FormProps>()

const {
  bind: { schema }
} = reactiveComputed(() =>
  useSchema(mode === 'Create' ? $IdentificationInput : $IdentificationUpdate)
)
</script>

<style scoped lang="scss"></style>
