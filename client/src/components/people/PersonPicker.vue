<template>
  <v-autocomplete
    v-model="model"
    :multiple="multiple"
    :label="label"
    :items="_items"
    item-title="full_name"
    clearable
    chips
    closable-chips
    auto-select-first
    clear-on-select
    counter
  >
  </v-autocomplete>
</template>

<script setup lang="ts">
import { PeopleService, Person, PersonUser } from '@/api'
import { handleErrors } from '@/api/responses'

const props = defineProps<{
  multiple?: boolean
  label: string
  items?: Person[]
}>()

const _items =
  props.items ??
  (await PeopleService.listPersons().then(
    handleErrors((err) => {
      console.error('Failed to fetch persons: ', err)
    })
  ))

const model = defineModel<PersonUser[] | PersonUser>()
</script>

<style lang="scss" scoped></style>
