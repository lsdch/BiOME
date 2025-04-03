<template>
  <v-text-field
    v-model.trim="model"
    label="DOI"
    :loading="isFetching"
    :error-messages="error?.detail"
  >
    <template #append-inner>
      <v-btn
        color="primary"
        icon="mdi-send"
        variant="tonal"
        size="small"
        @click="refetch()"
        :loading="isFetching"
      />
    </template>
  </v-text-field>
</template>

<script setup lang="ts">
import { ErrorModel } from '@/api'
import { crossRefOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { ArticleFormModel } from '@/models/article'
import { useQuery } from '@tanstack/vue-query'
import { watch } from 'vue'

const model = defineModel<string | null>({ default: '' })

const { data, error, isFetching, refetch } = useQuery({
  ...crossRefOptions({ query: { doi: model.value ?? '' } }),
  enabled: false
})

const emit = defineEmits<{
  fetched: [article: ArticleFormModel]
  error: [error: ErrorModel]
}>()

watch(data, (d) => {
  if (error.value) {
    emit('error', error.value)
  } else if (d) {
    const { message } = d
    const article = {
      doi: model.value!,
      authors: message?.author?.map(({ family, given }) => `${family} ${given}`) ?? [],
      title: message?.title?.[0],
      journal: message?.publisher,
      year: message?.published?.['date-parts']?.[0]?.[0]
    }
    emit('fetched', article)
  }
})
</script>

<style scoped lang="scss"></style>
