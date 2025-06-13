<template>
  <LayerOptionsCard
    title="Markers"
    :subtitle="layer.filterType"
    v-model="layer.active"
    prepend-icon="mdi-circle-multiple-outline"
  >
    <template #before-switch>
      <v-menu>
        <template #activator="{ props }">
          <v-btn icon="mdi-dots-vertical" variant="text" :rounded="100" v-bind="props" />
        </template>
        <v-card>
          <v-list density="compact">
            <v-list-item title="Reset" prepend-icon="mdi-restore" @click="emit('reset')" />
            <v-list-item title="Delete" prepend-icon="mdi-delete" @click="emit('delete')" />
          </v-list>
        </v-card>
      </v-menu>
    </template>
    <template #prepend>
      <SvgCircle
        :size="20"
        :fill-color="layer.config.fillColor"
        :stroke-color="withOpacity(layer.config.color)"
      />
    </template>
    <template #title>
      <v-text-field
        v-model="layer.name"
        variant="plain"
        density="compact"
        hide-details
        placeholder="Marker layer"
      />
    </template>
    <template #header>
      <DataFeedPicker
        v-model="layer.dataFeedID"
        class="mx-3 mb-2"
        density="compact"
        hide-details
        clearable
      />
    </template>
    <div class="bg-main">
      <v-list>
        <v-list-item>
          <SiteSamplingStatusFilter
            class="my-1"
            density="compact"
            v-model="layer.filterType"
            hide-details
          />
        </v-list-item>
        <ListItemInput label="Clustered" subtitle="Aggregate marker clusters">
          <v-switch v-model="layer.clustered" hide-details />
        </ListItemInput>
        <ListItemInput label="Radius">
          <v-slider
            :min="1"
            :max="20"
            :step="0.5"
            v-model="layer.config.radius"
            hide-details
            :width="250"
            thumb-label
          />
        </ListItemInput>
        <ListItemInput label="Stroke color" subtitle="Hue and opacity">
          <ColorPickerMenu v-model="layer.config.color" hide-details />
        </ListItemInput>
        <ListItemInput label="Fill color" subtitle="Hue and opacity">
          <ColorPickerMenu v-model="layer.config.fillColor" hide-details />
        </ListItemInput>
        <ListItemInput label="Stroke width">
          <v-slider
            :min="1"
            :max="5"
            v-model="layer.config.weight"
            hide-details
            :width="250"
            thumb-label
          />
        </ListItemInput>
      </v-list>
    </div>
    <template #footer>
      <v-divider />
    </template>
  </LayerOptionsCard>
</template>

<script setup lang="ts">
import { MarkerLayerDefinition } from '@/components/maps/map-layers'
import SvgCircle from '@/components/toolkit/ui/SvgCircle.vue'
import { withOpacity } from '@/functions/color_brewer'
import LayerOptionsCard from '@/views/location/LayerOptionsCard.vue'
import DataFeedPicker from '../occurrence/DataFeedPicker.vue'
import ColorPickerMenu from '../toolkit/ui/ColorPickerMenu.vue'
import ListItemInput from '../toolkit/ui/ListItemInput.vue'
import SiteSamplingStatusFilter from '@/views/location/SiteSamplingStatusFilter.vue'

const layer = defineModel<MarkerLayerDefinition>({ required: true })

const emit = defineEmits<{
  delete: []
  reset: []
}>()
</script>

<style scoped lang="scss"></style>
