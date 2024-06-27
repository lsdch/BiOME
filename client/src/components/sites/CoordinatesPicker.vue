<template>
  <v-row>
    <v-col cols="12" sm="4">
      <CoordPrecisionPicker v-model="model.precision" v-bind="schema('precision')" />
    </v-col>
    <v-col cols="6" sm="4">
      <NumberInput
        v-model="model.latitude"
        label="Latitude"
        float
        :step="0.001"
        v-bind="schema('latitude')"
      />
    </v-col>
    <v-col cols="6" sm="4">
      <NumberInput
        v-model="model.longitude"
        label="Longitude"
        float
        :step="0.001"
        v-bind="schema('longitude')"
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { $Coordinates, Coordinates } from '@/api'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'
import { useForm, useSchema } from '../toolkit/forms/form'

const model = defineModel<Coordinates>({ required: true })

const { schema } = useSchema($Coordinates)
</script>

<style scoped></style>
