<template>
  <v-text-field v-if="(samplings?.length ?? 0) > 10" v-model="search.term" class="mx-5 mt-1" hide-details label="Search"
    clearable density="compact" />
  <CRUDTable :items="samplings" entity-name="Sampling" :headers :search filter-mode="some">
    <template #item.coordinates.latitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.coordinates.longitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.performed_on="{ value }">
      <span class="font-monospace">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
    <template #expanded-row-inject="{ item: { occurrences, ...sampling } }: { item: SamplingWithOccurrences }">
      <OccurrencesAtSiteList v-if="occurrences?.length"
        :occurrences="occurrences?.map((o) => ({ ...o, sampling: sampling }))" />
    </template>
    <template v-for="name in Object.keys(slots)" :key="name" v-slot:[name]="slotProps">
      <slot :name="name" v-bind="slotProps" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { BaseOccurrence, CodeIdentifier, DateWithPrecision, Sampling } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrencesAtSiteList from '@/features/occurrences/components/OccurrencesAtSiteList.vue'
import { HeaderExtension, mergeHeaders } from '@/features/occurrences/components/tables/headers'
import { computed, ref, useSlots } from 'vue'

type SamplingWithOccurrences = Sampling & { occurrences?: BaseOccurrence[] }

const { samplings, ...props } = defineProps<{
  samplings?: SamplingWithOccurrences[]
  extendHeaders?: HeaderExtension<Sampling>[]
}>()

const search = ref({ term: undefined, owned: undefined })

const _headers: CRUDTableHeader<Sampling>[] = [
  {
    title: 'Site',
    value: 'site.name',
    sortable: true,
    // filter(value, query, item) {
    //   if (!query) return true
    //   return (
    //     item.raw.site?.code.toLowerCase().includes(query.toLowerCase()) ||
    //     !!item.raw.site?.name?.toLowerCase().includes(query.toLowerCase())
    //   )
    // }
  },
  {
    title: 'Latitude',
    value: 'coordinates.latitude',
    width: 0,
    sortable: true,
    // Disable filtering for numeric fields
    filter: () => false
  },
  {
    title: 'Longitude',
    value: 'coordinates.longitude',
    sortable: true,
    width: 0,
    // Disable filtering for numeric fields
    filter: () => false
  },
  {
    title: 'Date',
    value: 'performed_on',
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare,
    filter(value, query, item) {
      if (!query) return true
      return DateWithPrecision.format(value).toLowerCase().includes(query.toLowerCase())
    }
  },
  {
    title: 'Occurrences',
    key: 'occurrences',
    value: (item) => {
      return item.occurrences?.length || 0
    },
    cellProps: { class: 'font-monospace' },
    sortable: true,
    width: 0,
    align: 'end'
  }
] as const

const headers = computed(() => {
  return mergeHeaders(_headers, props.extendHeaders) satisfies CRUDTableHeader<Sampling>[]
})

const slots = useSlots()
</script>

<style scoped lang="scss"></style>
