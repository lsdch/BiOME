<template>
  <v-list-item title="Targets" prepend-icon="mdi-bullseye">
    <v-chip
      v-if="sampling.target.kind === 'Unknown'"
      class="ma-1"
      :text="sampling.target.kind"
      prepend-icon="mdi-help-circle"
    />
    <v-chip
      v-else-if="sampling.target.kind === 'Community'"
      class="ma-1"
      :text="sampling.target.kind"
      prepend-icon="mdi-grid"
    />
    <TaxonChip
      v-else-if="sampling.target.kind === 'Taxa'"
      class="ma-1"
      v-for="taxon in sampling.target.taxa"
      :taxon
      rounded
    />
  </v-list-item>

  <v-list-item title="Duration" prepend-icon="mdi-update">
    <code>
      {{ Duration.fromObject({ minutes: sampling.duration }).toFormat("hh'h' mm'm'") }}
    </code>
  </v-list-item>
  <v-list-item title="Fixatives" prepend-icon="mdi-snowflake">
    <v-chip v-for="f in sampling.fixatives" :text="f.label" />
    <v-list-item-subtitle v-if="sampling.fixatives?.length" class="font-italic">
      Unknown
    </v-list-item-subtitle>
  </v-list-item>
  <v-list-item title="Methods" prepend-icon="mdi-hook">
    <v-chip v-for="m in sampling.methods" :text="m.label" />
    <v-list-item-subtitle v-if="sampling.methods?.length" class="font-italic">
      Unknown
    </v-list-item-subtitle>
  </v-list-item>
  <v-list-item prepend-icon="mdi-image-filter-hdr-outline">
    <v-list-item title="Habitat" class="px-0 no-padding">
      <v-list-item-subtitle v-if="sampling.habitats?.length" class="font-italic">
        Unknown
      </v-list-item-subtitle>
      <v-chip v-for="h in sampling.habitats" class="ma-1" :text="h.label" :title="h.description" />
    </v-list-item>
    <v-list-item title="Access points" class="px-0 no-padding">
      <v-list-item-subtitle v-if="sampling.access_points?.length" class="font-italic">
        Unknown
      </v-list-item-subtitle>
      <v-chip v-for="access in sampling.access_points" class="ma-1" :text="access" />
    </v-list-item>
  </v-list-item>
  <v-list-item
    title="Comments"
    v-if="sampling.comments"
    prepend-icon="mdi-note-edit-outline"
    :subtitle="sampling.comments"
  />
</template>

<script setup lang="ts">
import { Sampling } from '@/api'
import { Duration } from 'luxon'
import TaxonChip from '../taxonomy/TaxonChip.vue'

const { sampling } = defineProps<{ sampling: Sampling }>()
</script>

<style scoped lang="scss">
.v-list-item.no-padding {
  padding-inline-start: 0px !important;
}
.v-list-item__prepend > .v-icon ~ .v-list-item__spacer {
  width: 16px !important;
}
</style>
