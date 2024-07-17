<template>
  <v-card class="h-100 d-flex flex-column">
    <template #title>
      {{ name }}
    </template>
    <template #subtitle>
      {{ authorship }}
    </template>
    <template #append>
      <v-tooltip location="top" text="Update from GBIF">
        <template #activator="{ isActive, props }">
          <v-btn
            v-bind="props"
            variant="text"
            :icon="isActive ? 'mdi-refresh' : 'mdi-pin'"
            size="small"
            :color="isActive ? 'primary' : 'warning'"
          />
        </template>
      </v-tooltip>
      <LinkIconGBIF v-if="GBIF_ID" :GBIF_ID="GBIF_ID" variant="text" />
    </template>
    <template #text>
      <div class="flex-grow-0">
        <v-chip :text="rank" />
      </div>
    </template>
    <v-divider />
    <v-card-actions class="flex-column align-start">
      <ItemDateChip v-if="meta.created" icon="created" :date="meta.created" color="grey" />
      <ItemDateChip v-if="meta.modified" icon="updated" :date="meta.modified" color="grey" />
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { Taxon } from '@/api'
import ItemDateChip from '../../toolkit/ItemDateChip.vue'
import LinkIconGBIF from '../LinkIconGBIF.vue'

defineProps<Taxon>()
</script>

<style scoped></style>
