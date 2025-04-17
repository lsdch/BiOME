<template>
  <v-menu
    transition="scale-transition"
    location="right center"
    origin="left center"
    :offset="10"
    :close-on-content-click="false"
    attach=".map"
  >
    <template #activator="{ props }">
      <v-card theme="light" flat rounded="0" elevation="5" :width="30" :height="30">
        <v-btn
          icon="mdi-cog"
          variant="text"
          color=""
          :rounded="0"
          v-bind="props"
          size="small"
          :width="30"
          :height="30"
        />
      </v-card>
    </template>
    <v-card
      theme="light"
      :min-width="400"
      title="Marker settings"
      class="small-card-title opacity-80"
      :max-height="300"
      @click.stop
      @mousemove.stop
      @wheel.stop
      @touchmove.stop
      @scroll.stop
      :ripple="false"
    >
      <template #append>
        <v-btn-toggle
          mandatory
          :rounded="10"
          v-model="marker"
          density="compact"
          border="sm"
          color="success"
        >
          <v-btn text="Clusters" size="x-small" density="compact" value="cluster"></v-btn>
          <v-btn text="Hex grid" size="x-small" density="compact" value="hexgrid"></v-btn>
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
