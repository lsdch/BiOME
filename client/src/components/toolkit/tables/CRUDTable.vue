<template>
  <div>
    <v-data-table
      id="table"
      :headers="processedHeaders"
      :items="filteredItems"
      :loading="loading"
      :search="search.term"
      :filter-keys="filterKeys"
      v-model="selected"
      v-model:sort-by="sortBy"
      v-bind="$attrs"
      show-expand
      must-sort
      fixed-header
      fixed-footer
      hover
      :mobile="mobile ?? xs"
      :density="mobile ? 'compact' : undefined"
      :items-per-page-options="[5, 10, 15, 20]"
      style="position: relative"
    >
      <!-- Toolbar -->
      <template #top v-if="toolbar">
        <TableToolbar
          ref="toolbar"
          id="table-toolbar"
          v-model:search="search.term"
          v-bind="toolbar"
          @reload="loadItems().then(() => feedback.show('Data reloaded'))"
        >
          <template #extension>
            <slot name="toolbar-extension" />
          </template>
          <template #[`prepend-actions`]>
            <slot name="toolbar-prepend-actions" />
          </template>
          <template #actions>
            <!-- Toggle item creation form -->
            <v-btn
              v-if="
                !!currentUser &&
                UserRole.isGranted(currentUser, 'Maintainer') &&
                hasSlotContent($slots['form'])
              "
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
              v-if="!toolbar?.noSort"
              sort-key="meta.last_updated"
              :sort-by="sortBy"
              @click="toggleSort('meta.last_updated')"
            />
          </template>

          <!-- Searchbar -->
          <template #search="props">
            <slot name="search" v-bind="props" :toggleMenu :menu-open="menu">
              <CRUDTableSearchBar v-model="search.term" v-if="smAndUp" />

              <v-badge
                dot
                :color="
                  Object.values(search).some((v) => v !== undefined) ? 'success' : 'transparent'
                "
                class="mx-1"
              >
                <v-btn
                  color="primary"
                  variant="tonal"
                  icon="mdi-menu"
                  @click="toggleMenu(true)"
                  :active="menu"
                  size="small"
                />
              </v-badge>
            </slot>
          </template>
        </TableToolbar>
      </template>

      <template #body.prepend="{ columns }">
        <tr v-if="!loading && error">
          <td :colspan="columns.length">
            <v-alert color="error" icon="mdi-alert" class="my-3">Failed to retrieve items</v-alert>
          </td>
        </tr>

        <v-menu
          id="search-menu"
          v-model="menu"
          location="bottom"
          target="#table-toolbar"
          attach="#table table"
          :close-on-content-click="false"
        >
          <v-card rounded="t-0">
            <v-card-text>
              <v-inline-search-bar v-model="search.term" label="Search term" />
            </v-card-text>
            <slot name="menu" :toggleMenu :menuOpen="menu"> </slot>
            <v-divider> </v-divider>
            <v-list-item>
              <template #title>
                <v-switch
                  v-model="search.owned"
                  label="Owned items"
                  color="primary"
                  hint="Restrict the list to elements you contributed"
                  persistent-hint
                  class="ml-2"
                  density="compact"
                />
              </template>
            </v-list-item>
            <v-divider></v-divider>
            <v-card-actions>
              <v-btn color="primary" text="OK" @click="toggleMenu(false)"></v-btn>
              <v-spacer></v-spacer>
              <v-btn color="" text="Clear" @click="search = {}"></v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>
      </template>

      <!-- Expose VDataTable slots -->
      <template v-for="(id, index) of slotNames" #[id]="slotData" :key="index">
        <slot :name="id" v-bind="{ ...slotData, actions }" />
      </template>

      <!-- <template v-for="header in processedHeaders" #[`header.${header.key}`]="data">
        <slot :name="`header.${header.key}`" v-bind="data" />
      </template> -->

      <!-- <template v-for="header in processedHeaders" #[`item.${header.key}`]="data">
        <slot :name="`item.${header.key}`" v-bind="{ ...data, actions }" />
      </template> -->

      <!-- Table footer -->
      <template #[`footer.prepend`]>
        <div class="d-flex align-center flex-grow-1">
          <slot name="footer.prepend-actions"></slot>
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
                  <slot name="expanded-row-inject" :item> </slot>
                  <v-divider v-show="$slots['expanded-row-inject']" />
                </div>
                <slot name="expanded-row-footer" :item>
                  <div class="d-flex flex-wrap align-center">
                    <MetaChip v-if="item.meta" :meta="item.meta" class="ma-1" />
                    <v-btn
                      prepend-icon="mdi-content-copy"
                      text="UUID"
                      variant="plain"
                      size="small"
                      rounded="sm"
                      class="ma-1 text-caption font-monospace"
                      @click="copyUUID(item)"
                    />
                    <v-spacer />

                    <!-- Item actions -->
                    <template
                      v-if="
                        !!currentUser &&
                        (UserRole.isGranted(currentUser, 'Maintainer') ||
                          User.isOwner(currentUser, item))
                      "
                    >
                      <v-btn
                        v-if="actions.edit"
                        text="Edit"
                        color="primary"
                        variant="tonal"
                        size="small"
                        class="ma-1"
                        prepend-icon="mdi-pencil"
                        @click="actions.edit(item)"
                      />
                      <v-btn
                        v-if="actions.delete"
                        text="Delete"
                        color="error"
                        variant="tonal"
                        size="small"
                        class="ma-1"
                        prepend-icon="mdi-delete"
                        @click="actions.delete(item)"
                      />
                    </template>
                  </div>
                </slot>
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

