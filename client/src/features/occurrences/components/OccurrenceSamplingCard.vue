<template>
  <v-card
    title="Sampling"
    variant="elevated"
    class="small-card-title"
    prepend-icon="mdi-package-down"
    :subtitle="DateWithPrecision.format(item.performed_on, undefined, 'Date unspecified')"
  >
    <!-- <template v-if="isGranted('Maintainer') || (isGranted('Maintainer') && isOwner(item))" #append>
      <v-btn icon="mdi-pencil" variant="tonal" size="small" @click="emit('edit')" />
    </template> -->
    <template #append v-if="item.site.country">
      <div class="d-flex align-center ga-1">
        <CoordinatesChip :coordinates="item.coordinates" label />
        <CoordPrecisionChip :precision="item.coordinates.precision" label />
      </div>
    </template>
    <v-divider />
    <v-list-item
      prepend-icon="mdi-map-marker-outline"
      :title="item.site.name || item.site.code"
      :subtitle="item.site.locality"
    >
      <template #title>
        {{ item.site.name || item.site.code }}
      </template>
      <template #append v-if="item.site.country">
        <CountryChip :country="item.site.country" size="small" />
      </template>
    </v-list-item>
    <ItemLocationMap :item="item" :height="500" :exclude-ids="[item.id]" />
    <v-divider />
    <v-list>
      <v-list-item prepend-icon="mdi-account-multiple">
        <v-chip
          v-for="person in item.performed_by"
          :key="person"
          :text="person"
          size="small"
          class="ma-1"
        />
        <span v-if="!item.performed_by" class="text-muted">Unknown</span>
        <template #append>
          <span class="text-muted text-caption">Sampled by</span>
        </template>
      </v-list-item>
      <v-divider />
      <v-list-group value="Details" prepend-icon="mdi-text-box">
        <template #activator="{ props }">
          <v-list-item v-bind="props" title="Details" lines="two" />
        </template>
        <SamplingListItems :sampling="item" />
      </v-list-group>

      <v-divider />

      <!-- <v-list-item prepend-icon="mdi-package-variant ">
        <v-chip
          v-for="{ id, code, identification } in samples"
          :variant="id === item.id ? 'outlined' : 'tonal'"
          :text="Identification.intersperseConfer(identification.taxon, identification.confer)"
          :class="['ma-1', { 'text-muted': id === item.id }]"
          :to="id !== item.id ? { name: 'occurrence-item', params: { code: code } } : undefined"
          label
          v-tooltip="
            id === item.id
              ? {
                  location: 'start',
                  origin: 'start',
                  openOnClick: true,
                  text: 'Currently viewed occurrence'
                }
              : undefined
          "
        />
        <template #append>
          <span class="text-muted text-caption">Reported occurrences</span>
        </template>
      </v-list-item> -->
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { DateWithPrecision, Identification, SamplingWithDetails } from '@/api/adapters'
import ItemLocationMap from '@/features/cartography/components/ItemLocationMap.vue'
import SamplingListItems from '@/features/occurrences/components/sampling/SamplingListItems.vue'
import CountryChip from '@/features/site/components/CountryChip'
import { useUserStore } from '@/stores/user'
import { useSorted } from '@vueuse/core'
import { computed } from 'vue'
import CoordinatesChip from './CoordinatesChip'
import CoordPrecisionChip from '@/features/site/components/CoordPrecisionChip'

const { item } = defineProps<{
  item: { id: string } & SamplingWithDetails
}>()
const emit = defineEmits<{
  edit: []
}>()

const { isGranted, isOwner } = useUserStore()

// const samples = useSorted(
//   computed(() => item.occurrences ?? []),
//   (a, b) => {
//     if (a.id === item.id) return -1
//     else return a.identification.taxon.name.localeCompare(b.identification.taxon.name)
//   }
// )
</script>

<style lang="scss"></style>
