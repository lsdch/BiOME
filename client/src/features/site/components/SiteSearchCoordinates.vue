<template>
  <div>
    <div class="d-flex">
      <v-number-input
        v-model="coords.latitude"
        label="Latitude"
        decimal-separator="."
        :precision="4"
        :min="-90"
        :max="90"
        :step="0.0001"
        rounded="e-0"
      >
      </v-number-input>
      <v-number-input
        v-model="coords.longitude"
        label="Longitude"
        decimal-separator="."
        :precision="4"
        :min="-180"
        :max="180"
        :step="0.0001"
        rounded="0"
      >
      </v-number-input>
      <v-number-input
        v-model="coords.radius"
        label="Radius (km)"
        :step="1"
        :min="1"
        :max="200"
        rounded="s-0"
        class="flex-grow-0"
        :min-width="150"
        :placeholder="DEFAULT_RADIUS.toString()"
        persistent-placeholder
      >
      </v-number-input>
    </div>
    <v-select
      v-model="model"
      return-object
      :items="sitesAtProximity"
      label="Sites at proximity"
      item-title="name"
      :loading="Coordinates.isValidCoordinates(coords) && isPending"
      :placeholder
      persistent-placeholder
      :disabled="!Coordinates.isValidCoordinates(coords) || isPending"
    >
      <template #item="{ item: site, props }">
        <v-list-item
          :title="site.name ?? site.code"
          v-bind="props"
          :disabled="disabledCodes?.includes(site.code)"
        >
          <template #subtitle>
            <v-list-item-subtitle>
              {{ site.locality ?? 'Unspecified locality' }}
              <country-chip v-if="site.country" :country="site.country" size="small" />
            </v-list-item-subtitle>
          </template>
          <template #append>
            <div class="d-flex flex-column align-end ga-1">
              <v-chip :text="site.code" class="font-monospace" size="small" />
              <div class="text-label-small text-muted font-monospace">
                {{ Math.round(site.distance) }}m
              </div>
            </div>
          </template>
        </v-list-item>
      </template>
    </v-select>
  </div>
</template>

<script setup lang="ts">
import { SiteItem, SiteWithDistance } from '@/api'
import { sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { Coordinates } from '@/features/cartography/coordinates'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'
import CountryChip from './CountryChip'

type CoordsSearchQuery = {
  latitude?: number
  longitude?: number
  radius?: number
}

const coords = ref<CoordsSearchQuery>({})

const model = defineModel<SiteItem>()

const { disabledCodes } = defineProps<{
  disabledCodes?: string[]
}>()

const DEFAULT_RADIUS = 100
const QUERY_LIMIT = 20

const placeholder = computed(() => {
  if (!Coordinates.isValidCoordinates(coords.value)) {
    return 'Waiting for valid coordinates...'
  } else if (isPending.value) {
    return 'Loading sites at proximity...'
  } else {
    const count = sitesAtProximity.value?.length ?? 0
    const radius = (coords.value as CoordsSearchQuery).radius ?? DEFAULT_RADIUS
    return `${count}${count === QUERY_LIMIT ? '+' : ''} sites within ${radius}km radius`
  }
})

const { data: sitesAtProximity, isPending } = useQuery(
  computed(() => ({
    enabled: Coordinates.isValidCoordinates(coords.value),
    ...sitesProximityOptions({
      query: {
        radius: (coords.value.radius ?? DEFAULT_RADIUS) * 1000,
        latitude: coords.value.latitude,
        longitude: coords.value.longitude,
        limit: QUERY_LIMIT
      }
    })
  }))
)
</script>

<style scoped lang="scss"></style>
