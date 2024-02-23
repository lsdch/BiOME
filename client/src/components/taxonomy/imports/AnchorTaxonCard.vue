<template>
  <v-card class="h-100 d-flex flex-column">
    <v-toolbar density="compact" color="surface">
      <v-toolbar-title :title="name">
        {{ name }}
      </v-toolbar-title>
      <v-tooltip location="top" text="Drop anchor">
        <template v-slot:activator="{ isActive, props }">
          <v-btn
            v-bind="props"
            :icon="isActive ? 'mdi-close' : 'mdi-anchor'"
            size="small"
            :color="isActive ? 'red' : 'orange'"
          />
        </template>
      </v-tooltip>
      <LinkIconGBIF v-if="GBIF_ID" :GBIF_ID="GBIF_ID" />
    </v-toolbar>
    <v-card-subtitle class="d-flex justify-space-between mb-3">
      <span>{{ rank }}</span>
      <v-spacer />
      <v-chip v-if="authorship" label rounded="xl">{{ authorship }}</v-chip>
    </v-card-subtitle>
    <v-divider />
    <v-card-actions class="justify-space-between">
      <ItemDateChip v-if="meta.created" icon="created" :date="meta.created" color="grey" />
      <ItemDateChip v-if="meta.modified" icon="updated" :date="meta.modified" color="grey" />
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { TaxonDB } from '@/api'
import ItemDateChip from '../../toolkit/ItemDateChip.vue'
import LinkIconGBIF from '../LinkIconGBIF.vue'

defineProps<TaxonDB>()
</script>

<style scoped></style>
