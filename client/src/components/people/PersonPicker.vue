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
    :itemValue
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

<script setup lang="ts" generic="ModelValue extends unknown[] | null | undefined">
import { PeopleService, Person, UserRole } from '@/api'
import { handleErrors } from '@/api/responses'
import { computed, onMounted, ref } from 'vue'
import { isGranted } from './userRole'

const model = defineModel<ModelValue>()
const loading = defineModel<boolean>('loading', { default: true })

const props = defineProps<{
  multiple?: boolean
  label: string
  itemValue?: keyof Person
  // Filter items by user role or account assignation
  users?: boolean | UserRole
}>()

const allPersons = ref<Person[]>(await fetch())

async function fetch() {
  loading.value = true
  return (
    (await PeopleService.listPersons()
      .then(
        handleErrors((err) => {
          console.error('Failed to fetch persons: ', err)
        })
      )
      .finally(() => (loading.value = false))) ?? []
  )
}

const items = computed(() => {
  if (props.users === undefined) {
    return allPersons.value
  }
  return props.users
    ? allPersons.value?.filter(
        ({ user }) => user && (props.users === true || isGranted(user, props.users as UserRole))
      )
    : allPersons.value?.filter(({ user }) => !user)
})
</script>

<style lang="scss" scoped></style>
