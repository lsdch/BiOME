<template>
  <div class="d-flex flex-column fill-height overflow-auto">
    <v-toolbar flat dark density="compact">
      <v-tabs v-model="tab" class="">
        <v-tab value="tree" :rounded="0" prepend-icon="mdi-graph" text="Taxonomy tree" />
        <v-tab value="table" :rounded="0" prepend-icon="mdi-table" text="Taxonomy table" />
        <v-tab value="GBIF" :rounded="0" :prepend-icon="IconGBIF" text="GBIF imports" />
      </v-tabs>
    </v-toolbar>
    <v-window v-model="tab" class="flex-grow-1 fill-height overflow-auto">
      <v-window-item value="tree" class="fill-height">
        <TaxonomyExplorer />
      </v-window-item>
      <v-window-item value="table" class="fill-height" style="max-height: 100%">
        <TaxaTable class="fill-height" />
      </v-window-item>
      <v-window-item value="GBIF">
        <ImportGBIF />
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import IconGBIF from '@/components/icons/IconGBIF.vue'
import ImportGBIF from '@/components/taxonomy/imports/ImportGBIF.vue'
import TaxaTable from '@/components/taxonomy/TaxaTable.vue'
import TaxonomyExplorer from '@/components/taxonomy/TaxonomyExplorer.vue'

type TaxonomyTab = 'tree' | 'table' | 'GBIF'
const tab = ref<TaxonomyTab>('tree')
</script>

<style></style>
