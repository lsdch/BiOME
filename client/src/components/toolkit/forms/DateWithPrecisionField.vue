<template>
  <div class="d-flex">
    <v-select
      v-model="model.precision"
      label="Precision"
      :items="$EventDatePrecision.enum"
      rounded="e-0"
      :max-width="150"
      :min-width="100"
      v-bind="schema('precision')"
    />
    <CompositeDateField v-model="model.date" :precision="model.precision" rounded="s-0" />
  </div>
</template>

<script setup lang="ts">
import { $EventDatePrecision } from '@/api'
import { useSchemaBinding } from '@/composables/schema'
import { $schema, DateWithPrecisionModel } from '@/models/date_with_precision'
import { reactive } from 'vue'
import CompositeDateField from './CompositeDateField.vue'

const model = defineModel<DateWithPrecisionModel>({
  default: () => reactive({ date: {}, precision: 'day' })
})

const { schema } = useSchemaBinding($schema)
</script>

<style lang="scss"></style>
