<template>
  <div>
    <v-toolbar class="bg-surface">
      <template #prepend>
        <v-btn
          variant="plain"
          icon="mdi-arrow-left"
          title="Back to taxonomy"
          :to="{ name: 'taxonomy' }"
        ></v-btn>
      </template>
      <template #title> Import clades from GBIF database </template>
    </v-toolbar>
    <ImportGBIF />
  </div>
</template>

<script setup lang="ts">
import ImportGBIF from '@/components/taxonomy/imports/ImportGBIF.vue'
</script>

<style scoped></style>
