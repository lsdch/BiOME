<template>
  <v-select label="Delimiter" :items="delimiters" v-model="model" item-value="value">
    <template #selection="{ item }">
      {{ item.raw.subtitle }}
      &nbsp;
      <v-chip :text="item.title" label size="small" class="font-monospace"></v-chip>
    </template>
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #title>
          <v-chip :text="item.title" label size="small" class="font-monospace"></v-chip>
        </template>
        <template #append>
          <span class="text-caption text-muted">{{ item.raw.subtitle }}</span>
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { Delimiter } from '@/composables/data_exports'

const model = defineModel<Delimiter>('\t')

const delimiters = [
  { title: 'Tab', subtitle: 'Tabulation', value: '\t' },
  { title: ',', subtitle: 'Comma', value: ',' },
  { title: ';', subtitle: 'Semicolon', value: ';' }
]
</script>

<style scoped lang="scss"></style>
