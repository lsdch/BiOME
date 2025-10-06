<template>
  <div class="d-flex flex-column fill-height">
    <v-progress-linear v-if="hasValidCoords && isPending" indeterminate />
    <div class="flex-grow-1">
      <BaseMap
        ref="map"
        :marker="hasValidCoords ? ({ coordinates: coords } as Geocoordinates) : undefined"
        :marker-layers="
          hasValidCoords
            ? [
                {
                  name: 'Proximal sites',
                  config: { radius: 4, fillColor: 'orangered', weight: 2, color: '#e41a1c' },
                  active: true,
                  clustered: true,
                  data: data?.filter(
                    ({ distance, code }) =>
                      distance <= proximityRadius && !omitCodes?.includes(code)
                  )
                }
              ]
            : undefined
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
      </BaseMap>
    </div>
    <ProximityRadiusSlider
      class="flex-grow-0"
      @update:radius="(radius) => (proximityRadius = radius)"
    />
  </div>
</template>

<script setup lang="ts">
import { LatLongCoords } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { LCircle, LMarker } from '@vue-leaflet/vue-leaflet'
import { useDebounceFn, useMousePressed } from '@vueuse/core'
import { LatLngLiteral } from 'leaflet'
import { computed, nextTick, ref, useTemplateRef, watch } from 'vue'
import { Coordinates, Geocoordinates, MaybeCoordinates } from '../maps'
import BaseMap from '../maps/BaseMap.vue'
import ProximityRadiusSlider from '../maps/ProximityRadiusSlider.vue'
import { hasEventListener } from '../toolkit/vue-utils'
import SitePopup from './SitePopup.vue'

const coords = defineModel<MaybeCoordinates>({ required: true })
const { omitCodes } = defineProps<{ omitCodes?: string[] }>()

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
