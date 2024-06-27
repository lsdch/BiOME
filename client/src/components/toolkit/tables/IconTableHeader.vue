<template>
  <div class="d-flex flex-no-wrap">
    <template v-if="expanded">
      {{ column.title }}
    </template>
    <v-icon
      v-else
      :title="column.title"
      :color="color ?? (isSorted(column) ? 'primary' : 'secondary')"
      :icon="icon"
    />
    <v-icon v-if="isSorted(column)" :icon="getSortIcon(column)" />
  </div>
</template>

<script setup lang="ts">
type Column = { title?: string }
defineProps<{
  color?: string
  icon: string
  column: Column
  expanded?: boolean
  isSorted(column: Column): boolean
  getSortIcon(column: Column): IconValue
}>()
</script>

<style scoped></style>
