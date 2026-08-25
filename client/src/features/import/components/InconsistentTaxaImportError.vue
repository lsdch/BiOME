<template>
  <v-list class="pa-0">
    <v-list-group :title="error.title">
      <template #activator="{ props }">
        <v-list-item
          :title="error.title"
          :subtitle="error.detail"
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
          :items="error.taxa"
          :headers="[
            { title: 'Taxon', value: 'name' },
            { title: 'Authorships', value: 'authorships' },
            { title: 'Ranks', value: 'ranks' }
          ]"
          v-model:page="page"
          v-model:items-per-page="itemsPerPage"
        >
          <template #item.authorships="{ value, index }: { value: string[]; index: number }">
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
</template>

<script setup lang="ts">
import { TaxonDefinition } from '@/api'
import { ref } from 'vue'

export type InconsistentTaxon = {
  name: string
  authorships: string[]
  ranks: string[]
}
export type InconsistentTaxaError = {
  title: string
  detail: string
  taxa: InconsistentTaxon[]
}

const itemsPerPage = ref(10)
const page = ref(1)

const { error } = defineProps<{
  error: InconsistentTaxaError
}>()

const taxonDefinitions = defineModel<TaxonDefinition[]>({ default: () => [] })
</script>

<style scoped lang="scss"></style>
