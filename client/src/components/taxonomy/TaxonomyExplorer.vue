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

    <!-- TREE -->
    <div class="taxonomy-explorer bg-surface" :style="{ 'grid-template-columns': templateColumns }">
      <!-- HEADERS -->
      <div
        v-for="{ rank } in headers.filter(
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
          @click="rank == 'Subspecies' ? unfold(parentRank(rank)!) : toggleFold(rank)"
        >
          {{ countsByRank[rank] }}
          <template #append>
            <v-icon
              v-if="rank !== 'Subspecies'"
              class="ml-1"
              :icon="isFolded(rank) ? 'mdi-plus-box-outline' : 'mdi-minus-box-outline'"
            />
          </template>
        </v-chip>
      </div>

      <!-- INNER TREE -->
      <div class="taxonomy-tree">
        <v-progress-linear v-if="loading" class="loading" indeterminate />
        <FTaxaNestedList
          v-if="filteredItems?.children"
          :items="filteredItems?.children"
          rank="Kingdom"
        />
        <div style="grid-column: start / span end; grid-row: -1"></div>
      </div>
    </div>

    <!-- FOOTER -->
    <div class="taxonomy-footer bg-surface pa-3 border-t-thin d-flex">
      <v-spacer />
      <v-btn color="primary" variant="plain" text="GBIF Imports" :to="{ name: 'import-GBIF' }">
        <template #prepend>
          <IconGBIF />
        </template>
      </v-btn>
    </div>

    <!-- MODALS -->
    <TaxonCard
      v-if="selected"
      v-model:open="showTaxonCard"
      v-model="selected"
      @add-child="addDescendant"
      @navigate="(target) => (selected = target)"
      @deleted="({ parent }) => update(parent?.id)"
    />
    <TaxonFormDialog v-model="formDialog" :parent="parentTaxon" @success="onTaxonCreated" />
  </div>
</template>

<script setup lang="ts">
import {
  $TaxonRank,
  GetTaxonomyData,
  Taxon,
  Taxonomy,
  TaxonomyService,
  TaxonRank,
  TaxonStatus
} from '@/api'
import { handleErrors } from '@/api/responses'
import { refDebounced, useLocalStorage } from '@vueuse/core'
import { computed, onMounted, provide, ref, watch } from 'vue'
import { MaxRankInjection, useRankFoldState, useTaxonFoldState, useTaxonSelection } from '.'
import { FTaxaNestedList } from './functionals'
import { isAscendant, isDescendant, parentRank } from './rank'
import StatusPicker from './StatusPicker.vue'
import TaxonCard from './TaxonCard.vue'
import TaxonFormDialog from './TaxonFormDialog.vue'
import TaxonRankPicker from './TaxonRankPicker.vue'
import IconGBIF from '../icons/IconGBIF.vue'

const formDialog = ref(false)
const parentTaxon = ref<Taxon>()
const showTaxonCard = ref(false)
function addDescendant(taxon: Taxon) {
  formDialog.value = true
  showTaxonCard.value = false
  parentTaxon.value = taxon
}

const { selected } = useTaxonSelection()
watch(selected, (taxon) => {
  showTaxonCard.value = taxon !== undefined
})

type Header = { rank: TaxonRank }

const headers: Header[] = [
  { rank: 'Kingdom' },
  { rank: 'Phylum' },
  { rank: 'Class' },
  { rank: 'Order' },
  { rank: 'Family' },
  { rank: 'Genus' },
  { rank: 'Species' },
  { rank: 'Subspecies' }
]

const maxRank = useLocalStorage<TaxonRank>('max-taxon-rank', 'Kingdom')
provide(MaxRankInjection, maxRank)

const { toggleFold, isFolded, unfold } = useRankFoldState()

const loading = ref(false)
const items = ref<Taxonomy>()
onMounted(async () => (items.value = await fetch()))

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
  if (!items.value || (!filterStatus.value && !debouncedSearchTerm.value)) return items.value
  const filters = {
    term: debouncedSearchTerm.value ? new RegExp(debouncedSearchTerm.value, 'i') : undefined,
    status: filterStatus.value
  }
  return {
    ...items.value,
    children: items.value.children?.map(matchSearch(filters)).filter((t) => t !== undefined)
  }
})

async function fetch(query?: GetTaxonomyData['query']) {
  loading.value = true
  const taxonomy = await TaxonomyService.getTaxonomy({ query }).then(
    handleErrors((err) => console.error(err))
  )
  loading.value = false
  return taxonomy
}

function find(subtree: Taxonomy, taxonID: string) {
  const match = subtree.children?.find(({ id }) => id === taxonID)
  if (match) return match
  return subtree.children?.reduce<Taxonomy | undefined>((acc, item): Taxonomy | undefined => {
    if (acc !== undefined) return acc
    if (!item.children) return acc
    return find(item, taxonID)
  }, undefined)
}

async function update(taxonID: string | undefined) {
  if (!taxonID) {
    items.value = await fetch()
    return
  }
  if (!items.value) return
  const toUpdate = find(items.value, taxonID)
  if (toUpdate === undefined) {
    console.error('Failed to find taxon with ID: ', taxonID)
    return
  }
  const subtree = await fetch({ identifier: taxonID })
  Object.assign(toUpdate, subtree)
}

async function onTaxonCreated(taxon: Taxonomy) {
  await update(taxon.parent?.id)
  const { show } = useTaxonFoldState(taxon)
  show()
}

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

function _countsByRank(acc: RanksCount, taxonomy: Taxonomy | undefined) {
  taxonomy?.children?.forEach((child) => {
    acc[child.rank] += 1
    if (child.children_count > 0) _countsByRank(acc, child)
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
