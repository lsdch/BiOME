<template>
  <div
    v-if="!isAscendant(item.rank, maxRank)"
    class="taxon-item-container"
    ref="container"
    :style="{ 'grid-column': item.rank }"
    :id="item.name"
  >
    <div class="taxon-item">
      <LinkIconGBIF v-if="item.GBIF_ID" :GBIF_ID="item.GBIF_ID" size="small" variant="plain" />
      <span class="mr-3 text-no-wrap cursor-pointer" @click="select(omitChildren(item))">
        {{ item.name }}
      </span>
      <v-spacer></v-spacer>
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
  <!-- <v-expand-transition @transitionend="scrollTo"> -->
  <TaxaNestedList
    v-if="expanded && item.children_count > 0"
    :items="item.children"
    :rank="item.children[0].rank"
  />
  <!-- </v-expand-transition> -->
</template>

<script setup lang="ts">
import { Taxon, Taxonomy } from '@/api'
import { useElementVisibility, useEventBus, useToggle } from '@vueuse/core'
import { inject, nextTick, ref, watch } from 'vue'
import { FoldEvent, MaxRankInjection, useTaxonSelection } from '.'
import LinkIconGBIF from './LinkIconGBIF.vue'
import { isAscendant } from './rank'
import TaxaNestedList from './TaxaNestedList.vue'

const maxRank = inject(MaxRankInjection) ?? ref('Kingdom')
watch(maxRank, (rank) => {
  if (isAscendant(props.item.rank, rank)) {
    expanded.value = true
  }
})

const props = defineProps<{ item: Taxonomy }>()

const { select } = useTaxonSelection()

function omitChildren(item: Taxonomy): Taxon {
  const { children: _, ...rest } = item
  return rest
}

const expanded = ref(true)
const toggle = useToggle(expanded)
const waitingForScroll = ref(false)

const { on: onFold } = useEventBus<FoldEvent>('fold')

onFold(({ rank, state }) => {
  if (rank == props.item.rank && state !== undefined) {
    expanded.value = state
  }
})

function toggleAndScroll() {
  const toggled = toggle()
  waitingForScroll.value = !toggled
  nextTick(scrollTo)
}

function scrollTo() {
  if (waitingForScroll.value && !containerVisible.value)
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
