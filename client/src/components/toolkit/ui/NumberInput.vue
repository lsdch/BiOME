<template>
  <v-text-field
    v-model.number="model"
    v-bind="$attrs"
    type="number"
    @keydown="filterInput"
    :min="min"
    :max="max"
  />
</template>

<script setup lang="ts">
const model = defineModel<number>()
const props = defineProps<{ float?: boolean | number; min?: number; max?: number }>()

function filterInput(event: KeyboardEvent) {
  if (event.key.length === 1 && isNaN(Number(event.key)) && !(props.float && event.key == '.')) {
    event.preventDefault()
  }
}
</script>

<style scoped></style>
