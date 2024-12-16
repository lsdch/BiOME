<template>
  <v-card
    title="Sampling"
    variant="elevated"
    prepend-icon="mdi-package-down"
    :subtitle="DateWithPrecision.format(item.event.performed_on)"
  >
    <template #append>
      <v-btn icon="mdi-pencil" variant="tonal" size="small" @click="emit('edit')" />
    </template>
    <v-card-text>
      <v-list density="compact">
        <v-list-item
          class="text-primary"
          prepend-icon="mdi-map-marker-outline"
          :title="item.event.site.name"
          subtitle="Locality, CC"
          :to="{ name: 'site-item', params: { code: item.event.site.code } }"
        ></v-list-item>
        <v-list-group value="Details" prepend-icon="mdi-text-box">
          <template #activator="{ props }">
            <v-list-item v-bind="props" title="Details" lines="two"></v-list-item>
          </template>
          <SamplingListItems :sampling="item.sampling" />
        </v-list-group>
        <v-divider></v-divider>
        <v-list-item title="Samples" prepend-icon="mdi-package-variant ">
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
                class="ma-1"
                :to="id !== item.id ? { name: 'biomat-item', params: { code: code } } : undefined"
                label
                v-bind="id === item.id ? props : undefined"
              />
            </template>
          </v-tooltip>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { EventInner, Sampling } from '@/api'
import { DateWithPrecision, OccurrenceCategory } from '@/api/adapters'
import SamplingListItems from '../events/SamplingListItems.vue'
import { useSorted } from '@vueuse/core'

const { item } = defineProps<{ item: { id: string; sampling: Sampling; event: EventInner } }>()
const emit = defineEmits<{
  edit: []
}>()

const samples = useSorted(item.sampling.samples, (a, b) => {
  if (a.id === item.id) return -1
  else return a.identification.taxon.name.localeCompare(b.identification.taxon.name)
})
</script>

<style scoped lang="scss"></style>
