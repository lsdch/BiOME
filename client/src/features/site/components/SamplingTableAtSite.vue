<template>
  <CRUDTable :items="samplings" :headers entityName="Sampling">
    <template #item.number="{ value }">
      <code>
        {{ value }}
      </code>
    </template>
    <template #item.performed_on="{ value, item }">
      <SamplingCardDialog :sampling="item">
        <template #activator="{ props }">
          <a :class="['font-monospace cursor-pointer', { 'text-muted': !value }]" v-bind="props">
            {{ DateWithPrecision.format(value) }}
          </a>
        </template>
      </SamplingCardDialog>
    </template>
    <template #item.target_taxa="{ value }">
      <div class="d-flex align-center ga-1">
        <TaxonChip v-for="taxon in value" :taxon />
      </div>
    </template>
    <template #item.occurrences="{ value, toggleExpand, internalItem }">
      <v-chip
        v-if="value"
        inline
        :text="value"
        color="success"
        @click="toggleExpand(internalItem)"
      />
      <span v-else class="text-muted font-italic">None</span>
    </template>

    <template #expanded-row-inject="{ item }: { item: SamplingAtSite }">
      <v-list density="compact">
        <v-list-item
          v-for="occurrence in item.occurrences"
          :key="occurrence.code"
          :to="{
            name: 'occurrence-item',
            params: { code: occurrence.code }
          }"
          @click.self
          ><div class="d-flex align-center">
            <OccurrenceIcon :item="occurrence" class="mx-2" />
            <a class="font-monospace">{{ CodeIdentifier.textWrap(occurrence.code) }}</a>
          </div>
          <template #append>
            <IdentificationChip
              :identification="occurrence.identification"
              size="small"
              class="ma-1"
              @click.stop
            />
          </template>
        </v-list-item>
      </v-list>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { CodeIdentifier, DateWithPrecision, SamplingAtSite } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import OccurrenceIcon from '@/features/occurrences/components/OccurrenceIcon'
import SamplingCardDialog from '@/features/occurrences/components/sampling/SamplingCardDialog.vue'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonChip from '@/features/taxonomy/components/TaxonChip'

const { samplings } = defineProps<{ samplings: SamplingAtSite[] }>()

const headers: CRUDTableHeader<SamplingAtSite>[] = [
  { title: 'Number', key: 'number', width: 0 },
  {
    title: 'Date',
    key: 'performed_on',
    sort: DateWithPrecision.compare
  },
  { title: 'Target', key: 'target_taxa', sortable: false },
  {
    title: 'Occurrences',
    key: 'occurrences',
    value(item, fallback) {
      return item.occurrences?.length ?? fallback
    },
    align: 'center',
    width: 0
  }
]
</script>

<style scoped lang="scss"></style>
