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
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #prepend="{ isSelected }" v-if="multiple">
          <v-checkbox :modelValue="isSelected" hide-details density="compact" class="mx-1" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="ModelValue extends unknown | unknown[] | null | undefined">
import { Dataset } from '@/api'
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'

const model = defineModel<ModelValue>()

defineProps<{
  multiple?: boolean
  label: string
  itemValue?: keyof Dataset
}>()

const { data: items, isPending: loading, error } = useQuery(listDatasetsOptions())
</script>

<style lang="scss" scoped></style>
