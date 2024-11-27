<template>
  <v-select
    v-model="model"
    :items
    :loading
    item-title="label"
    item-subtitle="description"
    prepend-inner-icon="mdi-snowflake"
    v-bind="$attrs"
  ></v-select>
</template>

<script setup lang="ts">
import { Fixative, SamplingService } from '@/api'
import { useFetchItems } from '@/composables/fetch_items'
import { ref } from 'vue'

const model = ref<Fixative[]>([])
const { items, loading } = useFetchItems(SamplingService.listFixatives)
</script>

<style scoped lang="scss"></style>
