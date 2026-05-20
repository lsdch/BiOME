<template>
  <v-text-field
    v-if="(sites?.length ?? 0) > 10"
    v-model="search.term"
    class="mx-5 mt-1"
    hide-details
    label="Search"
    clearable
    density="compact"
  />
  <CRUDTable :items="sites" entity-name="Sites" :headers :search>
    <template #item.code="{ value, item }: { value: string; item: SiteWithOccurrences }">
      <div class="d-flex flex-column">
        <RouterLink
          :to="{
            name: 'site-item',
            params: { code: value }
          }"
          target="_blank"
        >
          <span class="text-wrap font-monospace">
            {{ CodeIdentifier.textWrap(value) }}
          </span>
        </RouterLink>
        <span class="text-muted" v-if="item.name">
          {{ item.name }}
        </span>
      </div>
    </template>
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
    <template #item.coordinates.precision="{ value }">
      <CoordPrecisionChip :precision="value" size="small" />
    </template>
    <template #item.country="{ value }">
      <CountryChip :country="value" size="small" />
    </template>
    <template #item.last_visited="{ value, item }">
      <span class="font-monospace">
        {{ item.samplings.length ? DateWithPrecision.format(value) : 'Never' }}
      </span>
    </template>
    <template #expanded-row-inject="{ item }: { item: Data }">
      <OccurrenceAtSiteList
        v-if="item.samplings.length"
        :occurrences="
          item.samplings.flatMap((sampling) =>
            sampling.occurrences.map((occurrence) => ({
              ...occurrence,
              date: sampling.date
            }))
          )
        "
      />
    </template>
    <template v-for="name in Object.keys(slots)" :key="name" v-slot:[name]="slotProps">
      <slot :name="name" v-bind="slotProps" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts" generic="Data extends SiteWithOccurrences">
import { CodeIdentifier, DateWithPrecision, SiteWithOccurrences } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrenceAtSiteList from '@/features/occurrences/components/OccurrencesAtSiteList.vue'
import { HeaderExtension, mergeHeaders } from '@/features/occurrences/components/tables/headers'
import CoordPrecisionChip from '@/features/site/components/CoordPrecisionChip'
import CountryChip from '@/features/site/components/CountryChip'
import { lastSamplingDate } from '@/lib/dates'
import { computed, ref, useSlots } from 'vue'

const { sites, extendHeaders } = defineProps<{
  sites?: Data[]
  extendHeaders?: HeaderExtension<Data>[]
}>()

const baseHeaders = [
  {
    title: 'Code',
    value: 'code',
    sortable: true,
    filter(value, query, item) {
      if (!query) return true
      return (
        item.raw.code.toLowerCase().includes(query.toLowerCase()) ||
        !!item.raw.name?.toLowerCase().includes(query.toLowerCase())
      )
    }
  },
  {
    title: 'Latitude',
    value: 'coordinates.latitude',
    width: 0,
    sortable: true
  },
  { title: 'Longitude', value: 'coordinates.longitude', sortable: true, width: 0 },
  { title: 'Precision', value: 'coordinates.precision', sortable: true, width: 0 },
  { title: 'Locality', value: 'locality', sortable: true },
  {
    title: 'Country',
    value: 'country',
    sortable: true,
    width: 0,
    sort: (a, b) => a?.code?.localeCompare(b?.code)
  },
  {
    title: 'Last visited',
    key: 'last_visited',
    value: (item: Data) => {
      return lastSamplingDate(item)
    },
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare
  },
  {
    title: 'Occurrences',
    key: 'occurrences',
    align: 'end',
    width: 0,
    value(item: SiteWithOccurrences) {
      return item.samplings.reduce((sum, s) => sum + (s.occurrences?.length ?? 0), 0)
    },
    sortable: true
  }
] as const satisfies CRUDTableHeader<Data>[]

const headers = computed(() => {
  return mergeHeaders(baseHeaders, extendHeaders)
})

const slots = useSlots()

const search = ref({
  term: ''
})
</script>

<style scoped lang="scss"></style>
