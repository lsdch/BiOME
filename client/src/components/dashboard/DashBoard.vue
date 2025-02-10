<template>
  <v-container class="fill-height align-start" fluid>
    <v-row class="w-100">
      <v-col cols="12" md="6">
        <OccurrenceSunburst></OccurrenceSunburst>
      </v-col>
      <v-col>
        <v-card
          class="h-100"
          prepend-icon="mdi-sign-caution"
          title="Dashboard in progress"
        ></v-card>
      </v-col>
    </v-row>
  </v-container>
  <DatasetsDrawer />
</template>

<script setup lang="ts">
import DatasetsDrawer from './DatasetsDrawer.vue'
import OccurrenceSunburst from './OccurrenceSunburst.vue'
</script>

<style scoped></style>
