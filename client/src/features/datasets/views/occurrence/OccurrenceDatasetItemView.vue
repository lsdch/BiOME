<template>
  <DatasetItemView :slug :dataset :map-on-side>
    <template #map>
      <OccurrenceDatasetMap :sites="dataset?.sites ?? []" />
    </template>
    <template #details>
      <CenteredSpinner
        v-if="isPending"
        :height="300"
        size="large"
        color="primary"
        class="flex-grow-1"
      />
      <PageErrors v-else-if="error" :error class="flex-grow-1" />
      <v-card v-else-if="dataset" class="flex-grow-1 d-flex flex-column">
        <v-tabs v-model="tab" mandatory>
          <v-tab v-show="!mapOnSide" value="map" prepend-icon="mdi-map"> Map </v-tab>
          <v-tab value="sites" prepend-icon="mdi-map-marker">
            Sites
            <v-chip class="mx-1" :text="dataset.sites?.length.toString()" density="compact" />
          </v-tab>
          <v-tab value="occurrences" prepend-icon="mdi-crosshairs-gps">
            Occurrences
            <v-chip class="mx-1" :text="occurrences?.length.toString()" density="compact" />
          </v-tab>
          <v-tab value="bibliography" prepend-icon="mdi-newspaper">
            Bibliography
            <v-chip
              class="mx-1"
              :text="dataset?.bibliography?.length.toString()"
              density="compact"
            />
          </v-tab>
          <v-spacer></v-spacer>
          <v-tab
            value="administration"
            prepend-icon="mdi-cog"
            v-if="isUserMaintainer || userStore.isGranted('Admin')"
          >
          </v-tab>
        </v-tabs>
        <v-tabs-window v-model="tab" class="fill-height" crossfade>
          <v-tabs-window-item value="map" key="map" :transition="false" id="map-tab">
            <OccurrenceDatasetMap :sites="dataset.sites" style="min-height: 600px" />
          </v-tabs-window-item>
          <v-tabs-window-item value="sites" key="sites">
            <SitesTable :sites="dataset.sites.map(({ samplings, ...s }) => s)" />
          </v-tabs-window-item>
          <v-tabs-window-item value="occurrences" key="occurrences">
            <OccurrencesTable :with-site="true" :occurrences />
          </v-tabs-window-item>
          <v-tabs-window-item value="bibliography" key="bibliography">
            <OccurrenceDatasetBibliography :dataset />
          </v-tabs-window-item>
          <v-tabs-window-item value="administration" key="administration">
            <OccurrenceDatasetAdministration :dataset @refresh="refetch()" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
    </template>
    <template #subtitle>
      <div class="d-flex ga-1">
        <v-chip label text="Occurrences dataset" size="small" prepend-icon="mdi-crosshairs-gps" />
        <OccurrenceDatasetParticipants :dataset v-if="dataset?.contributors" />
      </div>
    </template>
  </DatasetItemView>
</template>

<script setup lang="tsx">
import { getOccurrenceDatasetOptions } from '@/api/gen/@tanstack/vue-query.gen'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import OccurrenceDatasetBibliography from '@/features/datasets/components/OccurrenceDatasetBibliography.vue'
import OccurrenceDatasetMap from '@/features/datasets/components/OccurrenceDatasetMap.vue'
import OccurrenceDatasetParticipants from '@/features/datasets/components/OccurrenceDatasetParticipants.vue'
import DatasetItemView from '@/features/datasets/views/DatasetItemView.vue'
import OccurrencesTable from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import SitesTable from '@/features/site/components/SitesTable.vue'
import { useUserStore } from '@/stores/user'
import { useQuery } from '@tanstack/vue-query'
import { computed, nextTick, ref, watch, watchEffect } from 'vue'
import { useDisplay } from 'vuetify/lib/composables/display.mjs'
import OccurrenceDatasetAdministration from '../../components/OccurrenceDatasetAdministration.vue'

type Tab = 'map' | 'sites' | 'occurrences' | 'bibliography' | 'administration'
const tab = ref<Tab>('map')

const { xlAndUp: mapOnSide } = useDisplay()

const userStore = useUserStore()
const isUserMaintainer = computed(() => {
  return dataset.value?.maintainers?.some((m) => m.id === userStore.user?.id)
})

watch(
  mapOnSide,
  (val) => {
    if (val && tab.value === 'map') tab.value = 'sites'
  },
  { immediate: true }
)

const { slug } = defineProps<{
  slug: string
}>()

const {
  data: dataset,
  error,
  refetch,
  isPending
} = useQuery(getOccurrenceDatasetOptions({ path: { slug } }))

watch(
  dataset,
  (d) => {
    nextTick(() => {
      document.title = d ? d.label : 'Dataset'
    })
  },
  { immediate: true }
)

// Occurrences table data
const occurrences = computed(() => {
  return dataset.value?.sites?.flatMap(({ samplings, ...site }) => {
    return samplings.flatMap(({ occurrences, date }) =>
      occurrences.map((o) => ({
        ...o,
        sampling_date: date,
        site
      }))
    )
  })
})
</script>

<style lang="scss">
.v-list-item.empty .v-list-item-subtitle {
  font-style: italic;
}
</style>
