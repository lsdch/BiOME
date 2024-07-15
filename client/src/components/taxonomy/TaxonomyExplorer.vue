<template>
  <div class="taxonomy-container d-flex flex-column">
    <v-toolbar flat dense>
      <template #prepend>
        <v-btn icon="mdi-graph" variant="outlined" color="secondary" size="small" />
      </template>

      <v-toolbar-title style="min-width: 150px" class="flex-grow-0">Taxonomy</v-toolbar-title>

      <div class="flex-grow-1 mx-3 d-flex justify-center">
        <v-text-field
          v-model="searchTerm"
          class="mr-2"
          label="Search"
          hide-details
          density="compact"
          clearable
          prepend-inner-icon="mdi-magnify"
          color="primary"
          max-width="500px"
        />
        <StatusPicker
          max-width="220px"
          v-model="filterStatus"
          density="compact"
          hide-details
          clearable
        />
      </div>

      <template #append>
        <TaxonRankPicker
          v-model="maxRank"
          label="Truncate above"
          hide-details
          density="compact"
          min-width="200px"
        />
      </template>
    </v-toolbar>
    <div class="taxonomy-explorer bg-surface" :style="{ 'grid-template-columns': templateColumns }">
      <div
        v-for="{ rank, noFold } in headers.filter(
          ({ rank }) => rank == maxRank || isDescendant(rank, maxRank)
        )"
        :key="rank"
        :style="{ 'grid-column': rank }"
        class="taxonomy-header bg-surface"
      >
        <span class="text-overline">
          {{ rank }}
        </span>
        <v-chip
          size="small"
          rounded="100"
          color="primary"
          @click="noFold ? unfold(parentRank(rank)!) : toggleFold(rank)"
        >
          {{ countsByRank[rank] }}
          <template #append>
            <v-icon
              v-if="!noFold"
              class="ml-1"
              :icon="isFolded(rank) ? 'mdi-plus-box-outline' : 'mdi-minus-box-outline'"
            />
          </template>
        </v-chip>
      </div>
      <div class="taxonomy-tree">
        <v-progress-linear v-if="loading" class="loading" indeterminate />
        <FTaxaNestedList :items="filteredItems" rank="Kingdom" />
        <div style="grid-column: start / span end; grid-row: -1"></div>
      </div>
    </div>
    <div class="taxonomy-footer bg-surface pa-3 border-t-thin"></div>

    <TaxonCard
      v-if="selected"
      v-model:open="showTaxon"
      v-model="selected"
      @add-child="addDescendant"
      @navigate="(target) => (selected = target)"
    />
    <TaxonFormDialog v-model="formDialog" :parent="parentTaxon" />
  </div>
</template>

<script setup lang="ts">
import { $TaxonRank, Taxon, Taxonomy, TaxonomyService, TaxonRank, TaxonStatus } from '@/api'
import { handleErrors } from '@/api/responses'
import { refDebounced, useLocalStorage } from '@vueuse/core'
import { computed, onMounted, provide, ref, watch } from 'vue'
import { MaxRankInjection, useFoldState, useTaxonSelection } from '.'
import { FTaxaNestedList } from './functionals'
import { isAscendant, isDescendant, parentRank } from './rank'
import StatusPicker from './StatusPicker.vue'
import TaxonCard from './TaxonCard.vue'
import TaxonFormDialog from './TaxonFormDialog.vue'
import TaxonRankPicker from './TaxonRankPicker.vue'

const formDialog = ref(false)
const parentTaxon = ref<Taxon>()
function addDescendant(taxon: Taxon) {
  formDialog.value = true
  showTaxon.value = false
  parentTaxon.value = taxon
}

const showTaxon = ref(false)
const { selected } = useTaxonSelection()
watch(selected, (taxon) => {
  showTaxon.value = taxon !== undefined
})

type Header = { rank: TaxonRank & string; noFold?: boolean }

const headers: Header[] = [
  { rank: 'Kingdom' },
  { rank: 'Phylum' },
  { rank: 'Class' },
  { rank: 'Order' },
  { rank: 'Family' },
  { rank: 'Genus' },
  { rank: 'Species' },
  { rank: 'Subspecies', noFold: true }
]

