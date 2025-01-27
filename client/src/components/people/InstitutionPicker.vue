<template>
  <v-autocomplete
    label="Institutions (optional)"
    v-model="model"
    :items="filteredItems"
    item-color="primary"
    :chips
    :closable-chips
    :multiple
    :loading
    :error-messages="error?.detail"
    :item-props="({ code, name }: Institution) => ({ title: code, subtitle: name })"
    :item-value
    :return-object
    prepend-inner-icon="mdi-domain"
    v-bind="$attrs"
  >
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #prepend="{ isSelected }">
          <v-checkbox :model-value="isSelected" hide-details />
        </template>
        <template #append>
          <InstitutionKindChip :kind="item.raw.kind" />
        </template>
      </v-list-item>
    </template>
    <template #chip="{ item, props }">
      <InstitutionKindChip :kind="item.raw.kind" v-bind="props" size="small">
        {{ item.raw.code }}
      </InstitutionKindChip>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Institution, InstitutionKind } from '@/api'
import { listInstitutionsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import InstitutionKindChip from './InstitutionKindChip.vue'

const { kinds } = defineProps<{
  multiple?: boolean
  chips?: boolean
  closableChips?: boolean
  itemValue?: keyof Institution
  returnObject?: boolean
  kinds?: InstitutionKind[]
}>()

const model = defineModel<any>()

const { data: items, isPending: loading, error } = useQuery(listInstitutionsOptions())

const filteredItems = computed(() =>
  kinds ? items.value?.filter(({ kind }) => kinds?.includes(kind)) : items.value
)
</script>

<style scoped></style>
