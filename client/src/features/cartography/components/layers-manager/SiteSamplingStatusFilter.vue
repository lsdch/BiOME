<template>
  <ListItemInput title="Include sites">
    <v-chip-group v-model="model">
      <v-chip
        v-for="option in options"
        :key="option"
        :value="option"
        color="success"
        variant="tonal"
        rounded="md"
        v-tooltip="{ text: tooltip(option), maxWidth: 300 }"
      >
        {{ option }}
      </v-chip>
    </v-chip-group>
  </ListItemInput>
</template>

<script setup lang="ts">
import { SiteSamplingStatus } from '@/api'
import ListItemInput from '@/components/toolkit/ui/ListItemInput.vue'

const options: SiteSamplingStatus[] = ['Occurrences', 'Sampled', 'All'] as const

const model = defineModel<SiteSamplingStatus>({
  required: true
})

function tooltip(value: SiteSamplingStatus) {
  switch (value) {
    case 'Occurrences':
      return 'Only include sites having occurrence records after filtering. This is the default option and the one to use in most cases.'
    case 'Sampled':
      return 'Only include sites having sampling events after filtering, regardless of whether they have occurrence records or not. Only relevant '
    case 'All':
      return 'All sites, regardless of their sampling status'
  }
}
</script>

<style scoped lang="scss"></style>
