<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerms"
    :items="filteredItems"
    :loading
    item-value="id"
    :multiple
    :chips="multiple"
    :closable-chips="multiple"
    no-filter
    clear-on-select
    placeholder="Enter search terms..."
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template #prepend-item>
      <div v-if="items && items.length" class="mx-3 text-caption text-center">
        {{
          searchTerms
            ? `${filteredItems.length} item(s) out of ${items.length} total`
            : `${items.length} items total`
        }}
      </div>
    </template>
    <template #chip="{ item, props }">
      <ArticleChip :article="item.obj" v-bind="props" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        v-bind="props"
        :title="Article.shortAuthors(item.obj.authors)"
        class="fuzzy-search-item"
      >
        <template #title>
          <v-list-item-title
            v-html="
              highlight(item, 'authors', {
                baseValue: Article.shortAuthors(item.obj.authors)
              })
            "
          />
        </template>
        <template #prepend="{ isSelected }">
          <v-checkbox :model-value="isSelected"></v-checkbox>
        </template>
        <template #subtitle>
          <v-list-item-subtitle style="max-width: 100%">
            <div class="d-flex flex-column" v-if="item.obj.title">
              <span v-html="highlight(item, 'title')"> </span>
              <div class="text-caption" v-html="highlight(item, 'journal')"></div>
            </div>
            <div v-else class="text-caption" v-html="highlight(item, 'verbatim')"></div>
          </v-list-item-subtitle>
        </template>
        <template #append>
          <v-chip size="small">
            <span v-html="highlight(item, 'year')"></span>
          </v-chip>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Article } from '@/api'
import { listArticlesOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useFuzzyItemsFilter } from '@/composables/fuzzy_search'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import ArticleChip from './ArticleChip'

const { threshold = 0.7, limit = 10 } = defineProps<{
  multiple?: boolean
  threshold?: number
  limit?: number
}>()

const { data: items, isPending: loading, error } = useQuery(listArticlesOptions())

const model = defineModel<Article[]>()
const searchTerms = ref<string>('')

const keys = [
  'year',
  'title',
  'journal',
  'verbatim',
  { key: 'authors', fn: (obj: Article) => Article.shortAuthors(obj.authors) }
] as const

const { highlight, filteredItems } = useFuzzyItemsFilter(keys, searchTerms, items, {
  threshold,
  limit
})
</script>

<style lang="scss"></style>
