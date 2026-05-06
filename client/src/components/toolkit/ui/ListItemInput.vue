<template>
  <v-list-item :title="label" :subtitle :lines>
    <template #append>
      <slot></slot>
    </template>
  </v-list-item>
</template>

<script setup lang="ts">
const { label } = defineProps<{
  label?: string
  subtitle?: string
  lines?: false | 'one' | 'two' | 'three'
}>()
</script>

<style scoped lang="scss"></style>
