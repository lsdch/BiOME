<template>
  <div class="taxonomy-container d-flex flex-column">
    <div class="taxonomy-explorer bg-surface" :style="{ 'grid-template-columns': templateColumns }">
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
        <v-btn
          color="primary"
          variant="plain"
          density="compact"
          :icon="foldState[rank] ? 'mdi-minus-box-outline' : 'mdi-plus-box-outline'"
          @click="toggleFold(rank)"
        />
      </div>
      <div class="taxonomy-tree">
        <v-progress-linear v-if="loading" class="loading" indeterminate />
        <TaxaNestedList :items rank="Kingdom" />
      </div>
    </div>
    <div class="taxonomy-footer bg-surface pa-3 border-t-thin">
      <TaxonRankPicker v-model="maxRank" label="Truncate above" hide-details density="compact" />
    </div>
    <v-bottom-sheet v-model="showTaxon" :inset="mdAndUp" content-class="rounded-0">
      <TaxonCard v-if="selected" :taxon="selected" />
    </v-bottom-sheet>
  </div>
</template>

<script setup lang="ts">
import { $TaxonRank, Taxonomy, TaxonomyService, TaxonRank } from '@/api'
import { handleErrors } from '@/api/responses'
import { useEventBus, useLocalStorage } from '@vueuse/core'
import { computed, onMounted, provide, ref, watch } from 'vue'
import TaxaNestedList from './TaxaNestedList.vue'
import TaxonRankPicker from './TaxonRankPicker.vue'
import { FoldEvent, MaxRankInjection, useTaxonSelection } from '.'
import { isAscendant, isDescendant } from './rank'
import TaxonCard from './TaxonCard.vue'
import { useDisplay } from 'vuetify'

const { mdAndUp } = useDisplay()

const showTaxon = ref(false)
const { selected } = useTaxonSelection()
watch(selected, (taxon) => {
  showTaxon.value = taxon !== undefined
})

type Header = { rank: TaxonRank & string }

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

const foldState = ref<{ [k in TaxonRank]: boolean | undefined }>({
  Kingdom: true,
  Phylum: true,
  Class: true,
  Order: true,
  Family: true,
  Genus: true,
  Species: true,
  Subspecies: true
})

const { emit } = useEventBus<FoldEvent>('fold')

function toggleFold(rank: TaxonRank) {
  foldState.value[rank] = !foldState.value[rank]
  emit({ rank, state: foldState.value[rank] })
}

const loading = ref(false)
const items = ref<Taxonomy[]>([])

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
</script>

<style scoped lang="scss">
.taxonomy-container {
  height: 0px;
  min-height: 100%;
}

.taxonomy-explorer {
  flex-grow: 1;
  display: grid;
  // grid-template-columns: dynamically defined in component
  grid-template-rows: 0fr auto;
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

  .taxonomy-footer {
    position: sticky;
    bottom: 0px;
    left: 0px;
  }
}
</style>
