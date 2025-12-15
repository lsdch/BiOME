<template>
  <CRUDTable entity-name="Sampled taxa" :items="occurrences" :headers :group-by :show-expand="false">
    <template #group-header="{ item, columns, toggleGroup, isGroupOpen, index }">
      <tr class="v-data-table-group-header-row" :style="`--v-data-table-group-header-row-depth: ${item.depth};`">
        <td :colspan="columns.length"
          class="v-data-table__td v-data-table-column--align-start v-data-table-group-header-row__column" v-ripple
          @click="toggleGroup(item)">
          <div class="d-flex align-center">
            <v-btn :icon="isGroupOpen(item) ? '$expand' : '$next'" color="primary" density="comfortable" size="small"
              variant="plain" />

            <span class="d-flex justify-space-between ms-4 w-100">
              <span>
                {{ item.value }}
                <v-badge v-if="item.depth === 0" inline :content="item.items.length" />
              </span>
              <v-chip v-if="item.depth === 0" size="small">
                {{
                  item.items.reduce((acc, curr: Group<any>) => acc + curr.items.length, 0)
                }} occurrences
              </v-chip>
              <v-badge v-else inline :content="item.items.length" />
            </span>
          </div>
        </td>
      </tr>
    </template>

    <template #item.code="{ item, value }: { item: OccurrenceAtSite; value: string }">
      <RouterLink :to="{
        name: 'occurrence-item',
        params: { code: item.code }
      }" target="_blank">
        <span class="text-wrap">{{ CodeIdentifier.textWrap(value) }}</span>
      </RouterLink>
    </template>
    <template #item.site.code="{ value }: { value: string }">
      <RouterLink :to="{
        name: 'site-item',
        params: { code: value }
      }" target="_blank">
        <span class="text-wrap font-monospace">
          {{ CodeIdentifier.textWrap(value) }}
        </span>
      </RouterLink>
    </template>
    <template #item.sampling_date="{ value }">
      <span :class="['font-monospace', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { CodeIdentifier, DateWithPrecision, OccurrenceAtSite } from '@/api';
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue';
import { OccurrenceTableItem } from '@/features/occurrences/components/tables/OccurrencesTable.vue';
import { Group } from 'vuetify/lib/components/VDataTable/composables/group.mjs';

defineProps<{
  occurrences: OccurrenceTableItem<true>[];
}>();

const headers = [
  { key: 'code', title: 'Code' },
  { key: 'site.code', title: 'Site' },
  { key: 'sampling_date', title: 'Sampling date' }
]

const groupBy = [
  { key: 'identification.taxon.rank' },
  { key: 'identification.taxon.name', order: 'asc' }
]

</script>

<style scoped lang="scss"></style>