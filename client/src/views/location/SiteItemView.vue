<template>
  <div class="bg-surface fill-height">
    <SiteEditForm v-if="editing" :site />
    <v-container v-else class="align-start fill-height" fluid>
      <v-row class="fill-height">
        <v-col cols="12" md="6">
          <div class="text-h5 d-flex align-center">
            <v-icon icon="mdi-map-marker-radius" />
            {{ site.name }}
            <v-btn
              class="ml-auto"
              color="primary"
              icon="mdi-pencil"
              variant="tonal"
              @click="toggleEdit(true)"
            />
          </div>
          <v-divider class="my-3" />
          <v-list>
            <v-list-item title="Code" :subtitle="site.code" />
            <v-list-item title="Coordinates" class="pr-0">
              <template #subtitle>
                <code class="text-wrap">
                  {{ site.coordinates.latitude }},
                  {{ site.coordinates.longitude }}
                </code>
                <v-chip
                  class="ml-5"
                  variant="outlined"
                  density="compact"
                  prepend-icon="mdi-crosshairs-question"
                  :text="site.coordinates.precision"
                />
              </template>
              <template #append>
                <v-btn
                  v-if="smAndDown"
                  class="ml-3"
                  variant="tonal"
                  icon="mdi-map-marker"
                  color="primary"
                  @click="toggleMap(true)"
                />
              </template>
            </v-list-item>
            <v-list-item title="Locality">
              <template #subtitle>
                {{ site.locality }}, {{ site.country.name }}
                <v-chip
                  class="ml-2 text-overline"
                  variant="outlined"
                  density="compact"
                  :text="site.country.code"
                />
              </template>
            </v-list-item>
            <v-list-item title="Description">
              <v-list-item-subtitle
                :class="{ 'font-italic': !site.description }"
                v-text="site.description ?? 'No description'"
              />
            </v-list-item>
          </v-list>
          <v-divider class="my-3" />
          <v-expansion-panels>
            <v-expansion-panel class="dataset-expansion-panel">
              <template #title>
                <div class="d-flex w-100 align-center mr-5">
                  Datasets
                  <v-chip
                    :text="site.datasets.length.toLocaleString()"
                    rounded
                    class="mx-2"
                    color="purple"
                  />
                  <v-btn icon="mdi-plus" density="compact" variant="tonal" class="ml-auto" />
                </div>
              </template>
              <template #text>
                <v-list
                  nav
                  lines="one"
                  density="compact"
                  :items="site.datasets"
                  item-title="label"
                  item-value="slug"
                >
                  <template #item="{ props: { title, value: slug } }">
                    <v-list-item
                      link
                      :to="$router.resolve({ name: 'dataset-item', params: { slug } })"
                      :title
                    />
                  </template>
                </v-list>
              </template>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
        <v-col cols="12" md="6">
          <ResponsiveDialog :as-dialog="smAndDown" v-model:open="mapActive">
            <v-card
              :class="['d-flex flex-column', lgAndUp ? 'fill-height' : 'h-50']"
              :rounded="!mapActive"
            >
              <SitesMap
                :items="[site]"
                regions
                :fitPad="0.3"
                :closable="mapActive"
                @close="toggleMap(false)"
              />
              <template #actions>
                <!-- :href="`https://www.google.com/maps/place/${site.coordinates.latitude}+${site.coordinates.longitude}/@${site.coordinates.latitude},${site.coordinates.longitude},10z`" -->
                <v-btn
                  prepend-icon="mdi-google-maps"
                  color="primary"
                  :href="`https://earth.google.com/web/search/${site.coordinates.latitude},${site.coordinates.longitude}`"
                  text="See in Google Maps"
                />
              </template>
            </v-card>
          </ResponsiveDialog>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup lang="ts">
import { LocationService, SiteUpdate } from '@/api'
import { handleErrors } from '@/api/responses'
import SitesMap from '@/components/maps/SitesMap.vue'
import SiteEditForm from '@/components/sites/SiteEditForm.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useToggle } from '@vueuse/core'
import { useRouteParams } from '@vueuse/router'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { useDisplay } from 'vuetify'

const { smAndDown, mdAndUp, lgAndUp } = useDisplay()

const [mapActive, toggleMap] = useToggle(false)
const [editing, toggleEdit] = useToggle(false)

const route = useRoute()
const code = route.params['code'] as string

const site = ref(
  await LocationService.getSite({ path: { code } }).then(
    handleErrors((err) => {
      console.error('Failed to fetch site:', err)
    })
  )
)

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

<style lang="scss">
.dataset-expansion-panel {
  .v-expansion-panel-text__wrapper {
    padding: 0px;
  }
}
</style>
