<template>
  <v-container class="fill-height align-start" fluid>
    <v-row class="w-100">
      <v-col cols="12" lg="6">
        <OccurrenceSunburst></OccurrenceSunburst>
      </v-col>
      <v-col cols="12" lg="6">
        <CountryOverview></CountryOverview>
      </v-col>
    </v-row>
  </v-container>
  <DatasetsDrawer />
</template>

<script setup lang="ts">
import CountryOverview from './CountryOverview.vue'
import DatasetsDrawer from './DatasetsDrawer.vue'
import OccurrenceSunburst from './OccurrenceSunburst.vue'
</script>

<style scoped></style>
