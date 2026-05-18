<template>
  <v-list>
    <v-list-item>
      <SiteAutocomplete :model-value="site" @update:model-value="addSite" class="pt-2" />
    </v-list-item>
    <v-divider />
    <v-list-item
      v-for="(marker, index) in markers"
      :key="marker.data.code"
      :title="marker.data.name ?? marker.data.code"
      :subtitle="marker.data.name ? marker.data.code : undefined"
      @click="emit('focusSite', marker.data)"
    >
      <template #append>
        <v-icon icon="mdi-map-marker" :color="marker.options?.color"></v-icon>
        <v-btn
          variant="plain"
          icon="mdi-close-circle"
          size="small"
          color=""
          @click="markers.splice(index, 1)"
        ></v-btn>
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { SiteItem, SiteWithDistance, SiteWithScore } from '@/api'
import SiteAutocomplete from '@/features/site/components/SiteAutocomplete.vue'
import { Hsluv } from 'hsluv'
import { computed, reactive, ref } from 'vue'
import { PinMarker } from './layers-manager/map-layers'

const site = ref<SiteWithDistance | SiteWithScore>()

const markers = defineModel<PinMarker<SiteItem>[]>({ default: reactive([]) })

const siteColors = computed(() => {
  return new Set(markers.value.map((s, i) => s.options?.color))
})

const emit = defineEmits<{
  focusSite: [site: SiteItem]
}>()

function addSite(s?: SiteItem) {
  if (!s) return
  site.value = undefined
  markers.value.push({ data: s, options: { color: generateColor() }, coordinates: s.coordinates })
}

function generateColor(index: number = 0): string {
  const conv = new Hsluv()
  conv.hsluv_h = (index * 137.508) % 360 // use golden angle approximation
  conv.hsluv_s = 90
  conv.hsluv_l = 65
  conv.hsluvToHex()
  return siteColors.value.has(conv.hex) ? generateColor(index + 1) : conv.hex
}
</script>

<style scoped lang="scss"></style>