const maxRank = useLocalStorage<TaxonRank>('max-taxon-rank', 'Kingdom')
provide(MaxRankInjection, maxRank)

const { toggleFold, isFolded, unfold } = useFoldState()

const loading = ref(false)
const items = ref<Taxonomy[]>([])

const filterStatus = ref<TaxonStatus>()
const searchTerm = ref<string>()
const debouncedSearchTerm = refDebounced(searchTerm, 200)

type SearchFilters = {
  status?: TaxonStatus
  term?: RegExp
}

function taxonMatches(taxon: Taxonomy, filters: SearchFilters) {
  return (
    (filters.status ? taxon.status === filters.status : true) &&
    (filters.term ? taxon.name.match(filters.term) : true)
  )
}

function matchSearch(filters: SearchFilters) {
  return (taxon: Taxonomy): Taxonomy | undefined => {
    if (taxonMatches(taxon, filters)) {
      if (isFolded(taxon.rank)) unfold(taxon.rank)
      return taxon
    }
    const matchingChildren =
      taxon.children?.map(matchSearch(filters)).filter((t) => t !== undefined) ?? []
    return matchingChildren.length > 0 ? { ...taxon, children: matchingChildren } : undefined
  }
}

const filteredItems = computed(() => {
  if (!filterStatus.value && !debouncedSearchTerm.value) return items.value
  const filters = {
    term: debouncedSearchTerm.value ? new RegExp(debouncedSearchTerm.value, 'i') : undefined,
    status: filterStatus.value
  }
  return items.value.map(matchSearch(filters)).filter((t) => t !== undefined)
})

async function fetch() {
  loading.value = true
  const taxa = await TaxonomyService.getTaxonomy({
    query: { 'max-depth': undefined }
  }).then(handleErrors((err) => console.error(err)))
  loading.value = false
  return taxa
}

onMounted(async () => (items.value = await fetch()))

const templateColumns = computed(() => {
  return $TaxonRank.enum
    .reduce((acc, rank) => {
      const name = `[${rank}${rank == 'Kingdom' ? ' start' : ''}]`
      return `${acc} ${name} ${isAscendant(rank, maxRank.value) ? '0px' : 'auto'}`
    }, '')
    .concat(' [end]')
})

type RanksCount = {
  [k in TaxonRank]: number
}

function _countsByRank(acc: RanksCount, taxa: Taxonomy[]) {
  taxa.forEach(({ rank, children, children_count }) => {
    acc[rank] += 1
    if (children_count > 0) _countsByRank(acc, children)
  })
  return acc
}

const countsByRank = computed(() => {
  const acc: RanksCount = {
    Kingdom: 0,
    Phylum: 0,
    Class: 0,
    Order: 0,
    Family: 0,
    Genus: 0,
    Species: 0,
    Subspecies: 0
  }

  return _countsByRank(acc, items.value)
})
</script>

<style lang="scss">
.taxonomy-container {
  height: 0px;
  min-height: 100%;
}

.taxonomy-explorer {
  flex-grow: 1;
  display: grid;
  // grid-template-columns: dynamically defined in component
  grid-template-rows: 0fr auto 1fr;
  border-collapse: collapse;
  overflow: scroll;

  > .taxonomy-tree {
    display: grid;
    grid-column: start / span end;
    grid-template-columns: subgrid;
    grid-template-rows: auto;
    .loading {
      grid-column: start / span end;
    }
  }

  .taxonomy-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 15px;
    border-right: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
    border-bottom: 3px solid rgba(var(--v-border-color), var(--v-border-opacity));
    margin-bottom: -1px;
    z-index: 200;
    position: sticky;
    top: 0;
    height: 60px;
  }

  .taxa-list {
    display: grid;
    grid-template-columns: subgrid;
    grid-template-rows: auto;
    grid-column: start / span end;
  }

  .taxonomy-footer {
    position: sticky;
    bottom: 0px;
    left: 0px;
  }
}
</style>
