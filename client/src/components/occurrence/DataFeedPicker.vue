<template>
  <v-select
    v-model="model"
    :items="registry"
    item-title="name"
    prepend-inner-icon="mdi-database-arrow-right-outline"
    item-value="id"
    v-bind="$attrs"
  >
    <template #append-inner="{}">
      <v-progress-circular
        v-if="model && remotes.get(model)?.isPending.value"
        indeterminate
        size="small"
        color="warning"
      />
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { type UUID } from 'crypto'
import { useDataFeeds } from './data_feeds'
import { watch } from 'vue'

const { registry, remotes } = useDataFeeds()

const model = defineModel<UUID>()

const { mandatory } = defineProps<{
  mandatory?: boolean
}>()

watch(
  remotes,
  (newRemotes) => {
    if (model.value && !newRemotes.has(model.value)) {
      model.value = mandatory ? registry.value[0].id : undefined
    }
  },
  { immediate: true }
)
</script>

<style scoped lang="scss"></style>
