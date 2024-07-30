<template>
  <v-autocomplete
    v-model="model"
    :multiple="multiple"
    :label="label"
    :items="items"
    item-title="full_name"
    clearable
    chips
    closable-chips
    auto-select-first
    clear-on-select
    counter
    :loading="loading"
    :itemValue
    v-bind="$attrs"
  >
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData" />
    </template>
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #prepend="{ isSelected }">
          <v-checkbox :modelValue="isSelected" hide-details density="compact" class="mx-1" />
        </template>
        <template #append>
          <v-chip :text="item.raw.user.role" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="ModelValue extends unknown[]">
import { PeopleService, Person } from '@/api'
import { handleErrors } from '@/api/responses'
import { onMounted, ref } from 'vue'

const model = defineModel<ModelValue>()
const loading = defineModel<boolean>('loading', { default: true })

const props = defineProps<{
  multiple?: boolean
  label: string
  items?: Person[]
  itemValue: keyof Person
}>()

const items = ref(props.items)

onMounted(async () => {
  if (items.value === undefined) {
    items.value = await PeopleService.listPersons().then(
      handleErrors((err) => {
        console.error('Failed to fetch persons: ', err)
      })
    )
  }
  loading.value = false
})
</script>

<style lang="scss" scoped></style>
