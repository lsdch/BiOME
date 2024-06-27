<template>
  <v-row>
    <v-col>
      <v-select
        label="Delimiter"
        v-model="model.delimiter"
        item-value="value"
        item-title="label"
        placeholder="auto"
        persistent-placeholder
        clearable
        :items="[
          { value: '\t', label: 'Tab', symbol: '⇥' },
          { value: ',', label: 'Comma', symbol: ',' },
          { value: ';', label: 'Semicolon', symbol: ';' }
        ]"
      >
        <template #item="{ item, props }">
          <v-list-item v-bind="props" :title="item.title">
            <template #append>
              {{ item.raw.symbol }}
            </template>
          </v-list-item>
        </template>
      </v-select>
    </v-col>
    <v-col>
      <v-select
        label="Quote"
        v-model="model.quoteChar"
        item-value="value"
        item-title="label"
        placeholder="auto"
        clearable
        persistent-placeholder
        :items="[
          { value: '&quot;', label: '&quot;' },
          { value: '\'', label: '\'' }
        ]"
      />
    </v-col>
    <v-col>
      <v-select
        label="New line"
        v-model="model.newline"
        clearable
        item-value="value"
        item-title="label"
        placeholder="auto"
        persistent-placeholder
        :items="[
          { value: '\n', label: '\\n' },
          { value: '\r', label: '\\r' },
          { value: '\r\n', label: '\\r\\n' }
        ]"
      />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ParseConfig } from 'papaparse'

const model = defineModel<ParseConfig<any, undefined>>({ required: true })
</script>

<style scoped></style>
