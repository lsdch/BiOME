<template>
  <v-row>
    <v-col>
      <HabitatPicker />
    </v-col>
  </v-row>
  <v-row>
    <v-col cols="12" md="6">
      <v-skeleton-loader class="mx-auto border" max-width="300" type="card-avatar, actions" />
    </v-col>

    <v-col cols="12" md="6">
      <v-skeleton-loader class="mx-auto border" max-width="300" type="image, article" />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import HabitatPicker from '../habitat/HabitatPicker.vue'
</script>

<style scoped></style>
