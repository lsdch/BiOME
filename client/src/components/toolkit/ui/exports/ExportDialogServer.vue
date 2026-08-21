<template>
  <FormDialog v-model="dialog" :title btn-text="Download" @submit="emit('submit', model, suffix)">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <v-card-text class="py-3">
      <v-select :items="items" v-model="model.format" item-value="value" label="Format" />
      <v-text-field v-model="model.filename" label="Filename" :suffix="suffix" required />
      <ExportFormCSV v-if="model.format === 'csv'" v-model="model.csvOptions" />
    </v-card-text>
  </FormDialog>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import FormDialog from '../../forms/FormDialog.vue'
import ExportFormCSV from './ExportFormCSV.vue'
import { CsvDelimiter, CsvQuoteChar } from '@/api/adapters.ts'

const { title = 'Data export' } = defineProps<{ title?: string }>()

const dialog = defineModel<boolean>('dialog')

const items = [
  { title: 'JSON', value: 'json' },
  { title: 'CSV / TSV', value: 'csv' }
]

type ExportFormat = 'json' | 'csv'

export type ExportOptions = {
  format: ExportFormat
  filename: string
  csvOptions: {
    delimiter: CsvDelimiter
    quoteChar: CsvQuoteChar
  }
}

const model = defineModel<ExportOptions>({
  default: () =>
    reactive({
      format: 'json',
      filename: '',
      csvOptions: {
        delimiter: '\t',
        quoteChar: '"'
      }
    })
})

const emit = defineEmits<{
  submit: [options: ExportOptions, suffix: string]
}>()

const suffix = computed(() => {
  switch (model.value.format) {
    case 'json':
      return '.json'
    case 'csv':
      switch (model.value.csvOptions.delimiter) {
        case '\t':
          return '.tsv'
        case ',':
        case ';':
          return '.csv'
        default:
          throw new Error(`Unsupported delimiter: ${model.value.csvOptions.delimiter}`)
      }
  }
})
</script>

<style scoped lang="scss"></style>
