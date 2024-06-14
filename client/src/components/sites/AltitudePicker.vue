<template>
  <v-row>
    <v-col>
      <v-select v-model="mode" label="Type" :items="Modes"></v-select>
    </v-col>
    <v-col>
      <NumberInput v-show="mode != 'NA'" label="Min" />
    </v-col>
    <v-col>
      <NumberInput v-show="mode == 'Range'" label="Max" prefix="m" />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import { AltitudeRange } from '@/api'

const Modes = ['NA', 'Point', 'Range'] as const
type Mode = (typeof Modes)[number]

const mode = ref<Mode>('Point')

const model = defineModel<AltitudeRange | undefined>({ required: true })
</script>

<style scoped></style>
