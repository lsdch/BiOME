<template>
  <v-card
    title="Sampling"
    variant="elevated"
    class="small-card-title"
    prepend-icon="mdi-package-down"
    :subtitle="DateWithPrecision.format(item.event.performed_on)"
  >
    <template #append>
      <v-btn icon="mdi-pencil" variant="tonal" size="small" @click="emit('edit')" />
    </template>

    <v-divider />
    <v-list-item
      class="text-primary"
      prepend-icon="mdi-map-marker-outline"
      :title="item.event.site.name || item.event.site.code"
      :subtitle="item.event.site.locality"
      :to="{ name: 'site-item', params: { code: item.event.site.code } }"
    >
      <template #append v-if="item.event.site.country">
        <CountryChip :country="item.event.site.country" size="small" />
      </template>
    </v-list-item>
    <v-sheet :height="300">
      <SitesMap
        :marker="item.event.site"
        :items="nearbySites?.filter(({ distance }) => distance <= proximityRadius)"
        :auto-fit="proximityRadius"
        clustered
        regions
      >
        <template #popup="{ item }">
          <KeepAlive>
            <SitePopup :item :options="{ keepInView: false }" />
          </KeepAlive>
        </template>
        <LCircle
          v-if="proximityRadius"
          :lat-lng="Geocoordinates.LatLng(item.event.site)"
          :radius="proximityRadius"
          :interactive="false"
        />
      </SitesMap>
    </v-sheet>
    <v-list-item>
      <ProximityRadiusSlider @update:radius="(radius) => (proximityRadius = radius)" />
    </v-list-item>
    <v-divider />
    <v-list density="compact">
      <v-list-item prepend-icon="mdi-account-multiple">
        <PersonChip v-for="person in item.event.performed_by" :person size="small" class="ma-1" />
        <span v-if="!item.event.performed_by" class="text-muted">Unknown</span>
        <template #append>
          <span class="text-muted text-caption">Sampled by</span>
        </template>
      </v-list-item>
      <v-divider />
      <v-list-group value="Details" prepend-icon="mdi-text-box">
        <template #activator="{ props }">
          <v-list-item v-bind="props" title="Details" lines="two" />
        </template>
        <SamplingListItems :sampling="item.sampling" />
      </v-list-group>

      <v-divider />

      <v-list-item prepend-icon="mdi-package-variant ">
        <v-tooltip location="start" origin="start" open-on-click>
          The currently viewed bio material
          <template #activator="{ props }">
            <v-chip
              v-for="{ id, code, category, identification } in samples"
              :variant="id === item.id ? 'outlined' : 'tonal'"
              :text="identification.taxon.name"
              :title="category"
              :color="OccurrenceCategory.props[category].color"
              :prepend-icon="OccurrenceCategory.icon(category)"
              :class="['ma-1', { 'text-muted': id === item.id }]"
              :to="id !== item.id ? { name: 'biomat-item', params: { code: code } } : undefined"
              label
              v-bind="id === item.id ? props : undefined"
            />
          </template>
        </v-tooltip>
        <template #append>
          <span class="text-muted text-caption">Samples bundle </span>
        </template>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { Sampling } from '@/api'
import { DateWithPrecision, EventWithParticipants, OccurrenceCategory } from '@/api/adapters'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { LCircle } from '@vue-leaflet/vue-leaflet'
import { useSorted } from '@vueuse/core'
import { computed, ref } from 'vue'
import SamplingListItems from '../events/SamplingListItems.vue'
import { Geocoordinates } from '../maps'
import ProximityRadiusSlider from '../maps/ProximityRadiusSlider.vue'
import SitesMap from '../maps/SitesMap.vue'
import PersonChip from '../people/PersonChip'
import CountryChip from '../sites/CountryChip'
import SitePopup from '../sites/SitePopup.vue'

const { item } = defineProps<{
  item: { id: string; sampling: Sampling; event: EventWithParticipants }
}>()
const emit = defineEmits<{
  edit: []
}>()

const samples = useSorted(
  computed(() => item.sampling.samples ?? []),
  (a, b) => {
    if (a.id === item.id) return -1
    else return a.identification.taxon.name.localeCompare(b.identification.taxon.name)
  }
)

const proximityRadius = ref(0)

const {
  data: nearbySites,
  isPending,
  error
} = useQuery(
  sitesProximityOptions({
    body: {
      latitude: item.event.site.coordinates.latitude,
      longitude: item.event.site.coordinates.longitude,
      radius: 100_000,
      exclude: [item.event.site.code]
    }
  })
)
</script>

<style lang="scss"></style>
