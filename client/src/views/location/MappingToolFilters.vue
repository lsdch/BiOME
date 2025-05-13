<template>
  <v-navigation-drawer v-model="drawer" :location="$vuetify.display.xs ? 'top' : 'left'">
    <v-list-subheader class="px-2"> Filters </v-list-subheader>
    <v-divider></v-divider>
    <div class="mt-3 px-2">
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
      />
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
      <TaxonPicker
        v-model="filters.taxa"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
        hide-details
        :ranks="sampledTaxa.rank"
        :sampled-only="sampledTaxa.sampledOnly"
      >
        <template #prepend-item>
          <div class="d-flex px-4 align-start position-sticky ga-3">
            <v-switch
              v-model="sampledTaxa.sampledOnly"
              label="Sampled only"
              color="primary"
              hide-details
              density="compact"
            />
            <v-spacer></v-spacer>
            <TaxonRankPicker
              v-model="sampledTaxa.rank"
              label="Rank"
              hide-details
              density="compact"
            />
          </div>
          <v-divider class="my-1" />
        </template>
      </TaxonPicker>
      <v-switch
        class="px-2"
        label="Use clade"
        v-model="filters.whole_clade"
        color="primary"
        hint="Include descendant taxa"
        density="compact"
      >
        <template #append>
          <InlineHelp>
            When enabled, all occurrences of descendant taxa will be included.
          </InlineHelp>
        </template>
      </v-switch>
      <v-divider class="my-2" />
      <SamplingTargetKindFilter density="compact" v-model="filters.sampling_target_kinds" />
      <TaxonPicker
        v-model="filters.sampling_target_taxa"
        label="Targeted taxa"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
        clearable
        hide-details
      />
      <v-divider class="my-2" />
      <HabitatPicker
        label="Habitats"
        v-model="filters.habitats"
        item-value="name"
        density="compact"
        multiple
        chips
        closable-chips
      />
    </div>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { HabitatRecord, OccurrencesBySiteData } from '@/api'
import DatasetPicker from '@/components/datasets/DatasetPicker.vue'
import HabitatPicker from '@/components/occurrence/habitat/HabitatPicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import TaxonRankPicker from '@/components/taxonomy/TaxonRankPicker'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'
import { useToggle } from '@vueuse/core'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { reactive, ref } from 'vue'
import SamplingTargetKindFilter from './SamplingTargetKindFilter.vue'

export type MappingFilters = Overwrite<
  NonNullable<OccurrencesBySiteData['query']>,
  { habitats?: HabitatRecord[] }
>

const [drawer, toggleDrawer] = useToggle(true)

const sampledTaxa = ref({
  rank: undefined,
  sampledOnly: true
})

const filters = defineModel<MappingFilters>({ default: () => reactive({}) })
</script>

<style scoped lang="scss"></style>
