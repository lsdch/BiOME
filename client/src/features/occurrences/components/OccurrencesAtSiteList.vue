<template>
  <v-list density="compact">
    <v-list-item v-for="{ id, code, identification, sampling } in occurrences">
      <template #title>
        <RouterLink :to="{
          name: 'occurrence-item',
          params: { id, code: code }
        }" class="font-monospace text-label-medium" target="_blank">
          {{ code }}
        </RouterLink>
      </template>
      <template #append>
        <IdentificationChip :identification size="small" />
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { DateWithPrecision, Occurrence } from '@/api';
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip';

defineProps<{
  occurrences: Occurrence[]
}>()
</script>

<style scoped lang="scss"></style>
