<template>
  <HexLayerCard :layer="hexLayer" flat class="ma-1" />
  <v-toolbar density="compact" class="bg-surface">
    <template #prepend v-if="!!markerLayers.length">
      <ConfirmDialog
        title="Delete all marker layers?"
        @confirm="markerLayers.splice(0, markerLayers.length)"
      >
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            prepend-icon="mdi-delete"
            text="Delete marker layers"
            variant="text"
            color="primary"
            rounded="md"
            size="small"
            class="text-label-small"
          ></v-btn>
        </template>
      </ConfirmDialog>
    </template>

    <template #append>
      <v-btn
        prepend-icon="mdi-plus-circle-multiple"
        text="Add marker layer"
        variant="elevated"
        color="primary"
        rounded="md"
        size="small"
        class="text-label-small"
        @click="addLayer()"
      ></v-btn>
    </template>
  </v-toolbar>
  <!-- <v-card
    title="Context"
    class="small-card-title"
    subtitle="Default settings for all new data feeds"
    flat
    :rounded="0"
  >
    <template #append>
      <v-switch v-model="contextEnabled" color="primary" hide-details></v-switch>
    </template>
    <v-expand-transition>
      <div v-if="contextEnabled" class="bg-main">
        <DataFeedsContextPicker />
      </div>
    </v-expand-transition>
  </v-card>
  <v-divider class="mb-3" :thickness="4" /> -->
  <div class="bg-main fill-height">
    <v-item-group selected-class="ma-1" v-model="expandedCard">
      <VueDraggable
        v-model="markerLayers"
        item-key="id"
        :disabled="markerLayers.length <= 1"
        handle=".handle"
        class="d-flex flex-column ga-2 px-1 py-2 bg-main fill-height"
        :animation="150"
      >
        <v-item
          v-for="(_, i) in markerLayers"
          #="{ isSelected, toggle, selectedClass }"
          :key="markerLayers[i].id"
        >
          <MarkerLayerCard
            v-model:layer="markerLayers[i]"
            :key="markerLayers[i].id"
            :index="i"
            :class="selectedClass"
            :expanded="isSelected"
            @delete="markerLayers.splice(i, 1)"
            @update:expanded="toggle?.()"
            @draghandle-down="expandedCard = undefined"
          />
        </v-item>
      </VueDraggable>
    </v-item-group>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import { HexLayerSpec, makeHexLayer, makeMarkerLayer, MarkerLayerSpec } from './map-layers'
import MarkerLayerCard from './MarkerLayerCard.vue'
import HexLayerCard from './HexLayerCard.vue'
import ConfirmDialog from '@/components/toolkit/ui/ConfirmDialog.vue'

const hexLayer = defineModel<HexLayerSpec>('hex-layer', {
  default: () => {
    return reactive(makeHexLayer())
  }
})
const markerLayers = defineModel<MarkerLayerSpec[]>('marker-layers', {
  default: () => reactive([])
})

function addLayer() {
  markerLayers.value.push(makeMarkerLayer())
}

const expandedCard = ref<unknown[]>()
</script>

<style scoped lang="scss"></style>
