<template>
  <v-menu location="top start" origin="top start" transition="scale-transition">
    <template #activator="{ props }">
      <v-chip :text="gene.code" class="font-monospace" v-bind="{ ...props, ...$attrs }" />
    </template>
    <v-card :title="gene.label" class="bg-surface-light small-card-title" density="compact">
      <template #subtitle>
        <code>{{ gene.code }}</code>
      </template>
      <template #prepend>
        <v-badge icon="mdi-dna" color="transparent" class="mr-3">
          <v-icon icon="mdi-tag" size="small"></v-icon>
        </v-badge>
      </template>
      <v-card-text v-if="gene.description">
        {{ gene.description }}
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Gene } from '@/api'

defineProps<{ gene: Gene }>()
</script>

<style scoped lang="scss"></style>
