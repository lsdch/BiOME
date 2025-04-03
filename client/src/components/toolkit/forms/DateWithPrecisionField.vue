<template>
  <div class="d-flex">
    <v-select
      v-model="model.precision"
      :label="model.precision === 'Unknown' ? 'Date' : 'Precision'"
      :items="$DatePrecision.enum"
      :rounded="model.precision === 'Unknown' ? undefined : 'e-0'"
      :max-width="150"
      v-bind="schema('precision')"
    />
    <CompositeDateField
      v-if="model.precision !== 'Unknown'"
      v-model="model.date"
      :precision="model.precision"
      rounded="s-0"
    />
  </div>
</template>

<script setup lang="ts">
import { $DatePrecision, $DateWithPrecisionInput, DateWithPrecisionInput } from '@/api'
import { useSchema } from '@/composables/schema'
import CompositeDateField from './CompositeDateField.vue'

const model = defineModel<DateWithPrecisionInput>({
  default: { date: undefined, precision: 'Day' }
})

const {
  bind: { schema }
} = useSchema($DateWithPrecisionInput)
</script>

<style lang="scss"></style>
