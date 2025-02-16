<template>
  <v-menu
    transition="scale-transition"
    location="bottom left"
    origin="bottom left"
    :close-on-content-click="false"
  >
    <template #activator="{ props }">
      <v-card theme="light">
        <v-btn icon="mdi-cog" variant="plain" color="" :rounded="30" v-bind="props"></v-btn>
      </v-card>
    </template>
    <v-card theme="light" :min-width="300" title="Marker settings">
      <template #append>
        <v-btn-toggle
          mandatory
          :rounded="10"
          v-model="marker"
          density="compact"
          border="sm"
          color="success"
        >
          <v-btn text="Clusters" size="small" density="compact" value="cluster"></v-btn>
          <v-btn text="Hex grid" size="small" density="compact" value="hexgrid"></v-btn>
        </v-btn-toggle>
      </template>
      <v-tabs-window v-model="marker">
        <v-tabs-window-item value="cluster">
          <v-card-text>
            <!-- <ClusterConfigPanel v-model="clusterConfig" /> -->
          </v-card-text>
        </v-tabs-window-item>
        <v-tabs-window-item value="hexgrid">
          <v-divider />
          <HexgridConfigPanel v-model="hexgridConfig" />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import HexgridConfigPanel, { HexgridConfig } from './HexgridConfigPanel.vue'

export type MarkerLayer = 'cluster' | 'hexgrid'
const marker = defineModel<MarkerLayer>({ required: true })
const hexgridConfig = defineModel<HexgridConfig>('hexgrid', { required: true })
</script>

<style scoped lang="scss"></style>
