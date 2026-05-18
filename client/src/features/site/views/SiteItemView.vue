<template>
  <v-container
    fluid
    id="item-view-container"
    class="responsive-container"
    min-height="100%"
    :class="['align-start w-100', { large: xlAndUp, small: mdAndDown }]"
  >
    <v-card
      id="info-container"
      :title="site?.name || code"
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
      <CenteredSpinner v-if="isPending" size="large" :height="300" />
      <v-alert v-else-if="error" color="error" icon="mdi-alert"> Failed to load site </v-alert>
      <v-list v-else-if="site">
        <v-list-item prepend-icon="mdi-crosshairs-gps">
          <div class="d-flex align-center ga-1">
            <CoordinatesChip :coordinates="site.coordinates" />
            <CoordPrecisionChip :precision="site.coordinates.precision" />
          </div>
          <template #append>
            <span class="text-muted text-caption">Coordinates</span>
          </template>
        </v-list-item>
        <v-list-item prepend-icon="mdi-arrow-up-down">
          <span :class="site.altitude ? 'font-monospace' : 'text-muted'">
            {{ site.altitude || 'Unknown' }}
          </span>
          <template #append>
            <span class="text-muted text-caption">Altitude (m)</span>
          </template>
        </v-list-item>
        <v-list-item prepend-icon="mdi-town-hall">
          <span :class="{ 'text-muted': !site.locality }">
            {{ site.locality || 'Unknown locality' }}
          </span>
          <CountryChip
            v-if="site.country"
            class="mx-2 text-overline"
            density="compact"
            :country="site.country"
          />
          <template #append>
            <span class="text-muted text-caption">Locality</span>
          </template>
        </v-list-item>
        <v-list-item>
          <span
            :class="['text-muted', { 'font-italic': !site.description }]"
            v-text="site.description ?? 'No description'"
          />
          <template #append>
            <span class="text-muted text-caption">Description</span>
          </template>
        </v-list-item>
        <v-divider />
        <v-list-item>
          <span class="text-muted font-italic" v-if="!targeted_taxa?.length">Unknown</span>
          <TaxonChip v-else v-for="taxon in targeted_taxa" class="ma-1" :taxon size="small" />
          <template #append>
            <span class="text-muted text-caption">Targeted taxa</span>
          </template>
        </v-list-item>
        <v-list-item>
          <span class="text-muted font-italic" v-if="!occurring_taxa?.length">None</span>
          <TaxonChip v-else v-for="taxon in occurring_taxa" class="ma-1" :taxon size="small" />
          <template #append>
            <span class="text-muted text-caption">Sampled taxa</span>
          </template>
        </v-list-item>
        <v-list-item>
          <span :class="['mx-1', last_sampled ? 'font-monospace' : 'text-muted']">
            {{
              last_sampled === null
                ? 'Never'
                : DateWithPrecision.format(last_sampled, undefined, 'Unknown date')
            }}
          </span>
          <template #append>
            <span class="text-muted text-caption">Last sampled</span>
          </template>
        </v-list-item>
      </v-list>
    </v-card>

    <div id="map-container">
      <v-card v-if="!site" height="100%">
        <CenteredSpinner height="100%" :size="50" />
      </v-card>
      <SiteItemMap v-else :site />
    </div>

    <SiteFormDialog v-model:item="site" v-model:dialog="editDialog"></SiteFormDialog>

    <SiteItemDetailTabs id="panels" :site> </SiteItemDetailTabs>
    <!-- <div id="panels">
      <v-expansion-panels :disabled="isPending">
        <v-expansion-panel>
          <template #title>
            Occurrences
            <v-badge color="primary" inline :content="site?.samplings?.length ?? 0" />
          </template>

          <template #text>
            <OccurrencesAtSiteTable v-if="site" :samplings="site.samplings ?? []" />
          </template>
        </v-expansion-panel>
        <v-expansion-panel class="dataset-expansion-panel">
          <template #title>
            <div class="d-flex w-100 align-center mr-5">
              Datasets
              <v-badge
                color="purple"
                inline
                :content="site?.datasets?.length ?? 0"
                variant="tonal"
              />
              <v-btn icon="mdi-plus" density="compact" variant="tonal" class="ml-auto" />
            </div>
          </template>
          <template v-if="site" #text>
            <v-list nav lines="one" density="compact" item-title="label" item-value="slug">
              <v-list-item
                v-for="{ category, slug, label } in site.datasets"
                :title="label"
                link
                :to="
                  $router.resolve({
                    name: `${category.toLocaleLowerCase()}-dataset-item`,
                    params: { slug }
                  })
                "
              />
            </v-list>
          </template>
        </v-expansion-panel>
      </v-expansion-panels>
    </div> -->
    <!-- <div id="misc">
      <v-card>
        <v-list-item
          title="Abiotic measurements"
          class="py-2"
          prepend-icon="mdi-chart-line"
          :link="!isPending"
        >
          <AbioticChartsDialog
            :abiotic_data="Object.values(abiotic_measurements)"
            activator="parent"
          />
          <template #append>
            <v-chip :text="Object.keys(abiotic_measurements).length.toString() || 'None'" />
          </template>
        </v-list-item>
      </v-card>
    </div> -->
  </v-container>
</template>

<script setup lang="ts">
import { DateWithPrecision, Taxon } from '@/api'
import { getSiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SiteFormDialog from '@/components/forms/SiteFormDialogMutation.vue'
import CoordPrecisionChip from '@/features/site/components/CoordPrecisionChip'
import CountryChip from '@/features/site/components/CountryChip'
import TaxonChip from '@/features/taxonomy/components/TaxonChip'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { computed, nextTick } from 'vue'
import { useDisplay } from 'vuetify'
import SiteItemDetailTabs from '../components/SiteItemDetailTabs.vue'
import SiteItemMap from '../components/ResponsiveSiteItemMap.vue'
import CoordinatesChip from '@/features/occurrences/components/CoordinatesChip'

const { mdAndDown, xlAndUp } = useDisplay()

const [editDialog, toggleEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()
nextTick(() => {
  document.title = code
})

const { data: site, error, isPending } = useQuery(getSiteOptions({ path: { code } }))

const last_sampled = computed(() => {
  if (!site.value?.samplings) return null
  return site.value?.samplings?.reduce<DateWithPrecision | undefined>((acc, { performed_on }) => {
    if (!acc) return performed_on
    return DateWithPrecision.compare(acc, performed_on) > 0 ? acc : performed_on
  }, undefined)
})

const targeted_taxa = computed(() => {
  return Object.values(
    site.value?.samplings?.reduce<Record<string, Taxon>>((acc, { target_taxa }) => {
      target_taxa?.forEach((t) => {
        acc[t.name] = t
      })
      return acc
    }, {}) ?? {}
  )
})

const occurring_taxa = computed(() => {
  return site.value?.samplings
    ?.reduce<Map<string, Taxon>>((acc, { occurrences }) => {
      occurrences?.forEach(({ identification: { taxon: t } }) => {
        acc.set(t.name, t)
      })
      return acc
    }, new Map<string, Taxon>())
    ?.values()
    .toArray()
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
      'panels panels';
  }
  &.small {
    grid-template-areas:
      'info info'
      'map map'
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
