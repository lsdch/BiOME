<template>
  <v-list class="map-tool-filters mt-3">
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
      <TaxonFilterPicker
        v-model:taxa="filters.taxa"
        v-model:whole-clade="filters.whole_clade"
        density="compact"
        multiple
        chips
        closable-chips
        hide-details
      />
    </v-list-item>
    <v-divider class="my-2" />

    <v-list-item>
      <TaxonFilterPicker
        v-model:taxa="filters.sampling_target_taxa"
        v-model:whole-clade="filters.sampling_target_whole_clade"
        label="Targeted taxa"
        item-value="name"
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
      <HabitatPicker
        label="Habitats"
        v-model="filters.habitats"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
      />
    </v-list-item>
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
      />
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { HabitatRecord, OccurrencesBySiteData } from '@/api'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import HabitatPicker from '@/features/registries/components/HabitatPicker.vue'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { reactive, ref } from 'vue'

export type MappingFilters = Overwrite<
  NonNullable<OccurrencesBySiteData['query']>,
  { habitats?: HabitatRecord[] }
>

const filters = defineModel<MappingFilters>({ default: () => reactive({}) })
</script>

<style lang="scss">
div.v-list.map-tool-filters .v-list-item .v-list-item__content .v-input {
  margin-top: 5px;
}
</style>
