<template>
  <v-row>
    <v-col cols="12" sm="">
      <CSVDelimiterPicker v-model="model.delimiter" />
    </v-col>
    <v-col cols="12" sm="">
      <CSVQuotePicker v-model="model.quoteChar" />
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
// import { Delimiter, QuoteChar } from '@/composables/data_exports'
import { CsvDelimiter, CsvQuoteChar } from '@/api/adapters.ts'
import { reactive } from 'vue'
import CSVDelimiterPicker from './CSVDelimiterPicker.vue'
import CSVQuotePicker from './CSVQuotePicker.vue'

type ExportCSVOptions = {
  delimiter: CsvDelimiter
  quoteChar: CsvQuoteChar
}

const model = defineModel<ExportCSVOptions>({
  default: () =>
    reactive({
      delimiter: '\t',
      quoteChar: '"'
    })
})
</script>

<style scoped lang="scss"></style>
