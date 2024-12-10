<template>
  <CRUDTable
    class="fill-height"
    :headers
    :fetch-items="ReferencesService.listArticles"
    :delete="({ code }: Article) => ReferencesService.deleteArticle({ path: { code } })"
    entity-name="Article"
    :toolbar="{ icon: 'mdi-newspaper-variant-multiple', title: 'Bibliography' }"
    append-actions
  >
    <template #item.authors="{ value }: { value: string[] }">
      {{
        value.length == 1 ? value[0] : value.length == 2 ? value.join(' & ') : `${value[0]} et. al`
      }}
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
import { Article, ReferencesService } from '@/api'
import ArticleFormDialog from '@/components/references/ArticleFormDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
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
