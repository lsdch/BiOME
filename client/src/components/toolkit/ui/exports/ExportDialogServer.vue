<template>
  <FormDialog v-model="dialog" :title btn-text="Download" @submit="emit('submit', model, suffix)">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <v-card-text class="py-3">
      <v-select :items="items" v-model="model.format" item-value="value" label="Format" />
      <v-text-field
        :model-value="!!model.filename.length ? model.filename : filename"
        @update:model-value="(f) => (model.filename = f)"
        label="Filename"
        :suffix="suffix"
        required
      />
      <ExportFormCSV v-if="model.format === 'csv'" v-model="model.csvOptions" />
    </v-card-text>
  </FormDialog>
</template>

<script setup lang="ts">
import { ExportOptions } from '@/lib/exports.ts'
import { computed, reactive } from 'vue'
import FormDialog from '../../forms/FormDialog.vue'
import ExportFormCSV from './ExportFormCSV.vue'

const { title = 'Data export', filename = '' } = defineProps<{
  title?: string
  filename?: string
}>()

const dialog = defineModel<boolean>('dialog')

const items = [
  { title: 'JSON', value: 'json' },
  { title: 'CSV / TSV', value: 'csv' }
]

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

const suffix = computed<string>(() => {
  switch (model.value.format) {
    case 'json':
      return '.json'
    case 'csv':
      switch (model.value.csvOptions.delimiter) {
        case '\t':
          return '.tsv'
        default:
          return '.csv'
      }
    default:
      throw new Error(`Unknown format: ${model.value.format}`)
  }
})
</script>

<style scoped lang="scss"></style>
