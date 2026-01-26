<template>
  <v-text-field v-model="search.term" label="Search" class="ma-2" clearable />
  <CRUDTable
    entity-name="Bibliography"
    :items="dataset?.bibliography ?? []"
    :search
    :headers="[
      {
        key: 'authors',
        title: 'Authors',
        filter(value, query, item) {
          return value.some((author: string) => author.toLowerCase().includes(query.toLowerCase()))
        }
      },
      {
        key: 'year',
        title: 'Year',
        cellProps: { class: 'font-monospace' },
        width: 0,
        align: 'end'
      },
      {
        key: 'title',
        title: 'Title / Verbatim',
        value(item, fallback) {
          return item.title ?? item.verbatim
        }
      }
    ]"
  >
    <template #item.authors="{ value: authors }: { value: string[] }">
      <v-menu location="top left" origin="top left">
        <template #activator="{ props }">
          <v-chip v-bind="props">{{ Article.shortAuthors(authors) }}</v-chip>
        </template>
        <v-card>
          <v-card-text>
            {{ authors.join(', ') }}
          </v-card-text>
        </v-card>
      </v-menu>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { Article, OccurrenceDataset } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { ref } from 'vue'
const { dataset } = defineProps<{ dataset: OccurrenceDataset }>()
const search = ref({ term: undefined })
</script>

<style scoped lang="scss"></style>
