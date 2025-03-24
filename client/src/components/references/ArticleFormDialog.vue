<template>
  <CreateUpdateForm v-model="item" :create :update @success="dialog = false">
    <template #default="{ model, field, mode, loading, submit, setModel }">
      <FormDialog
        v-model="dialog"
        :title="`${mode} article`"
        :loading="loading.value"
        @submit="submit"
      >
        <v-container fluid>
          <v-row>
            <v-col cols="12" md="9">
              <DoiInputFetcher v-model="model.doi" @fetched="(article) => setModel(article)" />
            </v-col>
            <v-col cols="12" md="3">
              <v-switch label="Manual input" v-model="manualMode" color="primary"></v-switch>
            </v-col>
          </v-row>
          <template v-if="manualMode">
            <v-combobox
              v-model="model.authors"
              label="Authors"
              multiple
              chips
              closable-chips
              v-bind="field('authors')"
            />

            <v-text-field
              v-model="model.title"
              v-bind="field('title')"
              label="Title"
            ></v-text-field>
            <v-row>
              <v-col cols="12" md="9">
                <v-text-field v-model="model.journal" v-bind="field('journal')" label="Journal" />
              </v-col>
              <v-col cols="12" md="3">
                <v-number-input v-model="model.year" v-bind="field('year')" label="Year" />
              </v-col>
            </v-row>
            <v-textarea
              v-model="model.verbatim"
              v-bind="field('verbatim')"
              label="Verbatim"
            ></v-textarea>
            <v-textarea
              v-model="model.comments"
              v-bind="field('comments')"
              label="Comments"
            ></v-textarea>
          </template>
          <v-card
            v-else-if="crossrefResult"
            class="article-info"
            :title="model.title ?? 'Untitled'"
            :subtitle="model.journal ?? 'Unknown journal'"
          >
            <template #append>
              <v-chip :text="model.year?.toString()" label></v-chip>
            </template>
            <v-card-text>
              <div class="mb-3" v-if="model.doi">
                <a
                  :href="`https://doi.org/${model.doi}`"
                  :text="`https://doi.org/${model.doi}`"
                  target="_blank"
                />
              </div>
              <v-chip v-for="author in model.authors" :text="author" class="ma-1" size="small" />
            </v-card-text>
          </v-card>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $ArticleInput, $ArticleUpdate, Article, ArticleLocalInput, ArticleUpdate } from '@/api'
import { createArticleMutation, updateArticleMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import { ref } from 'vue'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import DoiInputFetcher from './DoiInputFetcher.vue'

const manualMode = ref(false)
const dialog = defineModel<boolean>('dialog')
const item = defineModel<Article>()

const crossrefResult = ref<ArticleLocalInput>()

const initial: ArticleLocalInput = {
  authors: []
}

function updateTransformer({ id, meta, $schema, ...rest }: Article): ArticleUpdate {
  return rest
}

const create = defineFormCreate(createArticleMutation(), {
  schema: $ArticleInput,
  initial,
  requestData: (input: ArticleLocalInput) => ({
    body: {
      ...input,
      year: input.year! // will be caught by server-side validation if not provided
    }
  })
})

const update = defineFormUpdate(updateArticleMutation(), {
  schema: $ArticleUpdate,
  itemToModel: updateTransformer,
  requestData: ({ code }) => ({ path: { code } })
})
</script>

<style lang="scss">
.article-info {
  .v-card-title {
    font-size: 1rem;
  }
}
</style>
