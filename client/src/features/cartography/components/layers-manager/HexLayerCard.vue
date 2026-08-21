<template>
  <v-card class="small-card-title" v-bind="$attrs">
    <template #prepend>
      <v-icon icon="mdi-hexagon-multiple" class="text-muted" />
    </template>
    <template #title>
      <v-text-field
        v-model="layer.name"
        density="compact"
        hide-details
        :placeholder="title"
        variant="underlined"
      />
    </template>
    <template #append>
      <v-progress-circular
        v-if="remote.isPending.value"
        indeterminate
        size="small"
        color="warning"
      ></v-progress-circular>
      <v-btn
        v-else
        icon="mdi-reload"
        size="small"
        variant="text"
        color="grey"
        v-tooltip="`Reload data`"
        @click="remote.refetch()"
      ></v-btn>
      <v-switch v-model="layer.active" color="primary" hide-details></v-switch>
      <v-btn
        :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
        color=""
        variant="text"
        :rounded="100"
        size="small"
        @click="() => (expanded = !expanded)"
      />
    </template>
    <v-expand-transition>
      <div v-show="expanded" class="">
        <v-divider />
        <v-tabs v-model="tab" density="compact">
          <v-tab prepend-icon="mdi-database" value="data">Data</v-tab>
          <v-tab prepend-icon="mdi-hexagon-multiple-outline" value="layer">Style</v-tab>
          <v-tab prepend-icon="mdi-circle-multiple-outline" value="markers">Markers</v-tab>
          <v-tab prepend-icon="mdi-chart-bar" value="stats" :disabled="remote.isFetching.value">
            <v-icon icon="mdi-download" size="small" />
          </v-tab>
          <!-- <v-tab
            slim
            prepend-icon="mdi-download"
            value="export"
            :disabled="remote.isFetching.value"
          >
          </v-tab> -->
        </v-tabs>
        <v-divider></v-divider>
        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="data">
            <!-- <SiteSamplingStatusFilter v-model="layer.include_sites" /> -->
            <!-- <v-confirm-edit v-model="layer.filters" @save="register()">
              <template #default="{ model, actions, isPristine }"> -->
            <LayerDataFeed v-model="layer.filters" v-model:mode="layer.mode" class="pa-2" />
            <!-- <div v-if="!isPristine" class="d-flex justify-end">
                  <component :is="actions"></component>
                </div> -->
            <!-- </template> -->
            <!-- </v-confirm-edit> -->
          </v-tabs-window-item>
          <v-tabs-window-item value="layer">
            <HexgridLayerStylePanel v-model="layer" />
          </v-tabs-window-item>
          <v-tabs-window-item value="markers">
            <MarkerLayerStylePanel v-model="layer.markers">
              <template #prepend-item>
                <!-- <v-list-subheader title="Markers can show up when zooming in on the hexgrid layer.">
                </v-list-subheader> -->
                <ListItemInput
                  label="Zoom threshold"
                  :subtitle="
                    layer.markers.minZoomMode === 'manual'
                      ? layer.markers.minZoom === 0
                        ? 'Always'
                        : `Zoom ${layer.markers.minZoom}+`
                      : layer.markers.minZoomMode === 'auto'
                        ? 'Auto'
                        : 'Never'
                  "
                >
                  <v-chip-group v-model="layer.markers.minZoomMode" color="primary">
                    <v-chip value="auto">Auto</v-chip>
                    <v-chip value="never">Never</v-chip>
                    <v-chip value="manual">Manual</v-chip>
                  </v-chip-group>
                  <InlineHelp>
                    The minimum zoom level at which cells are faded to display their content as
                    markers.
                  </InlineHelp>
                </ListItemInput>
                <v-list-item v-if="layer.markers.minZoomMode === 'manual'">
                  <v-slider
                    v-model="layer.markers.minZoom"
                    class="py-1 px-3"
                    :min="0"
                    :max="18"
                    :step="1"
                    hide-details
                    :track-size="2"
                    density="compact"
                  >
                  </v-slider>
                </v-list-item>
                <v-divider />
              </template>
            </MarkerLayerStylePanel>
          </v-tabs-window-item>
          <v-tabs-window-item value="stats">
            <OccurrencesStats :cells="remote.data.value" />
            <v-divider></v-divider>
            <v-card-text>
              <div class="d-flex justify-end">
                <ExportDialogServer
                  @submit="(options, suffix) => exportData(options)"
                  :disabled="remote.isFetching.value"
                >
                  <template #activator="{ props }">
                    <v-btn
                      prepend-icon="mdi-download"
                      text="Export data"
                      variant="tonal"
                      color="primary"
                      v-bind="props"
                    ></v-btn>
                  </template>
                  <!-- <v-btn icon="mdi-download" variant="text" color="primary" @click="exportData()" /> -->
                </ExportDialogServer>
              </div>
            </v-card-text>
          </v-tabs-window-item>
        </v-tabs-window>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script setup lang="ts">
