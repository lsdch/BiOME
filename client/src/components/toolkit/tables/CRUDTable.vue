<template>
  <div>
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
      <!-- Toolbar -->
      <template v-slot:top>
        <TableToolbar
          ref="toolbar"
          v-model="items"
          v-model:search="searchTerm"
          v-bind="toolbar"
          @create-item="actions.create"
        >
          <!-- Right toolbar actions -->
          <template v-slot:append>
            <SortLastUpdatedBtn
              sort-key="meta.last_updated"
              :sort-by="sortBy"
              @click="toggleSort('meta.last_updated')"
            />
          </template>

          <!-- Searchbar -->
          <template v-slot:search>
            <slot name="search">
              <CRUDTableSearchBar v-model="searchTerm" />
            </slot>
          </template>
        </TableToolbar>
      </template>

      <!-- Actions column -->
      <template v-if="props.showActions" v-slot:[`header.actions`]>
        <v-icon title="Actions" color="secondary">mdi-cog </v-icon>
      </template>
      <template v-if="props.showActions" v-slot:[`item.actions`]="{ item }">
        <v-icon v-if="$slots.form" color="primary" icon="mdi-pencil" @click="actions.edit(item)" />
        <v-icon color="primary" icon="mdi-delete" @click="actions.delete(item)" />
      </template>

      <!-- Expose VDataTable slots -->
      <template v-for="(name, index) of slotNames" v-slot:[name]="slotData" :key="index">
        <slot :name="name" v-bind="slotData || {}" />
      </template>

      <!-- Table footer -->
      <template v-slot:[`footer.prepend`]>
        <div class="flex-grow-1">
          <v-btn
            variant="plain"
            size="small"
            prepend-icon="mdi-download"
            text="Export"
            :loading="exportDialog.loading"
            @click="exportTSV"
          />
        </div>
      </template>

      <!-- Expanded row -->
      <template v-slot:expanded-row="{ columns, item, ...others }">
        <slot name="expanded-row" v-bind="{ columns, item, ...others }">
          <tr class="expanded">
            <td :colspan="columns.length" class="px-0">
              <div class="d-flex flex-column h-auto">
                <div class="d-flex flex-wrap">
                  <ItemDateChip
                    v-if="item.meta?.created"
                    icon="created"
                    :date="item.meta.created"
                  />
                  <ItemDateChip
                    v-if="item.meta?.modified"
                    icon="updated"
                    :date="item.meta.modified"
                  />
                  <v-spacer></v-spacer>
                  <v-btn
                    prepend-icon="mdi-identifier"
                    variant="plain"
                    class="text-caption"
                    @click="copyUUID(item)"
                    :text="item.id"
                  />
                </div>
                <div class="flex-grow-1">
                  <v-divider v-show="$slots['expanded-row-inject']"></v-divider>
                  <slot name="expanded-row-inject" v-bind="{ item }"> </slot>
                </div>
              </div>
            </td>
          </tr>
        </slot>
      </template>
    </v-data-table>

    <!-- Form dialog with form slot -->
    <FormDialog v-model="formDialog" v-if="$slots.form" :mode="formMode" :entityName="entityName">
      <slot name="form"></slot>
    </FormDialog>

    <!-- Confirm item deletion dialog -->
    <ConfirmDialog v-model="deleteDialog.open" v-bind="deleteDialog.props" />

    <!-- Feedback snackbar -->
    <CRUDFeedback v-model="feedback.model" v-bind="feedback.props" />

    <!-- CSV export dialog -->
    <ExportDialog
      v-model="exportDialog.open"
      v-bind="exportDialog.props"
      @ready="exportDialog.loading = false"
    />
  </div>
</template>

<script setup lang="ts" generic="ItemType extends { id: string; meta?: Meta }">
import { CancelablePromise, Meta } from '@/api'
import { Ref, computed, onMounted, ref, useSlots } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { type VDataTable } from 'vuetify/components'
import { TableEmitEvents, TableProps, useTable } from '.'
import CRUDFeedback from '../CRUDFeedback.vue'
import ConfirmDialog from '../ConfirmDialog.vue'
import ExportDialog from '../ExportDialog.vue'
import FormDialog from '../FormDialog.vue'
import ItemDateChip from '../ItemDateChip.vue'
import SortLastUpdatedBtn from '../ui/SortLastUpdatedBtn.vue'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import TableToolbar from './TableToolbar.vue'

type Props = TableProps<ItemType, () => CancelablePromise<ItemType[]>> & {
  filter?: (item: ItemType) => boolean
  filterKeys?: string | string[]
}
type SortItem = {
  key: string
  order?: boolean | 'asc' | 'desc'
}

const slots = useSlots()
// Assert type here to prevent errors in template when exposing VDataTable slots
const slotNames = Object.keys(slots) as 'default'[]

const loading = ref(true)
const searchTerm = defineModel<string>('search')
const props = defineProps<Props>()
const emit = defineEmits<TableEmitEvents<ItemType>>()

defineSlots<
  VDataTable['$slots'] & {
    search(): any
    'expanded-row-inject': (props: { item: ItemType }) => any
    form(): any
  }
>()

const { items, actions, deleteDialog, formDialog, feedback, formMode } = useTable(props, emit)

const sortBy = ref<SortItem[]>([])

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

const exportDialog: Ref<{
  open: boolean
  loading: boolean
  props: ComponentProps<typeof ExportDialog>
}> = ref({
  open: false,
  loading: false,
  props: {
    items: [],
    namePrefix: props.entityName
  }
})

function exportTSV() {
  exportDialog.value = {
    open: true,
    loading: true,
    props: {
      items: items.value,
      namePrefix: props.entityName
    }
  }
}

async function copyUUID(item: ItemType) {
  try {
    await navigator.clipboard.writeText(item.id)
    feedback.value.show('UUID copied to clipboard', 'primary')
  } catch (err) {
    feedback.value.show('Failed to copy UUID to clipboard', 'warning')
  }
}
</script>

<style lang="less">
tr.expanded td {
  border-left: 1px solid rgb(16, 113, 176);
}

tr:has(+ tr.expanded) {
  td:first-child {
    border-left: 3px solid rgb(16, 113, 176);
  }
}
</style>
