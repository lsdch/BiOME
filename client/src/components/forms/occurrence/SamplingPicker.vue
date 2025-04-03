<template>
  <v-select v-model="model" :items item-title="code" v-bind="$attrs"></v-select>
</template>

<script setup lang="ts">
import { SamplingInner } from '@/api'

const model = defineModel<SamplingInner>()

defineProps<{ items: SamplingInner[] }>()
</script>

<style scoped lang="scss"></style>
