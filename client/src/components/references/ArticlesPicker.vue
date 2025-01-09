<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerms"
    :items="filteredItems"
    :loading
    item-value="code"
    :multiple
    :chips="multiple"
    :closable-chips="multiple"
    no-filter
    clear-on-select
    placeholder="Enter search terms..."
    v-bind="$attrs"
  >
    <template #prepend-item>
      <div v-if="items.length" class="mx-3 text-caption text-center">
        {{
          searchTerms
            ? `${filteredItems.length} item(s) out of ${items.length} total`
            : `${items.length} items total`
        }}
      </div>
    </template>
    <template #chip="{ item, props }">
      <ArticleChip :article="item.raw.obj" v-bind="props" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        v-bind="props"
        :title="Article.shortAuthors(item.raw.obj.authors)"
        class="article-picker-item"
      >
        <template #title>
          <v-list-item-title
            v-html="
              highlight(item.raw, 'authors', {
                baseValue: Article.shortAuthors(item.raw.obj.authors)
              })
            "
          />
        </template>
        <template #prepend="{ isSelected }">
          <v-checkbox :model-value="isSelected"></v-checkbox>
        </template>
        <template #subtitle>
          <v-list-item-subtitle style="max-width: 100%">
            <div class="d-flex flex-column">
              <span v-html="highlight(item.raw, 'title')"> </span>
              <div class="text-caption" v-html="highlight(item.raw, 'journal')"></div>
            </div>
          </v-list-item-subtitle>
        </template>
        <template #append>
          <v-chip size="small">
            <span v-html="highlight(item.raw, 'year')"></span>
          </v-chip>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Article, ReferencesService } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'
import { useFuzzyItemsFilter } from '@/composables/fuzzy_search'
import { ref } from 'vue'
import ArticleChip from './ArticleChip.vue'

const { threshold, limit } = withDefaults(
  defineProps<{ multiple?: boolean; threshold?: number; limit?: number }>(),
  {
    limit: 10,
    threshold: 0.7
  }
)

const { items, loading } = useFetchItems(ReferencesService.listArticles)

const model = defineModel<string | string[] | null>()
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

<style lang="scss">
@use 'vuetify';
.article-picker-item b {
  color: rgb(var(--v-theme-primary));
}
</style>
