<template>
  <CRUDTable
    ref="table"
    class="fill-height"
    :headers
    :fetch-items="listSiteDatasetsOptions"
    entity-name="Site dataset"
    :toolbar="{ title: 'Site datasets', icon: 'mdi-map-marker-multiple' }"
  >
    <template #item.label="{ item }: { item: SiteDataset }">
      <RouterLink
        :to="{ name: 'site-dataset-item', params: { slug: item.slug } }"
        class="font-weight-bold"
        :text="item.label"
      />
    </template>
    <template #item.description="{ value }">
      <LineClampedText :title="value" :text="value" :lines="3" />
    </template>

    <template #header.pin="props">
      <IconTableHeader icon="mdi-pin" v-bind="props" />
    </template>
    <template #item.pin="{ item, index }: { item: SiteDataset }">
      <DatasetPinButton
        :model-value="item"
        @update:model-value="
          ({ pinned }) => {
            table?.updateItem<SiteDataset>(index, { ...item, pinned })
          }
        "
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { SiteDataset } from '@/api'
import { listSiteDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { LineClampedText } from '@/components/toolkit/ui/LineClampedText'
import DatasetPinButton from './DatasetPinButton.vue'
import { useTemplateRef } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { useUserStore } from '@/stores/user'
import IconTableHeader from '@/components/toolkit/tables/IconTableHeader.vue'

const table = useTemplateRef<ComponentExposed<typeof CRUDTable>>('table')

const { isGranted } = useUserStore()

const headers: CRUDTableHeader<SiteDataset>[] = [
  { key: 'label', title: 'Label' },
  {
    key: 'description',
    title: 'Description'
  },
  {
    key: 'sites',
    title: 'Sites',
    width: 0,
    value(item, fallback) {
      return item.sites.length
    }
  },
  ...(isGranted('Admin')
    ? [
        {
          key: 'pin',
          title: 'Pin status',
          width: 0
        }
      ]
    : [])
]
</script>

<style scoped></style>
