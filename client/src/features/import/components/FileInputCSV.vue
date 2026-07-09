<template>
  <v-file-input v-model="csv.file" :label></v-file-input>
  <div class="d-flex ga-3">
    <CSVQuotePicker v-model="csv.quotes"></CSVQuotePicker>
    <CSVDelimiterPicker v-model="csv.separator"></CSVDelimiterPicker>
  </div>
</template>

<script setup lang="ts">
import CSVDelimiterPicker from '@/components/toolkit/ui/exports/CSVDelimiterPicker.vue'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'

export type ImportDataCSV = {
  file?: File
  separator: string
  quotes: string
}

const csv = defineModel<ImportDataCSV>({
  default: () => ({ file: null, separator: ',', quotes: '"' })
})

defineProps<{
  label?: string
}>()
</script>

<style scoped lang="scss"></style>
