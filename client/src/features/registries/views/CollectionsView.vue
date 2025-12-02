<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="listCollectionsOptions()"
    :delete="{
      mutation: deleteCollectionMutation,
      params: ({ code }: Collection) => ({ path: { code } })
    }"
    entity-name="Collection"
    :toolbar="{ icon: 'mdi-newspaper-variant-multiple', title: 'Bibliography' }"
    append-actions
    v-model:search="search"
  >
    <template #item.personal="{ value }: { value: boolean }">
      <v-icon v-if="value" icon="mdi-account"></v-icon>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { type Collection } from '@/api'
import { deleteCollectionMutation, listCollectionsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { ref } from 'vue'

type CollectionFilters = {
  term?: string
}

const search = ref<CollectionFilters>({})

const headers: CRUDTableHeader<Collection>[] = [
  { key: 'label', title: 'Label' },
  { key: 'code', title: 'Code', width: 0, cellProps: { class: 'font-monospace' } },
  { key: 'location', title: 'Location' },
  { key: 'personal', title: 'Personal', width: 0, align: 'center' }
]
</script>

<style lang="scss">
.article-details .v-card-title {
  font-size: 1rem;
}
</style>
