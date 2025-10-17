<template>
  <v-list density="compact">
    <v-list-item
      v-for="{ code, category, taxon, date } in occurrences"
      :subtitle="DateWithPrecision.format(date)"
      :to="{
        name: 'occurrence-item',
        params: { code: code }
      }"
      target="_blank"
    >
      <template #prepend>
        <v-icon
          icon="mdi-package-variant"
          :color="category === 'Internal' ? 'primary' : 'warning'"
        />
      </template>
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
        <TaxonChip :taxon="taxon" size="small" />
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { DateWithPrecision, OccurrenceAtSite } from '@/api'

defineProps<{
  occurrences: (OccurrenceAtSite & { date: DateWithPrecision })[]
}>()
</script>

<style scoped lang="scss"></style>
