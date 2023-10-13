<template>
  <v-tooltip
    text="Go to original GBIF record"
    location="top"
    :disabled="!tooltip"
    :open-delay="700"
  >
    <template v-slot:activator="{ props }">
      <v-btn
        v-bind="{ ...$attrs, ...props }"
        :size="size"
        icon
        :href="`https://www.gbif.org/species/${GBIF_ID}`"
        target="_blank"
      >
        <IconGBIF size="50%"></IconGBIF>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import IconGBIF from '../icons/IconGBIF.vue'

withDefaults(
  defineProps<{
    GBIF_ID: number
    size?: string
    tooltip?: boolean
  }>(),
  {
    size: 'small',
    tooltip: true
  }
)
</script>

<style scoped></style>
