<template>
  <v-sheet :height>
    <BaseMap
      :markers="[site]"
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
    </BaseMap>
  </v-sheet>
  <v-list-item>
    <ProximityRadiusSlider v-model="proximitySliderValue" v-model:radius="proximityRadius" />
    <template #prepend>
      <OccurrencesOverviewDialog
        :data="nearbySites?.filter(({ distance }) => distance <= proximityRadius) ?? []"
      >
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
            variant="tonal"
            v-tooltip="`See list of sites and occurrences within radius`"
            v-bind="props"
          />
        </template>
      </OccurrencesOverviewDialog>
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
import { CoordinatesPrecision, SiteWithDistance } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import OccurrencesOverviewDialog from '@/features/occurrences/components/tables/OccurrencesOverviewDialog.vue'
import SitePopup from '@/features/site/components/SitePopup.vue'
import SiteRadius from '@/features/site/components/SiteRadius'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { Geocoordinates } from '.'
import BaseMap from './BaseMap.vue'
import { MarkerLayer } from './map-layers'
import ProximityRadiusSlider from './ProximityRadiusSlider.vue'
import SiteProximityRadius from './SiteProximityRadius.vue'

const { site, height = 350 } = defineProps<{
  site: Geocoordinates & { code: string }
  height?: number
}>()

const proximitySliderValue = ref<number>(0)
const proximityRadius = ref<number>(0)

const {
  data: nearbySites,
  isPending,
  error
} = useQuery(
  sitesProximityOptions({
    query: {
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
    clustered: true,
    maxClusterRadius: 5,
    data: nearbySites.value?.filter(({ distance }) => distance <= proximityRadius.value),
    config: {
      radius: 6,
      color: '#FF0000BB',
      fillColor: '#FF000055',
      weight: 2
    }
  }
})
</script>

<style scoped lang="scss"></style>
