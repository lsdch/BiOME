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
    <DateFiltersListItem v-model="filters.date" />
    <!-- <v-menu>
      <template #activator="{ props }">
        <v-list-item title="Sampling date">
          <template #append>
            <v-switch hide-details density="compact" color="primary"></v-switch>
            <v-btn v-bind="props" icon="mdi-calendar" />
          </template>
        </v-list-item>
      </template>
      <v-card>
        <div class="d-flex ga-3 align-center">
          <v-select
            v-model="datePrecision"
            :items="['day', 'month', 'year']"
            hide-details
            density="compact"
          ></v-select>
          <v-select
            :items="[
              { title: 'Fixed', value: false },
              { title: 'Range', value: true }
            ]"
            v-model="dateRange"
            hide-details
            density="compact"
          ></v-select>
        </div>
        <v-list-item v-if="!dateRange">
          <DateFilter :precision="datePrecision" />
        </v-list-item>
        <template v-else>
          <v-list-item>
            <template #prepend>
              <span class="text-muted" :style="{ width: '50px' }">From:</span>
            </template>
            <DateFilter label="" :precision="datePrecision" />
          </v-list-item>
          <v-list-item>
            <template #prepend>
              <span class="text-muted" :style="{ width: '50px' }">To:</span>
            </template>
            <DateFilter label="" :precision="datePrecision" />
          </v-list-item>
        </template>
      </v-card>
    </v-menu> -->
    <!-- <v-list-item
      title="Sampling date"
      lines="two"
      :subtitle="
        filters.date.enabled
          ? `
        [${
          (filters.date.from
            ? CompositeDate.toString(filters.date.from, filters.date.precision)
            : undefined) ?? ' -'
        } ; ${
          (filters.date.to
            ? CompositeDate.toString(filters.date.to, filters.date.precision)
            : undefined) ?? '- '
        }]`
          : undefined
      "
      class="text-muted"
    >
      <template #append>
        <div class="d-flex ga-3 align-center">
          <v-select
            v-if="filters.date.enabled"
            v-model="filters.date.precision"
            label="Precision"
            :items="['day', 'month', 'year']"
            hide-details
            density="compact"
          ></v-select>
          <v-select
            v-if="filters.date.enabled"
            label="Mode"
            :items="[
              { title: 'Fixed', value: false },
              { title: 'Range', value: true }
            ]"
            v-model="filters.date.is_range"
            hide-details
            density="compact"
          ></v-select>
          <v-switch
            v-model="filters.date.enabled"
            hide-details
            density="compact"
            color="primary"
          ></v-switch>
        </div>
      </template>
    </v-list-item>
    <template v-if="filters.date.enabled">
      <v-list-item>
        <template #prepend>
          <span class="text-muted" :style="{ width: '50px' }">From:</span>
        </template>
        <DateFilter
          label=""
          :max-date="filters.date.is_range ? filters.date.to : undefined"
          v-model="filters.date.from"
          :precision="filters.date.precision ?? 'day'"
        />
      </v-list-item>
      <v-list-item v-if="filters.date.is_range">
        <template #prepend>
          <span class="text-muted" :style="{ width: '50px' }">To:</span>
        </template>
        <DateFilter
          label=""
          v-model="filters.date.to"
          :min-date="filters.date.from"
          :precision="filters.date.precision ?? 'day'"
        />
      </v-list-item>
      <ListItemInput title="Include unknown dates" class="text-muted">
        <v-switch
          v-model="filters.date.include_unknown"
          color="primary"
          hide-details
          density="compact"
        ></v-switch>
      </ListItemInput>
    </template> -->
  </v-list>
</template>

<script setup lang="ts">
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import DatasetPicker from '@/features/datasets/components/DatasetPicker.vue'
import HabitatPicker from '@/features/registries/components/HabitatPicker.vue'
import TaxonFilterPicker from '@/features/taxonomy/components/TaxonFilterPicker.vue'
import { reactive, ref } from 'vue'
import { MappingFilters } from './map-layers'
import ImportBatchPicker from '@/features/import/components/ImportBatchPicker.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'
import DateFilter from '@/components/toolkit/ui/DateFilter.vue'
import { CompositeDate } from '@/api'
import DateFiltersListItem from '@/components/toolkit/ui/DateFiltersListItem.vue'

const filters = defineModel<MappingFilters>({ default: () => reactive({}) })

type QueryMode = 'occurrences' | 'samplings'
const queryMode = defineModel<QueryMode>('mode', { default: 'occurrences' })
</script>

<style lang="scss">
div.v-list.map-tool-filters .v-list-item .v-list-item__content .v-input {
  margin-top: 5px;
}
</style>
