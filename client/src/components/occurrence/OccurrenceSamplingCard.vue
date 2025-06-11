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
    <ItemLocationMap :site="item.event.site" :height="300" />
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
import { useSorted } from '@vueuse/core'
import { computed } from 'vue'
import SamplingListItems from '../events/SamplingListItems.vue'
import ItemLocationMap from '../maps/ItemLocationMap.vue'
import PersonChip from '../people/PersonChip'
import CountryChip from '../sites/CountryChip'

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
</script>

<style lang="scss"></style>
