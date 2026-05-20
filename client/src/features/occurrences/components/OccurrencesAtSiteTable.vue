<template>
  <div v-if="items?.length" class="d-flex align-center ga-2">
    <v-text-field
      v-model="search.term"
      class="mx-5"
      hide-details
      label="Search"
      clearable
      density="compact"
    />
  </div>
  <CRUDTable :items :headers entityName="Occurrence" :search filter-mode="some">
    <template #item.code="{ item, value }: { item: OccurrenceAtSite; value: string }">
      <div class="d-flex justify-space-between align-center">
        <RouterLink
          :to="{
            name: 'occurrence-item',
            params: { code: item.code }
          }"
        >
          <span class="text-wrap font-monospace">{{ CodeIdentifier.textWrap(value) }}</span>
        </RouterLink>
      </div>
    </template>
    <template #item.sampling.performed_on="{ value, item }">
      <div class="d-flex flex-column">
        <span :class="['font-monospace', { 'text-muted': !value }]">
          {{ DateWithPrecision.format(value) }}
        </span>
        <span class="text-caption text-muted font-monospace"> #{{ item.sampling.number }} </span>
      </div>
    </template>
    <template #item.identification="{ value: identification }">
      <IdentificationChip :identification size="small" />
    </template>
  </CRUDTable>
</template>

<script setup lang="tsx">
import {
  CodeIdentifier,
  DateWithPrecision,
  Identification,
  OccurrenceAtSite,
  SamplingAtSite
} from '@/api/adapters'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import { computed, ref } from 'vue'

const { samplings } = defineProps<{ samplings: SamplingAtSite[] }>()

const search = ref({
  term: undefined,
  owned: undefined
})

type OccurrenceTableItem = {
  sampling: SamplingAtSite
} & OccurrenceAtSite

const items = computed(
  () =>
    samplings.reduce<OccurrenceTableItem[]>((acc, { occurrences, ...s }) => {
      occurrences?.forEach((o) => {
        acc.push({ sampling: s, ...o })
      })
      return acc
    }, []) ?? []
)

const headers: CRUDTableHeader<OccurrenceTableItem>[] = [
  { title: 'Code', key: 'code' },
  {
    title: 'Sampling',
    key: 'sampling.performed_on',
    sort: DateWithPrecision.compare
  },
  Identification.tableHeader({ key: 'identification' })
]
</script>

<style scoped></style>
