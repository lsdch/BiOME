<template>
  <v-chip-group v-model="searchBy" mandatory color="primary">
    <v-chip value="name" text="By name / code"></v-chip>
    <v-chip value="coordinates" text="By coords"></v-chip>
  </v-chip-group>
  <SiteAutocomplete v-if="searchBy === 'name'" v-model="model" :disabled-codes class="pt-2" />
  <SiteSearchCoordinates v-else v-model="model" :disabled-codes class="pt-2" />
</template>

<script setup lang="ts">
import { Site } from '@/api'
import SiteAutocomplete from './SiteAutocomplete.vue'
import SiteSearchCoordinates from './SiteSearchCoordinates.vue'
import { ref } from 'vue'

const searchBy = ref<'name' | 'coordinates'>('name')

const model = defineModel<Site>()

const { disabledCodes } = defineProps<{
  disabledCodes?: string[]
}>()
</script>

<style scoped lang="scss"></style>
