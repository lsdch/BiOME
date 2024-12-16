<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip :text="gene.code" class="font-monospace" v-bind="{ ...props, ...$attrs }" />
    </template>
    <v-card :title="gene.label" class="bg-surface-light" density="compact">
      <template #subtitle>
        <code>{{ gene.code }}</code>
      </template>
      <template #prepend>
        <v-badge icon="mdi-dna" color="transparent">
          <v-icon icon="mdi-tag" size="small"></v-icon>
        </v-badge>
      </template>
      <v-card-text>
        <span class="font-italic">
          {{ gene.description ?? 'No description provided' }}
        </span>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Gene } from '@/api'

defineProps<{ gene: Gene }>()
</script>

<style scoped lang="scss"></style>
