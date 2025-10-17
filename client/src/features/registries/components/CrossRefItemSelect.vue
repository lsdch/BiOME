<template>
  <v-select
    :items
    v-model="model"
    return-object
    hide-details
    variant="plain"
    class="crossref-select"
  >
    <template #selection="{ item }">
      <v-list-item width="100%">
        <template #title>
          {{
            item.raw.author
              ? Article.shortAuthors(item.raw.author.map(({ given, family }) => `${family}`))
              : 'Unknown'
          }}
          <v-chip size="small" :text="item.raw.issued?.['date-parts']?.[0][0].toString()" />
        </template>
        <template #subtitle>
          <v-list-item-subtitle style="max-width: 100%">
            <div class="d-flex flex-column">
              <!-- <span v-for="title in item.raw.title"> {{ title }} </span> -->
              <div v-for="journal in item.raw['container-title']" class="text-caption">
                {{ journal }}
              </div>
            </div>
          </v-list-item-subtitle>
        </template>
      </v-list-item>
    </template>
    <template #item="{ item, props }">
      <v-list-item
        v-bind="props"
        :title="
          item.raw.author
            ? Article.shortAuthors(item.raw.author.map(({ given, family }) => `${family}`))
            : 'Unknown'
        "
      >
        <template #subtitle>
          <v-list-item-subtitle style="max-width: 100%">
            <div class="d-flex flex-column">
              <span v-for="title in item.raw.title"> {{ title }} </span>
              <div v-for="journal in item.raw['container-title']" class="text-caption">
                {{ journal }}
              </div>
            </div>
          </v-list-item-subtitle>
        </template>
        <template #prepend>
          <v-chip
            size="small"
            :text="item.raw.score?.toFixed(0).toString()"
            class="mr-2"
            :color="
              item.raw.score === undefined
                ? ''
                : item.raw.score > 120
                  ? 'success'
                  : item.raw.score > 80
                    ? 'warning'
                    : 'error'
            "
          ></v-chip>
        </template>
        <template #append>
          <v-chip size="small" :text="item.raw.issued?.['date-parts']?.[0][0].toString()" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { Article, BibSearchResults, Message } from '@/api'

const model = defineModel<Message>()

const { items, total } = defineProps<Omit<BibSearchResults, '$schema'>>()
// : { score, author, title, URL, doi, 'container-title': journal, issued }
</script>

<style lang="scss">
@use 'vuetify';
.crossref-select {
  .v-field__input {
    padding-top: var(--v-field-input-padding-bottom);
  }
  .v-field__append-inner {
    padding-top: 0px;
  }
}
</style>
