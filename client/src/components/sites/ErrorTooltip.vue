<template>
  <v-tooltip :disabled="error === undefined">
    <template #activator="{ props }">
      <div v-bind="props" :class="{ 'text-error': error }">
        <slot> </slot>
      </div>
    </template>
    {{ error }}
  </v-tooltip>
</template>

<script setup lang="ts">
defineProps<{ error?: string }>()
</script>

<style scoped></style>
