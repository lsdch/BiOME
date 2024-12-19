<template>
  <v-select v-model="model" :items :loading v-bind="$attrs" item-title="code">
    <template #item="{ item, props }">
      <v-list-item :subtitle="item.raw.label" v-bind="props" />
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { SequencesService } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'

const model = defineModel<any>()

const { items, loading } = useFetchItems(SequencesService.listGenes)
</script>

<style scoped lang="scss"></style>
