<template>
  <div class="bg-surface fill-height">
    <v-container class="align-start" fluid>
      <SiteEditForm v-if="editing" :site />
      <div v-else>
        <v-row>
          <v-col class="text-h5 d-flex align-center justify-space-between">
            {{ site.name }}
            <v-btn color="primary" icon="mdi-pencil" variant="plain" @click="toggleEdit(true)" />
          </v-col>
        </v-row>
        <v-divider class="my-3" />

        <v-row>
          <v-col cols="12" md="6">
            <v-list>
              <v-list-item title="Code" :subtitle="site.code"></v-list-item>
              <v-list-item title="Description">
                <v-list-item-subtitle
                  :class="{ 'font-italic': !site.description }"
                  v-text="site.description ?? 'No description'"
                />
              </v-list-item>
            </v-list>
          </v-col>
          <v-col cols="12" md="6" style="height: 50vh; min-height: 500px">
            <v-card class="fill-height d-flex flex-column">
              <template #prepend>
                <v-icon icon="mdi-crosshairs-gps"></v-icon>
              </template>
              <template #subtitle>
                <div class="d-flex justify-space-between flex-wrap">
                  <div>
                    Coordinates: {{ site.coordinates.latitude }},
                    {{ site.coordinates.longitude }}
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
      </div>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { LocationService, SiteUpdate } from '@/api'
import { handleErrors } from '@/api/responses'
import SitesMap from '@/components/maps/SitesMap.vue'
import SiteEditForm from '@/components/sites/SiteEditForm.vue'
import { useToggle } from '@vueuse/core'
import { useRouteQuery } from '@vueuse/router'
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

const [editing, toggleEdit] = useToggle()

async function submit() {
  const { name, code, description, coordinates, locality, altitude } = site.value
  const updatedSite: SiteUpdate = {
    name,
    code,
    coordinates,
    locality,
    altitude,
    description,
    country_code: site.value.country.code
  }
  return await LocationService.updateSite({ body: updatedSite, path: { code } })
}
</script>

<style scoped></style>
