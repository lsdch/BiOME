<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip :text="program.label" v-bind="{ ...props, ...$attrs }" />
    </template>
    <v-card
      :max-width="400"
      :title="program.label"
      :subtitle="program.code"
      class="bg-surface-light"
      density="compact"
    >
      <v-card-text>{{ program.description }}</v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { ProgramInner } from '@/api'

defineProps<{ program: ProgramInner }>()
</script>

<style scoped lang="scss"></style>
