<template>
  <v-text-field
    v-model="searchTerm"
    id="searchbar"
    density="compact"
    name="search"
    label="Search"
    single-line
    hide-details
    clearable
    color="primary"
    v-bind="$attrs"
  >
    <template v-slot:[`prepend-inner`]="{ isActive }">
      <v-icon :color="isActive.value ? 'primary' : 'grey'">mdi-magnify</v-icon>
    </template>
  </v-text-field>
</template>

<script setup lang="ts">
const searchTerm = defineModel()
</script>

<style scoped></style>
