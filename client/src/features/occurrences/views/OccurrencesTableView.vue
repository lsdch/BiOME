<template>
  <v-data-table-server
    id="table"
    class="crud-table fill-height"
    :headers="headers"
    :items="data?.items"
    :items-length="data?.total_count ?? 0"
    item-key="id"
    :loading
    :items-per-page-options="[5, 10, 15, 25, 50]"
    v-model:items-per-page="pagination.itemsPerPage"
    v-model:page="pagination.page"
    @update:sort-by="updateSortBy"
    @update:page="prefetchNext"
  >
    <!-- Toolbar -->
    <template #top>
      <TableToolbar
        title="Occurrences"
        icon="mdi-package-variant"
        ref="toolbar"
        id="table-toolbar"
        v-model:search="filters.search_term"
        @reload="reload()"
      >
        <template #extension>
          <slot name="toolbar-extension" />
        </template>
        <template #[`prepend-actions`]>
          <slot name="toolbar-prepend-actions" />
        </template>
        <template #[`append-actions`]>
          <slot name="toolbar-append-actions" />
        </template>

        <!-- Right toolbar actions -->
        <!-- <template #append>
          <SortLastUpdatedBtn v-if="!toolbar?.noSort" sort-key="meta.last_updated" :sort-by @click="
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
            " />
        </template> -->

        <!-- Searchbar -->
        <template #search="props">
          <slot name="search" v-bind="props" :toggleMenu :menu-open="menu">
            <CRUDTableSearchBar v-model="filters.search_term" v-if="$vuetify.display.smAndUp" />

            <v-badge
              dot
              :color="
                Object.values(filters)
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
            <v-inline-search-bar v-model="filters.search_term" label="Search term" />
          </v-card-text>
          <slot name="menu" :toggleMenu :menuOpen="menu">
            <v-row class="ma-0">
              <v-col cols="12" md="6">
                <v-list>
                  <v-list-item prepend-icon="mdi-folder-table">
                    <DatasetPicker
                      v-model="filters.datasets"
                      label="Datasets"
                      class="mt-2"
                      item-value="label"
                      clearable
                      multiple
                      chips
                      closable-chips
                      density="compact"
                      hide-details
                    />
                  </v-list-item>
                  <v-list-item prepend-icon="mdi-file-table">
                    <ImportBatchPicker
                      v-model="filters.batches"
                      label="Import batches"
                      class="mt-2"
                      item-value="id"
                      clearable
                      multiple
                      chips
                      closable-chips
                      density="compact"
                      hide-details
                    />
                  </v-list-item>
                  <DateFiltersListItem v-model="filters.date" />
                  <!-- <v-list-item prepend-icon="mdi-dna">
                    <ClearableSwitch v-model="filters.has_sequences" class="pl-2" label="Sequences available"
                      color-true="primary" color-false="red" hint="Show only bio material having registered sequences"
                      persistent-hint density="compact" />
                  </v-list-item> -->
                  <!-- <v-list-item prepend-icon="mdi-calendar" v-if="yearRange?.min && yearRange?.max">
              <div class="d-flex align-start mt-5">
                <v-checkbox class="" density="compact" :model-value="useYearRange" @update:model-value="toggleYearRange"
                  hide-details color="primary"></v-checkbox>
                <v-range-slider :disabled="!useYearRange" class="mx-5 mt-1" :min="yearRange.min" :max="yearRange.max"
                  :step="1" :width="400" thumb-label="always" clearable
                  :model-value="[filters.year ?? yearRange.min, filters.year_end ?? yearRange.max]" @update:model-value="
                    ([start, end]: [number, number]) => {
                      filters.year = start
                      filters.year_end = end
                    }
                  ">
                </v-range-slider>
              </div>
              <template #append>
                <InlineHelp class="mr-4"
                  text="Filters occurrences by sampling date. If sampling date is not available, identification date is used as fallback. If neither sampling nor identification dates are available, the earliest bibliographic reference year is used." />
              </template>
  </v-list-item> -->
                </v-list>
              </v-col>
              <v-col cols="12" md="6">
                <v-list density="compact">
                  <v-list-item prepend-icon="mdi-family-tree">
                    <TaxonFilterPicker
                      v-model:taxa="filters.taxa"
                      v-model:whole-clade="filters.whole_clade"
                      item-value="id"
                      label="Assigned taxon"
                      density="compact"
                      class="mt-1"
                      clearable
                      multiple
                      hide-details
                      chips
                      closable-chips
                    />
                    <ClearableSwitch
                      v-model="filters.confer"
                      class="pl-2"
                      label="Confer (cf.)"
                      base-color="grey"
                      color-true="primary"
                      color-false="red"
                      :hint="
                        filters.confer
                          ? 'Show only occurrences with a confer identification'
                          : filters.confer !== undefined
                            ? 'Show only occurrences without a confer identification'
                            : undefined
                      "
                      density="compact"
                    />
                    <div class="d-flex align-center ga-3 flex-wrap">
                      <v-select
                        v-model="filters.rank"
                        :items="$TaxonRank.enum"
                        label="Rank"
                        density="compact"
                        hide-details
                        clearable
                        :min-width="200"
                      />
                      <v-select
                        v-model="filters.status"
                        :items="$TaxonStatus.enum"
                        label="Status"
                        density="compact"
                        hide-details
                        clearable
                        :min-width="200"
                      />
                    </div>
                  </v-list-item>

                  <v-list-item prepend-icon="mdi-star-four-points">
                    <TypeStatusPicker
                      v-model="filters.type_status"
                      label="Type status"
                      clearable
                      :multiple="false"
                      density="compact"
                      hide-details
                    />
                  </v-list-item>
                </v-list>
              </v-col>
            </v-row>
          </slot>
          <v-divider />
          <v-list-item v-if="currentUser">
            <!-- <template #title>
              <v-switch v-model="genericFilters.owned" label="Owned items" color="primary"
                hint="Restrict the list to elements you contributed" persistent-hint class="ml-2" density="compact" />
            </template> -->
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
    <!-- <template v-for="(id, index) of slotNames" #[id]="slotData" :key="index">
      <slot :name="id" v-bind="{ ...slotData }" />
    </template> -->
    <!-- <slot :name="id" v-bind="{ ...slotData, actions }" /> -->

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
                  <!-- <MetaChip v-if="item.meta" :meta="item.meta" class="ma-1" /> -->
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
                </div>
              </slot>
            </div>
          </td>
        </tr>
      </slot>
    </template>

    <template #item.code="{ value, item }: { value: string; item: Occurrence }">
      <span class="d-flex justify-space-between align-center">
        <RouterLink
          :text="value"
          :to="{ name: 'occurrence-item', params: { id: item.id, code: value } }"
        />
        <span class="d-flex align-center ga-2 justify-end">
          <v-icon
            v-if="item.type_status"
            icon="mdi-star-four-points"
            size="small"
            v-tooltip="item.type_status"
            density="compact"
          />
          <!-- <v-icon v-if="item.has_sequences" size="small" icon="mdi-dna" v-tooltip="`Sequence(s) available`" /> -->
        </span>
      </span>
    </template>
    <template #item.sampling.site="{ value: { code, name, country } }: { value: Site }">
      <div class="d-flex justify-space-between">
        <span class="font-size-small">{{ name }}</span>
        <CountryChip v-if="country" :country size="small" class="flex-shrink-0" />
      </div>
    </template>
    <template #item.sampling.performed_on="{ value }: { value?: DateWithPrecision }">
      <span :class="['font-monospace text-caption', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>

    <template #item.identification="{ value: identification }: { value: Identification }">
      <IdentificationChip :identification size="small" short />
    </template>
    <template #item.identification.identified_by="{ value: identificators }: { value?: string[] }">
      <div v-if="identificators?.length" class="d-flex flex-wrap ga-1">
        <v-chip v-for="name in identificators" :text="name" size="small"></v-chip>
      </div>
      <span v-else class="text-muted text-caption">Unknown</span>
    </template>
    <template #item.identification.identified_on="{ value }: { value?: DateWithPrecision }">
      <span :class="['font-monospace text-caption', { 'text-muted': !value }]">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
  </v-data-table-server>
