<template>
  <LCircle
    v-if="site && proximityRadius > 0"
    :lat-lng="Geocoordinates.LatLng(site)"
    :radius="proximityRadius"
    :interactive="false"
    dash-array="1px 8px"
    dash-offset="5px"
  />
</template>

<script setup lang="ts">
import { LCircle } from '@vue-leaflet/vue-leaflet'
import { Geocoordinates } from '.'

defineProps<{
  site?: Geocoordinates
  proximityRadius: number
}>()
</script>

<style scoped lang="scss"></style>
