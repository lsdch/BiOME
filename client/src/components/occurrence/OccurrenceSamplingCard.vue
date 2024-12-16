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
        <v-list-item title="Other samples" prepend-icon="mdi-package-variant">
          <v-chip
            v-for="sample in item.sampling.samples.filter(({ id }) => id !== item!.id)"
            :text="sample.identification.taxon.name"
            :title="sample.category"
            class="ma-1"
            :to="{ name: 'biomat-item', params: { code: sample.code } }"
            label
          />
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { EventInner, Sampling } from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import SamplingListItems from '../events/SamplingListItems.vue'

defineProps<{ item: { id: string; sampling: Sampling; event: EventInner } }>()
const emit = defineEmits<{
  edit: []
}>()
</script>

<style scoped lang="scss"></style>
