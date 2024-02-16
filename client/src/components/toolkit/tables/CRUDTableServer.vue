<template>
  <v-data-table-server
    :headers="headers"
    :items="items"
    :items-length="totalItems"
    item-key="id"
    search="search"
    :loading="loading"
    :items-per-page="itemsPerPage"
    @update:options="loadItems"
  >
    <template v-slot:top>
      <TableToolbar
        ref="toolbar"
        v-model="items"
        :delete-request="crud.delete"
        v-bind="toolbarProps"
      >
        <template v-slot:search>
          <slot name="search"></slot>
        </template>
      </TableToolbar>
    </template>
    <template v-if="props.showActions" v-slot:[`item.actions`]="{ item }">
      <v-icon size="small" color="primary" icon="mdi-pencil" @click="toolbar?.editItem(item)" />
      <v-icon size="small" color="primary" icon="mdi-delete" @click="toolbar?.deleteItem(item)" />
    </template>
    <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
      <slot :name="name" v-bind="slotData || {}" />
    </template>
  </v-data-table-server>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }, ItemInputType, FilterType">
import { CancelablePromise } from '@/api'
import { Ref, ref, useSlots } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { VDataTable } from 'vuetify/components'
import TableToolbar from './TableToolbar.vue'
import { TableProps } from '.'

const slots = useSlots()
// Assert type here to prevent errors in template
const slotNames = Object.keys(slots) as 'default'[]
defineSlots<VDataTable['$slots'] & { search(): any }>()

const toolbar = ref<ComponentExposed<typeof TableToolbar<ItemType>> | null>(null)

type ItemsQuery = { page: number; itemsPerPage: number; sortBy: string; filters?: FilterType }

type FetchList = (query: ItemsQuery) => CancelablePromise<{
  items: ItemType[]
  totalItems: number
}>

type Props = TableProps<ItemType, ItemInputType, FetchList>
const props = defineProps<Props>()

const itemsPerPage = ref(10)
const items: Ref<ItemType[]> = ref([])
const totalItems = ref(0)
const loading = ref(true)

async function loadItems(query: ItemsQuery) {
  loading.value = true
  const response = await props.crud.list(query)
  items.value = response.items
  totalItems.value = response.totalItems
  loading.value = false
}
</script>

<style scoped></style>