</template>

<script setup lang="ts">
import {
  $TaxonRank,
  $TaxonStatus,
  CompositeDate,
  EventDatePrecision,
  ListOccurrencesData,
  OccurrenceSortKey,
  OccurrenceTypeStatus,
  Site,
  TaxonRank
} from '@/api'

import { DateWithPrecision, Identification, Occurrence, TaxonStatus } from '@/api'
import { listOccurrencesOptions, listOccurrencesQueryKey } from '@/api/gen/@tanstack/vue-query.gen'
import CRUDTableSearchBar from '@/components/toolkit/tables/CRUDTableSearchBar.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import ClearableSwitch from '@/components/toolkit/ui/ClearableSwitch.vue'
import DateFiltersListItem from '@/components/toolkit/ui/DateFiltersListItem.vue'
import {
  DateFilters,
  mappingFiltersToQuery
} from '@/features/cartography/components/layers-manager/map-layers'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import ImportBatchPicker from '@/features/import/components/ImportBatchPicker.vue'
import TypeStatusPicker from '@/features/occurrences/components/TypeStatusPicker'
import CountryChip from '@/features/site/components/CountryChip'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { keepPreviousData, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computedAsync, promiseTimeout, useToggle, useUrlSearchParams } from '@vueuse/core'
import { storeToRefs } from 'pinia'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { computed, onMounted, ref, toRef, watch } from 'vue'
import { DataTableSortItem, FilterMatch } from 'vuetify'

