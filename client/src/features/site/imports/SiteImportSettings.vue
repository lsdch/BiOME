<template>
  <v-row>
    <v-col>
      <v-select
        label="Existing codes"
        v-model="model.existing"
        item-value="label"
        item-title="label"
        :items="existing"
      >
        <template #item="{ item, props }">
          <v-list-item :title="item.raw.label" :subtitle="item.raw.description" v-bind="props" />
        </template>
      </v-select>
    </v-col>
    <v-col>
      <CoordPrecisionPicker
        v-model="model.defaultPrecision"
        label="Default coord precision"
        clearable
        placeholder="None"
        persistent-placeholder
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { CoordinatesPrecision } from '@/api'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'

const existing = [
  {
    label: 'Restrict',
    description: 'Alert when a code is already in use'
  },
  {
    label: 'Omit',
    description: 'Omit existing codes from the dataset'
  },
  {
    label: 'Include',
    description: 'Include existing codes in the dataset'
  }
] as const

export type ExistingSetting = (typeof existing)[number]['label']

export type ImportSettings = {
  existing: ExistingSetting
  defaultPrecision?: CoordinatesPrecision
}

const model = defineModel<ImportSettings>({
  default: { existing: 'Include' }
})
</script>

<style scoped></style>
