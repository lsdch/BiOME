<template>
  <v-autocomplete
    label="Institutions (optional)"
    v-model="model"
    :items
    item-color="primary"
    :chips
    :closable-chips
    :multiple
    :loading
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
import { Institution, InstitutionKind, PeopleService } from '@/api'
import { handleErrors } from '@/api/responses'
import { useToggle } from '@vueuse/core'
import { onMounted, ref } from 'vue'
import InstitutionKindChip from './InstitutionKindChip.vue'

const props = defineProps<{
  multiple?: boolean
  chips?: boolean
  closableChips?: boolean
  itemValue?: keyof Institution
  returnObject?: boolean
  kinds?: InstitutionKind[]
}>()

const model = defineModel<any>()

const items = ref<Institution[]>([])

onMounted(fetch)

const [loading, toggleLoading] = useToggle(true)
async function fetch() {
  toggleLoading(true)
  const data = await PeopleService.listInstitutions().then(
    handleErrors((err) => console.error('Failed to fetch institutions: ', err))
  )
  items.value = props.kinds ? data.filter(({ kind }) => props.kinds?.includes(kind)) : data
  toggleLoading(false)
}
</script>

<style scoped></style>
