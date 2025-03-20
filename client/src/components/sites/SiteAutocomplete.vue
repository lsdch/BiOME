<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchValue"
    :items="items ?? []"
    item-title="name"
    item-value="code"
    :loading="isFetching ? 'primary' : false"
    no-filter
    chips
    clear-on-select
    clearable
    return-object
    label="Site search"
    placeholder="WGS84 coords or search term..."
    persistent-placeholder
    v-bind="$attrs"
  >
    <template #item="{ item: { raw: site }, props }">
      <v-list-item :title="site.name" v-bind="props">
        <template #subtitle>
          <v-list-item-subtitle>
            {{ site.locality ?? 'Unspecified locality' }}
            <CountryChip v-if="site.country" :country="site.country" size="small" />
          </v-list-item-subtitle>
        </template>
        <template #append>
          <div class="d-flex flex-column align-end">
            <v-chip :text="site.code" class="font-monospace" size="small" />
            <div v-if="'distance' in site" class="text-caption font-monospace">
              {{ Math.round(site.distance) }}m
            </div>
          </div>
        </template>
      </v-list-item>
    </template>
    <template #no-data>
      <v-list-item
        class="text-muted"
        :title="
          isFetching
            ? 'Loading...'
            : isSearchCoordinates
              ? `No sites found within ${radius}m radius`
              : (searchValue?.length ?? 0) > 2
                ? 'No matching sites'
                : 'Waiting for query...'
        "
      >
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { searchSitesOptions, sitesProximityOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { reactiveComputed, toRefs } from '@vueuse/core'
import { computed, ref } from 'vue'
import CountryChip from './CountryChip'
import { SiteWithDistance, SiteWithScore } from '@/api'

const model = defineModel<SiteWithDistance | SiteWithScore>()
const searchValue = ref<string | undefined>('')
const radius = ref<number>(20_000)

const isSearchCoordinates = computed(
  () => !!searchValue.value && /^\s*\d+(\.\d+)?\s*,\s*\d+(\.\d+)?\s*$/.test(searchValue.value)
)

function asCoordinates(value: string): { latitude: number; longitude: number } {
  const [latitude, longitude] = value.split(',').map((v) => parseFloat(v.trim()))
  return { latitude, longitude }
}

const sitesMatching = useQuery(
  computed(() => ({
    enabled: !isSearchCoordinates.value && (searchValue.value?.length ?? 0) > 2,
    ...searchSitesOptions({ query: { query: searchValue.value } })
  }))
)

const sitesAtProximity = useQuery(
  computed(() => ({
    enabled: isSearchCoordinates.value,
    ...sitesProximityOptions({
      body: {
        radius: radius.value,
        ...(searchValue.value ? asCoordinates(searchValue.value) : { latitude: 0, longitude: 0 })
      }
    })
  }))
)

const { data: items, isFetching } = toRefs(
  reactiveComputed(() => (isSearchCoordinates.value ? sitesAtProximity : sitesMatching))
)
</script>

<style scoped lang="scss"></style>
