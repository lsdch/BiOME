<template>
  <v-list-item prepend-icon="mdi-family-tree">
    <template #append>
      <span class="text-caption text-muted">Target</span>
    </template>
    <v-chip
      v-if="sampling.target.kind === 'Unknown'"
      class="ma-1"
      :text="sampling.target.kind"
      prepend-icon="mdi-help-circle"
      size="small"
    />
    <v-chip
      v-else-if="sampling.target.kind === 'Community'"
      class="ma-1"
      :text="sampling.target.kind"
      prepend-icon="mdi-grid"
      size="small"
    />
    <TaxonChip
      v-else-if="sampling.target.kind === 'Taxa'"
      class="ma-1"
      v-for="taxon in sampling.target.taxa"
      :taxon
      rounded
      size="small"
    />
  </v-list-item>

  <v-list-item prepend-icon="mdi-update">
    <template #append>
      <span class="text-caption text-muted">Duration</span>
    </template>
    <code v-if="sampling.duration">
      {{ Duration.fromObject({ minutes: sampling.duration }).toFormat("hh'h' mm'm'") }}
    </code>
    <v-list-item-subtitle v-else class="font-italic"> Unknown </v-list-item-subtitle>
  </v-list-item>

  <v-list-item prepend-icon="mdi-snowflake">
    <template #append>
      <span class="text-caption text-muted">Fixative(s)</span>
    </template>
    <v-chip v-for="f in sampling.fixatives" :text="f.label" size="small" />
    <v-list-item-subtitle v-if="!sampling.fixatives?.length" class="font-italic">
      Unknown
    </v-list-item-subtitle>
  </v-list-item>

  <v-list-item prepend-icon="mdi-hook">
    <template #append>
      <span class="text-caption text-muted">Methods</span>
    </template>
    <v-chip v-for="m in sampling.methods" :text="m.label" size="small" />
    <v-list-item-subtitle v-if="!sampling.methods?.length" class="font-italic">
      Unknown
    </v-list-item-subtitle>
  </v-list-item>
  <v-list-item prepend-icon="mdi-image-filter-hdr-outline">
    <v-list-item class="px-0 no-padding">
      <template #append>
        <span class="text-caption text-muted">Habitat</span>
      </template>
      <v-list-item-subtitle v-if="!sampling.habitats?.length" class="font-italic">
        Unknown
      </v-list-item-subtitle>
      <v-chip
        v-for="h in sampling.habitats"
        class="ma-1"
        :text="h.label"
        :title="h.description"
        size="small"
      />
    </v-list-item>
    <v-list-item class="px-0 no-padding">
      <template #append>
        <span class="text-caption text-muted">Access points</span>
      </template>
      <v-list-item-subtitle v-if="!sampling.access_points?.length" class="font-italic">
        Unknown
      </v-list-item-subtitle>
      <v-chip v-for="access in sampling.access_points" class="ma-1" :text="access" size="small" />
    </v-list-item>
  </v-list-item>
  <v-list-item
    v-if="sampling.comments"
    prepend-icon="mdi-note-edit-outline"
    title="Comments"
    :subtitle="sampling.comments"
  />
</template>

<script setup lang="ts">
import { HabitatRecord, Sampling } from '@/api'
import { Duration } from 'luxon'
import TaxonChip from '../taxonomy/TaxonChip'
import { Optional } from 'ts-toolbelt/out/Object/Optional'

const { sampling } = defineProps<{
  sampling: Omit<Optional<Sampling, 'id' | 'meta' | 'code'>, 'habitats'> & {
    habitats?: Array<HabitatRecord>
  }
}>()
</script>

<style scoped lang="scss">
.v-list-item.no-padding {
  padding-inline-start: 0px !important;
}
.v-list-item__prepend > .v-icon ~ .v-list-item__spacer {
  width: 16px !important;
}
</style>
