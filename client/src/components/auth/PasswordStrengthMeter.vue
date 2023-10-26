<template>
  <v-progress-linear
    :model-value="strength + 1"
    max="5"
    :color="colors[strength]"
    height="15px"
    v-bind="$attrs"
  >
    <template v-slot:default>
      <span class="meter-label">{{ labels[strength] }}</span>
    </template>
  </v-progress-linear>
</template>

<script setup lang="ts">
defineProps<{ strength: number }>()

const colors = ['red', 'orange', 'yellow', 'lime', 'green']
const labels = ['Very weak', 'Weak', 'Acceptable', 'Good', 'Very good']
</script>

<style scoped>
.meter-label {
  font-size: 10px;
}
</style>
