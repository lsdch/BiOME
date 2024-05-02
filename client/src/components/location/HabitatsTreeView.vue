<template>
  <v-treeview activatable> </v-treeview>
</template>

<script setup lang="ts">
import { LocationService } from '@/api'

const habitats = await LocationService.listHabitats()

console.log(habitats)
</script>

<style scoped></style>
