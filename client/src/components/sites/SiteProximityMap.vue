<template>
  <v-progress-linear v-if="hasValidCoords && isPending" indeterminate />
  <div class="flex-grow-1">
    <SitesMap
      ref="map"
      :marker="hasValidCoords ? ({ coordinates: coords } as Geocoordinates) : undefined"
      :items="
        hasValidCoords ? data?.filter(({ distance }) => distance <= proximityRadius) : undefined
      "
      :auto-fit="proximityRadius"
      clustered
      :center="hasValidCoords ? [coords!.latitude!, coords!.longitude!] : undefined"
      :zoom="hasValidCoords ? 10 : 0"
      :min-zoom="1"
      hide-marker-control
    >
      <!-- Coordinates marker -->
      <template #marker="{ latLng }">
        <LMarker
          v-if="latLng"
          :lat-lng
          :draggable="hasModelBinding"
          @update:latLng="updateFromMarkerCoords"
        />
      </template>

      <!-- Proximity radius indicator -->
      <LCircle
        v-if="hasValidCoords && proximityRadius > 0"
        :lat-lng="[coords!.latitude!, coords!.longitude!]"
        :radius="proximityRadius"
        :interactive="false"
      />

      <!-- Proximal sites popup -->
      <template #popup="{ item }">
        <SitePopup :item />
      </template>
    </SitesMap>
  </div>
  <ProximityRadiusSlider
    class="flex-grow-0"
    @update:radius="(radius) => (proximityRadius = radius)"
  />
</template>

<script setup lang="ts">
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { LCircle, LMarker } from '@vue-leaflet/vue-leaflet'
import { computed, nextTick, ref, useTemplateRef, watch } from 'vue'
import SitesMap from '../maps/SitesMap.vue'
import SitePopup from './SitePopup.vue'
import { Coordinates, Geocoordinates, MaybeCoordinates } from '../maps'
import { LatLngLiteral } from 'leaflet'
import { useDebounceFn, useMousePressed, watchOnce } from '@vueuse/core'
import { hasEventListener } from '../toolkit/vue-utils'
import { LatLongCoords } from '@/api'
import ProximityRadiusSlider from '../maps/ProximityRadiusSlider.vue'

const coords = defineModel<MaybeCoordinates>({ required: true })

const hasModelBinding = hasEventListener('onUpdate:modelValue')

const hasValidCoords = computed(() => Coordinates.isValidCoordinates(coords.value))

const proximityRadius = ref(0)

const proximityFetchOptions = computed(() => ({
  enabled: hasValidCoords.value,
  ...sitesProximityOptions({
    body: {
      latitude: coords.value!.latitude!,
      longitude: coords.value!.longitude!,
      radius: 100_000
    }
  })
}))

const { data, isPending } = useQuery(proximityFetchOptions)

const map = useTemplateRef<HTMLElement>('map')
const mouse = useMousePressed({ target: map })

const draggingCoords = ref<LatLongCoords>()

watch(mouse.pressed, (pressed, wasPressing) => {
  console.log('pressed', pressed, wasPressing)
  if (wasPressing && !pressed) {
    nextTick(
      useDebounceFn(() => {
        coords.value.latitude = draggingCoords.value?.latitude
        coords.value.longitude = draggingCoords.value?.longitude
      }, 100)
    )
  }
})

const updateFromMarkerCoords = ({ lat, lng }: LatLngLiteral) => {
  draggingCoords.value = {
    latitude: Number(lat.toFixed(4)),
    longitude: Number(lng.toFixed(4))
  }
}
</script>

<style scoped lang="scss"></style>
