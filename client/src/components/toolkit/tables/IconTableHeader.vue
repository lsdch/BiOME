<template>
  <div class="d-flex flex-no-wrap justify-center">
    <template v-if="expanded">
      {{ column.title }}
    </template>
    <v-tooltip v-else location="top">
      <template #activator="{ props }">
        <v-icon
          v-bind="props"
          :title="column.title"
          :color="color ?? (isSorted(column) ? 'primary' : undefined)"
          :icon="icon"
        />
      </template>
      {{ column.title }}
    </v-tooltip>
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
