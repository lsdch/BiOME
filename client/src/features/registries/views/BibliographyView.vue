<template>
  <CRUDTable
    class="fill-height"
    :headers
    entity-name="Publication"
    :toolbar="{ icon: 'mdi-newspaper-variant-multiple', title: 'Bibliography' }"
    append-actions
    v-model:search="search"
    :filter
    :items
    :error
    :loading
    @reload="refetch()"
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
    <template #item.title="{ value, item }">
      <span v-if="value">{{ value }}</span>
      <span v-else class="text-muted">{{ item.verbatim }}</span>
    </template>
    <template #expanded-row-inject="{ item }">
      <v-card
        class="article-details"
        flat
        :title="item.title ?? 'Untitled'"
        :subtitle="item.journal ?? 'Journal unknown'"
      >
        <template #append v-if="item.year">
          <v-chip label :text="item.year.toString()"></v-chip>
        </template>
        <v-card-text>
          <v-list-item prepend-icon="mdi-account-multiple">
            {{ item.authors?.join(', ') }}
          </v-list-item>
          <v-list-item class="font-monospace" prepend-icon="mdi-text-box">
            <template #subtitle>
              <code class="text-caption font-monospace">{{ item.verbatim }}</code>
            </template>
          </v-list-item>
        </v-card-text>
      </v-card>
    </template>
    <template #expanded-row-footer-append="{ item }">
      <DeleteBtn
        title="Delete publication ?"
        message="Deleted publication will be removed for all associated occurrences"
        @confirm="deleteMutation({ path: { id: item.id } }).then(() => refetch())"
      />
    </template>
    <!-- <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <ArticleFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :item="editItem"
        @close="onClose"
        @success="onSuccess"
      />
    </template> -->
    <!-- <template #footer.prepend-actions>
      <ArticlesImportDialog v-model="importDialog" />
      <v-btn
        color="primary"
        text="Import"
        variant="plain"
        prepend-icon="mdi-upload"
        size="small"
        @click="toggleImportDialog(true)"
      />
    </template> -->
  </CRUDTable>
</template>

<script setup lang="ts">
import { Article, Publication } from '@/api/adapters'
import {
  deletePublicationMutation,
  listPublicationsOptions
} from '@/api/gen/@tanstack/vue-query.gen'
// import ArticleFormDialogMutation from '@/components/forms/ArticleFormDialogMutation.vue'
// import ArticleFormDialogMutation from '@/components/forms/ArticleFormDialogMutation.vue'
// import ArticlesImportDialog from '@/features/registries/components/ArticlesImportDialog.vue'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import DeleteBtn from '@/components/toolkit/ui/DeleteBtn.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { ref } from 'vue'

type PubFilters = {
  term?: string
  year?: number
  author?: string
}

const [importDialog, toggleImportDialog] = useToggle(false)

const { data: items, refetch, error, isPending: loading } = useQuery(listPublicationsOptions())
const { mutateAsync: deleteMutation } = useMutation(deletePublicationMutation())

const search = ref<PubFilters>({})

function filter({ authors, year }: Publication) {
  const { author, year: searchYear } = search.value
  return (
    (!author || authors!.some((a) => a.toLowerCase().includes(author.toLowerCase()))) &&
    (!searchYear || year === searchYear)
  )
}

const headers: CRUDTableHeader<Publication>[] = [
  { key: 'authors', title: 'Authors' },
  { key: 'year', title: 'Year', width: 0 },
  {
    key: 'title',
    title: 'Title / Verbatim',
    sortRaw(i: Publication, j: Publication) {
      return (i.title ?? i.verbatim ?? '').localeCompare(j.title ?? j.verbatim ?? '')
    }
  }
]
</script>

<style lang="scss">
.article-details .v-card-title {
  font-size: 1rem;
}
</style>
