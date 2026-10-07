<template>
  <v-card
    title="Location"
    class="small-card-title border-md"
    flat
    prepend-icon="mdi-crosshairs-gps"
  >
    <v-card-text>
      <CoordinatesInput
        v-model="sampling.coordinates"
        @update-coords="(v) => (sampling.altitude = v)"
      />
      <v-number-input
        v-model.number="sampling.coordinates.precision"
        label="Coordinates precision radius (m)"
        placeholder="Unknown"
        persistent-placeholder
        :min="0"
      />
      <v-number-input
        v-model.number="sampling.altitude"
        label="Altitude (m)"
        placeholder="Unknown"
        persistent-placeholder
      />
    </v-card-text>
    <!-- <v-divider />
            <v-card-text>
              <v-text-field label="Site name"></v-text-field>
            </v-card-text> -->
    <!-- v-bind="schema('altitude')" -->
  </v-card>
  <v-card title="Date" class="small-card-title border-sm" flat prepend-icon="mdi-calendar">
    <template #append>
      <v-checkbox
        label="Unknown"
        hide-details
        density="compact"
        v-model="unknownSamplingDate"
      ></v-checkbox>
    </template>
    <v-card-text>
      <DateWithPrecisionField v-model="sampling.performed_on" v-if="!unknownSamplingDate" />
      <span v-else class="text-subtitle-1 text-warning">
        Providing a sampling date is optional but <b>highly recommended</b>.
      </span>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import CoordinatesInput from '@/components/toolkit/forms/CoordinatesInput.vue'
import DateWithPrecisionField from '@/components/toolkit/forms/DateWithPrecisionField.vue'

import { SamplingCoreModel } from '@/models/sampling'
import { watch } from 'vue'

const sampling = defineModel<SamplingCoreModel>('sampling', {
  default: (): SamplingCoreModel => ({
    coordinates: {},
    altitude: undefined,
    performed_on: undefined
  })
})

const unknownSamplingDate = defineModel<boolean>('unknownDate', { default: false })

watch(unknownSamplingDate, (v) => {
  if (v) {
    sampling.value.performed_on = undefined
  }
})
</script>

<style scoped lang="scss"></style>
