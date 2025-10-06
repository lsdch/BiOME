<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="listDataSourcesOptions()"
    entity-name="Data source"
    :toolbar="{
      title: 'Data sources',
      icon: 'mdi-database-sync'
    }"
    append-actions
  >
    <template #item.code="{ value, item }">
      <a v-if="item.url" :href="item.url" target="_blank">{{ value }}</a>
      <template v-else>{{ value }}</template>
    </template>

    <template #expanded-row-inject="{ item }">
      <v-list>
        <v-list-item
          title="Description"
          :subtitle="item.description ?? 'No description provided'"
        ></v-list-item>
        <v-list-item title="Link template ">
          <template #subtitle>
            <code v-if="item.link_template">{{ item.link_template }}</code>
            <div v-else class="text-wrap overflow-auto">
              No template provided: direct link generation using accession numbers won't be
              available.
            </div>
          </template>
        </v-list-item>
      </v-list>
    </template>
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <DataSourceFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :model-value="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { DataSource } from '@/api'
import { listDataSourcesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import DataSourceFormDialogMutation from '@/components/forms/DataSourceFormDialogMutation.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'

const headers: CRUDTableHeader<DataSource>[] = [
  { key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } },
  { key: 'label', title: 'Label' }
]
</script>

<style scoped lang="scss"></style>
