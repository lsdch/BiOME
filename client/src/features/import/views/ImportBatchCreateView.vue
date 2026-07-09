<template>
  <v-card
    title="Import occurrences batch"
    class="d-flex flex-column h-100"
    prepend-icon="mdi-file-upload"
  >
    <v-divider />
    <v-form class="bg-main flex-grow-1" @submit.prevent="submit">
      <v-container class="overflow-y-auto">
        <v-card>
          <v-card-text>
            <v-text-field v-model="label" label="Batch label"></v-text-field>
            <div class="d-flex ga-3">
              <v-file-input label="Occurrences TSV file" v-model="csv.file"></v-file-input>
              <CSVQuotePicker v-model="csv.quotes" class="flex-grow-0"></CSVQuotePicker>
            </div>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              text="Submit dataset"
              variant="tonal"
              size="large"
              prepend-icon="mdi-upload"
              type="submit"
            ></v-btn>
          </v-card-actions>
        </v-card>
      </v-container>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import { importOccurrencesCsvMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation } from '@tanstack/vue-query'
import { ref } from 'vue'
import { ImportDataCSV } from '../components/FileInputCSV.vue'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'

const label = ref<string>()
const csv = ref<ImportDataCSV>({ file: undefined, separator: '\t', quotes: '"' })

const { mutateAsync } = useMutation(importOccurrencesCsvMutation())

async function submit() {
  if (!label.value || !csv.value.file) {
    throw new Error('Label and CSV file are required')
  }

  await mutateAsync({
    body: {
      label: label.value,
      file: csv.value.file,
      separator: csv.value.separator,
      quotes: csv.value.quotes
    }
  })
}
</script>

<style scoped lang="scss"></style>
