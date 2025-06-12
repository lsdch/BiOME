<template>
  <CardDialog :title v-bind="props">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <CRUDTable :items="sites" entity-name="Sites" :headers>
      <template #item.code="{ value }: { value: string }">
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
      <template #item.country="{ value }">
        <CountryChip :country="value" size="small" />
      </template>
      <template #item.last_visited="{ value }">
        <span class="font-monospace">
          {{ DateWithPrecision.format(value) }}
        </span>
      </template>
      <template #expanded-row-inject="{ item }: { item: SiteWithOccurrences }">
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
    </CRUDTable>
  </CardDialog>
</template>

<script setup lang="ts">
import { CodeIdentifier, DateWithPrecision, SiteWithOccurrences } from '@/api'
import CountryChip from '@/components/sites/CountryChip'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrenceAtSiteList from './OccurrenceAtSiteList.vue'

const {
  sites,
  title = 'Sites',
  ...props
} = defineProps<
  {
    sites?: SiteWithOccurrences[]
  } & CardDialogProps
>()

const headers: CRUDTableHeader<SiteWithOccurrences>[] = [
  {
    title: 'Code',
    value: 'code',
    sortable: true
  },
  {
    title: 'Latitude',
    value: 'coordinates.latitude',
    width: 0,
    sortable: true
  },
  { title: 'Longitude', value: 'coordinates.longitude', sortable: true, width: 0 },
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
    value: 'last_visited',
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
      return item.samplings.reduce((sum, s) => sum + s.occurrences?.length, 0)
    },
    sortable: true
  }
]
</script>

<style scoped lang="scss"></style>
