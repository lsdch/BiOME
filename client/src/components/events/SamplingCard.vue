<template>
  <v-card class="sampling-card border-s-lg border-success border-opacity-100" rounded="sb-0">
    <template #prepend>
      <code class="top-tag bg-success"> {{ cornerTag }} </code>
    </template>
    <v-card-text>
      <v-list density="compact">
        <v-list-item title="Targets" prepend-icon="mdi-bullseye">
          <v-chip
            v-if="sampling.target.kind === 'Unknown'"
            :text="sampling.target.kind"
            prepend-icon="mdi-help-circle"
          />
          <v-chip
            v-else-if="sampling.target.kind === 'Community'"
            :text="sampling.target.kind"
            prepend-icon="mdi-grid"
          />
          <v-chip
            v-else-if="sampling.target.kind === 'Taxa'"
            v-for="taxon in sampling.target.target_taxa"
            :text="taxon.name"
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
        </v-list-item>
        <v-list-item title="Methods" prepend-icon="mdi-hook">
          <v-chip v-for="m in sampling.methods" :text="m.label" />
        </v-list-item>
        <v-list-item title="Habitat">
          <v-chip v-for="h in sampling.habitats" :text="h.label" :title="h.description" />
        </v-list-item>
        <v-list-item title="Access points">
          <v-chip v-for="access in sampling.access_points" :text="access" />
        </v-list-item>
        <v-list-item title="Comments" :subtitle="sampling.comments" />
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Sampling } from '@/api'
import { Duration } from 'luxon'

defineProps<{ sampling: Sampling; cornerTag: string }>()
</script>

<style lang="scss">
.sampling-card {
  .v-card-item {
    padding-top: 0px;
    padding-left: 0px;
    .top-tag {
      height: 45px;
      padding: 10px;
      border-bottom-right-radius: 25%;
      font-weight: bold;
    }
  }
}
</style>
