<template>
  <FormDialog
    v-bind="props"
    v-model="dialog"
    :title="title ?? `${mode} article`"
    @submit="emit('submit', model)"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container fluid>
      <v-row>
        <v-col cols="12" md="9">
          <DoiInputFetcher v-model="model.doi" @fetched="(article) => (model = article)" />
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
          v-bind="schema('authors')"
        />

        <v-text-field v-model="model.title" v-bind="schema('title')" label="Title"></v-text-field>
        <v-row>
          <v-col cols="12" md="9">
            <v-text-field v-model="model.journal" v-bind="schema('journal')" label="Journal" />
          </v-col>
          <v-col cols="12" md="3">
            <v-number-input v-model="model.year" v-bind="schema('year')" label="Year" />
          </v-col>
        </v-row>
        <v-textarea
          v-model="model.verbatim"
          v-bind="schema('verbatim')"
          label="Verbatim"
        ></v-textarea>
        <v-textarea
          v-model="model.comments"
          v-bind="schema('comments')"
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

<script setup lang="ts">
import { $ArticleInput, $ArticleUpdate } from '@/api'
import DoiInputFetcher from '@/components/references/DoiInputFetcher.vue'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { useSchema } from '@/composables/schema'
import { FormProps } from '@/functions/mutations'
import { ArticleModel } from '@/models'
import { reactiveComputed } from '@vueuse/core'
import { ref } from 'vue'

const dialog = defineModel<boolean>('dialog')
const model = defineModel<ArticleModel.ArticleFormModel>({
  default: ArticleModel.initialModel
})

const manualMode = ref(false)
const crossrefResult = ref<ArticleModel.ArticleFormModel>()

const { mode = 'Create', ...props } = defineProps<FormProps & FormDialogProps>()

const emit = defineEmits<{
  submit: [model: ArticleModel.ArticleFormModel | undefined]
}>()

const {
  bind: { schema }
} = reactiveComputed(() => useSchema(mode === 'Create' ? $ArticleInput : $ArticleUpdate))
</script>

<style scoped lang="scss"></style>
