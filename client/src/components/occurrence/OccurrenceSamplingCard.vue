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
        <v-list-item title="Samples" prepend-icon="mdi-package-variant">
          <v-chip
            v-for="sample in item.sampling.samples"
            variant="tonal"
            :text="sample.identification.taxon.name"
            :title="sample.category"
            :prepend-icon="OccurrenceCategory.icon(sample.category)"
            class="ma-1"
            :color="sample.id === item.id ? 'success' : 'primary'"
            :active="sample.id === item.id"
            :disabled="sample.id === item.id"
            :to="
              sample.id !== item.id
                ? { name: 'biomat-item', params: { code: sample.code } }
                : undefined
            "
            label
          />
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { EventInner, Sampling } from '@/api'
import { DateWithPrecision, OccurrenceCategory } from '@/api/adapters'
import SamplingListItems from '../events/SamplingListItems.vue'

defineProps<{ item: { id: string; sampling: Sampling; event: EventInner } }>()
const emit = defineEmits<{
  edit: []
}>()
</script>

<style scoped lang="scss"></style>
