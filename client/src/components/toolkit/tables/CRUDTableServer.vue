<template>
  <v-data-table-server
    id="table"
    class="crud-table"
    :headers="processedHeaders"
    :items="data?.items"
    :items-length="data?.total_count ?? 0"
    item-key="id"
    :loading
    :items-per-page-options="[5, 10, 15, 25, 50]"
    v-model:items-per-page="pagination.itemsPerPage"
    v-model:page="pagination.page"
    v-model:sort-by="sortBy"
  >
    <!-- Toolbar -->
    <template #top v-if="toolbar">
      <TableToolbar
        ref="toolbar"
        id="table-toolbar"
        v-model:search="genericFilters.search"
        v-bind="toolbar"
        @reload="refetch().then(() => feedback({ message: 'Data reload' }))"
      >
        <template #extension>
          <slot name="toolbar-extension" />
        </template>
        <template #[`prepend-actions`]>
          <slot name="toolbar-prepend-actions" />
        </template>
        <template #actions>
          <!-- Toggle item creation form -->
          <!-- <v-btn
            v-if="
              !!currentUser &&
              UserRole.isGranted(currentUser, 'Maintainer') &&
              hasSlotContent($slots['form'])
            "
            style="min-width: 30px"
            variant="text"
            color="primary"
            :icon="$vuetify.display.xs"
            size="small"
            @click="actions.create"
          >
            <v-tooltip v-if="$vuetify.display.xs" left activator="parent" text="New item" />
            <v-icon v-if="$vuetify.display.xs" icon="mdi-plus" size="small" />
            <span v-else>New Item</span>
          </v-btn> -->
        </template>
        <template #[`append-actions`]>
          <slot name="toolbar-append-actions" />
        </template>

        <!-- Right toolbar actions -->
        <template #append>
          <SortLastUpdatedBtn
            v-if="!toolbar?.noSort"
            sort-key="meta.last_updated"
            :sort-by
            @click="
              sortBy = [
                {
                  key: 'meta.last_updated',
                  remoteKey: 'last_updated',
                  order:
                    sortBy?.[0]?.key === 'meta.last_updated'
                      ? sortBy[0].order === 'asc'
                        ? 'desc'
                        : 'asc'
                      : 'asc'
                }
              ]
            "
          />
        </template>

        <!-- Searchbar -->
        <template #search="props">
          <slot name="search" v-bind="props" :toggleMenu :menu-open="menu">
            <CRUDTableSearchBar v-model="genericFilters.search" v-if="$vuetify.display.smAndUp" />

            <v-badge
              dot
              :color="
                Object.values(genericFilters)
                  .concat(Object.values(filters))
                  .some((v) => v !== undefined && v !== null && v !== '')
                  ? 'success'
                  : 'transparent'
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
            <v-inline-search-bar v-model="genericFilters.search" label="Search term" />
          </v-card-text>
          <slot name="menu" :toggleMenu :menuOpen="menu" />
          <v-divider />
          <v-list-item v-if="currentUser">
            <template #title>
              <v-switch
                v-model="genericFilters.owned"
                label="Owned items"
                color="primary"
                hint="Restrict the list to elements you contributed"
                persistent-hint
                class="ml-2"
                density="compact"
              />
            </template>
          </v-list-item>
          <v-divider v-if="currentUser" />
          <v-card-actions>
            <v-btn color="primary" text="OK" @click="toggleMenu(false)" />
            <v-spacer />
            <v-btn color="" text="Clear" @click="resetFilters()" />
          </v-card-actions>
        </v-card>
      </v-menu>
    </template>

    <template #body.prepend="{ columns }">
      <tr v-if="!loading && error">
        <td :colspan="columns.length">
          <v-alert color="error" icon="mdi-alert" class="my-3">Failed to retrieve items</v-alert>
        </td>
      </tr>
    </template>

    <!-- Expose VDataTable slots -->
    <template v-for="(id, index) of slotNames" #[id]="slotData" :key="index">
      <slot :name="id" v-bind="{ ...slotData }" />
      <!-- <slot :name="id" v-bind="{ ...slotData, actions }" /> -->
    </template>

    <!-- Table footer -->
    <!-- <template #[`footer.prepend`]>
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
    </template> -->

    <!-- Expanded row -->
    <template #expanded-row="{ columns, item, ...others }">
      <slot name="expanded-row" v-bind="{ columns, item, ...others }">
        <tr class="expanded">
          <td :colspan="columns.length" class="px-0">
            <div class="d-flex flex-column h-auto">
              <div class="flex-grow-1">
                <slot name="expanded-row-inject" :item> </slot>
              </div>
              <slot name="expanded-row-footer" :item>
                <div class="d-flex flex-wrap align-center">
                  <MetaChip v-if="item.meta" :meta="item.meta" class="ma-1" />
                  <!-- <v-btn
                    prepend-icon="mdi-content-copy"
                    text="UUID"
                    variant="plain"
                    size="small"
                    rounded="sm"
                    class="ma-1 text-caption font-monospace"
                    @click="copyUUID(item)"
                  /> -->
                  <v-spacer />

                  <!-- Item actions -->
                  <!-- <template
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
                  </template> -->
                </div>
              </slot>
            </div>
          </td>
        </tr>
      </slot>
    </template>
  </v-data-table-server>
