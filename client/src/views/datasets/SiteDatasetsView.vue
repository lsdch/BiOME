<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="DatasetsService.listSiteDatasets"
    entity-name="Site dataset"
    :toolbar="{ title: 'Site datasets', icon: 'mdi-map-marker-multiple' }"
  >
    <template #item.label="{ item }: { item: SiteDataset }">
      <RouterLink
        :to="{ name: 'site-dataset-item', params: { slug: item.slug } }"
        :text="item.label"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { DatasetsService, SiteDataset } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<SiteDataset>[] = [
  { key: 'label', title: 'Label' },
  { key: 'description', title: 'Description' },
  {
    key: 'sites',
    title: 'Sites',
    width: 0,
    value(item, fallback) {
      return item.sites.length
    }
  }
]
</script>

<style scoped></style>
