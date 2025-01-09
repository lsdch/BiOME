<template>
  <v-menu
    location="top start"
    origin="top start"
    transition="scale-transition"
    :open-on-focus="false"
    open-on-click
  >
    <template #activator="{ props }">
      <v-chip
        :text="Article.toString(article)"
        v-bind="{ ...props, ...$attrs }"
        :color="article.original_source ? 'success' : undefined"
      />
    </template>
    <v-card
      :title="article.title ?? 'Untitled'"
      :subtitle="article.journal ?? 'Unknown journal'"
      class="small-card-title bg-surface-light"
      density="compact"
      :max-width="600"
    >
      <template #append>
        <v-chip label :text="article.year.toString()"></v-chip>
      </template>
      <v-card-text>
        {{ article.authors.join(', ') }}
      </v-card-text>
      <v-card-actions v-if="article.doi || article.original_source">
        <v-chip
          label
          v-if="article.original_source"
          title="This is the original reference reporting the occurrence"
          text="Original source"
        />
        <a v-if="article.doi" :href="Article.linkDOI(article)" :text="Article.linkDOI(article)"></a>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Article } from '@/api/adapters'

defineProps<{ article: Article }>()
</script>

<style scoped lang="scss"></style>
