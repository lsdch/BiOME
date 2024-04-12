<template>
  <div class="d-flex flex-no-wrap">
    <v-icon
      :title="column.title"
      :color="isSorted(column) ? 'primary' : 'secondary'"
      :icon="icon"
    />
    <v-icon v-if="isSorted(column)" :icon="getSortIcon(column)" />
  </div>
</template>

<script setup lang="ts">
type Column = { title?: string }
defineProps<{
  icon: string
  column: Column
  isSorted(column: Column): boolean
  getSortIcon(column: Column): IconValue
}>()
</script>

<style scoped></style>
