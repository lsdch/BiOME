<template>
  <v-tooltip text="Go to original GBIF record" location="top">
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="props"
        :size="size"
        icon
        :href="`https://www.gbif.org/species/${GBIF_ID}`"
        target="_blank"
      >
        <IconGBIF></IconGBIF>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import IconGBIF from '../icons/IconGBIF.vue'

type Props = {
  GBIF_ID: number
  size?: string
}
withDefaults(defineProps<Props>(), {
  size: 'small'
})
</script>

<style scoped></style>
