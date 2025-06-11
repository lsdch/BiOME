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
            <v-spacer />
            <TaxonRankPicker
              v-model="sampledTaxa.rank"
              label="Rank"
              hide-details
              density="compact"
            />
          </div>
          <v-divider />
        </template>
      </TaxonPicker>
    </v-list-item>
    <v-list-item>
      <v-switch
        class="px-2"
        label="Use clade"
        v-model="filters.whole_clade"
        color="primary"
        density="compact"
        :disabled="!filters.taxa?.length"
        hide-details
      >
        <template #append>
          <InlineHelp>
            When enabled, all occurrences of descendant taxa will be included.
          </InlineHelp>
        </template>
      </v-switch>
    </v-list-item>
    <v-divider class="my-2" />

    <v-list-item>
      <SamplingTargetKindFilter density="compact" v-model="filters.sampling_target_kinds">
        <template #chip="{ item, props }">
          <v-chip
            v-bind="props"
            :title="item.title"
            :closable="item.value !== 'Taxa' || !filters.sampling_target_taxa?.length"
          />
        </template>
        <template #item="{ item, props }">
          <v-list-item
            v-bind="props"
            :title="item.title"
            density="compact"
            lines="one"
            :disabled="item.value === 'Taxa' && !!filters.sampling_target_taxa?.length"
          >
            <template #prepend="{ isSelected, select }">
              <v-checkbox
                :model-value="isSelected"
                @update:model-value="(v) => select(v ?? false)"
                color=""
                hide-details
              />
            </template>
          </v-list-item>
        </template>
      </SamplingTargetKindFilter>
    </v-list-item>
    <v-list-item>
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
    </v-list-item>
    <v-list-item>
      <v-switch
        class="px-2"
        label="Use clade"
        v-model="filters.sampling_target_whole_clade"
        color="primary"
        density="compact"
        :disabled="!filters.sampling_target_taxa?.length"
        hide-details
      >
        <template #append>
          <InlineHelp>
            When enabled, all occurrences of descendant taxa will be included.
          </InlineHelp>
        </template>
      </v-switch>
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
import DatasetPicker from '@/components/datasets/DatasetPicker.vue'
import HabitatPicker from '@/components/occurrence/habitat/HabitatPicker.vue'
import TaxonPicker from '@/components/taxonomy/TaxonPicker.vue'
import TaxonRankPicker from '@/components/taxonomy/TaxonRankPicker'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { reactive, ref, watch } from 'vue'
import SamplingTargetKindFilter from '@/components/occurrence/SamplingTargetKindPicker.vue'
import SiteSamplingStatusFilter from '@/views/location/SiteSamplingStatusFilter.vue'

export type MappingFilters = Overwrite<
  NonNullable<OccurrencesBySiteData['query']>,
  { habitats?: HabitatRecord[] }
>

const sampledTaxa = ref({
  rank: undefined,
  sampledOnly: true
})

const filters = defineModel<MappingFilters>({ default: () => reactive({}) })

watch(
  () => filters.value.sampling_target_taxa,
  (newValue) => {
    if (newValue?.length && !filters.value.sampling_target_kinds?.includes('Taxa')) {
      filters.value.sampling_target_kinds = [...(filters.value.sampling_target_kinds ?? []), 'Taxa']
    }
  }
)
</script>

<style lang="scss">
div.v-list.map-tool-filters .v-list-item .v-list-item__content .v-input {
  margin-top: 5px;
}
</style>
