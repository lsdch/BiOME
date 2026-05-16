<template>
  <DeckGlMap :hexgrid="hexgrid">
    <template #cluster-popup="{ data }">
      <SiteClusterPopup :data />
    </template>
    <template #popup="{ item }">
      <SitePopupWithOccurrences :item />
    </template>
  </DeckGlMap>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import DeckGlMap from '@/features/cartography/components/DeckGlMap.vue'
import {
  HexgridLayer,
  makeHexLayer
} from '@/features/cartography/components/layers-manager/map-layers'
import SiteClusterPopup from '@/features/cartography/components/popups/SiteClusterPopup.vue'
import SitePopupWithOccurrences from '@/features/cartography/components/popups/SitePopupWithOccurrences.vue'
import { hexgridLayerFromSpec } from '@/features/cartography/composables/hexgrid-layer'
import { computed, ref } from 'vue'

const { sites } = defineProps<{
  sites: SiteWithOccurrences[]
}>()

const hexLayerSpec = ref(makeHexLayer())

const hexgrid = computed<HexgridLayer<SiteWithOccurrences>>(() => {
  return hexgridLayerFromSpec(hexLayerSpec.value, sites)
})
</script>

<style scoped lang="scss"></style>
