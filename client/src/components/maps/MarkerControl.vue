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
    <v-card theme="light">
      <v-list>
        <v-list-item title="Markers">
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
        </v-list-item>
      </v-list>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
export type MarkerLayer = 'cluster' | 'hexgrid'
const marker = defineModel<MarkerLayer>({ required: true })
</script>

<style scoped lang="scss"></style>
