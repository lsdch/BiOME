<template>
  <v-autocomplete
    v-model="model"
    :multiple="multiple"
    :chips="multiple"
    :closable-chips="multiple"
    :label="label"
    :items
    item-title="full_name"
    auto-select-first
    clear-on-select
    :loading="loading"
    :item-value
    v-bind="$attrs"
  >
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData" />
    </template>
    <template #item="{ item, props }">
      <v-list-item v-bind="props">
        <template #prepend="{ isSelected }" v-if="multiple">
          <v-checkbox :modelValue="isSelected" hide-details density="compact" class="mx-1" />
        </template>
        <template v-if="item.raw.user" #append>
          <v-chip :text="item.raw.user.role" />
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts" generic="ModelValue extends unknown | unknown[] | null | undefined">
import { PeopleService, Person, UserRole } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'
import { computed } from 'vue'

const model = defineModel<ModelValue>()

const { restrict } = defineProps<{
  multiple?: boolean
  label: string
  itemValue?: keyof Person
  // Filter items by user role or account assignation
  restrict?: 'users' | 'unregistered' | UserRole
}>()

const { items: allPersons, loading } = useFetchItems(PeopleService.listPersons)

const items = computed(() => {
  switch (restrict) {
    case undefined:
      return allPersons.value
    case 'users':
      return allPersons.value.filter(({ user }) => user)
    case 'unregistered':
      return allPersons.value.filter(({ user }) => !user)
    default:
      return allPersons.value.filter(({ user }) => user && UserRole.isGranted(user, restrict))
  }
})
</script>

<style lang="scss" scoped></style>
