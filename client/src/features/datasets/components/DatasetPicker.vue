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
    no-data-text="No datasets found"
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
    <template #append v-if="withForm && userStore.isGranted('contributor')">
      <dataset-create-form-dialog
        @created="
          (value) => {
            refetch()
            multiple
              ? (model = [...(model || []), itemValue ? value[itemValue] : value] as ModelValue)
              : (model = (itemValue ? value[itemValue] : value) as ModelValue)
          }
        "
      >
        <template #activator="{ props: activatorProps }">
          <v-btn v-bind="activatorProps" icon="mdi-plus" v-tooltip="'Create dataset'" />
        </template>
      </dataset-create-form-dialog>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="Multiple extends boolean, ReturnObject extends boolean">
import { Dataset, DatasetWithMaintainers } from '@/api'
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useUserStore } from '@/stores/user'
import { useQuery } from '@tanstack/vue-query'
import DatasetCreateFormDialog from './DatasetCreateFormDialog.vue'
import { Value } from 'vuetify/lib/components/VAutocomplete/VAutocomplete.mjs'

const userStore = useUserStore()

type ModelValue = Value<DatasetWithMaintainers, ReturnObject, Multiple>

const model = defineModel<ModelValue>()

const { itemValue = 'id' } = defineProps<{
  multiple?: boolean
  label: string
  itemValue?: keyof Dataset
  withForm?: boolean
}>()

const { data: items, isPending: loading, error, refetch } = useQuery(listDatasetsOptions())
</script>

<style lang="scss" scoped></style>
