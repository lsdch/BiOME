<template>
  <v-card class="small-card-title flex-grow-1" prepend-icon="mdi-map-marker" v-bind="$attrs">
    <template #title>
      <b v-if="site" class="font-monospace">
        {{ site.code }}
      </b>
      <v-card-title v-else>Site</v-card-title>
    </template>
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
    <template #append>
      <slot name="append" />
    </template>
    <v-divider />
    <v-expand-transition>
      <slot>
        <div v-if="site">
          <v-list>
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
        </div>
      </slot>
    </v-expand-transition>
    <slot name="below" />
    <template #actions v-if="$slots.actions">
      <slot name="actions" />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { SiteItem } from '@/api'
import { hasEventListener } from '../toolkit/vue-utils'
import CoordPrecisionChip from './CoordPrecisionChip'
import { SiteModel } from '@/models'

const props = defineProps<{ site?: SiteItem | SiteModel.SiteFormModel }>()
const emit = defineEmits<{ edit: [] }>()
const hasEmitBinding = hasEventListener('onEdit')
</script>

<style scoped lang="scss"></style>
