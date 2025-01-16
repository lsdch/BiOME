<template>
  <TaxonPicker v-model="model.taxon" item-value="name" />
  <PersonPicker v-model="model.identified_by" label="Curator" />
  <DateWithPrecisionField v-model="model.identified_on" />
</template>

<script setup lang="ts">
import { IdentificationInput, IdentificationUpdate } from '@/api'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import PersonPicker from '../people/PersonPicker.vue'
import DateWithPrecisionField from '../events/DateWithPrecisionField.vue'

const model = defineModel<IdentificationInput | IdentificationUpdate>({ required: true })
</script>

<style scoped lang="scss"></style>
