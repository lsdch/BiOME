<template>
  <v-card
    title="Previews"
    variant="flat"
    prepend-icon="mdi-image-filter-center-focus-strong-outline"
    v-bind="$attrs"
  >
    <v-slide-group show-arrows :direction="direction" center-active>
      <v-slide-group-item v-for="size in previewSizes" :key="size">
        <div
          :class="`d-flex align-center ${direction == 'vertical' ? 'w-100 flex-column' : 'h-100'}`"
        >
          <v-avatar class="ma-3" :size="size">
            <preview
              :width="size"
              :height="size"
              :image="result?.image"
              :coordinates="result?.coordinates"
            />
            <v-tooltip activator="parent" location="center"> {{ size }}px </v-tooltip>
          </v-avatar>
        </div>
      </v-slide-group-item>
    </v-slide-group>
  </v-card>
</template>

<script setup lang="ts">
import { CropperResult, Preview } from 'vue-advanced-cropper'

const previewSizes = [50, 100, 150]
export type Direction = 'horizontal' | 'vertical'

withDefaults(
  defineProps<{
    result?: CropperResult
    direction?: Direction
  }>(),
  {
    direction: 'horizontal'
  }
)
</script>

<style scoped></style>
