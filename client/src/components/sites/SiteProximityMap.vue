<template>
  <v-progress-linear v-if="hasValidCoords && isPending" indeterminate></v-progress-linear>
  <SitesMap
    ref="map"
    :items="hasValidCoords ? data : undefined"
    :auto-fit="false"
    clustered
    :center="hasValidCoords ? [coords!.latitude!, coords!.longitude!] : undefined"
    :zoom="hasValidCoords ? 10 : 0"
    :min-zoom="1"
    hide-marker-control
  >
    <!-- Coordinates marker -->
    <LMarker
      v-if="hasValidCoords"
      :lat-lng="{ lat: coords!.latitude!, lng: coords!.longitude! }"
      :draggable="hasModelBinding"
      @update:latLng="updateFromMarkerCoords"
    />

    <!-- Proximity radius indicator -->
    <LCircle
      v-if="hasValidCoords"
      :lat-lng="[coords!.latitude!, coords!.longitude!]"
      :radius
      :interactive="false"
    />

    <!-- Proximal sites popup -->
    <template #popup="{ item }">
      <SitePopup :item />
    </template>
  </SitesMap>
</template>

<script setup lang="ts">
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { LCircle, LMarker } from '@vue-leaflet/vue-leaflet'
import { computed, useTemplateRef } from 'vue'
import SitesMap from '../maps/SitesMap.vue'
import SitePopup from './SitePopup.vue'
import { Coordinates, MaybeCoordinates } from '../maps'
import { LatLngLiteral } from 'leaflet'
import { useMousePressed, watchOnce } from '@vueuse/core'
import { hasEventListener } from '../toolkit/vue-utils'

const coords = defineModel<MaybeCoordinates>({ required: true })

const hasModelBinding = hasEventListener('onUpdate:modelValue')

const hasValidCoords = computed(() => Coordinates.isValidCoordinates(coords.value))

const props = withDefaults(
  defineProps<{
    radius?: number
  }>(),
  { radius: 10000 }
)

const proximityFetchOptions = computed(() => ({
  enabled: hasValidCoords.value,
  ...sitesProximityOptions({
    body: {
      latitude: coords.value!.latitude!,
      longitude: coords.value!.longitude!,
      radius: props.radius
    }
  })
}))

const { data, isPending } = useQuery(proximityFetchOptions)

const map = useTemplateRef<HTMLElement>('map')
const mouse = useMousePressed({ target: map })

function updateFromMarkerCoords({ lat, lng }: LatLngLiteral) {
  if (!mouse.pressed.value) return
  // Update coordinates on mouse release
  watchOnce(
    () => mouse.pressed.value,
    () => {
      coords.value.latitude = Number(lat.toFixed(4))
      coords.value.longitude = Number(lng.toFixed(4))
    }
  )
}
</script>

<style scoped lang="scss"></style>
