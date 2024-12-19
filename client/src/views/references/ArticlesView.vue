<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="ReferencesService.listArticles"
    :delete="({ code }: Article) => ReferencesService.deleteArticle({ path: { code } })"
    entity-name="Article"
    :toolbar="{ icon: 'mdi-newspaper-variant-multiple', title: 'Bibliography' }"
    append-actions
    v-model:search="search"
    :filter
  >
    <template #menu>
      <v-divider class="mb-2"></v-divider>
      <v-row class="mb-2">
        <v-col cols="12" md="8">
          <v-list-item>
            <v-text-field
              class="mt-1"
              v-model="search.author"
              label="Author"
              density="compact"
              hide-details
              clearable
            />
          </v-list-item>
        </v-col>
        <v-col cols="12" md="4">
          <v-list-item>
            <v-number-input
              class="mt-1"
              v-model="search.year"
              label="Year"
              density="compact"
              hide-details
              clearable
            />
          </v-list-item>
        </v-col>
      </v-row>
    </template>
    <template #item.authors="{ value }: { value: string[] }">
      {{ Article.shortAuthors(value) }}
    </template>
    <template #expanded-row-inject="{ item }">
      <v-card
        class="article-details"
        flat
        :title="item.title ?? 'Untitled'"
        :subtitle="item.journal ?? 'Journal unknown'"
      >
        <template #append>
          <v-chip label :text="item.year.toString()"></v-chip>
        </template>
        <v-card-text>
          <v-list-item prepend-icon="mdi-account-multiple">
            {{ item.authors.join(', ') }}
          </v-list-item>
          <v-list-item class="font-monospace" prepend-icon="mdi-text-box" :title="item.code">
            <template #subtitle>
              <code class="text-caption font-monospace">{{ item.verbatim }}</code>
            </template>
          </v-list-item>
        </v-card-text>
      </v-card>
    </template>
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <ArticleFormDialog
        :model-value="dialog"
        @close="onClose"
        @success="onSuccess"
        :edit="editItem"
      />
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { ReferencesService } from '@/api'
import { Article } from '@/api/adapters'
import ArticleFormDialog from '@/components/references/ArticleFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { ref } from 'vue'

type ArticleFilters = {
  term?: string
  year?: number
  author?: string
}

const search = ref<ArticleFilters>({})

function filter({ authors, year }: Article) {
  const { author, year: searchYear } = search.value
  return (
    (!author || authors.some((a) => a.toLowerCase().includes(author.toLowerCase()))) &&
    (!searchYear || year === searchYear)
  )
}

const headers: CRUDTableHeader[] = [
  { key: 'authors', title: 'Authors' },
  { key: 'year', title: 'Year', width: 0 },
  { key: 'journal', title: 'Journal', mobile: false }
]
</script>

<style lang="scss">
.article-details .v-card-title {
  font-size: 1rem;
}
</style>
