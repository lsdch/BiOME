<template>
  <teleport :to="$vuetify.display.mdAndDown ? '#inline-map-container' : '#map-container'">
    <v-card
      class="d-flex flex-column fill-height"
      :min-height="400"
      :flat="$vuetify.display.mdAndDown"
      :rounded="$vuetify.display.mdAndDown ? 0 : undefined"
    >
      <ItemLocationMap :site />
      <template #actions v-if="site">
        <!-- :href="`https://www.google.com/maps/place/${site.coordinates.latitude}+${site.coordinates.longitude}/@${site.coordinates.latitude},${site.coordinates.longitude},10z`" -->
        <v-menu>
          <template #activator="{ props }">
            <v-btn prepend-icon="mdi-map" color="primary" v-bind="props" text="Open in" />
          </template>
          <v-list density="compact">
            <v-list-item
              title="Google Maps"
              prepend-icon="mdi-google-maps"
              :href="`http://maps.google.com/maps?&z=15&mrt=yp&t=k&q=${site.coordinates.latitude}+${site.coordinates.longitude}`"
              target="_blank"
            />
            <v-list-item
              title="Google Earth"
              prepend-icon="mdi-google-earth"
              :href="`https://earth.google.com/web/search/${site.coordinates.latitude},${site.coordinates.longitude}`"
              target="_blank"
            />
          </v-list>
        </v-menu>
      </template>
    </v-card>
  </teleport>
</template>

<script setup lang="ts">
import { Site } from '@/api'
import ItemLocationMap from '@/components/maps/ItemLocationMap.vue'

const { site } = defineProps<{
  site: Site
}>()
</script>

<style scoped lang="scss"></style>
