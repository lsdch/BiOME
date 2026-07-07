<template>
  <v-card
    title="Import occurrences batch"
    class="d-flex flex-column h-100"
    prepend-icon="mdi-file-upload"
  >
    <template #append>
      <v-btn color="primary" variant="tonal" text="Submit" prepend-icon="mdi-upload"></v-btn>
    </template>
    <v-divider />
    <v-container fluid class="bg-main flex-grow-1 overflow-y-auto"> </v-container>
  </v-card>
</template>

<script setup lang="ts"></script>

<style scoped lang="scss"></style>
