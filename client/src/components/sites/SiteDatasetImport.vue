<template>
  <v-form>
    <v-text-field label="Dataset name" />
    <PersonPicker label="Maintainers" multiple></PersonPicker>
    <v-textarea label="Description" variant="outlined" />
  </v-form>
</template>

<script setup lang="ts">
import PersonPicker from '../people/PersonPicker.vue'
</script>

<style scoped></style>
