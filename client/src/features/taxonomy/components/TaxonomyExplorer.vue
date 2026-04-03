<template>
  <div class="taxonomy-container d-flex flex-column">
    <TableToolbar
      title="Taxonomy"
      icon="mdi-family-tree"
      :togglable-search="smAndDown"
      @reload="refetch()"
    >
      <template #search>
        <v-container>
          <v-row>
            <v-col cols="12" sm="6" lg="8">
              <v-text-field
                v-model="searchTerm"
                label="Search"
                hide-details
                density="compact"
                clearable
                prepend-inner-icon="mdi-magnify"
                color="primary"
              />
            </v-col>
            <v-col cols="12" sm="6" lg="4">
              <StatusPicker v-model="filterStatus" density="compact" hide-details clearable />
            </v-col>
          </v-row>
        </v-container>
      </template>
      <template #append>
        <TaxonRankPicker
          class="ml-3"
          v-model="maxRankDisplay"
          :exclude="['Subgenus']"
          label="Truncate above"
          hide-details
          density="compact"
          min-width="200px"
        />
      </template>
    </TableToolbar>

    <!-- TREE -->
    <div class="taxonomy-explorer bg-surface" :style="{ 'grid-template-columns': templateColumns }">
      <!-- HEADERS -->
      <div
        v-for="{ rank } in headers.filter(
          ({ rank }) => rank == maxRankDisplay || TaxonRank.isDescendant(rank, maxRankDisplay)
        )"
        :key="rank"
        :style="{ 'grid-column': rank }"
        class="taxonomy-header bg-surface"
      >
        <span class="text-overline">
          {{ rank === 'Species' ? 'Species / Subgenus' : rank }}
        </span>
        <v-chip
          size="small"
          rounded="100"
          color="primary"
          @click="rank == 'Subspecies' ? unfold(TaxonRank.parentRank(rank)!) : toggleFold(rank)"
        >
          {{ countsByRank[rank] + (rank == "Species" ? countsByRank["Subgenus"] : 0) }}
          <template #append>
            <v-progress-circular
              v-if="isRankFetching(rank)"
              class="ml-1"
              size="14"
              width="2"
              indeterminate
            />
            <v-icon
              v-else-if="rank !== 'Subspecies'"
              class="ml-1"
              :icon="isFolded(rank) ? 'mdi-plus-box-outline' : 'mdi-minus-box-outline'"
            />
          </template>
        </v-chip>
      </div>

      <!-- INNER TREE -->
      <div class="taxonomy-tree">
        <v-progress-linear v-if="loading" class="loading" indeterminate />
        <v-container v-if="error" style="grid-column: start / span end">
          <v-alert type="error" icon="mdi-alert"> Failed to load taxonomy </v-alert>
        </v-container>
        <FTaxaNestedList
          v-else-if="filteredItems?.length"
          :items="filteredItems"
          rank="Kingdom"
        />
        <div v-else class="mx-auto my-5" style="grid-column: start / span end">
          {{ loading ? 'Loading...' : 'Nothing to display' }}
        </div>
        <div style="grid-column: start / span end; grid-row: -1"></div>
      </div>
    </div>

    <!-- FOOTER -->
    <div class="taxonomy-footer bg-surface pa-3 border-t-thin d-flex">
      
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
    <TaxonFormDialogMutation
      v-if="parentTaxon"
      v-model:dialog="formDialog"
      v-model:parent="parentTaxon"
      @success="onTaxonCreated"
    />
  </div>
</template>

<script setup lang="ts">
import {
  $TaxonRank,
  Taxon,
  TaxonRank,
  TaxonStatus
} from '@/api'
import { getTaxonomyAtRankOptions } from '@/api/gen/@tanstack/vue-query.gen'
import TaxonFormDialogMutation from '@/components/forms/TaxonFormDialogMutation.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { useQuery } from '@tanstack/vue-query'
import { refDebounced } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useDisplay } from 'vuetify'
import {
  maxRankDisplay,
  useRankFoldState,
  useTaxonFoldState,
  useTaxonSelection
} from '../composables'
import { FTaxaNestedList } from './functionals'
import StatusPicker from './StatusPicker.vue'
import TaxonCard from './TaxonCard.vue'
import { TaxonomyElement } from './TaxonomyItem.vue'
import TaxonRankPicker from './TaxonRankPicker'

const { smAndDown } = useDisplay()

const formDialog = ref(false)
const parentTaxon = ref<Taxon>()
const showTaxonCard = ref(false)
function addDescendant(taxon: Taxon) {
  formDialog.value = true
  showTaxonCard.value = false
  parentTaxon.value = taxon
}


const { selected, onSelect, select } = useTaxonSelection()
onSelect((_taxon) => {
  showTaxonCard.value = true
})
 
type Header = { rank: TaxonRank.NoSubgenus }

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

const { toggleFold, isFolded, unfold } = useRankFoldState()

