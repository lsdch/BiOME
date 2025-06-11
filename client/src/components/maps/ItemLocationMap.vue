<template>
  <v-sheet :height>
    <SitesMap
      :marker="site"
      :marker-layers="[proximalSitesMarkers]"
      :auto-fit="proximityRadius || CoordinatesPrecision.radius(site?.coordinates.precision)"
      clustered
      regions
    >
      <template #popup="{ item }">
        <KeepAlive>
          <SitePopup :item :options="{ keepInView: false }" />
        </KeepAlive>
      </template>
      <SiteRadius v-if="site" :site />
      <SiteProximityRadius :site :proximity-radius />
    </SitesMap>
  </v-sheet>
  <v-list-item>
    <ProximityRadiusSlider @update:radius="(radius) => (proximityRadius = radius)" />
  </v-list-item>
</template>

<script setup lang="ts">
import { CoordinatesPrecision, SiteWithDistance } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { Geocoordinates } from '.'
import SiteRadius from '../sites/SiteRadius'
import ProximityRadiusSlider from './ProximityRadiusSlider.vue'
import SiteProximityRadius from './SiteProximityRadius.vue'
import SitesMap, { MarkerLayer } from './SitesMap.vue'

const { site, height = 300 } = defineProps<{
  site: Geocoordinates & { code: string }
  height?: number
}>()

const proximityRadius = ref<number>(0)

const {
  data: nearbySites,
  isPending,
  error
} = useQuery(
  sitesProximityOptions({
    body: {
      latitude: site.coordinates.latitude,
      longitude: site.coordinates.longitude,
      radius: 100_000,
      exclude: [site.code]
    }
  })
)

const proximalSitesMarkers = computed<MarkerLayer<SiteWithDistance>>(() => {
  return {
    active: true,
    data: nearbySites.value?.filter(({ distance }) => distance <= proximityRadius.value),
    config: {
      clustered: false,
      radius: 6,
      color: '#FF0000BB',
      fillColor: '#FF000055',
      weight: 2
    }
  }
})
</script>

<style scoped lang="scss"></style>
