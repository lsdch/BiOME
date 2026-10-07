<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerm"
    :items="habitats"
    label="Habitat tags"
    prepend-inner-icon="mdi-tag-multiple"
    color="primary"
    :multiple
    chips
    clearable
    closable-chips
    item-title="name"
    item-subtitle="description"
    auto-select-first
    clear-on-select
    placeholder="Start typing to search"
    :return-object
    :item-value
  >
    <template #item="{ item, props }">
      <v-list-item :title="item.name" v-bind="props">
        <template #title="{ title }">
          <span>{{ title }}</span>
        </template>
        <template #subtitle="{ subtitle }">
          <span class="text-caption">{{ subtitle }}</span>
        </template>
        <template #append>
          <v-chip class="text-overline" color="primary" :text="item.facet" />
        </template>
      </v-list-item>
    </template>
    <!-- <template #append-item>
      <v-list v-if="quickSelect.length > 0">
        <v-list-subheader> Quick select </v-list-subheader>
        <v-list-item v-for="item in quickSelect" :key="item.id" @click="addWithDependencies(item)">
          <template v-for="dep in item.upstream ?? []" :key="dep.id">
            <v-chip class="" :text="dep.name" variant="text" />
            <v-icon icon="mdi-chevron-right" />
          </template>
          <v-chip class="ma-1" :text="item.name" color="primary" />
          <template #append>
            <v-chip class="text-overline" color="primary" :text="getHabitat(item.id)?.group.name" />
          </template>
        </v-list-item>
      </v-list>
    </template> -->
  </v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean, ReturnObject extends boolean">
import { Habitat } from '@/api'
import { listHabitatsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import type { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'

export type HabitatPickerProps<
  Multiple extends boolean = false,
  ReturnObject extends boolean = false
> = {
  multiple?: Multiple
  returnObject?: ReturnObject
  itemValue?: 'id' | 'name'
}

const {
  multiple = false,
  returnObject = false,
  itemValue = 'id'
} = defineProps<HabitatPickerProps<Multiple, ReturnObject>>()

const model = defineModel<Value<Habitat, ReturnObject, Multiple> | null>()
const searchTerm = ref<string | undefined>(undefined)

const { data: habitats } = useQuery(listHabitatsOptions())
</script>

<style scoped></style>
