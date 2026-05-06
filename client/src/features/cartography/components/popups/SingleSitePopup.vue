<template>
  <v-card class="text-label-small" :max-width="300">
    <v-list-item
      class="text-no-wrap"
      :title="item.name"
      :subtitle="item.code"
      :to="{ name: 'site-item', params: { code: item.code } }"
      target="_blank"
    >
      <template #title="{ title }">
        <span class="text-wrap text-label">{{ title }}</span>
      </template>
      <template #subtitle="{ subtitle }">
        <span class="font-monospace text-label-small">{{ subtitle }}</span>
      </template>
    </v-list-item>
    <v-divider />
    <v-list density="compact">
      <v-list-item prepend-icon="mdi-crosshairs-gps" @click.stop="copyCoordinates">
        <div class="coordinates font-monospace cursor-pointer">
          <span class="label"> Lat </span>
          {{ item.coordinates.latitude }}
          <span class="label"> Lng </span>
          {{ item.coordinates.longitude }}
        </div>
        <v-overlay
          v-model="hasCopied"
          class="align-center justify-center"
          contained
          content-class="w-100"
        >
          <v-alert
            icon="mdi-content-copy"
            color="success"
            variant="elevated"
            density="compact"
            text="Copied"
            width="100%"
          />
        </v-overlay>
        <template #append>
          <CoordPrecisionChip :precision="item.coordinates.precision" size="small" />
        </template>
      </v-list-item>
      <v-list-item prepend-icon="mdi-map-marker">
        <span :class="{ 'text-muted': !item.locality }">
          {{ item.locality ?? 'Locality unspecified' }}
        </span>
        <template #append v-if="item.country">
          <CountryChip :country="item.country" size="small" class="ml-2" />
        </template>
      </v-list-item>
      <slot name="append-items" :item />
    </v-list>
  </v-card>
</template>

<script setup lang="ts" generic="Item extends SiteItem">
import { SiteItem } from '@/api'
import CoordPrecisionChip from '@/features/site/components/CoordPrecisionChip'
import CountryChip from '@/features/site/components/CountryChip'
import { useClipboard, useTimeoutFn, useToggle } from '@vueuse/core'

const { item } = defineProps<{ item: Item }>()

const [hasCopied, toggleHasCopied] = useToggle(false)
const hasCopiedTimeout = useTimeoutFn(() => toggleHasCopied(false), 2000)

const { copy } = useClipboard()
function copyCoordinates() {
  hasCopiedTimeout.stop()
  copy(`${item.coordinates.latitude}, ${item.coordinates.longitude}`)
  toggleHasCopied(true)
  hasCopiedTimeout.start()
}
</script>

<style scoped lang="scss"></style>
