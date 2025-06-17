<template>
  <v-autocomplete
    v-model="model"
    :items
    :multiple="multiple"
    :chips="multiple"
    :closable-chips="multiple"
    :label="label"
    item-title="label"
    auto-select-first
    clear-on-select
    :loading="loading"
    :item-value
    v-bind="$attrs"
    :error-messages="error?.detail"
  >
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData" />
    </template>
    <template #prepend-item>
      <div class="d-flex w-100 justify-center">
        <v-btn-toggle
          rounded
          multiple
          v-model="categories"
          class="flex-grow-1"
          color="primary"
          variant="text"
          density="compact"
          divided
          border="sm"
        >
          <v-btn
            :icon="DatasetCategory.icon('Site')"
            size="small"
            value="Site"
            class="flex-grow-1"
            v-tooltip="{ text: 'Sites', location: 'top' }"
          />
          <v-btn
            :icon="DatasetCategory.icon('Occurrence')"
            size="small"
            value="Occurrence"
            class="flex-grow-1"
            v-tooltip="{ text: 'Occurrences', location: 'top' }"
          />
          <v-btn
            :icon="DatasetCategory.icon('Seq')"
            size="small"
            value="Seq"
            class="flex-grow-1"
            v-tooltip="{ text: 'Sequences', location: 'top' }"
          />
        </v-btn-toggle>
      </div>
      <v-divider class="mt-2" />
    </template>
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #prepend="{ isSelected }" v-if="multiple">
          <v-checkbox :modelValue="isSelected" hide-details density="compact" class="mx-1" />
        </template>
        <template #append>
          <DatasetCategoryIcon :category="item.raw.category" size="x-small" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="ModelValue extends unknown | unknown[] | null | undefined">
import { Dataset, DatasetCategory } from '@/api'
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import DatasetCategoryIcon from '@/components/datasets/DatasetCategoryIcon'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const model = defineModel<ModelValue>()
const categories = ref<DatasetCategory[]>(['Occurrence', 'Site', 'Seq'])

defineProps<{
  multiple?: boolean
  label: string
  itemValue?: keyof Dataset
}>()

const { data: datasets, isPending: loading, error } = useQuery(listDatasetsOptions({ query: {} }))

const items = computed(
  () => datasets.value?.filter(({ category }) => categories.value.includes(category)) ?? []
)
</script>

<style lang="scss" scoped></style>
