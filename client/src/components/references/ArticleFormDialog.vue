<template>
  <FormDialog v-model="dialog" :title="`${mode} article`" :loading @submit="submit">
    <v-container fluid>
      <v-row>
        <v-col cols="12" md="9">
          <v-text-field
            v-model.trim="model.doi"
            label="DOI"
            @keyup.enter="model.doi ? fetchFromDOI(model.doi) : undefined"
            @input="hasFetched = false"
            :loading
            :error-messages="crossRefError?.detail"
          >
            <template #append-inner>
              <v-btn
                color="primary"
                icon="mdi-send"
                variant="tonal"
                size="small"
                @click="model.doi ? fetchFromDOI(model.doi) : undefined"
                :loading
              />
            </template>
          </v-text-field>
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

        <v-text-field v-model="model.title" v-bind="field('title')" label="Title"></v-text-field>
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
        v-else-if="hasFetched"
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

<script setup lang="ts">
import {
  Article,
  ArticleInput,
  ArticleUpdate,
  $ArticleInput,
  $ArticleUpdate,
  ReferencesService,
  ErrorModel
} from '@/api'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'
import { ref } from 'vue'

const manualMode = ref(false)
const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Article>>()
const emit = defineEmits<FormEmits<Article>>()

const initial: Omit<ArticleInput, 'year'> & { year?: number } = {
  authors: []
}

const { model, mode, makeRequest, reset } = useForm(props, {
  initial,
  updateTransformer({ id, meta, $schema, ...rest }): ArticleUpdate {
    return rest
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $ArticleInput : $ArticleUpdate)
)

const [loading, toggleLoading] = useToggle(false)

const crossRefError = ref<ErrorModel>()
const hasFetched = ref(false)
async function fetchFromDOI(doi: string) {
  toggleLoading(true)
  const { data, error } = await ReferencesService.crossref({ query: { doi } })
  if (error) {
    crossRefError.value = error
    reset()
  } else {
    crossRefError.value = undefined
    const { message } = data
    hasFetched.value = true

    model.value = {
      doi,
      authors: message?.author?.map(({ family, given }) => `${family} ${given}`),
      title: message?.title?.[0],
      journal: message?.publisher,
      year: message?.published?.['date-parts']?.[0]?.[0]
    }
  }
  toggleLoading(false)
}

async function submit() {
  console.log('SUBMIt')
  toggleLoading(true)
  return await makeRequest({
    create: ({ body }) => ReferencesService.createArticle({ body: { ...body, year: body.year! } }),
    edit: ({ code }, body) => ReferencesService.updateArticle({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style lang="scss">
.article-info {
  .v-card-title {
    font-size: 1rem;
  }
}
</style>