const RANKS_TO_LOAD = $TaxonRank.enum
const TAXONOMY_STALE_MS = 5 * 60 * 1000
const TAXONOMY_GC_MS = 30 * 60 * 1000


// Helper to assemble tree structure from flat rank-filtered lists
function assembleTreeFromRanks(rankData: Record<string, TaxonomyElement[]>): TaxonomyElement[] {
  // Build maps for each rank for O(1) lookup
  const rankMap: Record<string, Map<string, TaxonomyElement>> = RANKS_TO_LOAD.reduce(
    (acc, rank) => {
      const taxaMap = new Map<string, TaxonomyElement>()
      const rankTaxa = rankData[rank] ?? []
      rankTaxa.forEach((taxon) => {
        taxaMap.set(taxon.name, { ...taxon, children: [] } as TaxonomyElement)
      })
      acc[rank] = taxaMap
      return acc
    },
    {} as Record<string, Map<string, TaxonomyElement>>
  )

  // Link parents to children
  RANKS_TO_LOAD.forEach((rank) => {
    const rankMapForRank = rankMap[rank]
    if (!rankMapForRank) return
    
    const parentRank = TaxonRank.parentRank(rank)
    rankMapForRank.values().forEach((taxon) => {
      if (taxon.parent?.name) {
        try {
          const parent = parentRank ? rankMap[parentRank]?.get(taxon.parent.name) : null
          if (parent) {
            parent.children = parent.children || []
            parent.children.push(taxon)
          }
        } catch (e) {
          console.log("error:", e)
          // Continue if parent rank resolution fails
        }
      }
    })
  })

  return Array.from(rankMap.Kingdom?.values() ?? []) 
}

// Create queries for each rank in parallel
const rankQueries = RANKS_TO_LOAD.map((rank) => {
  const options = getTaxonomyAtRankOptions({
    path: { rank }
  })
  return useQuery({
    ...options,
    staleTime: TAXONOMY_STALE_MS,
    gcTime: TAXONOMY_GC_MS
  })
})

const loading = computed(() => rankQueries.some((q) => q.isPending.value))
const error = computed(() => rankQueries.find((q) => q.error.value)?.error.value)

function isRankFetching(rank: TaxonRank.NoSubgenus): boolean {
  const rankIndex = RANKS_TO_LOAD.indexOf(rank)
  if (rankIndex === -1) return false
  const query = rankQueries[rankIndex]
  return Boolean(query?.isPending.value || query?.isFetching.value)
}

// Unified data from all ranks assembled into tree
const items = computed(() => {
  const rankData: Record<string, TaxonomyElement[]> = {}
  RANKS_TO_LOAD.forEach((rank, index) => {
    const query = rankQueries[index]
    if (!query) return

    const data = query.data.value
    const rankKey = String(rank)    
    rankData[rankKey] = data ?? []
  })
  return assembleTreeFromRanks(rankData)
})

const refetch = async () => {
  await Promise.all(rankQueries.map((q) => q.refetch?.()))
}

const filterStatus = ref<TaxonStatus>()
const searchTerm = ref<string>()
const debouncedSearchTerm = refDebounced(searchTerm, 200)

type SearchFilters = {
  status?: TaxonStatus
  term?: RegExp
}

function taxonMatches(taxon: TaxonomyElement, filters: SearchFilters) {
  return (
    (filters.status ? taxon.status === filters.status : true) &&
    (filters.term ? taxon.name.match(filters.term) : true)
  )
}

function matchSearch(filters: SearchFilters) {
  return (taxon: TaxonomyElement): TaxonomyElement | undefined => {
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
  return items.value?.map(matchSearch(filters)).filter((t) => t !== undefined)
})

async function update(taxonID: string | undefined) {
  if (!taxonID) {
    await refetch()
    return
  }
  await refetch()
}

async function onTaxonCreated(taxon: TaxonomyElement) {
  await update(taxon.parent?.id)
  const { show } = useTaxonFoldState(taxon)
  show()
}

// CSS grid template columns based on ranks and maxRankDisplay
const templateColumns = computed(() => {
  return $TaxonRank.enum
    .reduce((acc, rank) => {
      if (rank === 'Subgenus') return acc
      const name = `[${rank}${rank == 'Kingdom' ? ' start' : ''}]`
      return `${acc} ${name} ${TaxonRank.isAscendant(rank, maxRankDisplay.value) ? '0px' : 'auto'}`
    }, '')
    .concat(' [end]')
})

type RanksCount = {
  [k in TaxonRank]: number
}

const countsByRank = computed(() => {
  const acc: RanksCount = {
    Kingdom: 0,
    Phylum: 0,
    Class: 0,
    Order: 0,
    Family: 0,
    Genus: 0,
    Subgenus: 0,
    Species: 0,
    Subspecies: 0
  }
  
  // Count items from parallel rank queries
  RANKS_TO_LOAD.forEach((rank, index) => {
    acc[rank as TaxonRank] = rankQueries[index]?.data.value?.length ?? 0
  })

  return acc
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
