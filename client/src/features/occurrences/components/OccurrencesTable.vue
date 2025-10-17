<template>
  <CRUDTable :items entity-name="Occurrences" :headers :search>
    <template #item.code="{ item, value }: { item: OccurrenceAtSite; value: string }">
      <div class="d-flex justify-space-between align-center">
        <RouterLink
          :to="{
            name: 'occurrence-item',
            params: { code: item.code }
          }"
          target="_blank"
        >
          <span class="text-wrap">{{ CodeIdentifier.textWrap(value) }}</span>
        </RouterLink>
        <OccurrenceIcon :item class="mx-1"></OccurrenceIcon>
      </div>
    </template>
    <template #item.site.code="{ value }: { value: string }">
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
    </template>
    <template #item.site.coordinates.latitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.site.coordinates.longitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.sampling_date="{ value }">
      <span :class="['font-monospace', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
    <template #item.taxon="{ value }">
      <TaxonChip :taxon="value" size="small" />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import { CodeIdentifier, DateWithPrecision, OccurrenceAtSite, SiteItem } from '@/api'
import TaxonChip from '@/features/taxonomy/components/TaxonChip'
import { computed, ref } from 'vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrenceIcon from './OccurrenceIcon'

type Occurrence = OccurrenceAtSite & { sampling_date?: DateWithPrecision } & (WithSite extends true
    ? { site: SiteItem }
    : {})

const { occurrences, ...props } = defineProps<{
  /**
   * If true, the site information is included with each occurrence.
   * If false, only the occurrence information is included.
   */
  withSite: WithSite
  occurrences?: Occurrence[]
}>()

const items = computed(() => {
  return occurrences?.filter((o) => {
    return (
      (occurrenceTypes.value.has('internal') && o.category === 'Internal') ||
      (occurrenceTypes.value.has('external') && o.category === 'External')
    )
  })
})

const occurrenceTypes = computed(() => new Set(search.value.occurrenceTypes))

const search = ref({
  term: undefined,
  owned: undefined,
  occurrenceTypes: ['internal', 'external', 'sequence']
})

const headersWithSites: CRUDTableHeader<Occurrence>[] = [
  {
    title: 'Code',
    value: 'code',
    sortable: true
  },
  {
    title: 'Site',
    value: 'site.code',
    width: 0,
    sortable: true
  },
  {
    title: 'Latitude',
    value: 'site.coordinates.latitude',
    width: 0,
    sortable: true
  },
  {
    title: 'Longitude',
    value: 'site.coordinates.longitude',
    sortable: true,
    width: 0
  },
  {
    title: 'Sampl. date',
    value: 'sampling_date',
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare
  },
  {
    title: 'Taxon',
    key: 'taxon',
    sortable: true,
    align: 'start',
    sort: (a, b) => a.name.localeCompare(b.name),
    filter(value, query, item) {
      return value.name.toLowerCase().includes(query.toLowerCase())
    }
  }
]

const headers = computed(() =>
  headersWithSites.filter(
    (header) =>
      props.withSite !== false ||
      typeof header.value != 'string' ||
      !header.value.startsWith('site.')
  )
)
</script>

<style scoped lang="scss"></style>
