<template>
  <div>
    <v-data-table
      v-bind="$attrs"
      :headers="processedHeaders"
      :items="filteredItems"
      :loading="loading"
      :search="searchTerm"
      :filter-keys="filterKeys"
      show-expand
      v-model="selected"
      v-model:sort-by="sortBy"
      must-sort
      fixed-header
      fixed-footer
      hover
      :items-per-page-options="[5, 10, 15, 20]"
    >
      <!-- Toolbar -->
      <template #top>
        <TableToolbar
          ref="toolbar"
          v-model:search="searchTerm"
          v-bind="toolbar"
          @reload="loadItems().then(() => feedback.show('Data reloaded'))"
        >
          <template #[`prepend-actions`]>
            <slot name="toolbar-prepend-actions" />
          </template>
          <template #actions>
            <!-- Toggle item creation form -->
            <v-btn
              style="min-width: 30px"
              variant="text"
              color="primary"
              :icon="xs"
              size="small"
              @click="actions.create"
            >
              <v-tooltip v-if="xs" left activator="parent" text="New item" />
              <v-icon v-if="xs" icon="mdi-plus" size="small" />
              <span v-else>New Item</span>
            </v-btn>
          </template>
          <template #[`append-actions`]>
            <slot name="toolbar-append-actions" />
          </template>

          <!-- Right toolbar actions -->
          <template #append>
            <SortLastUpdatedBtn
              v-if="!toolbar.noSort"
              sort-key="meta.last_updated"
              :sort-by="sortBy"
              @click="toggleSort('meta.last_updated')"
            />
            <TableFilterMenu v-if="!toolbar.noFilters" v-model="tableFilters" :user="currentUser" />
          </template>

          <!-- Searchbar -->
          <template #search>
            <slot name="search">
              <CRUDTableSearchBar v-model="searchTerm" />
            </slot>
          </template>
        </TableToolbar>
      </template>

      <!-- Actions column -->
      <template v-if="props.showActions" #[`header.actions`]>
        <v-icon title="Actions" color="secondary">mdi-cog </v-icon>
      </template>
      <template v-if="props.showActions && currentUser !== undefined" #[`item.actions`]="{ item }">
        <slot name="actions" v-bind="{ actions, show: showActions, currentUser, item }">
          <CRUDItemActions
            :show="showActions"
            :item="item"
            :actions="actions"
            :user="currentUser"
          />
        </slot>
      </template>

      <!-- Expose VDataTable slots -->
      <template v-for="(id, index) of slotNames" #[id]="slotData" :key="index">
        <slot :name="id" v-bind="slotData || {}" />
      </template>

      <!-- <template v-for="header in processedHeaders" #[`header.${header.key}`]="data">
        <slot :name="`header.${header.key}`" v-bind="data" />
      </template>

      <template v-for="header in processedHeaders" #[`item.${header.key}`]="data">
        <slot :name="`item.${header.key}`" v-bind="data" />
      </template> -->

      <!-- Table footer -->
      <template #[`footer.prepend`]>
        <div class="d-flex align-center flex-grow-1">
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
      <template #expanded-row="{ columns, item, ...others }">
        <slot name="expanded-row" v-bind="{ columns, item, ...others }">
          <tr class="expanded">
            <td :colspan="columns.length" class="px-0">
              <div class="d-flex flex-column h-auto">
                <div class="flex-grow-1">
                  <slot name="expanded-row-inject" v-bind="{ item }"> </slot>
                  <v-divider v-show="$slots['expanded-row-inject']" />
                </div>
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
                  <v-spacer />
                  <v-btn
                    prepend-icon="mdi-identifier"
                    variant="plain"
                    class="text-caption"
                    @click="copyUUID(item)"
                    :text="item.id"
                  />
                </div>
              </div>
            </td>
          </tr>
        </slot>
      </template>
    </v-data-table>

    <slot name="form" v-bind="form" />

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
import { Meta } from '@/api'
import { useClipboard } from '@vueuse/core'
import { Ref, UnwrapRef, computed, ref, useSlots } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { type VDataTable } from 'vuetify/components'
import { TableProps, useTable } from '.'
import CRUDFeedback from '../CRUDFeedback.vue'
import ExportDialog from '../ExportDialog.vue'
import ItemDateChip from '../ItemDateChip.vue'
import { isOwner } from '../meta'
import SortLastUpdatedBtn from '../ui/SortLastUpdatedBtn.vue'
import CRUDItemActions from './CRUDItemActions.vue'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import TableFilterMenu from './TableFilterMenu.vue'
import TableToolbar from './TableToolbar.vue'
import { useDisplay } from 'vuetify'

type Props = TableProps<ItemType> & {
  filter?: (item: ItemType) => boolean
  filterKeys?: string | string[]
}

const { xs } = useDisplay()

const slots = useSlots()
// Assert type here to prevent errors in template when exposing VDataTable slots
const slotNames = Object.keys(slots) as 'default'[]

const items = defineModel<ItemType[]>({ default: [] })
const selected = defineModel<string[]>('selected', { default: [] })
const searchTerm = defineModel<string>('search')
const props = defineProps<Props>()

const { currentUser, actions, feedback, form, processedHeaders, loading, loadItems } = useTable(
  items,
  props
)

defineSlots<
  VDataTable['$slots'] & {
    actions(bind: {
      actions: typeof actions
      show: typeof props.showActions
      item: ItemType
      currentUser: typeof currentUser
    }): any
    search(): any
    'expanded-row-inject': (props: { item: ItemType }) => any
    'toolbar-prepend-actions': () => any
    'toolbar-append-actions': () => any
    form(props: UnwrapRef<typeof form>): any
  }
>()

const sortBy = ref<SortItem[]>([])

function toggleSort(sortKey: string) {
  const sortMeta = sortBy.value?.find(({ key }) => key === sortKey)
  let order: 'desc' | 'asc' = 'asc'
  if (sortMeta?.order === 'asc') {
    order = 'desc'
  }
  sortBy.value?.splice(0, sortBy.value.length)
  sortBy.value?.push({ key: sortKey, order })
}

const tableFilters = ref({
  ownedItems: false
})

function ownedItemFilter(item: ItemType) {
  return tableFilters.value.ownedItems && currentUser !== undefined
    ? isOwner(currentUser, item)
    : true
}

const filteredItems = computed(() => {
  return props.filter || tableFilters.value.ownedItems
    ? items.value.filter((item) => {
        return (props.filter ? props.filter(item) : true) && ownedItemFilter(item)
      })
    : items.value
})

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

const { copy } = useClipboard()
async function copyUUID(item: ItemType) {
  if (item.id === undefined) return
  try {
    await copy(item.id)
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
