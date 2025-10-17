<template>
  <LLayerGroup
    name="PolygonLayer"
    pane="overlayPane"
    @click="
      (e: LeafletMouseEvent) => {
        console.log('click', e.latlng)
        if (active) {
          addPolylinePoint(e.latlng)
        }
      }
    "
  >
    <LPolygon
      v-if="!active && polyline.length > 2"
      :lat-lngs="polyline"
      color="orangered"
      :weight="2"
      fill
      no-clip
      fill-rule=""
      :fill-opacity="0.3"
    />
    <LPolyline
      v-if="active && polyline.length"
      :lat-lngs="[...polyline, ...(cursorCoordinates ? [cursorCoordinates] : [])]"
      color="orangered"
      :weight="2"
      fill
      no-clip
      fill-rule=""
      :fill-opacity="0.3"
      :interactive="false"
    />
    <LCircleMarker
      v-if="active || polyline.length > 1"
      v-for="(latLng, i) in polyline"
      interactive
      :lat-lng
      :radius="i === 0 || i === polyline.length - 1 ? 6 : 3"
      fill
      :fill-opacity="1"
      :fillColor="i === 0 ? 'green' : 'orangered'"
      :color="i === 0 ? 'green' : 'orangered'"
      @click="
        (ev) => {
          console.log('click', i)
          if (active && i === 0) {
            if (polyline.length == 2) {
              clearPolyline()
            }
            active = false
          }
        }
      "
    />
  </LLayerGroup>
</template>

<script setup lang="ts">
import { LCircleMarker, LLayerGroup, LPolygon, LPolyline } from '@vue-leaflet/vue-leaflet'
import { onKeyStroke } from '@vueuse/core'
import { LatLngExpression, LatLngLiteral, LeafletMouseEvent } from 'leaflet'
import { ref } from 'vue'

const { cursorCoordinates } = defineProps<{
  cursorCoordinates?: LatLngLiteral
}>()

const active = defineModel<boolean>('active', { default: false })
const polyline = ref<LatLngExpression[]>([])
function addPolylinePoint(latlng: LatLngExpression) {
  polyline.value = [...polyline.value, latlng]
}
function clearPolyline() {
  polyline.value = []
}

onKeyStroke('Escape', () => {
  console.log('Escape pressed')
  if (active.value) {
    clearPolyline()
    active.value = false
  }
})
</script>

<style scoped lang="scss"></style>
