<template>
  <div class="bg-surface fill-height">
    <TableToolbar :title="site.name" icon="mdi-map-marker"></TableToolbar>
    <v-container class="align-start" fluid>
      <v-row>
        <v-col cols="12" md="6">
          <v-list>
            <v-list-item title="Code" :subtitle="site.code"></v-list-item>
            <v-list-item title="Description" :subtitle="site.description"></v-list-item>
          </v-list>
        </v-col>
        <v-col cols="12" md="6" style="height: 50vh">
          <v-card class="fill-height d-flex flex-column">
            <template #prepend>
              <v-icon icon="mdi-crosshairs-gps"></v-icon>
            </template>
            <template #subtitle>
              <div class="d-flex justify-space-between flex-wrap">
                <div>
                  Coordinates: {{ site.coordinates.latitude }}, {{ site.coordinates.longitude }}
                  <v-chip :text="site.coordinates.precision"></v-chip>
                </div>
                <div>
                  {{ [site.locality, site.country.name].filter((e) => e).join(', ') }}
                  <v-chip :text="site.country.code"></v-chip>
                </div>
              </div>
            </template>
            <SitesMap :items="[site]" regions :fitPad="0.3"></SitesMap>
            <template #actions>
              <!-- :href="`https://www.google.com/maps/place/${site.coordinates.latitude}+${site.coordinates.longitude}/@${site.coordinates.latitude},${site.coordinates.longitude},10z`" -->
              <v-btn
                prepend-icon="mdi-google-maps"
                color="primary"
                :href="`https://earth.google.com/web/search/${site.coordinates.latitude},${site.coordinates.longitude},10z`"
                text="See in Google Maps"
              />
            </template>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { LocationService } from '@/api'
import { handleErrors } from '@/api/responses'
import SitesMap from '@/components/maps/SitesMap.vue'
import TableToolbar from '@/components/toolkit/tables/TableToolbar.vue'
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const code = route.params['code'] as string

const site = ref(
  await LocationService.getSite({ path: { code } }).then(
    handleErrors((err) => {
      console.error('Failed to fetch site:', err)
    })
  )
)
</script>

<style scoped></style>
