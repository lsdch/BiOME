<template>
  <v-list density="compact">
    <v-list-item
      v-for="{ code, identification, date } in occurrences"
      :subtitle="DateWithPrecision.format(date)"
      :to="{
        name: 'occurrence-item',
        params: { code: code }
      }"
      target="_blank"
    >
      <template #title>
        <RouterLink
          :to="{
            name: 'occurrence-item',
            params: { code: code }
          }"
          class="font-monospace text-caption"
          target="_blank"
        >
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
import { DateWithPrecision, OccurrenceAtSite } from '@/api'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'

defineProps<{
  occurrences: (OccurrenceAtSite & { date?: DateWithPrecision })[]
}>()
</script>

<style scoped lang="scss"></style>
