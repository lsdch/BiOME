<template>
  <v-select
    v-model="model"
    label="Palette"
    :items="palettes"
    item-title="name"
    item-value="name"
    density="compact"
  >
    <template #append-inner>
      <div style="width: 100px">
        <v-sparkline
          v-if="model"
          :gradient="palette(model)"
          gradient-direction="left"
          :model-value="Array.from({ length: 10 }, (_, i) => 0)"
          :line-width="100"
          :height="50"
        />
      </div>
    </template>
    <template #item="{ item, props }">
      <v-list-item :title="item.title" v-bind="props">
        <template #append>
          <div style="width: 100px">
            <v-sparkline
              :gradient="item.raw.palette"
              gradient-direction="left"
              :model-value="Array.from({ length: 10 }, (_, i) => 0)"
              :line-width="50"
            />
          </div>
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import {
  brewerPalettes,
  ColorBrewerPalette,
  ColorBrewerPaletteKey,
  palette
} from '@/functions/color_brewer'

const model = defineModel<keyof typeof brewerPalettes>({
  default: 'Viridis'
})

const palettes = Object.entries(brewerPalettes).map(([key, value]) => ({
  name: key as ColorBrewerPaletteKey,
  palette: value as ColorBrewerPalette
}))
</script>

<style scoped lang="scss"></style>
