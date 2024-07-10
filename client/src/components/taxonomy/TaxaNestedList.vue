<template>
  <div class="pa-0 taxa-list bg-surface" :style="{ 'grid-column': `${rank} / span end` }">
    <template v-for="k in items" :key="k.id">
      <TaxonomyItem :item="k"></TaxonomyItem>
    </template>
  </div>
</template>

<script setup lang="ts">
import { Taxonomy, TaxonRank } from '@/api'
import TaxonomyItem from './TaxonomyItem.vue'

const props = defineProps<{ items: Taxonomy[]; rank: TaxonRank }>()
</script>

<style scoped lang="scss">
@use 'vuetify';

.taxa-list {
  display: grid;
  grid-template-columns: subgrid;
  grid-template-rows: auto;
  grid-column: start / span end;
}
</style>
