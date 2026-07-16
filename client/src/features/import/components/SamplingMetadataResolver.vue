<template>
  <div class="d-flex">
    <v-tabs v-model="tab" direction="vertical">
      <v-tab prepend-icon="mdi-hook" value="methods">Methods</v-tab>
      <v-tab prepend-icon="mdi-snowflake" value="fixatives">Fixatives</v-tab>
      <v-tab prepend-icon="mdi-map-marker-up" value="access-points">Access points</v-tab>
      <v-tab prepend-icon="mdi-image-filter-hdr" value="habitats">Habitats</v-tab>
    </v-tabs>
    <v-divider vertical></v-divider>
    <v-tabs-window v-model="tab" class="flex-grow-1 fill-height" direction="vertical">
      <v-tabs-window-item value="methods">
        <SamplingMethodResolver :import_id="import_id" />
      </v-tabs-window-item>
      <v-tabs-window-item class="fill-height" value="fixatives">
        <SamplingFixativeResolver :import_id="import_id" />
      </v-tabs-window-item>
      <v-tabs-window-item value="access-points">
        <AccessPointsPicker />
      </v-tabs-window-item>
    </v-tabs-window>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import SamplingMethodResolver from './SamplingMethodResolver.vue'
import AccessPointsPicker from '@/features/occurrences/components/sampling/AccessPointsPicker.vue';
import SamplingFixativeResolver from './SamplingFixativeResolver.vue';

type Tabs = 'methods' | 'fixatives' | 'access-points'
const tab = ref<Tabs>('methods')

const { import_id } = defineProps<{
  import_id: UUID
}>()
</script>

<style scoped lang="scss"></style>
