<template>
  <v-container>
    <HabitatGroupBox v-for="group in habitatGraph" :key="group.id" :group="group" open />
  </v-container>
</template>

<script setup lang="ts">
import HabitatGroupBox from '@/features/registries/components/HabitatGroupBox.vue'
import { useHabitats } from '../composables/habitats'
const { habitatGraph } = useHabitats()
</script>

<style scoped></style>
