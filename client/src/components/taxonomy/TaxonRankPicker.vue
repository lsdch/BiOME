<template>
  <v-select :items v-model="model" v-bind="$attrs"></v-select>
</template>

<script setup lang="ts">
import { $TaxonRank, TaxonRank } from '@/api'

const items = $TaxonRank.enum
const model = defineModel<TaxonRank>()
</script>

<style scoped></style>