// import { occurrencesBySiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { ExportSamplingsWithOccurrencesData, ListOccurrencesH3Data } from '@/api/adapters.ts'
import {
  listOccurrencesH3Options,
  listSamplingsH3Options
} from '@/api/gen/@tanstack/vue-query.gen.ts'
import { client } from '@/api/gen/client.gen.ts'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'
import OccurrencesStats from '@/features/occurrences/components/OccurrencesStats.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, onMounted, ref, toValue, watch } from 'vue'
import HexgridLayerStylePanel from './HexgridLayerStylePanel.vue'
import { useLayerData } from './layer-data'
import LayerDataFeed from './LayerDataFeed.vue'
import { automaticResolution, HexLayerSpec, mappingFiltersToQuery } from './map-layers'
import MarkerLayerStylePanel from './MarkerLayerStylePanel.vue'
import InlineHelp from '@/components/toolkit/ui/InlineHelp.vue'
import ExportDialogServer, {
  ExportOptions
} from '@/components/toolkit/ui/exports/ExportDialogServer.vue'

const layer = defineModel<HexLayerSpec>('layer', { required: true })
const { zoom } = defineProps<{ zoom: number }>()

const title = computed(() => (layer.value.name?.length ? layer.value.name : 'Cells'))

const tab = ref<'data' | 'layer' | 'markers' | 'stats'>('data')

const expanded = defineModel<boolean>('expanded', { default: false })

const { layerData, registerLayer, deleteLayer } = useLayerData()

const emit = defineEmits<{
  delete: []
  'draghandle-down': []
}>()

watch(
  [() => zoom, () => [layer.value.resolutionMode, layer.value.markers.resolutionMode]],
  (newValues) => {
    const [newZoom] = newValues
    if (layer.value.resolutionMode === 'auto') {
      layer.value.resolution = automaticResolution(layer.value, newZoom)
    }
    if (layer.value.markers.resolutionMode === 'auto') {
      layer.value.markers.resolution = automaticResolution(layer.value.markers, newZoom)
    }
  },
  { immediate: true }
)

async function exportData(options: ExportOptions) {
  const queryClient = client
  const requestUrl = queryClient.buildUrl<ExportSamplingsWithOccurrencesData>({
    url: '/occurrences/export',
    query: {
      filename: options.filename,
      format: options.format,
      delimiter: options.csvOptions.delimiter,
      quoteChar: options.csvOptions.quoteChar,
      taxa: layer.value.filters.taxa,
      whole_clade: layer.value.filters.whole_clade,
      countries: layer.value.filters.countries,
      batches: layer.value.filters.batches
    }
  })
  window.location.assign(requestUrl)
}

const remote = useQuery(
  computed(() => {
    switch (layer.value.mode) {
      case 'samplings':
        const samplingOptions = listSamplingsH3Options({
          path: { resolution: layer.value.resolution },
          query: mappingFiltersToQuery(layer.value.filters, layer.value.mode)
        })
        return {
          enabled: layer.value.active,
          initialData: [],
          ...samplingOptions
        }
      case 'occurrences':
        const options = listOccurrencesH3Options({
          path: { resolution: layer.value.resolution },
          query: mappingFiltersToQuery(
            layer.value.filters,
            layer.value.mode
          ) as ListOccurrencesH3Data['query']
        })
        return {
          enabled: true, //layer.value.active,
          initialData: [],
          ...options
          //   taxa: layer.value.filters.taxa,
          //   whole_clade: layer.value.filters.whole_clade,
          //   countries: layer.value.filters.countries,
          //   batches: layer.value.filters.batches
          // }
        }
    }
  })
)

watch(
  () => layer.value.id,
  (current, previous) => {
    if (previous !== current) {
      deleteLayer(previous)
      register()
    }
  }
)

function register() {
  console.debug(`Registering hex layer ${layer.value.id}`)
  registerLayer(layer.value, remote)
}

onMounted(() => {
  register()
})
</script>

<style scoped lang="scss"></style>
