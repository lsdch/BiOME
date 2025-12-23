<template>
  <v-container class="fill-height align-start" fluid>
    <v-row class="w-100">
      <v-col cols="12" lg="6">
        <OccurrenceSunburst :items :error :loading="isPending"></OccurrenceSunburst>
      </v-col>
      <v-col cols="12" lg="6">
        <CountryOverview></CountryOverview>
      </v-col>
    </v-row>
  </v-container>
  <DatasetsDrawer />
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import CountryOverview from '../components/CountryOverview.vue'
import DatasetsDrawer from '../components/DatasetsDrawer.vue'
import OccurrenceSunburst from '../components/OccurrenceSunburst.vue'
import { occurrenceOverviewOptions } from '@/api/gen/@tanstack/vue-query.gen'

const { data: items, error, isPending } = useQuery(occurrenceOverviewOptions())
</script>

<style scoped></style>
