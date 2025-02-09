<template>
  <v-autocomplete
    label="Organisations (optional)"
    v-model="model"
    :items="filteredItems"
    item-color="primary"
    :chips
    :closable-chips
    :multiple
    :loading
    :error-messages="error?.detail"
    :item-props="({ code, name }: Organisation) => ({ title: code, subtitle: name })"
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
          <OrgKindChip :kind="item.raw.kind" />
        </template>
      </v-list-item>
    </template>
    <template #chip="{ item, props }">
      <OrgKindChip :kind="item.raw.kind" v-bind="props" size="small">
        {{ item.raw.code }}
      </OrgKindChip>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { Organisation, OrgKind } from '@/api'
import { listOrganisationsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import OrgKindChip from './OrgKindChip.vue'

const { kinds } = defineProps<{
  multiple?: boolean
  chips?: boolean
  closableChips?: boolean
  itemValue?: keyof Organisation
  returnObject?: boolean
  kinds?: OrgKind[]
}>()

const model = defineModel<any>()

const { data: items, isPending: loading, error } = useQuery(listOrganisationsOptions())

const filteredItems = computed(() =>
  kinds ? items.value?.filter(({ kind }) => kinds?.includes(kind)) : items.value
)
</script>

<style scoped></style>
