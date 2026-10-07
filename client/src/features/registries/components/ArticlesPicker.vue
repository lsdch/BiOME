<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerms"
    :items="items"
    :loading
    item-value="id"
    :multiple
    :return-object
    :chips="multiple"
    :closable-chips="multiple"
    no-filter
    clear-on-select
    placeholder="Enter search terms..."
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template #chip="{ item, props }">
      <PublicationChip :publication="item" v-bind="props" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        v-bind="props"
        :title="Publication.shortAuthors(item.authors ?? [])"
        class="fuzzy-search-item"
      >
        <template #prepend="{ isSelected }">
          <v-checkbox :model-value="isSelected"></v-checkbox>
        </template>
        <template #subtitle>
          <v-list-item-subtitle style="max-width: 100%">
            <div class="d-flex flex-column" v-if="item.title">
              <span> {{ item.title }} </span>
              <div class="text-caption">{{ item.journal }}</div>
            </div>
            <div v-else class="text-caption">{{ item.verbatim }}</div>
          </v-list-item-subtitle>
        </template>
        <template #append v-if="item.year">
          <v-chip size="small">
            {{ item.year }}
          </v-chip>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean, ReturnObject extends boolean">
import { Publication } from '@/api'
import { listPublicationsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import type { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'
import PublicationChip from './PublicationChip'

export type ArticlesPickerProps<
  Multiple extends boolean = false,
  ReturnObject extends boolean = false
> = {
  multiple?: Multiple
  returnObject?: ReturnObject
}

const { multiple = false, returnObject = false } =
  defineProps<ArticlesPickerProps<Multiple, ReturnObject>>()

const { data: items, isPending: loading, error } = useQuery(listPublicationsOptions())

const model = defineModel<Value<Publication, ReturnObject, Multiple> | null>()
const searchTerms = ref<string>('')
</script>

<style lang="scss"></style>
