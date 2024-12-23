<template>
  <v-container class="fill-height bg-surface" fluid>
    <v-row>
      <v-col cols="12" md="6">
        <OccurrenceSunburst />
      </v-col>
      <v-col cols="12" md="6">
        <v-skeleton-loader class="mx-auto border" max-width="300" type="card-avatar, actions" />
      </v-col>

      <v-col cols="12" md="6">
        <v-skeleton-loader class="mx-auto border" max-width="300" type="image, article" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { useInstanceSettings } from '../settings'
import OccurrenceSunburst from './OccurrenceSunburst.vue'

const { settings } = useInstanceSettings()
</script>

<style scoped></style>
