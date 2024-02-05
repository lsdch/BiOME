<template>
  <v-data-table
    v-bind="$attrs"
    :headers="headers"
    :items="filteredItems"
    :loading="loading"
    :search="searchTerm"
    :filter-keys="filterKeys"
    show-expand
    v-model:sort-by="sortBy"
    must-sort
    :items-per-page-options="[5, 10, 20]"
  >
    <template v-slot:top>
      <TableToolbar
        ref="toolbar"
        v-model="items"
        v-model:search="searchTerm"
        v-bind="toolbarProps"
        :delete-request="crud.delete"
      >
        <template v-slot:append>
          <v-tooltip bottom activator="parent" :open-delay="400">
            Sort by last update
            <template v-slot:activator="{ props }">
              <v-btn
                variant="plain"
                v-bind="{ ...sortLastUpdateIcon, ...props }"
                @click="() => toggleSort('meta.last_updated')"
              />
            </template>
          </v-tooltip>
        </template>
        <template v-slot:search>
          <slot name="search">
            <CRUDTableSearchBar v-model="searchTerm" />
          </slot>
        </template>
      </TableToolbar>
    </template>
    <template v-if="props.showActions" v-slot:[`item.actions`]="{ item }">
      <v-icon color="primary" icon="mdi-pencil" @click="toolbar?.editItem(item)" />
      <v-icon color="primary" icon="mdi-delete" @click="toolbar?.deleteItem(item)" />
    </template>
    <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
      <slot :name="name" v-bind="slotData || {}" />
    </template>
    <template v-slot:[`footer.prepend`]>
      <div class="flex-grow-1">
        <v-btn variant="plain" size="small" prepend-icon="mdi-download">Export</v-btn>
      </div>
    </template>
    <template v-slot:expanded-row="{ columns, item, ...others }">
      <slot name="expanded-row" v-bind="{ columns, item, ...others }">
        <tr>
          <td :colspan="columns.length">
            <div class="d-flex">
              <div class="d-flex flex-column flex-grow-0 mr-3">
                <ItemDateChip v-if="item.meta?.created" icon="created" :date="item.meta.created" />
                <ItemDateChip
                  v-if="item.meta?.modified"
                  icon="updated"
                  :date="item.meta.modified"
                />
              </div>
              <v-divider vertical />
              <slot name="expanded-row-inject" v-bind="{ item }"> </slot>
            </div>
          </td>
        </tr>
      </slot>
    </template>
  </v-data-table>
</template>

<script
  setup
  lang="ts"
  generic="ItemInputType extends {}, ItemType extends { id: string; meta?: Meta }"
>
import { CancelablePromise, Meta } from '@/api'
import { Ref, computed, onMounted, ref, useSlots } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { type VDataTable } from 'vuetify/components'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import ItemDateChip from './ItemDateChip.vue'
import TableToolbar from './TableToolbar.vue'
import { TableProps } from './table'

type Props = TableProps<ItemType, ItemInputType, () => CancelablePromise<ItemType[]>> & {
  filter?: (item: ItemType) => boolean
  filterKeys?: string | string[]
}
type SortItem = {
  key: string
  order?: boolean | 'asc' | 'desc'
}

const slots = useSlots()
// Assert type here to prevent errors in template
const slotNames = Object.keys(slots) as 'default'[]

const toolbar = ref<ComponentExposed<typeof TableToolbar<ItemType>> | null>(null)

const loading = ref(true)
const searchTerm = defineModel<string>('search')
const items: Ref<ItemType[]> = ref([])

const sortBy = ref<SortItem[]>([])
const sortLastUpdateIcon = computed(() => {
  const sortMeta = sortBy.value.find(({ key }) => key === 'meta.last_updated')
  if (sortMeta) {
    return sortMeta.order === 'asc'
      ? {
          icon: 'mdi-sort-calendar-ascending',
          color: 'primary'
        }
      : { icon: 'mdi-sort-calendar-descending', color: 'primary' }
  }
  return { icon: 'mdi-sort-calendar-ascending', color: 'secondary' }
})

const props = defineProps<Props>()

defineSlots<
  VDataTable['$slots'] & {
    search(): any
    'expanded-row-inject': (props: { item: ItemType }) => any
  }
>()

function toggleSort(sortKey: string) {
  const sortMeta = sortBy.value?.find(({ key }) => key === sortKey)
  let order: 'desc' | 'asc' = 'asc'
  if (sortMeta?.order === 'asc') {
    order = 'desc'
  }
  sortBy.value.splice(0, sortBy.value.length)
  sortBy.value.push({ key: sortKey, order })
}

const filteredItems = computed(() => {
  return props.filter ? items.value.filter(props.filter) : items.value
})

async function loadItems() {
  loading.value = true
  items.value = await props.crud.list()
  loading.value = false
}

onMounted(loadItems)
</script>

<style scoped></style>
