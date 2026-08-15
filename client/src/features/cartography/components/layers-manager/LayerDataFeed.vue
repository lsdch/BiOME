<template>
  <v-list class="map-tool-filters mt-1">
    <ListItemInput label="Mode">
      <v-chip-group mandatory color="success" v-model="queryMode">
        <v-chip label value="occurrences">Occurrences</v-chip>
        <v-chip label value="samplings">Samplings</v-chip>
      </v-chip-group>
      <InlineHelp>
        <ul>
          <li>
            <b class="text-success"> Occurrences </b> mode displays all samplings with occurrences
            that match the defined filters, <i>excluding samplings with no occurrences</i>.
          </li>
          <li>
            <b class="text-info"> Samplings </b> mode displays all samplings that match the defined
            filters, <i>including</i> those with no occurrences.
          </li>
        </ul>
      </InlineHelp>
    </ListItemInput>
    <v-divider></v-divider>
    <v-list-item>
      <DatasetPicker
        v-model="filters.datasets"
        density="compact"
        item-value="slug"
        label="Datasets"
        multiple
        chips
        closable-chips
        clear-on-select
        clearable
        placeholder="All datasets"
        persistent-placeholder
        hide-details
      />
    </v-list-item>
    <v-list-item>
      <ImportBatchPicker
        v-model="filters.batches"
        item-value="id"
        label="Import batches"
        multiple
        clearable
        chips
        closable-chips
        density="compact"
        placeh
        older="All batches"
        persistent-placeholder
        placeholder="All batches"
        hide-details
      />
    </v-list-item>
    <v-list-item v-if="queryMode === 'occurrences'">
      <TaxonFilterPicker
        v-model:taxa="filters.taxa"
        v-model:whole-clade="filters.whole_clade"
        item-value="id"
        density="compact"
        multiple
        chips
        closable-chips
        clearable
        hide-details
      />
    </v-list-item>
    <v-divider class="my-2" />

    <v-list-item>
      <TaxonFilterPicker
        label="Sampling target"
        v-model:taxa="filters.target_taxa"
        v-model:whole-clade="filters.target_taxa_whole_clade"
        item-value="id"
        density="compact"
        multiple
        chips
        closable-chips
        clearable
        hide-details
      />
    </v-list-item>

    <!-- <v-list-item>
      <HabitatPicker
        label="Habitats"
        v-model="filters.habitats"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
      />
    </v-list-item> -->
    <v-list-item>
      <CountryPicker
        density="compact"
        multiple
        v-model="filters.countries"
        item-value="code"
        clear-on-select
        chips
        closable-chips
        clearable
        hide-details²
        :filter="(c) => c.sampling_count > 0"
      />
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import HabitatPicker from '@/features/registries/components/HabitatPicker.vue'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import { reactive } from 'vue'
import { MappingFilters } from './map-layers'
import ImportBatchPicker from '@/features/import/components/ImportBatchPicker.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'

const filters = defineModel<MappingFilters>({ default: () => reactive({}) })

type QueryMode = 'occurrences' | 'samplings'
const queryMode = defineModel<QueryMode>('mode', { default: 'occurrences' })
</script>

<style lang="scss">
div.v-list.map-tool-filters .v-list-item .v-list-item__content .v-input {
  margin-top: 5px;
}
</style>
