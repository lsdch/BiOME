<template>
  <v-icon v-bind="spec" :size="size"></v-icon>
</template>

<script setup lang="ts">
import { TaxonStatus } from '@/api'
import { computed } from 'vue'

const props = defineProps<{ status: TaxonStatus; size?: string }>()

const spec = computed(() => {
  switch (props.status) {
    case TaxonStatus.Accepted:
      return {
        icon: 'mdi-check',
        color: 'green'
      }
    case TaxonStatus.Synonym:
      return {
        icon: 'mdi-arrow-left-right',
        color: 'orange'
      }
    case TaxonStatus.Unclassified:
      return {
        icon: 'mdi-plus',
        color: 'blue'
      }
    default:
      throw new Error(`Unhandled TaxonStatus value : ${props.status}`)
  }
})
</script>

<style scoped></style>
