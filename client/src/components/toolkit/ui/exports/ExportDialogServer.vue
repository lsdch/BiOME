<template>
  <FormDialog v-model="dialog" :title btn-text="Download" @submit="emit('submit', options, suffix)">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <v-card-text class="py-3">
      <v-select :items="items" v-model="format" item-value="value" label="Format" />
      <v-text-field v-model="filename" label="Filename" :suffix="suffix" required />
      <ExportFormCSV v-if="format === 'csv'" v-model="csvOptions" />
    </v-card-text>
  </FormDialog>
</template>

<script setup lang="ts">
import { ExportOptions } from '@/lib/exports.ts'
import { computed, reactive } from 'vue'
import FormDialog from '../../forms/FormDialog.vue'
import ExportFormCSV from './ExportFormCSV.vue'
import { ExportFormat } from '@/api/adapters.ts'

const { title = 'Data export' } = defineProps<{
  title?: string
}>()

const dialog = defineModel<boolean>('dialog')

const items = [
  { title: 'JSON', value: 'json' },
  { title: 'CSV / TSV', value: 'csv' }
]

const format = defineModel<ExportFormat>('format', { default: 'json' })
const filename = defineModel<string>('filename')
const csvOptions = defineModel<ExportOptions['csvOptions']>('csvOptions', {
  default: () =>
    reactive({
      delimiter: '\t',
      quoteChar: '"'
    })
})

const options = computed<ExportOptions>(() => ({
  format: format.value,
  filename: filename.value ?? '',
  csvOptions: csvOptions.value
}))

const emit = defineEmits<{
  submit: [options: ExportOptions, suffix: string]
}>()

const suffix = computed<string>(() => {
  switch (format.value) {
    case 'json':
      return '.json'
    case 'csv':
      switch (csvOptions.value.delimiter) {
        case '\t':
          return '.tsv'
        default:
          return '.csv'
      }
    default:
      throw new Error(`Unknown format: ${format.value}`)
  }
})
</script>

<style scoped lang="scss"></style>
