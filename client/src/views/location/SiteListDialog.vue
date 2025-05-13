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
        <v-list density="compact">
          <v-list-item
            v-for="o in item.occurrences"
            :subtitle="DateWithPrecision.format(o.sampling_date)"
            :to="{
              name: o.element === 'Sequence' ? 'sequence' : 'biomat-item',
              params: { code: o.code }
            }"
          >
            <template #prepend>
              <v-icon
                :icon="o.element === 'Sequence' ? 'mdi-dna' : 'mdi-package-variant'"
                :color="o.category === 'Internal' ? 'primary' : 'warning'"
              />
            </template>
            <template #title>
              <RouterLink
                :to="{
                  name: o.element === 'Sequence' ? 'sequence' : 'biomat-item',
                  params: { code: o.code }
                }"
                class="font-monospace text-caption"
                >{{ o.code }}</RouterLink
              >
            </template>
            <template #append>
              <TaxonChip :taxon="o.taxon" size="small" />
            </template>
          </v-list-item>
        </v-list>
      </template>
    </CRUDTable>
  </CardDialog>
</template>

<script setup lang="ts">
import { CodeIdentifier, DateWithPrecision, SiteWithOccurrences } from '@/api'
import CountryChip from '@/components/sites/CountryChip'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'

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
      return item.occurrences?.length
    },
    sortable: true
  }
]
</script>

<style scoped lang="scss"></style>