<script
  setup
  lang="ts"
  generic="
    ItemType extends { id: string; meta?: Meta },
    ItemsQueryData extends {},
    ItemsDeleteData extends {},
    Filters extends { owned?: boolean; term?: string }
  "
>
import { Meta, User, UserRole } from '@/api'
import { useArrayFilter, useClipboard, useToggle } from '@vueuse/core'
import { Ref, UnwrapRef, reactive, ref, useSlots } from 'vue'
import { ComponentProps } from 'vue-component-type-helpers'
import { useDisplay } from 'vuetify'
import { VDataTable } from 'vuetify/components'
import { TableProps, useTable } from '.'
import CRUDFeedback from '../CRUDFeedback.vue'
import ExportDialog from '../ExportDialog.vue'
import MetaChip from '../MetaChip.vue'
import SortLastUpdatedBtn from '../ui/SortLastUpdatedBtn.vue'
import { hasSlotContent } from '../vue-utils'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import TableToolbar from './TableToolbar.vue'

type Props = TableProps<ItemType, ItemsQueryData, ItemsDeleteData> & {
  filter?: (item: ItemType) => boolean
  filterKeys?: string | string[]
  mobile?: boolean
}

const { xs, smAndUp } = useDisplay()

const tableSlots = useSlots()
// Assert type here to prevent errors in template when exposing VDataTable slots
const slotNames = Object.keys(tableSlots) as 'default'[]

const items = defineModel<ItemType[]>('items', { default: reactive([]) })
const selected = defineModel<string[]>('selected', { default: [] })
const search = defineModel<Partial<Filters>>('search', { default: {} })
const props = withDefaults(defineProps<Props>(), {})
const emit = defineEmits<{
  itemCreated: [item: ItemType, index: number]
  itemEdited: [item: ItemType, index: number]
}>()

const { currentUser, actions, feedback, form, processedHeaders, loading, loadItems, error } =
  useTable(items, props, emit)

const [menu, toggleMenu] = useToggle(false)

defineExpose({ form, actions })

const slots = defineSlots<
  VDataTable['$slots'] & {
    actions(bind: { actions: typeof actions; item: ItemType; currentUser: typeof currentUser }): any
    search(props: { toggleMenu: typeof toggleMenu; menuOpen: boolean }): any
    'toolbar-extension': () => any
    menu: (props: { toggleMenu: typeof toggleMenu; menuOpen: boolean }) => any
    'expanded-row-inject': (props: { item: ItemType }) => any
    'expanded-row-footer': (props: { item: ItemType }) => any
    'toolbar-prepend-actions': () => any
    'toolbar-append-actions': () => any
    'footer.prepend-actions': () => any
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

function ownedItemFilter(item: ItemType) {
  return search.value.owned && currentUser !== undefined ? User.isOwner(currentUser, item) : true
}

const filteredItems = useArrayFilter(items as Ref<ItemType[]>, (item) => {
  return (props.filter ? props.filter(item) : true) && ownedItemFilter(item)
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
      items: items.value ?? [],
      namePrefix: props.entityName
    }
  }
}

const { copy } = useClipboard()
async function copyUUID(item: ItemType) {
  if (item.id === undefined) return
  try {
    await copy(item.id)
    feedback.value.show(`UUID copied to clipboard\n${item.id}`, 'primary')
  } catch (err) {
    feedback.value.show('Failed to copy UUID to clipboard', 'warning')
  }
}
</script>

<style lang="less">
#search-menu .v-overlay__content {
  left: 0px !important;
}

tr.expanded td {
  border-left: 1px solid rgb(16, 113, 176);
}

tr:has(+ tr.expanded) {
  td:first-child {
    border-left: 3px solid rgb(16, 113, 176);
  }
}
</style>
