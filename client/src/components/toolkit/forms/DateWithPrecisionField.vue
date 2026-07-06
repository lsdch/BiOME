<template>
  <div class="d-flex">
    <v-select
      v-model="model.precision"
      label="Precision"
      :items
      :rounded="model.precision !== 'Unknown' ? 'e-0' : undefined"
      :max-width="150"
      v-bind="schema('precision')"
    />
    <CompositeDateField
      v-model="model.date"
      v-if="model.precision !== 'Unknown'"
      :precision="model.precision"
      rounded="s-0"
    />
  </div>
</template>

<script setup lang="ts">
import { $EventDatePrecision, DateWithPrecisionInput } from '@/api'
import { useSchema } from '@/composables/schema'
import { $schema, DateWithPrecisionModel } from '@/models/date_with_precision'
import CompositeDateField from './CompositeDateField.vue'

const model = defineModel<DateWithPrecisionModel>({
  default: { date: undefined, precision: 'Day' }
})

const {
  bind: { schema }
} = useSchema($schema)

const items: Array<DateWithPrecisionInput['precision'] | 'Unknown'> = [
  ...$EventDatePrecision.enum,
  'Unknown'
]
</script>

<style lang="scss"></style>
