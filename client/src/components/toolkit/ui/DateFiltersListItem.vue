<template>
  <v-list-item
    title="Sampling date"
    lines="two"
    :subtitle="
      model.enabled
        ? `
        [${
          (model.from ? CompositeDate.toString(model.from, model.precision) : undefined) ?? ' -'
        } ; ${(model.to ? CompositeDate.toString(model.to, model.precision) : undefined) ?? '- '}]`
        : undefined
    "
    class="text-muted"
  >
    <template #append>
      <div class="d-flex ga-3 align-center">
        <v-select
          v-if="model.enabled"
          v-model="model.precision"
          label="Precision"
          :items="['day', 'month', 'year']"
          hide-details
          density="compact"
        ></v-select>
        <v-select
          v-if="model.enabled"
          label="Mode"
          :items="[
            { title: 'Fixed', value: false },
            { title: 'Range', value: true }
          ]"
          v-model="model.is_range"
          hide-details
          density="compact"
        ></v-select>
        <v-switch v-model="model.enabled" hide-details density="compact" color="primary"></v-switch>
      </div>
    </template>
  </v-list-item>
  <template v-if="model.enabled">
    <v-list-item>
      <template #prepend>
        <span class="text-muted" :style="{ width: '50px' }">From:</span>
      </template>
      <DateFilter
        label=""
        :max-date="model.is_range ? model.to : undefined"
        v-model="model.from"
        :precision="model.precision ?? 'day'"
      />
    </v-list-item>
    <v-list-item v-if="model.is_range">
      <template #prepend>
        <span class="text-muted" :style="{ width: '50px' }">To:</span>
      </template>
      <DateFilter
        label=""
        v-model="model.to"
        :min-date="model.from"
        :precision="model.precision ?? 'day'"
      />
    </v-list-item>
    <ListItemInput title="Include unknown dates" class="text-muted">
      <v-switch
        v-model="model.include_unknown"
        color="primary"
        hide-details
        density="compact"
      ></v-switch>
    </ListItemInput>
  </template>
</template>

<script setup lang="ts">
import { CompositeDate, DateFilterParams } from '@/api'
import { DateFilters } from '@/features/cartography/components/layers-manager/map-layers'
import ListItemInput from './ListItemInput.vue'
import DateFilter from './DateFilter.vue'

const model = defineModel<DateFilters>({ default: () => ({}) })
</script>

<style scoped lang="scss"></style>
