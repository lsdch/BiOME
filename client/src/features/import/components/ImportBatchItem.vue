<template>
  <v-card
    :title="batch?.label"
    prepend-icon="mdi-file-table"
    class="fill-height d-flex flex-column"
    flat
  >
    <template #subtitle>
      <v-chip label prepend-icon="mdi-file-arrow-up-down"></v-chip>
    </template>
    <template #append>
      <v-btn
        prepend-icon="mdi-file-download"
        text="Download raw data"
        variant="tonal"
        @click="downloadRawFile()"
      ></v-btn>
    </template>
    <div class="d-flex flex-column flex-grow-1 min-h-0 overflow-y-auto">
      <!-- <v-container v-if="batch" class="bg-main flex-shrink-0" fluid>
        <v-row>
          <v-col cols="12">
            <v-card>
              <v-card-text>
                {{ batch.description ?? 'No description provided.' }}
              </v-card-text>
              <v-divider></v-divider>
              <v-list>
                <v-list-item>
                  <div class="d-flex ga-1">
                    <v-chip v-for="person in batch.assembled_by" :text="person"></v-chip>
                  </div>
                  <template #append>
                    <span class="text-muted">Assembled by</span>
                  </template>
                </v-list-item>
                <v-list-item
                  :subtitle="
                    DateTime.fromJSDate(batch.created_at!).toLocaleString(DateTime.DATETIME_SHORT)
                  "
                >
                  <template #title>
                    <UserChip :user="batch?.created_by_user!" size="small" />
                  </template>
                  <template #append>
                    <span class="text-muted">Submitted</span>
                  </template>
                </v-list-item>
                <v-list-item
                  :subtitle="
                    DateTime.fromJSDate(batch.completed_at!).toLocaleString(DateTime.DATETIME_SHORT)
                  "
                >
                  <template #title>
                    <UserChip :user="batch?.completed_by_user!" size="small" />
                  </template>
                  <template #append>
                    <span class="text-muted">Completed</span>
                  </template>
                </v-list-item>
              </v-list>
            </v-card>
          </v-col>
        </v-row>
      </v-container> -->
      <v-tabs v-model="tab" class="flex-shrink-0">
        <v-tab value="map">Map</v-tab>
        <v-tab value="samplings">Samplings</v-tab>
        <v-tab value="occurrences">Occurrences</v-tab>
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
        <v-tabs-window-item value="samplings" key="samplings">
          <SamplingWithOccurrencesTable :samplings="batchSamplings" />
        </v-tabs-window-item>
        <v-tabs-window-item value="occurrences" key="occurrences">
          <OccurrencesTable :occurrences="batchOccurrences?.items" />
        </v-tabs-window-item>
      </v-tabs-window>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { DownloadRawFileData } from '@/api'
import {
  getImportBatchWithContentOptions,
  listOccurrencesH3Options,
  listOccurrencesOptions,
  listSamplingsWithOccurrencesOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import { client } from '@/api/gen/client.gen'
import DeckGlMap from '@/features/cartography/components/DeckGlMap.vue'
import {
  automaticResolution,
  makeHexLayer
} from '@/features/cartography/components/layers-manager/map-layers'
import MultiSamplingsPopup from '@/features/cartography/components/popups/MultiSamplingsPopup.vue'
import { hexgridLayerFromSpec } from '@/features/cartography/composables/hexgrid-layer'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'

const { uuid } = defineProps<{
  uuid: UUID
}>()

const tab = ref<'map' | 'samplings' | 'occurrences'>('map')

const { data: batch } = useQuery(
  computed(() => ({
    ...getImportBatchWithContentOptions({ path: { id: uuid } })
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
        batches: [uuid]
      }
    })
  )
)
const hexLayer = computed(() => hexgridLayerFromSpec(hexgridLayerSpec.value, hexData.value ?? []))

const { data: batchSamplings } = useQuery(
  computed(() => ({
    ...listSamplingsWithOccurrencesOptions({
      query: { batches: [uuid], sort: 'event_date', sort_direction: 'desc' }
    })
  }))
)
const { data: batchOccurrences } = useQuery(
  computed(() => ({
    ...listOccurrencesOptions({
      query: { batches: [uuid], sort: 'code', sort_direction: 'asc' }
    })
  }))
)

function downloadRawFile() {
  const url = client.buildUrl<DownloadRawFileData>({
    path: { id: uuid },
    url: '/import-batches/{id}/raw'
  })
  window.open(url, '_blank')
}
</script>

<style lang="scss">
.tabs-window-fill {
  .v-window__container {
    height: 100%;
  }
}
</style>
