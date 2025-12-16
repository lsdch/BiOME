<template>
  <v-text-field
    v-if="(samplings?.length ?? 0) > 10"
    v-model="search.term"
    class="mx-5"
    hide-details
    label="Search"
    clearable
    density="compact"
  />
  <CRUDTable :items="samplings" entity-name="Sampling" :headers :search filter-mode="some">
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
    <template #item.date="{ value }">
      <span class="font-monospace">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
    <template #expanded-row-inject="{ item }: { item: SamplingEvent }">
      <OccurrencesAtSiteList
        v-if="item.occurrences.length"
        :occurrences="item.occurrences.map((o) => ({ ...o, date: item.date }))"
      />
    </template>
    <template v-for="name in Object.keys(slots)" :key="name" v-slot:[name]="slotProps">
      <slot :name="name" v-bind="slotProps" />
    </template>
  </CRUDTable>
</template>

<script
  setup
  lang="ts"
  generic="
    WithSite extends boolean,
    SamplingData extends SamplingDateWithOccurrences,
    Site extends SiteItem
  "
>
import { CodeIdentifier, DateWithPrecision, SamplingDateWithOccurrences, SiteItem } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrencesAtSiteList from '@/features/occurrences/components/OccurrencesAtSiteList.vue'
import { HeaderExtension, mergeHeaders } from '@/features/occurrences/components/tables/headers'
import { computed, ref, useSlots } from 'vue'

export type SamplingTableItem<
  S extends SamplingDateWithOccurrences,
  WithSite extends boolean,
  Site extends SiteItem
> = S & (WithSite extends true ? { site: Site } : {})

type SamplingEvent = SamplingTableItem<SamplingData, WithSite, Site>

const { samplings, ...props } = defineProps<{
  /**
   * If true, the site information is included with each occurrence.
   * If false, only the occurrence information is included.
   */
  withSite: WithSite
  samplings?: SamplingEvent[]
  extendHeaders?: HeaderExtension<SamplingEvent>[]
}>()

const search = ref({ term: undefined, owned: undefined })

const headersWithSite = [
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
    sortable: true,
    // Disable filtering for numeric fields
    filter: () => false
  },
  {
    title: 'Longitude',
    value: 'site.coordinates.longitude',
    sortable: true,
    width: 0,
    // Disable filtering for numeric fields
    filter: () => false
  },
  {
    title: 'Date',
    value: 'date',
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare,
    filter(value, query, item) {
      if (!query) return true
      return DateWithPrecision.format(item.raw.date).toLowerCase().includes(query.toLowerCase())
    }
  },
  {
    title: 'Occurrences',
    key: 'occurrences',
    value: (item: SamplingEvent) => {
      return item.occurrences.length
    },
    cellProps: { class: 'font-monospace' },
    sortable: true,
    width: 0,
    align: 'end'
  }
] as const satisfies CRUDTableHeader<SamplingEvent>[]

const headers = computed(() => {
  const h = headersWithSite.filter(
    (header) =>
      props.withSite !== false ||
      !(typeof header.value === 'string' && header.value?.startsWith('site.'))
  )
  return mergeHeaders(h, props.extendHeaders) satisfies CRUDTableHeader<SamplingEvent>[]
})

const slots = useSlots()
</script>

<style scoped lang="scss"></style>
