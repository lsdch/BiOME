<template>
  <div v-if="(occurrences?.length ?? 0) > 10" class="d-flex align-center ga-2">
    <v-text-field
      v-model="search.term"
      class="mx-5"
      hide-details
      label="Search"
      clearable
      density="compact"
    />
  </div>
  <CRUDTable :items entity-name="Occurrences" :headers :search>
    <template #item.code="{ item, value }: { item: OccurrenceAtSite; value: string }">
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
  generic="Data extends OccurrenceAtSite, WithSite extends boolean, Site extends SiteItem"
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

export type OccurrenceTableItem<
  Data extends OccurrenceAtSite,
  WithSite extends boolean,
  Site extends SiteItem
> = Data & {
  sampling_date?: DateWithPrecision
} & (WithSite extends true ? { site: Site } : {})

type Occurrence = OccurrenceTableItem<Data, WithSite, Site>
const { occurrences, ...props } = defineProps<{
  /**
   * If true, the site information is included with each occurrence.
   * If false, only the occurrence information is included.
   */
  withSite: WithSite
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
  Identification.tableHeader({ key: 'identification' })
] as const satisfies CRUDTableHeader<Occurrence>[]

const headers = computed(() => {
  const h = headersWithSites.filter(
    (header) =>
      props.withSite !== false ||
      !(typeof header.value === 'string' && header.value?.startsWith('site.'))
  )
  return mergeHeaders(h, props.extendHeaders) satisfies CRUDTableHeader<Occurrence>[]
})

const slots = useSlots()
</script>

<style scoped lang="scss"></style>
