<template>
  <v-list>
    <v-list-subheader title="Markers"> </v-list-subheader>
    <ListItemInput label="Marker tooltips" subtitle="Display site name or code on hover">
      <v-switch v-model="markerOptions.tooltips" color="primary" hide-details />
    </ListItemInput>
    <ListItemInput
      label="Cluster radius scale"
      subtitle="Scale factor applied to cluster radius based on the number of sites in the cluster"
      lines="three"
    >
      <v-slider
        color="primary"
        :min="0"
        :max="2"
        :step="0.1"
        :width="200"
        :ticks="[0, 0.5, 1, 2]"
        show-ticks="always"
        v-model="markerOptions.cluster.radiusScaleFactor"
      ></v-slider>
    </ListItemInput>
    <ListItemInput
      label="Cluster zoom threshold"
      subtitle="Zoom level at which marker clusters are displayed"
      lines="three"
    >
      <v-slider
        color="primary"
        :min="5"
        :max="15"
        :step="1"
        :width="200"
        v-model="markerOptions.cluster.labelZoomThreshold"
        :ticks="[5, 10, 15]"
        show-ticks="always"
      ></v-slider>
    </ListItemInput>
    <v-divider></v-divider>
    <v-list-subheader title="Presets"> </v-list-subheader>
    <ListItemInput
      label="Save map configuration"
      subtitle="Restore last map configuration on next visit"
    >
      <v-switch color="primary" hide-details />
    </ListItemInput>
    <CardDialog v-if="userStore.isGranted('Contributor')" title="Map presets">
      <template #append>
        <v-switch
          v-if="userStore.isGranted('Maintainer')"
          v-model="showAllPresets"
          label="Maintainer view"
          hide-details
          color="warning"
          v-tooltip="'Display all registered presets'"
        />
      </template>
      <template #activator="{ props }">
        <v-list-item
          title="Manage presets"
          prepend-icon="mdi-folder-star-multiple"
          v-bind="props"
        />
      </template>
      <MapPresetManager :all="showAllPresets" />
    </CardDialog>
  </v-list>
</template>

<script setup lang="ts">
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import MapPresetManager from './map-presets/MapPresetManager.vue'
import { useUserStore } from '@/stores/user'
import { reactive, ref } from 'vue'
import { GlobalMarkerOptions } from './DeckGlMap.vue'

const userStore = useUserStore()

const markerOptions = defineModel<GlobalMarkerOptions>('markerOptions', {
  default: reactive<GlobalMarkerOptions>({
    cluster: {
      radiusScaleFactor: 0.5,
      labelZoomThreshold: 8
    },
    tooltips: true
  })
})

const showAllPresets = ref(false)
</script>

<style scoped lang="scss"></style>
