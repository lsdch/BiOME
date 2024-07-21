<template>
  <div
    v-if="!isAscendant(item.rank, maxRankDisplay)"
    class="taxon-item-container"
    ref="container"
    :style="{ 'grid-column': item.rank }"
    :id="item.name"
  >
    <div class="taxon-item">
      <FTaxonStatusIndicator :status="item.status" />
      <span class="mr-3 text-no-wrap cursor-pointer" @click="select(item)">
        {{ item.name }}
      </span>
      <v-spacer></v-spacer>
      <v-icon v-if="item.anchor" icon="mdi-pin" size="x-small" color="warning" />
      <v-chip
        v-if="item.children_count > 0"
        :color="expanded ? 'success' : 'primary'"
        size="small"
        density="compact"
        @click="toggleAndScroll()"
        :rounded="100"
      >
        {{ item.children_count }}
      </v-chip>
    </div>
  </div>

  <FTaxaNestedList
    v-if="expanded && item.children"
    :items="item.children"
    :rank="item.children[0].rank"
  />
</template>

<script setup lang="ts">
import { Taxonomy } from '@/api'
import { useElementVisibility } from '@vueuse/core'
import { nextTick, ref, watch } from 'vue'
import { maxRankDisplay, useRankFoldState, useTaxonFoldState, useTaxonSelection } from '.'
import { FTaxaNestedList, FTaxonStatusIndicator } from './functionals'
import { isAscendant } from './rank'

watch(maxRankDisplay, (rank) => {
  if (isAscendant(props.item.rank, rank)) {
    expanded.value = true
  }
})

const props = defineProps<{ item: Taxonomy }>()

const { select } = useTaxonSelection()

const { onFold: onRankFold, onUnfold: onRankUnfold, isFolded: isRankFolded } = useRankFoldState()
const { expanded, toggleFold } = useTaxonFoldState(props.item, !isRankFolded(props.item.rank))
// const expanded = ref(!isRankFolded(props.item.rank))
// const toggleFold = useToggle(expanded)
onRankFold((rank) => {
  if (rank == props.item.rank) expanded.value = false
})
onRankUnfold((rank) => {
  if (rank == props.item.rank) expanded.value = true
})

async function toggleAndScroll() {
  const toggled = toggleFold()
  await nextTick()
  setTimeout(() => {
    if (!toggled && !containerVisible.value) scrollTo()
  }, 25)
}

function scrollTo() {
  document.getElementById(`${props.item.name}`)!.scrollIntoView({ block: 'center' })
}

const container = ref()
const containerVisible = useElementVisibility(container)
</script>

<style scoped lang="scss">
@use 'vuetify';
.taxon-item-container {
  padding: 0.3rem;
  /* border-right: thin solid rgba(var(--v-border-color), var(--v-border-opacity)); */
  border: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
  margin-top: -1px;
  margin-left: -1px;
  > div.taxon-item {
    position: sticky;
    top: 60px;
    display: flex;
    // justify-content: space-between;
    align-items: center;
  }
}
</style>