const { feedback } = useFeedback()
const { user: currentUser } = storeToRefs(useUserStore())

const [menu, toggleMenu] = useToggle(false)

const sortBy = ref<Overwrite<DataTableSortItem, { key: OccurrenceSortKey }>>({
  key: 'code',
  order: 'asc'
})

const sortKeys: Record<string, OccurrenceSortKey> = {
  code: 'code',
  'sampling.site': 'site_name',
  'sampling.performed_on': 'event_date',
  identification: 'taxon_name',
  'identification.identified_on': 'identified_on'
}

function updateSortBy(newSortBy: DataTableSortItem[]) {
  if (newSortBy.length) {
    sortBy.value = { ...newSortBy[0], key: sortKeys[newSortBy[0].key] as OccurrenceSortKey }
  } else {
    sortBy.value = { key: 'code', order: 'asc' }
  }
}

type Pagination = {
  itemsPerPage: number
  page: number
}

const pagination = ref<Pagination>({
  itemsPerPage: 15,
  page: 1
})

type TableUrlParams = {
  search_term?: string
  datasets?: string[]
  batches?: UUID[]
  type_status?: OccurrenceTypeStatus
  confer?: boolean
  whole_clade?: boolean
  rank?: TaxonRank
  status?: TaxonStatus
  taxa?: string[]
  date_from?: string
  date_to?: string
  date_is_range?: boolean
  date_precision?: EventDatePrecision
  date_include_unknown?: boolean
  date_buffer?: string
}

const urlParams = toRef(
  useUrlSearchParams<TableUrlParams>('history', {
    removeNullishValues: true,
    initialValue: {}
  })
)

type TableFilters = {
  search_term?: string
  datasets?: string[]
  batches?: UUID[]
  type_status?: OccurrenceTypeStatus
  confer?: boolean
  whole_clade?: boolean
  rank?: TaxonRank
  status?: TaxonStatus
  taxa?: string[]
  date: DateFilters
}

const filters = ref<TableFilters>({ date: { enabled: false, is_range: false, precision: 'year' } })

onMounted(() => {
  const {
    date_from,
    date_to,
    date_is_range,
    date_precision,
    date_include_unknown,
    date_buffer,
    ...rest
  } = urlParams.value

  const date: DateFilters =
    date_from || date_to
      ? {
          enabled: true,
          from: CompositeDate.parse(date_from),
          to: CompositeDate.parse(date_to),
          is_range: date_is_range ?? false,
          precision: (date_precision ?? 'year') as EventDatePrecision,
          include_unknown: date_include_unknown ?? false,
          buffer: date_buffer
        }
      : { enabled: false, is_range: false, precision: 'year' }

  filters.value = {
    ...rest,
    date
  }
})

