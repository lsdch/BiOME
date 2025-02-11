<template>
  <v-navigation-drawer v-if="!mobile" location="right" permanent>
    <v-tabs v-model="tab" density="compact">
      <v-tab prepend-icon="mdi-creation-outline" value="pinned" stacked size="small" />
      <v-tab prepend-icon="mdi-update" value="latest" stacked size="small" />
    </v-tabs>
    <v-divider />
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="pinned">
        <DatasetFrontDisplay
          qualifier="pinned"
          :query="{ pinned: true, limit: 10, orderBy: 'label' }"
        />
      </v-tabs-window-item>
      <v-tabs-window-item value="latest">
        <DatasetFrontDisplay
          qualifier="latest"
          :query="{ limit: 10, orderBy: 'meta.lastUpdated' }"
        />
      </v-tabs-window-item>
    </v-tabs-window>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDisplay } from 'vuetify'
import DatasetFrontDisplay from './DatasetFrontDisplay.vue'

const { mobile } = useDisplay()

const tab = ref<'pinned' | 'latest'>('pinned')
</script>

<style scoped lang="scss"></style>
