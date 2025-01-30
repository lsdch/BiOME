<template>
  <v-container
    fluid
    id="item-view-container"
    min-height="100%"
    :class="['align-start w-100', { large: xlAndUp, small: mdAndDown }]"
  >
    <v-card
      id="info-container"
      :title="site?.name ?? code"
      :subtitle="code"
      prepend-icon="mdi-map-marker-radius"
    >
      <template #append>
        <v-btn
          class="ml-auto"
          color="primary"
          size="small"
          icon="mdi-pencil"
          variant="tonal"
          @click="toggleEdit(true)"
        />
      </template>
      <v-divider class="my-3" />
      <v-list v-if="site">
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
              title="Precision"
            />
          </template>
          <template #append>
            <v-btn
              v-if="mapAsDialog"
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
        <v-list-item title="Targeted taxa">
          <TaxonChip v-for="taxon in targeted_taxa" class="ma-1" :taxon />
        </v-list-item>
        <v-list-item title="Sampled taxa">
          <TaxonChip v-for="taxon in occurring_taxa" class="ma-1" :taxon />
        </v-list-item>
      </v-list>
    </v-card>

    <SiteFormDialog :edit="site" v-model="editDialog"></SiteFormDialog>

    <div id="panels" v-if="site">
      <v-expansion-panels>
        <v-expansion-panel>
          <template #title>
            Events
            <v-badge color="primary" inline :content="site.events?.length ?? 0" />
          </template>

          <template #text>
            <SiteEventsTable :site />
          </template>
        </v-expansion-panel>
        <v-expansion-panel class="dataset-expansion-panel">
          <template #title>
            <div class="d-flex w-100 align-center mr-5">
              Datasets
              <v-badge
                color="purple"
                inline
                :content="site.datasets?.length ?? 0"
                variant="tonal"
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
    </div>
    <div id="misc">
      <v-card>
        <v-list-item title="Abiotic measurements" class="py-2" prepend-icon="mdi-chart-line" link>
          <AbioticChartsDialog
            :abiotic_data="Object.values(abiotic_measurements)"
            activator="parent"
          />
          <template #append>
            <v-chip :text="Object.keys(abiotic_measurements).length.toString() || 'None'" />
          </template>
        </v-list-item>
      </v-card>
    </div>
    <div id="map-container">
      <ResponsiveDialog :as-dialog="mapAsDialog" v-model:open="mapActive">
        <v-card class="d-flex flex-column fill-height" :rounded="!mapActive">
          <SitesMap
            :items="site ? [site] : []"
            regions
            :fitPad="0.3"
            :closable="mapActive"
            @close="toggleMap(false)"
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
      </ResponsiveDialog>
    </div>
  </v-container>
</template>

<script setup lang="ts">
import { Taxon } from '@/api'
import { getSiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SiteEventsTable from '@/components/events/SiteEventsTable.vue'
import SitesMap from '@/components/maps/SitesMap.vue'
import SiteFormDialog from '@/components/sites/SiteFormDialog.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import ResponsiveDialog from '@/components/toolkit/ui/ResponsiveDialog.vue'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { computed } from 'vue'
import { useDisplay } from 'vuetify'
import AbioticChartsDialog from './AbioticChartsDialog.vue'
import { AbioticData, AbioticDataPoint } from './AbioticLineChart.vue'

const { mdAndDown, xlAndUp } = useDisplay()

const mapAsDialog = mdAndDown

const [mapActive, toggleMap] = useToggle(false)
const [editDialog, toggleEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()

const { data: site, error } = useQuery(getSiteOptions({ path: { code } }))

const targeted_taxa = computed(() => {
  return Object.values(
    site.value?.events?.reduce<Record<string, Taxon>>((acc, event) => {
      event.samplings?.forEach(({ target }) => {
        target.taxa?.forEach((t) => {
          acc[t.name] = t
        })
      })
      return acc
    }, {}) ?? {}
  )
})

const occurring_taxa = computed(() => {
  return Object.values(
    site.value?.events?.reduce<Record<string, Taxon>>((acc, event) => {
      event.samplings?.forEach(({ occurring_taxa }) => {
        occurring_taxa?.forEach((t) => {
          acc[t.name] = t
        })
      })
      return acc
    }, {}) ?? {}
  )
})

const abiotic_measurements = computed(() => {
  return (
    site.value?.events?.reduce<Record<string, AbioticData>>(
      (acc, { performed_on, abiotic_measurements }) => {
        abiotic_measurements?.forEach(({ param, value }) => {
          if (performed_on.date === undefined) return
          acc[param.code] = {
            param,
            points: [{ y: value, date: performed_on.date }].concat(
              acc[param.code]?.points ?? Array<AbioticDataPoint>()
            )
          }
        })
        return acc
      },
      {}
    ) ?? {}
  )
})
</script>

<style lang="scss">
#map-container {
  grid-area: map;
  align-self: stretch;
}
#info-container {
  grid-area: info;
}
#panels {
  grid-area: panels;
}
#misc {
  grid-area: misc;
}

#item-view-container {
  display: grid;
  gap: 20px 2%;
  grid-template-columns: 49% 49%;
  grid-template-rows: auto auto auto 1fr;
  grid-template-areas:
    'info map'
    'panels panels'
    'misc misc';
  &.large {
    grid-template-areas:
      'info map'
      'panels misc';
  }
  &.small {
    grid-template-areas:
      'info info'
      'panels panels'
      'misc misc';
  }
}

.dataset-expansion-panel {
  .v-expansion-panel-text__wrapper {
    padding: 0px;
  }
}
</style>
