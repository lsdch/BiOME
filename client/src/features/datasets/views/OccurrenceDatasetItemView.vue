<template>
  <v-card
    :title="dataset?.label"
    subtitle="Occurrences dataset"
    prepend-icon="mdi-folder-table"
    class="fill-height d-flex flex-column"
    flat
  >
    <!-- <template #subtitle>
    </template> -->
    <div class="d-flex flex-column flex-grow-1 min-h-0 overflow-y-auto">
      <v-tabs v-model="tab" class="flex-shrink-0">
        <v-tab value="map" prepend-icon="mdi-map">Map</v-tab>
        <v-tab value="samplings" prepend-icon="mdi-map-marker">
          Samplings
          <v-badge
            v-if="datasetSamplings?.length"
            :content="datasetSamplings?.length"
            inline
            color="primary"
            class="ml-1"
          />
        </v-tab>
        <v-tab value="occurrences" prepend-icon="mdi-crosshairs">
          Occurrences
          <v-badge
            v-if="datasetOccurrences?.total_count"
            :content="datasetOccurrences?.total_count"
            inline
            color="success"
            class="ml-1"
          />
        </v-tab>
        <v-tab value="import-batches" prepend-icon="mdi-file-table">
          Import batches
          <v-badge
            v-if="importBatches?.length"
            :content="importBatches?.length"
            inline
            color="purple"
            class="ml-1"
          />
        </v-tab>
      </v-tabs>
      <v-tabs-window v-model="tab" class="flex-grow-1 flex-shrink-0 tabs-window-fill" crossfade>
        <v-tabs-window-item
          value="map"
          key="map"
          :transition="false"
          id="map-tab"
          height="100%"
          class="fill-height"
        >
          <v-sheet height="100%" :min-height="600">
            <DeckGlMap :hexgrid="hexLayer" v-model:zoom="zoom">
              <template #popup="{ selection, mapContainer }">
                <MultiSamplingsPopup
                  v-if="selection.type === 'hexagon' && !!selection.info.object"
                  :data="selection.info.object"
                  :resolution="selection.resolution"
                  :params="selection.params"
                  :attach="mapContainer"
                />

                <MultiSamplingsPopup
                  v-else-if="selection.type === 'marker' && !!selection.info.object"
                  :data="selection.info.object"
                  :resolution="selection.resolution"
                  :params="selection.params"
                  :attach="mapContainer"
                />
              </template>
            </DeckGlMap>
          </v-sheet>
        </v-tabs-window-item>
        <v-tabs-window-item value="samplings" key="samplings" class="fill-height">
          <SamplingWithOccurrencesTable :samplings="datasetSamplings" />
        </v-tabs-window-item>
        <v-tabs-window-item value="occurrences" key="occurrences" class="fill-height">
          <OccurrencesTable :occurrences="datasetOccurrences?.items" />
        </v-tabs-window-item>
        <v-tabs-window-item value="import-batches" key="import-batches" class="fill-height">
          <ImportBatchTable :data="importBatches" class="fill-height" />
        </v-tabs-window-item>
      </v-tabs-window>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import {
  getDatasetByIdOptions,
  listImportBatchesInDatasetOptions,
  listOccurrencesH3Options,
  listOccurrencesOptions,
  listSamplingsWithOccurrencesOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import DeckGlMap from '@/features/cartography/components/DeckGlMap.vue'
import {
  automaticResolution,
  makeHexLayer
} from '@/features/cartography/components/layers-manager/map-layers'
import MultiSamplingsPopup from '@/features/cartography/components/popups/MultiSamplingsPopup.vue'
import { hexgridLayerFromSpec } from '@/features/cartography/composables/hexgrid-layer'
import ImportBatchTable from '@/features/import/components/ImportBatchTable.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import { v } from 'vue-router/dist/index-BN0B0y8a.js'

const { ulid } = defineProps<{
  ulid: string
}>()

const tab = ref<'map' | 'samplings' | 'occurrences'>('map')

const { data: dataset } = useQuery(
  computed(() => ({
    ...getDatasetByIdOptions({ path: { ulid } })
  }))
)

const zoom = ref(0)
const hexgridLayerSpec = ref(makeHexLayer())
watch(zoom, (newZoom) => {
  hexgridLayerSpec.value.resolution = automaticResolution(hexgridLayerSpec.value, newZoom)
})

const { data: hexData } = useQuery(
  computed(() =>
    listOccurrencesH3Options({
      path: { resolution: hexgridLayerSpec.value.resolution },
      query: {
        datasets: [ulid]
      }
    })
  )
)

const hexLayer = computed(() => hexgridLayerFromSpec(hexgridLayerSpec.value, hexData.value ?? []))

const { data: datasetSamplings } = useQuery(
  computed(() => ({
    ...listSamplingsWithOccurrencesOptions({
      query: { datasets: [ulid], sort: 'event_date', sort_direction: 'desc' }
    })
  }))
)

const { data: datasetOccurrences } = useQuery(
  computed(() => ({
    ...listOccurrencesOptions({
      query: { datasets: [ulid], sort: 'code', sort_direction: 'asc' }
    })
  }))
)

const { data: importBatches } = useQuery(
  computed(() => ({
    ...listImportBatchesInDatasetOptions({
      path: { ulid: ulid }
    })
  }))
)
</script>

<style lang="scss">
.tabs-window-fill {
  .v-window__container {
    height: 100%;
  }
}
</style>
