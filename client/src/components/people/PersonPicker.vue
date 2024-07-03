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
  >
  </v-autocomplete>
</template>

<script setup lang="ts">
import { PeopleService, Person, PersonUser } from '@/api'
import { handleErrors } from '@/api/responses'
import { onMounted, ref } from 'vue'

const props = defineProps<{
  multiple?: boolean
  label: string
  items?: Person[]
}>()

const items = ref(props.items)
const loading = ref(true)

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

const model = defineModel<PersonUser[] | PersonUser>()
</script>

<style lang="scss" scoped></style>
