<template>
  <div class="d-flex">
    <v-select
      :label="model.precision === 'Unknown' ? 'Date' : 'Precision'"
      v-model="model.precision"
      :items="$DatePrecision.enum"
      :rounded="model.precision === 'Unknown' ? undefined : 'e-0'"
      v-bind="field('precision')"
      :max-width="150"
    />
    <CompositeDateField
      v-if="model.precision !== 'Unknown'"
      v-model="model.date"
      :precision="model.precision"
      rounded="s-0"
      v-bind="field('date')"
    />
  </div>
</template>

<script setup lang="ts">
import { $DatePrecision, $DateWithPrecision, DateWithPrecisionInput } from '@/api'
import { useSchema } from '../toolkit/forms/schema'
import CompositeDateField from './CompositeDateField.vue'

const model = defineModel<DateWithPrecisionInput>({
  default: { date: undefined, precision: 'Day' }
})

const { field } = useSchema($DateWithPrecision)
</script>

<style lang="scss"></style>
