<template>
  <v-sheet :height class="d-flex flex-column">
    <DeckGlMap
      class="flex-grow-1"
      ref="map"
      :zoom="0.1"
      :min-zoom="0.1"
      :pinMarkers="pinMarker ? [pinMarker] : undefined"
      :marker-layers="proximityRadius ? [proximalMarkers] : undefined"
      :auto-fit="proximityRadius || item?.coordinates?.precision"
      @drag-pin="(_index, { latitude, longitude }) => setCoordinates(latitude, longitude)"
    >
      <template #popup="{ selection }">
        <SiteClusterPopup
          v-if="selection.type === 'marker' && !!selection.info.object"
          :data="selection.info.object"
          :resolution="12"
          :params="selection.params"
        >
          <template #append-item v-if="'distance_meters' in selection.info.object">
            <v-divider />
            <v-list-item title="Distance" prepend-icon="mdi-ruler" class="text-muted">
              <template #append>
                <DistanceDisplay :distance="selection.info.object.distance_meters" />
              </template>
            </v-list-item>
          </template>
        </SiteClusterPopup>
      </template>
    </DeckGlMap>

    <v-list-item>
      <ProximityRadiusSlider v-model="proximitySliderValue" v-model:radius="proximityRadius" />
      <template #prepend>
        <v-progress-circular
          indeterminate
          color="primary"
          v-if="hasValidCoordinates && isPending"
          size="small"
          class="mr-2"
        />
        <OccurrencesOverviewDialog
          v-else
          :data="
            nearbyOccurrences?.filter(
              ({ distance_meters }) => distance_meters <= proximityRadius
            ) ?? []
          "
          :max-width="1400"
        >
          <template #samplings-table="{ samplings }">
            <SamplingWithOccurrencesTable
              :samplings="samplings"
              :extend-headers="[
                {
                  position: 1,
                  key: 'distance_meters',
                  title: 'Distance'
                }
              ]"
            >
              <template #item.distance_meters="{ value }">
                <DistanceDisplay :distance="value" />
              </template>
            </SamplingWithOccurrencesTable>
          </template>
          <template #occurrences-table="{ occurrences }">
            <OccurrencesTable
              with-site
              :occurrences="occurrences"
              :extend-headers="[
                {
                  position: 1,
                  key: 'sampling.distance_meters',
                  title: 'Distance'
                }
              ]"
            >
              <template #item.sampling.distance_meters="{ value }">
                <DistanceDisplay :distance="value" />
              </template>
            </OccurrencesTable>
          </template>
          <template #sampled-taxa-table="{ occurrences }">
            <SampledTaxaTable
              :occurrences="occurrences"
              :extend-headers="[
                {
                  position: 1,
                  key: 'distance_meters',
                  title: 'Distance'
                }
              ]"
            >
              <template #item.distance_meters="{ value }: { value: number }">
                <DistanceDisplay :distance="value" />
              </template>
            </SampledTaxaTable>
          </template>
          <template #prepend-body>
            <v-card-text>
              <ProximityRadiusSlider
                label="Radius"
                v-model="proximitySliderValue"
                v-model:radius="proximityRadius"
              />
            </v-card-text>
          </template>
          <template #activator="{ props }">
            <v-btn
              icon="mdi-list-box"
              size="small"
              variant="plain"
              v-tooltip="`See list of sites and occurrences within radius`"
              v-bind="props"
              :disabled="!hasValidCoordinates || nearbyCells?.length === 0 || proximityRadius === 0"
            />
          </template>
        </OccurrencesOverviewDialog>
      </template>
    </v-list-item>
  </v-sheet>
</template>

