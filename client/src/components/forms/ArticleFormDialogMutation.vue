<template>
  <ArticleFormDialog
    v-model="model"
    v-model:dialog="dialog"
    :mode
    :errors
    :title="`${mode} article`"
    :loading="loading || activeMutation.isPending.value"
    :fullscreen="fullscreen || $vuetify.display.mdAndDown"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
  </ArticleFormDialog>
</template>

<script setup lang="ts">
import { $ArticleInput, $ArticleUpdate, Article, ArticleUpdate } from '@/api'
import { createArticleMutation, updateArticleMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate, useMutationForm } from '@/functions/mutations'
import { ArticleModel } from '@/models'
import { useFeedback } from '@/stores/feedback'
import ArticleFormDialog from './ArticleFormDialog.vue'

const dialog = defineModel<boolean>('dialog')
const item = defineModel<Article>('item')

defineProps<FormDialogProps>()

const create = defineFormCreate(createArticleMutation(), {
  initial: ArticleModel.initialModel,
  schema: $ArticleInput,
  requestData: (model) => ({ body: ArticleModel.toRequestBody(model) })
})

const update = defineFormUpdate(updateArticleMutation(), {
  itemToModel: ArticleModel.fromArticle,
  schema: $ArticleUpdate,
  requestData: ({ code }, model) => ({
    path: { code },
    body: ArticleModel.toRequestBody(model)
  })
})

const { feedback } = useFeedback()

const { mode, model, activeMutation, submit, errors } = useMutationForm(item, {
  create,
  update,
  onSuccess(item, mode) {
    dialog.value = false
    feedback({
      type: 'success',
      message: mode === 'Create' ? `Article created` : `Article updated`
    })
  }
})
</script>

<style scoped lang="scss"></style>
