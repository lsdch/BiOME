<template>
  <teleport :to="$vuetify.display.mdAndDown ? '#inline-map-container' : '#map-container'">
    <v-card
      class="d-flex flex-column fill-height"
      :min-height="400"
      :flat="$vuetify.display.mdAndDown"
      :rounded="$vuetify.display.mdAndDown ? 0 : undefined"
    >
      <v-sheet class="flex-grow-1" :height="300">
        <SitesMap
          :marker="site"
          :items="nearbySites?.filter(({ distance }) => distance <= proximityRadius)"
          regions
          :auto-fit="proximityRadius"
          clustered
        >
          <template #default="{ zoom }">
            <SiteRadius v-if="site" :site :zoom />
            <LCircle
              v-if="site && proximityRadius > 0"
              :lat-lng="Geocoordinates.LatLng(site)"
              :radius="proximityRadius"
              :interactive="false"
              dash-array="1px 8px"
              dash-offset="5px"
            />
          </template>
          <template #popup="{ item }">
            <KeepAlive>
              <SitePopup
                v-if="item"
                :item
                :options="{ keepInView: false, autoPan: false }"
                :showRadius="false"
              />
            </KeepAlive>
          </template>
        </SitesMap>
      </v-sheet>
      <ProximityRadiusSlider
        @update:radius="(radius) => (proximityRadius = radius)"
        class="flex-grow-0 px-2"
      />
      <template #actions v-if="site">
        <!-- :href="`https://www.google.com/maps/place/${site.coordinates.latitude}+${site.coordinates.longitude}/@${site.coordinates.latitude},${site.coordinates.longitude},10z`" -->
        <v-menu>
          <template #activator="{ props }">
            <v-btn prepend-icon="mdi-map" color="primary" v-bind="props" text="Open in" />
          </template>
          <v-list density="compact">
            <v-list-item
              title="Google Maps"
              prepend-icon="mdi-google-maps"
              :href="`http://maps.google.com/maps?&z=15&mrt=yp&t=k&q=${site.coordinates.latitude}+${site.coordinates.longitude}+(${site.name})`"
            />
            <v-list-item
              title="Google Earth"
              prepend-icon="mdi-google-earth"
              :href="`https://earth.google.com/web/search/${site.coordinates.latitude},${site.coordinates.longitude}`"
            />
          </v-list>
        </v-menu>
      </template>
    </v-card>
  </teleport>
</template>

<script setup lang="ts">
import { Site } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { Geocoordinates } from '@/components/maps'
import ProximityRadiusSlider from '@/components/maps/ProximityRadiusSlider.vue'
import SitesMap from '@/components/maps/SitesMap.vue'
import SitePopup from '@/components/sites/SitePopup.vue'
import SiteRadius from '@/components/sites/SiteRadius'
import { useQuery } from '@tanstack/vue-query'
import { LCircle } from '@vue-leaflet/vue-leaflet'
import { computed, ref } from 'vue'

const { site } = defineProps<{
  site: Site
}>()

const proximityRadius = ref(0)

const {
  data: nearbySites,
  isPending: nearbySitesPending,
  error: nearbySitesError
} = useQuery(
  computed(() => ({
    ...sitesProximityOptions({
      body: {
        latitude: site.coordinates.latitude,
        longitude: site.coordinates.longitude,
        radius: 100_000,
        exclude: [site.code]
      }
    }),
    enabled: site !== undefined
  }))
)
</script>

<style scoped lang="scss"></style>