</template>

<script
  setup
  lang="ts"
  generic="
    ItemType extends { id: string; meta?: Meta },
    ItemFilters extends {},
    ItemSortKey extends string
  "
>
import { ErrorModel, Meta } from '@/api'
import { PaginatedList } from '@/api/responses'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { keepPreviousData, UndefinedInitialQueryOptions, useQuery } from '@tanstack/vue-query'
import { computedAsync, promiseTimeout, useToggle } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import { computed, ComputedRef, reactive, ref, useSlots } from 'vue'
import { TableSlots, ToolbarProps, useTableSort } from '.'
import MetaChip from '../MetaChip'
import SortLastUpdatedBtn from '../ui/SortLastUpdatedBtn.vue'
import CRUDTableSearchBar from './CRUDTableSearchBar.vue'
import TableToolbar from './TableToolbar.vue'
export type SortItemRemote<ItemSortKey> = SortItem<ItemSortKey> & {
  remoteKey: string
}
type ItemsQueryData = {
  limit?: number
  offset?: number
}

type Pagination = {
  itemsPerPage: number
  page: number
}

const pagination = ref<Pagination>({
  itemsPerPage: 15,
  page: 1
})

type Filters = {
  search?: string
  owned?: boolean
}
const genericFilters = ref<Filters>({
  search: '',
  owned: undefined
})
function resetFilters() {
  genericFilters.value = {}
}

const sortBy = ref<SortItemRemote<ItemSortKey>[]>()

const slots = useSlots()
// Assert type here to prevent errors in template
const slotNames = Object.keys(slots) as 'default'[]
defineSlots<TableSlots<ItemType>>()

const props = defineProps<{
  sortKeyTransform?: (key: ItemSortKey | undefined) => ItemSortKey | undefined
  filters: ItemFilters & Filters
  headers: CRUDTableHeader<ItemType>[]
  toolbar?: ToolbarProps | false
  fetchItems: (options?: {
    query: ItemsQueryData &
      ItemFilters &
      Filters & {
        order?: 'asc' | 'desc'
        sort?: ItemSortKey
      }
  }) => UndefinedInitialQueryOptions<
    PaginatedList<ItemType>,
    ErrorModel,
    PaginatedList<ItemType>,
    any
  > & {
    queryKey: any
  }
}>()

const { feedback } = useFeedback()
const { user: currentUser } = storeToRefs(useUserStore())

const [menu, toggleMenu] = useToggle(false)

const processedHeaders = computed((): CRUDTableHeader<ItemType>[] => {
  return props.headers.filter(({ hide }) => {
    return !hide?.value
  })
}) as ComputedRef<DataTableHeader[]>

const { data, error, isPending, isFetching, refetch } = useQuery(
  computed(() => ({
    ...props.fetchItems({
      query: {
        limit: pagination.value.itemsPerPage,
        offset: (pagination.value.page - 1) * pagination.value.itemsPerPage,
        ...genericFilters.value,
        ...props.filters,
        order: sortBy.value?.[0]?.order,
        sort: props.sortKeyTransform?.(sortBy.value?.[0]?.key) ?? sortBy.value?.[0]?.key
      }
    }),
    placeholderData: keepPreviousData
  }))
)

const loading = computedAsync(async () => {
  return (
    isPending.value ||
    (isFetching.value &&
      (await promiseTimeout(500).then(() => {
        return isFetching.value
      })))
  )
}, true)
</script>

<style lang="scss">
#table table {
  position: relative;
}
#search-menu .v-overlay__content {
  left: 0px !important;
}
</style>
