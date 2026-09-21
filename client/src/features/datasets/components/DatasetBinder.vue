<template>
  <v-menu :close-on-content-click="false">
    <template #activator="{ props: activatorProps }">
      <v-btn
        prepend-icon="mdi-folder-star-multiple"
        :text="selected.size"
        append-icon="mdi-chevron-down"
        variant="outlined"
        rounded="md"
        color=""
        v-bind="activatorProps"
      ></v-btn>
    </template>
    <v-list density="compact" :min-width="300">
      <v-text-field
        class="mx-2"
        placeholder="Search datasets..."
        clearable
        density="compact"
        hide-details
      ></v-text-field>
      <v-list-item v-for="dataset in sortedItems" :title="dataset.label">
        <template #prepend>
          <v-checkbox
            :disabled="
              !userStore.isGranted('maintainer') ||
              !(
                !!userStore.user &&
                (dataset.owner_id === userStore.user!.id ||
                  dataset.maintainers.find(({ id }) => id === userStore.user!.id))
              )
            "
            :model-value="selected.has(dataset.id)"
            @update:model-value="
              (value) => {
                if (value) {
                  emit('add', dataset.id)
                } else {
                  emit('remove', dataset.id)
                }
              }
            "
            hide-details
            density="compact"
          ></v-checkbox>
        </template>
        <template #append>
          <v-btn
            icon="mdi-open-in-new"
            :to="{ name: 'occurrence-dataset-item', params: { ulid: dataset.id } }"
            target="_blank"
            variant="plain"
            size="small"
            color=""
          />
        </template>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script setup lang="ts">
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useUserStore } from '@/stores/user'
import { useQuery } from '@tanstack/vue-query'
import { storeToRefs } from 'pinia'
import { computed, reactive } from 'vue'

const userStore = useUserStore()

const { model = [] } = defineProps<{
  model?: string[]
}>()

const emit = defineEmits<{
  add: [ulid: string]
  remove: [ulid: string]
}>()

const { data: items } = useQuery({
  ...listDatasetsOptions(),
  initialData: []
})

const selected = computed(() => {
  return new Set(model)
})

const sortedItems = computed(() => {
  return [...items.value].sort(
    (a, b) =>
      a.label.localeCompare(b.label) +
      (selected.value.has(b.id) ? 1 : 0) -
      (selected.value.has(a.id) ? 1 : 0)
  )
})
</script>

<style scoped lang="scss"></style>
