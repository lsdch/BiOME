<template>
  <v-select :items item-title="code" v-bind="$attrs"></v-select>
</template>

<script setup lang="ts">
import { SamplingInner } from '@/api'

defineProps<{ items: SamplingInner[] }>()
</script>

<style scoped lang="scss"></style>
