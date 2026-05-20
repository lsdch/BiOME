<template>
  <div v-if="(occurrences?.length ?? 0) > 10" class="d-flex align-center ga-2">
    <v-text-field
      v-model="search.term"
      class="ma-2"
      hide-details
      label="Search"
      density="compact"
      clearable
    />
  </div>
  <CRUDTable :items entity-name="Occurrences" :headers :search filter-mode="some">
    <template #item.code="{ item, value }: { item: Data; value: string }">
      <RouterLink
        :to="{
          name: 'occurrence-item',
          params: { code: item.code }
        }"
        target="_blank"
      >
        <span class="text-wrap">{{ CodeIdentifier.textWrap(value) }}</span>
      </RouterLink>
    </template>
    <template #item.site="{ value }: { value: SiteItem }" v-if="withSite">
      <RouterLink
        :to="{
          name: 'site-item',
          params: { code: value.code }
        }"
        target="_blank"
        v-tooltip="value.name"
      >
        <span class="text-wrap font-monospace">
          {{ CodeIdentifier.textWrap(value.code) }}
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
    <template #item.identification="{ value: identification }">
      <IdentificationChip :identification size="small" />
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
    Data extends OccurrenceAtSite | (OccurrenceAtSite & { site: Site }),
    Site extends SiteItem
  "
>
import {
  CodeIdentifier,
  DateWithPrecision,
  Identification,
  OccurrenceAtSite,
  SiteItem
} from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { HeaderExtension, mergeHeaders } from '@/features/occurrences/components/tables/headers'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import { computed, ref, useSlots } from 'vue'

export type OccurrenceTableItem<Site> = OccurrenceAtSite & {
  sampling_date?: DateWithPrecision
  site: Site
}

type Occurrence = Data

const {
  occurrences,
  withSite = false,
  ...props
} = defineProps<{
  /**
   * If true, the site information is included with each occurrence.
   * If false, only the occurrence information is included.
   */
  withSite?: boolean
  occurrences?: Occurrence[]
  extendHeaders?: HeaderExtension<Occurrence>[]
}>()

const items = computed(() => {
  return occurrences
})

const search = ref({
  term: undefined,
  owned: undefined
})

const headersWithSites = [
  {
    title: 'Code',
    value: 'code',
    sortable: true
  },
  {
    title: 'Site',
    value: 'site',
    width: 0,
    sortable: true,
    sort: (a, b) => a?.code.localeCompare(b?.code) || 0,
    filter(value: SiteItem, query, item) {
      if (!query) return true
      return (
        value.code.toLowerCase().includes(query.toLowerCase()) ||
        !!value.name?.toLowerCase().includes(query.toLowerCase())
      )
    }
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
  Identification.tableHeader({ key: 'identification' })
] as const satisfies CRUDTableHeader<Data>[]

const headers = computed(() => {
  const h = headersWithSites.filter(
    (header) => withSite || !(typeof header.value === 'string' && header.value?.startsWith('site.'))
  ) as CRUDTableHeader<Occurrence>[]
  return mergeHeaders(h, props.extendHeaders) satisfies CRUDTableHeader<Occurrence>[]
})

const slots = useSlots()
</script>

<style scoped lang="scss"></style>