<script setup lang="tsx" generic="Item extends DeepPartial<ItemWithCoordinates>">
import { H3CellWithRichnessAndDistance } from '@/api/adapters.ts'
import {
  listOccurrencesAtProximityOptions,
  listSamplingsH3AtProximityOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import OccurrencesOverviewDialog from '@/features/occurrences/components/tables/OccurrencesOverviewDialog.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SampledTaxaTable from '@/features/occurrences/components/tables/SampledTaxaTable.vue'
import SamplingWithOccurrencesTable from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import { formatDistance } from '@/lib/distances'
import { useQuery } from '@tanstack/vue-query'
import { circle } from '@turf/turf'
import { cellToLatLng } from 'h3-js'
import { Map, MarkerOptions } from 'maplibre-gl'
import { computed, ref, useTemplateRef, watch } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import { Coordinates, ItemWithCoordinates } from '../coordinates'
import DeckGlMap from './DeckGlMap.vue'
import { makeMarkerLayer, markerLayerFromSpec, PinMarker } from './layers-manager/map-layers'
import SiteClusterPopup from './popups/MultiSamplingsPopup.vue'
import ProximityRadiusSlider from './ProximityRadiusSlider.vue'

const DistanceDisplay = (props: { distance: number }) => {
  return <span class="font-monospace">{formatDistance(props.distance)}</span>
}

const item = defineModel<Item>('item')

const { height, markerOptions, excludeIds } = defineProps<{
  height?: number
  markerOptions?: MarkerOptions
  excludeIds?: UUID[]
}>()

function setCoordinates(latitude: number, longitude: number) {
  if (!item.value) return

  item.value.coordinates = {
    latitude,
    longitude,
    precision: item.value.coordinates?.precision
  }
}

const proximitySliderValue = ref<number>(0)
const proximityRadius = ref<number>(0)

const map = useTemplateRef<ComponentExposed<typeof DeckGlMap>>('map')

function displayProximityRadius(radius: number) {
  if (
    !map.value ||
    !proximityRadius.value ||
    !item.value ||
    !Coordinates.isValidCoordinates(item.value.coordinates)
  )
    return
  const c = circle([item.value.coordinates.longitude, item.value.coordinates.latitude], radius, {
    steps: 64,
    units: 'meters'
  })
  map.value.instance?.addSource(`proximity-radius`, {
    type: 'geojson',
    data: c
  })
  map.value?.instance?.addLayer({
    id: `proximity-radius-outline`,
    type: 'line',
    source: `proximity-radius`,
    paint: {
      'line-color': '#068DB1',
      'line-width': 3,
      'line-dasharray': [1, 1]
    }
  })
}

function removeProximityRadius(mapInstance: Map) {
  mapInstance.getLayer(`proximity-radius`) && mapInstance.removeLayer(`proximity-radius`)
  mapInstance.getLayer(`proximity-radius-outline`) &&
    mapInstance.removeLayer(`proximity-radius-outline`)
  mapInstance.getSource(`proximity-radius`) && mapInstance.removeSource(`proximity-radius`)
}

watch(
  proximityRadius,
  () => {
    if (!map.value?.instance) return
    removeProximityRadius(map.value.instance)
    displayProximityRadius(proximityRadius.value)
  },
  { immediate: true }
)

const hasValidCoordinates = computed(() => Coordinates.isValidCoordinates(item.value?.coordinates))

const {
  data: nearbyCells,
  isPending,
  error
} = useQuery(
  computed(() => ({
    enabled: hasValidCoordinates.value,
    ...listSamplingsH3AtProximityOptions({
      query: {
        latitude: item.value?.coordinates?.latitude ?? 0,
        longitude: item.value?.coordinates?.longitude ?? 0,
        radius_meters: 100_000,
        exclude_ids: excludeIds
      }
    })
  }))
)

const { data: nearbyOccurrences } = useQuery(
  computed(() => ({
    enabled: hasValidCoordinates.value,
    ...listOccurrencesAtProximityOptions({
      query: {
        latitude: item.value?.coordinates?.latitude ?? 0,
        longitude: item.value?.coordinates?.longitude ?? 0,
        radius_meters: 100_000,
        exclude_ids: excludeIds
      }
    })
  }))
)

const pinMarker = computed<PinMarker<Item> | undefined>(() =>
  item.value && Coordinates.isValidCoordinates(item.value.coordinates)
    ? {
        data: item.value,
        coordinates: item.value.coordinates,
        options: markerOptions
      }
    : undefined
)

type CellMarker = H3CellWithRichnessAndDistance & {
  coordinates: Coordinates
}

const proximalMarkersSpec = ref(
  makeMarkerLayer('Proximal sites', {
    ready: true
  })
)
const proximalMarkers = computed(() => {
  return markerLayerFromSpec(
    proximalMarkersSpec.value,
    nearbyCells.value
      ?.filter(({ distance_meters }) => distance_meters <= proximityRadius.value)
      .map<CellMarker>((cell) => {
        const [latitude, longitude] = cellToLatLng(cell.h3_index)
        return {
          coordinates: { latitude, longitude },
          ...cell
        }
      }) ?? [],
    {
      radius(item) {
        return item.occurrences_count
      },
      getText: (item) => (item.occurrences_count > 1 ? `${item.occurrences_count}` : ''),
      getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude]
    }
  )
})
</script>

<style scoped lang="scss"></style>
