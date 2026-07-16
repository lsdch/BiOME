<template>
  <v-card title="Import occurrences batch" class="d-flex flex-column h-100" prepend-icon="mdi-file-upload">
    <v-divider />
    <v-form class="bg-main flex-grow-1 overflow-y-auto" @submit.prevent="submit">
      <v-container class="">
        <v-card>
          <v-card-text>
            <v-text-field v-model="model.label" label="Batch label" v-bind="schema('label')"></v-text-field>
            <v-combobox v-model="model.assembled_by" label="Assembled by" v-bind="schema('assembled_by')" multiple chips
              closable-chips clearable>
            </v-combobox>
            <v-textarea v-model="model.description" label="Description" v-bind="schema('description')"></v-textarea>
            <div class="d-flex ga-3">
              <v-file-input label="Occurrences TSV file" v-model="csv.file" class="required"
                :rules="[(v) => !!v || 'File is required']"
                filter-by-type="text/tab-separated-values, text/csv"></v-file-input>
              <CSVQuotePicker v-model="csv.quotes" class="flex-grow-0"></CSVQuotePicker>
            </div>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn text="Submit dataset" variant="tonal" size="large" prepend-icon="mdi-upload" type="submit"></v-btn>
          </v-card-actions>
        </v-card>
        <v-card v-if="error">
          <v-card-text>
            <p class="text-error">{{ error }}</p>
          </v-card-text>
        </v-card>

        <v-card>
          <v-list>
            <v-list-item v-for="batch in imports" :key="batch.workflow.import_id">
              {{ batch }}
            </v-list-item>
          </v-list>
        </v-card>
      </v-container>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import {
  importOccurrencesCsvMutation,
  listImportsForCurrentUserOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import { ImportDataCSV } from '../components/FileInputCSV.vue'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'
import { useRouter } from 'vuetify/lib/composables/router.mjs'
import { ImportWorkflowInput } from '@/api/adapters.ts'
import { useSchema } from '@/composables/schema.ts'
import { $ImportWorkflowInput } from '@/api/index.ts'

const model = ref<ImportWorkflowInput>({ label: '', description: undefined, assembled_by: [] })
const { bind: { schema } } = useSchema($ImportWorkflowInput)
const csv = ref<ImportDataCSV>({ file: undefined, separator: '\t', quotes: '"' })

const { mutateAsync, error } = useMutation(importOccurrencesCsvMutation())

const { data: imports } = useQuery(listImportsForCurrentUserOptions())

const router = useRouter()

async function submit() {
  if (!model.value.label || !csv.value.file) {
    throw new Error('Label and CSV file are required')
  }

  await mutateAsync(
    {
      body: {
        workflow: model.value,
        file: csv.value.file,
        separator: csv.value.separator,
        quotes: csv.value.quotes
      }
    },
    {
      onSuccess(data) {
        router?.push({ name: 'import-batch-item', params: { uuid: data.import_id } })
      }
    }
  )
}
</script>

<style scoped lang="scss"></style>
