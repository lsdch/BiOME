<template>
  <v-card
    title="Import occurrences batch"
    class="d-flex flex-column h-100"
    prepend-icon="mdi-file-upload"
  >
    <v-divider />
    <v-tabs v-model="tab" class="flex-shrink-0">
      <v-tab value="new">New batch</v-tab>
      <v-tab value="existing">Existing batches</v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab" class="flex-grow-1 overflow-y-auto">
      <v-tabs-window-item value="new">
        <v-form class="bg-main flex-grow-1" @submit.prevent="submit">
          <v-container class="d-flex flex-column ga-3">
            <v-card>
              <v-card-text>
                <v-text-field
                  v-model="model.label"
                  label="Batch label"
                  v-bind="schema('label')"
                ></v-text-field>
                <v-combobox
                  v-model="model.assembled_by"
                  label="Assembled by"
                  v-bind="schema('assembled_by')"
                  multiple
                  chips
                  closable-chips
                  clearable
                >
                </v-combobox>
                <v-textarea
                  v-model="model.description"
                  label="Description"
                  v-bind="schema('description')"
                ></v-textarea>
                <GBIFKingdomPicker
                  v-model="model.taxonomic_scope"
                  label="Taxonomic scope"
                  hint="The kingdom to which belong occurrences in the imported batch. This is required to unambiguously resolve taxa from GBIF."
                  persistent-hint
                  item-value="key"
                  class="flex-grow-0"
                  v-bind="schema('taxonomic_scope')"
                ></GBIFKingdomPicker>
                <div class="d-flex ga-3 mt-5">
                  <v-file-input
                    label="Occurrences TSV file"
                    v-model="csv.file"
                    class="required"
                    :rules="[(v) => !!v || 'File is required']"
                    filter-by-type="text/tab-separated-values, text/csv"
                  ></v-file-input>
                  <CSVQuotePicker v-model="csv.quotes" class="flex-grow-0"></CSVQuotePicker>
                </div>
              </v-card-text>
              <v-alert v-if="error" type="error" :title="error.title" :text="error.detail" closable>
                <ul>
                  <li v-for="e in error.errors">
                    {{ e.message }}
                  </li>
                </ul></v-alert
              >
              <v-list v-if="inconsistentTaxaError" class="pa-0">
                <v-list-group :title="inconsistentTaxaError.title">
                  <template #activator="{ props }">
                    <v-list-item
                      :title="inconsistentTaxaError.title"
                      :subtitle="inconsistentTaxaError.detail"
                      lines="three"
                      color="error"
                      class="text-error"
                      prepend-icon="mdi-alert-circle"
                      v-bind="props"
                    ></v-list-item>
                  </template>
                  <template #default>
                    <v-divider></v-divider>
                    <v-data-table
                      :items="inconsistentTaxaError.taxa"
                      :headers="[
                        { title: 'Taxon', value: 'name' },
                        { title: 'Authorships', value: 'authorships' },
                        { title: 'Ranks', value: 'ranks' }
                      ]"
                      v-model:page="page"
                      v-model:items-per-page="itemsPerPage"
                    >
                      <template
                        #item.authorships="{ value, index }: { value: string[]; index: number }"
                      >
                        <v-chip-group
                          v-model="taxonDefinitions[(page - 1) * itemsPerPage + index].authorship"
                          mandatory
                          color="primary"
                          class="d-flex ga-2"
                        >
                          <v-chip
                            v-for="(authorship, index) in value"
                            :key="index"
                            :value="authorship"
                            :text="authorship"
                            filter
                            label
                            variant="tonal"
                          >
                          </v-chip>
                        </v-chip-group>
                      </template>
                      <template #item.ranks="{ value, index }: { value: string[]; index: number }">
                        <div class="d-flex ga-2">
                          <v-chip-group
                            v-model="taxonDefinitions[(page - 1) * itemsPerPage + index].rank"
                            mandatory
                            color="primary"
                            class="d-flex ga-2"
                          >
                            <v-chip
                              v-for="(rank, index) in value"
                              :key="index"
                              :text="rank"
                              :value="rank"
                              class="text-capitalize"
                              filter
                              label
                              variant="tonal"
                            >
                            </v-chip>
                          </v-chip-group>
                        </div>
                      </template>
                    </v-data-table>
                  </template>
                </v-list-group>
              </v-list>
              <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn
                  :loading="isImporting"
                  text="Submit dataset"
                  variant="tonal"
                  size="large"
                  prepend-icon="mdi-upload"
                  type="submit"
                ></v-btn>
              </v-card-actions>
            </v-card>
            <v-card v-if="error">
              {{ error }}
            </v-card>
          </v-container>
        </v-form>
      </v-tabs-window-item>
      <v-tabs-window-item value="existing">
        <ImportBatchesTable />
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup lang="ts">
import {
  importOccurrencesCsvMutation,
  listImportsForCurrentUserOptions
} from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { ref, watch } from 'vue'
import { ImportDataCSV } from '../components/FileInputCSV.vue'
import CSVQuotePicker from '@/components/toolkit/ui/exports/CSVQuotePicker.vue'
import { useRouter } from 'vuetify/lib/composables/router.mjs'
import { ImportBatchInput, TaxonDefinition, TaxonRank } from '@/api/adapters.ts'
import { useSchemaBinding } from '@/composables/schema.ts'
import { $ImportBatchInput } from '@/api/index.ts'
import ImportBatchesTable from '../components/ImportBatchesTable.vue'
import GBIFKingdomPicker from '../components/GBIFKingdomPicker.vue'

const tab = ref<'new' | 'existing'>('new')

const model = ref<ImportBatchInput>({
  label: '',
  description: undefined,
  assembled_by: [],
  taxonomic_scope: 1
})
const taxonDefinitions = ref<TaxonDefinition[]>([])
const itemsPerPage = ref(10)
const page = ref(1)

type InconsistentTaxon = {
  name: string
  authorships: string[]
  ranks: string[]
}
const inconsistentTaxaError = ref<{
  title: string
  detail: string
  taxa: InconsistentTaxon[]
}>()

const { schema } = useSchemaBinding($ImportBatchInput)
const csv = ref<ImportDataCSV>({ file: undefined, separator: '\t', quotes: '"' })

const { mutateAsync, error, isPending: isImporting } = useMutation(importOccurrencesCsvMutation())

const router = useRouter()

async function submit() {
  if (!model.value.label || !csv.value.file) {
    throw new Error('Label and CSV file are required')
  }

  await mutateAsync(
    {
      body: {
        batch: model.value,
        file: csv.value.file,
        separator: csv.value.separator,
        quotes: csv.value.quotes,
        taxon_definitions: taxonDefinitions.value
      }
    },
    {
      onSuccess(data) {
        router?.push({ name: 'import-batch-item', params: { uuid: data.id } })
      },
      onError(err) {
        if (err.code === 'inconsistent_taxa') {
          inconsistentTaxaError.value = {
            title: err.title!,
            detail: err.detail!,
            taxa: err.errors?.map(({ value }) => value as InconsistentTaxon) ?? []
          }
          taxonDefinitions.value = inconsistentTaxaError.value.taxa.map(
            ({ name, authorships, ranks }) => ({
              name,
              authorship: authorships.length > 0 ? authorships[0] : undefined,
              rank: ranks.length > 0 ? (ranks[0] as TaxonRank) : undefined
            })
          )
          tab.value = 'new'
        }
      }
    }
  )
}
</script>

<style scoped lang="scss"></style>