watch(
  filters,
  (f) => {
    if (!f) return
    console.log('Updating URL params with filtervalues:', f)
    urlParams.value.batches = f.batches
    urlParams.value.datasets = f.datasets
    urlParams.value.search_term = f.search_term
    urlParams.value.type_status = f.type_status
    urlParams.value.confer = f.confer
    urlParams.value.whole_clade = f.whole_clade
    urlParams.value.rank = f.rank
    urlParams.value.status = f.status
    urlParams.value.taxa = f.taxa
    if (f.date.enabled) {
      urlParams.value.date_is_range = f.date.is_range
      urlParams.value.date_precision = f.date.precision
      urlParams.value.date_from = f.date.from
        ? CompositeDate.toString(f.date.from, f.date.precision)
        : undefined
      urlParams.value.date_to = f.date.to
        ? CompositeDate.toString(f.date.to, f.date.precision)
        : undefined
      urlParams.value.date_include_unknown = f.date.include_unknown
      urlParams.value.date_buffer = f.date.buffer
    } else {
      urlParams.value.date_is_range = undefined
      urlParams.value.date_precision = undefined
      urlParams.value.date_from = undefined
      urlParams.value.date_to = undefined
      urlParams.value.date_include_unknown = undefined
      urlParams.value.date_buffer = undefined
    }
  },
  { deep: true }
)

const headers: DataTableHeader[] = [
  {
    title: 'Occurrence',
    children: [{ key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } }]
  },
  {
    title: 'Sampling',
    align: 'center',
    headerProps: { class: 'border-s' },
    children: [
      {
        key: 'sampling.site',
        title: 'Site'
      },
      { key: 'sampling.performed_on', title: 'Date', align: 'end' }
    ]
  },
  {
    key: 'identification',
    title: 'Identification',
    align: 'center',
    sortable: false,
    headerProps: { class: 'border-s' },
    children: [
      {
        title: 'Taxon',
        key: 'identification',
        sortable: true,
        align: 'start',
        sort: (a, b) => a.taxon.name.localeCompare(b.name),
        filter(value, query, item): FilterMatch {
          return (value as unknown as Identification).taxon.name
            .toLowerCase()
            .includes(query.toLowerCase())
        }
      },
      {
        key: 'identification.identified_on',
        title: 'Date',
        align: 'end'
      }
    ]
  }
] as const

function invalidateQuery() {
  queryClient.invalidateQueries({ queryKey: listOccurrencesQueryKey() })
}

onMounted(invalidateQuery)

function dateQuery(
  filters: DateFilters
): NonNullable<ListOccurrencesData['query']>['date'] | undefined {
  if (!filters.enabled || !filters.from) return undefined

  const from_date = filters.from
    ? CompositeDate.toString(filters.from, filters.precision)
    : undefined
  const to_date = filters.to ? CompositeDate.toString(filters.to, filters.precision) : undefined

  return {
    from: from_date,
    to: filters.is_range ? to_date : from_date,
    include_unknown: filters.include_unknown,
    buffer: filters.buffer
  }
}

async function reload() {
  await refetch()
  feedback({ message: 'Reloaded occurrences' })
}

const { data, error, isPending, isFetching, refetch } = useQuery(
  computed(() => ({
    staleTime: Infinity,
    ...listOccurrencesOptions({
      query: {
        limit: pagination.value.itemsPerPage,
        offset: (pagination.value.page - 1) * pagination.value.itemsPerPage,
        confer: filters.value.confer,
        datasets: filters.value.datasets,
        batches: filters.value.batches,
        taxon_rank: filters.value.rank,
        taxon_status: filters.value.status,
        taxa: filters.value.taxa,
        search_term: filters.value.search_term,
        type_status: filters.value.type_status,
        date: dateQuery(filters.value.date),
        sort: sortBy.value.key,
        sort_direction: sortBy.value.order as 'asc' | 'desc'
      }
    }),
    placeholderData: keepPreviousData
  }))
)

const queryClient = useQueryClient()
async function prefetchNext(currentPage: number) {
  await setTimeout(() => {}, 500)
  await queryClient.prefetchQuery({
    staleTime: Infinity,
    ...listOccurrencesOptions({
      query: {
        ...mappingFiltersToQuery(filters.value, 'occurrences'),
        ...{
          limit: pagination.value.itemsPerPage,
          offset: currentPage * pagination.value.itemsPerPage
        }
      }
    })
  })
}

onMounted(() => {
  prefetchNext(pagination.value.page)
})

const loading = computedAsync(async () => {
  return (
    isPending.value ||
    (isFetching.value &&
      (await promiseTimeout(1000).then(() => {
        return isFetching.value
      })))
  )
}, true)

function resetFilters() {
  filters.value = { date: { enabled: false, is_range: false, precision: 'year' } }
}
</script>

<style scoped lang="scss"></style>
