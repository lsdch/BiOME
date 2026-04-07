<template>
  <v-select label="Quote char" :items="quoteChars" v-model="model" item-value="value" :disabled>
    <template #selection="{ item }">
      {{ item.subtitle }}
      &nbsp;
      <v-chip :text="item.title" label size="small" class="font-monospace"></v-chip>
    </template>
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #title>
          <v-chip :text="item.title" label size="small" class="font-monospace"></v-chip>
        </template>
        <template #append>
          <span class="text-caption text-muted">{{ item.subtitle }}</span>
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { QuoteChar } from '@/composables/data_exports'

defineProps<{
  disabled?: boolean
}>()

const model = defineModel<QuoteChar>('"')
const quoteChars = [
  { title: '"', subtitle: 'Double', value: '"' },
  { title: "'", subtitle: 'Single', value: "'" }
]
</script>

<style scoped lang="scss"></style>
