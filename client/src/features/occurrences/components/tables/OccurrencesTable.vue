<template>
  <div v-if="(occurrences?.length ?? 0) > 10" class="d-flex align-center ga-2">
    <v-text-field v-model="search.term" class="ma-2" hide-details label="Search" density="compact" clearable />
  </div>
  <CRUDTable :items entity-name="Occurrences" :headers :search filter-mode="some">
    <template #item.code="{ item, value }: { item: Occurrence; value: string }">
      <RouterLink :to="{
        name: 'occurrence-item',
        params: { id: item.id, code: item.code }
      }" target="_blank">
        <span class="text-wrap">{{ CodeIdentifier.textWrap(value) }}</span>
      </RouterLink>
    </template>
    <template #item.sampling.coordinates.latitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.sampling.coordinates.longitude="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.sampling.performed_on="{ value }">
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

<script setup lang="ts" generic="">
import {
  CodeIdentifier,
  DateWithPrecision,
  Identification,
  Occurrence,
} from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { HeaderExtension, mergeHeaders } from '@/features/occurrences/components/tables/headers'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import { computed, ref, useSlots } from 'vue'
import { FilterMatch } from 'vuetify'


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

const _headers = [
  {
    title: 'Code',
    value: 'code',
    sortable: true
  },
  {
    title: 'Site',
    value: 'sampling.site.name',
    width: 0,
    sortable: true,
    sort: (a, b) => a?.code.localeCompare(b?.code) || 0,
    // filter(value: SiteItem, query, item) {
    //   if (!query) return true
    //   return (
    //     value.code.toLowerCase().includes(query.toLowerCase()) ||
    //     !!value.name?.toLowerCase().includes(query.toLowerCase())
    //   )
    // }
  },
  {
    title: 'Latitude',
    value: 'sampling.coordinates.latitude',
    width: 0,
    sortable: true
  },
  {
    title: 'Longitude',
    value: 'sampling.coordinates.longitude',
    sortable: true,
    width: 0
  },
  {
    title: 'Sampl. date',
    value: 'sampling.performed_on',
    sortable: true,
    align: 'end',
    sort: DateWithPrecision.compare
  },
  {
    title: 'Taxon',
    key: 'identification',
    sortable: true,
    align: 'start',
    sort: (a, b) => a.taxon.name.localeCompare(b.name),
    filter(value, query, item): FilterMatch {
      return (value as unknown as Identification).taxon.name
        .toLowerCase()
        .includes(query.toLowerCase())
    }
  }
] as const satisfies CRUDTableHeader<Occurrence>[]

const headers = computed(() => {
  return mergeHeaders(_headers, props.extendHeaders) satisfies CRUDTableHeader<Occurrence>[]
})

const slots = useSlots()
</script>

<style scoped lang="scss"></style>
