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
      <v-btn
        icon="mdi-list-box"
        size="small"
        variant="tonal"
        v-tooltip="`See list of sites and occurrences within radius`"
        @click="toggleTableDialog(true)"
      />
    </template>
  </v-list-item>
  <v-dialog v-model="tableDialog">
    <v-card max-height="80vh">
      <v-tabs v-model="tab">
        <v-tab value="sites">Sites</v-tab>
        <v-tab value="occurrences">Occurrences</v-tab>
        <v-tab value="occurrences">Occurring taxa</v-tab>
      </v-tabs>
      <v-card-text>
        <ProximityRadiusSlider
          label="Radius"
          v-model="proximitySliderValue"
          v-model:radius="proximityRadius"
        />
      </v-card-text>
      <v-tabs-window v-model="tab">
        <v-tabs-window-item value="sites">sites </v-tabs-window-item>
        <v-tabs-window-item value="occurrences">
          <OccurrencesTable
            with-site
            :occurrences="
              nearbySites
                ?.filter(({ distance }) => distance <= proximityRadius)
                .flatMap(({ samplings, ...site }) => {
                  return samplings.flatMap((s) => s.occurrences).map((x) => ({ site, ...x }))
                })
            "
          />
        </v-tabs-window-item>
      </v-tabs-window>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { CoordinatesPrecision, SiteWithDistance } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SitePopup from '@/features/site/components/SitePopup.vue'
import SiteRadius from '@/features/site/components/SiteRadius'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import { Geocoordinates } from '.'
import BaseMap from './BaseMap.vue'
import { MarkerLayer } from './map-layers'
import ProximityRadiusSlider from './ProximityRadiusSlider.vue'
import SiteProximityRadius from './SiteProximityRadius.vue'
import { useToggle } from '@vueuse/core'
import OccurrencesTable from '@/features/occurrences/components/OccurrencesTable.vue'

const tab = ref<'sites' | 'occurrences'>()

const { site, height = 350 } = defineProps<{
  site: Geocoordinates & { code: string }
  height?: number
}>()

const [tableDialog, toggleTableDialog] = useToggle(false)

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
