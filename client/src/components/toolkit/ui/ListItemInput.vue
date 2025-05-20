<template>
  <v-list-item :title="label" :subtitle>
    <template #append>
      <slot></slot>
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
const { label } = defineProps<{
  label?: string
  subtitle?: string
}>()
</script>

<style scoped lang="scss"></style>
