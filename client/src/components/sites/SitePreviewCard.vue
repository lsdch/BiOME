<template>
  <v-card
    class="small-card-title flex-grow-1"
    title="Site"
    prepend-icon="mdi-map-marker"
    v-bind="$attrs"
  >
    <template #subtitle>
      <v-chip
        v-if="site && 'id' in site"
        label
        size="x-small"
        color="purple"
        variant="flat"
        prepend-icon="mdi-link"
        text="From database"
      />
      <v-chip
        v-else-if="!!site"
        label
        size="x-small"
        color="success"
        variant="flat"
        prepend-icon="mdi-plus"
        text="New site"
      />
    </template>
    <template #append v-if="hasEmitBinding">
      <v-btn icon="mdi-pencil" size="small" variant="tonal" @click="emit('edit')" />
    </template>
    <v-divider />
    <v-list v-if="site">
      <v-list-item :title="site.name">
        <v-list-item-subtitle>
          {{ site.locality ?? 'Unspecified locality' }}
          <!-- <CountryChip v-if="site.country" :country="site.country" size="small" /> -->
        </v-list-item-subtitle>
        <template #append>
          <v-chip v-if="site" :text="site.code" class="font-monospace" size="small" />
        </template>
      </v-list-item>
      <v-list-item>
        <!-- <div class="d-flex flex-column"> -->
        <div class="coordinates font-monospace">
          <span class="label"> Lat </span>
          {{ site.coordinates.latitude }}
          <span class="label"> Lng </span>
          {{ site.coordinates.longitude }}
        </div>
        <template #append>
          <CoordPrecisionChip :precision="site.coordinates.precision" size="small" />
        </template>
        <!-- </div> -->
      </v-list-item>
    </v-list>
    <slot></slot>
  </v-card>
</template>

<script setup lang="ts">
import { SiteInput, SiteItem } from '@/api'
import { hasEventListener } from '../toolkit/vue-utils'
import CoordPrecisionChip from './CoordPrecisionChip'

const props = defineProps<{ site?: SiteItem | SiteInput }>()
const emit = defineEmits<{ edit: [] }>()
const hasEmitBinding = hasEventListener('onEdit')
</script>

<style scoped lang="scss"></style>
