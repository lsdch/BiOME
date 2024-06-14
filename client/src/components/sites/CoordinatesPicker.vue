<template>
  <v-row>
    <v-col>
      <v-select v-model="model.precision" :items="$CoordinatePrecision.enum"></v-select>
    </v-col>
    <v-col>
      <NumberInput
        v-model="model.latitude"
        label="Latitude"
        float
        v-show="model.precision != 'Unknown'"
      ></NumberInput>
    </v-col>
    <v-col>
      <NumberInput
        v-model="model.longitude"
        label="Longitude"
        float
        v-show="model.precision != 'Unknown'"
      ></NumberInput>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { $CoordinatePrecision, Coordinates } from '@/api'
import NumberInput from '../toolkit/ui/NumberInput.vue'

const model = defineModel<Coordinates>({ required: true })
</script>

<style scoped></style>
