<template>
  <v-select
    v-model="model"
    label="Habitat"
    multiple
    chips
    clearable
    :items="items"
    item-title="name"
    item-value="name"
    return-object
    density="compact"
  ></v-select>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type HabitatQualifier = {
  name: string
  depends?: string | string[]
  incompatible?: string | string[]
  color?: string
}

const habitats = ref<HabitatQualifier[]>([
  { name: 'Aquatic', incompatible: 'Terrestrial', color: 'primary' },
  { name: 'Lotic', depends: 'Aquatic', incompatible: 'Lentic' },
  { name: 'Lentic', depends: 'Aquatic', incompatible: 'Lotic' },
  { name: 'Freshwater', depends: 'Aquatic', incompatible: ['Saltwater', 'Brackish water'] },
  { name: 'Saltwater', depends: 'Aquatic', incompatible: ['Freshwater', 'Brackish water'] },
  { name: 'Brackish water', depends: 'Aquatic', incompatible: ['Saltwater', 'Freshwater'] },
  { name: 'Terrestrial', incompatible: 'Aquatic' },
  { name: 'Surface', incompatible: 'Subsurface' },
  { name: 'Subsurface', incompatible: 'Surface' },
  { name: 'Aquifer', depends: 'Subsurface' },
  { name: 'Saturated', depends: 'Aquifer', incompatible: 'Unsaturated' },
  { name: 'Unsaturated', depends: 'Aquifer', incompatible: 'Saturated' },
  {
    name: 'Alluvial',
    depends: 'Aquifer',
    incompatible: ['Porous', 'Fissured', 'Hyporheic zone', 'Karst']
  },
  {
    name: 'Hyporheic zone',
    depends: 'Aquifer',
    incompatible: ['Porous', 'Fissured', 'Hyporheic zone', 'Karst']
  },
  { name: 'Karst', depends: 'Aquifer', incompatible: ['Alluvial', 'Hyporheic zone'] },
  { name: 'Porous', depends: 'Aquifer', incompatible: ['Alluvial', 'Hyporheic zone'] },
  { name: 'Fissured', depends: 'Aquifer', incompatible: ['Alluvial', 'Hyporheic zone'] }
])

const model = ref<HabitatQualifier[]>([])

const items = computed(() => {
  return habitats.value.filter((h) => {
    const isCompatible = !model.value.find(
      ({ name }) => name === h.incompatible || h.incompatible?.includes(name)
    )
    return (
      isCompatible &&
      (h.depends === undefined ||
        model.value?.find(({ name }) => h.depends === name || h.depends?.includes(name)))
    )
  })
})
</script>

<style scoped></style>
