<template>
  <CardDialog :title v-bind="props">
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>
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
      <template #item.target="{ value: { kind, taxa } }: { value: SamplingTarget }">
        <v-chip :text="kind" size="small" label />
        <span v-if="taxa?.length" class="text-overline ml-1">|</span>
        <TaxonChip v-for="taxon in taxa" :taxon size="small" class="ma-1" />
      </template>
      <template #item.occurrences="{ value, toggleExpand, item }">
        <v-badge inline :content="value" color="success" @click="toggleExpand(item)" />
      </template>
      <template #expanded-row-inject="{ item }: { item: SamplingEvent }">
        <OccurrenceAtSiteList
          v-if="item.occurrences.length"
          :occurrences="item.occurrences.map((o) => ({ ...o, date: item.date }))"
        />
      </template>
    </CRUDTable>
  </CardDialog>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import {
  CodeIdentifier,
  DateWithPrecision,
  SamplingEventWithOccurrences,
  SamplingTarget,
  SiteItem
} from '@/api'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import { computed, ref } from 'vue'
import { ComponentSlots } from 'vue-component-type-helpers'
import OccurrenceAtSiteList from './OccurrenceAtSiteList.vue'

type SamplingEvent = SamplingEventWithOccurrences &
  (WithSite extends true ? { site: SiteItem } : {})

const {
  title = 'Sampling events',
  samplings,
  ...props
} = defineProps<
  {
    /**
     * If true, the site information is included with each occurrence.
     * If false, only the occurrence information is included.
     */
    withSite: WithSite
    samplings?: SamplingEvent[]
  } & CardDialogProps
>()

const search = ref({ term: undefined, owned: undefined })

const headersWithSite: CRUDTableHeader<SamplingEvent>[] = [
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
    title: 'Target',
    value: 'target',
    sort: (a: SamplingTarget, b: SamplingTarget) => a.kind.localeCompare(b.kind),
    filter(value, query, item) {
      if (!query) return true
      if (item.raw.target.kind === 'Taxa') {
        return (
          item.raw.target.taxa?.some((taxon) =>
            taxon.name.toLowerCase().includes(query.toLowerCase())
          ) ?? false
        )
      }
      return item.raw.target.kind.toLowerCase().includes(query.toLowerCase()) ?? false
    },
    sortable: true,
    align: 'start'
  },
  {
    title: 'Occurrences',
    key: 'occurrences',
    value: (item: SamplingEvent) => {
      return item.occurrences.length
    },
    sortable: true,
    width: 0,
    align: 'end'
  }
]

const headers = computed(() =>
  headersWithSite.filter(
    (header) =>
      props.withSite !== false ||
      !(typeof header.value === 'string' && header.value?.startsWith('site.'))
  )
)

defineSlots<ComponentSlots<typeof CardDialog>>()
</script>

<style scoped lang="scss"></style>
