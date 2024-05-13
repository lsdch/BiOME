<template>
  <v-tooltip>
    <template #activator="{ props }">
      <v-btn v-bind="{ ...$attrs, ...props }" />
    </template>
    {{ tooltip }}
  </v-tooltip>
</template>

<script setup lang="ts">
defineProps<{ tooltip: string }>()
</script>

<style scoped></style>
