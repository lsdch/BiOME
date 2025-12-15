<template>
  <CardDialog :title v-bind="props">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <SiteWithOccurrencesTable :sites />
  </CardDialog>
</template>

<script setup lang="ts">
import { SiteWithOccurrences } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import SiteWithOccurrencesTable from '@/features/occurrences/components/tables/SiteWithOccurrencesTable.vue'

const {
  sites,
  title = 'Sites',
  ...props
} = defineProps<
  {
    sites?: SiteWithOccurrences[]
  } & CardDialogProps
>()
</script>

<style scoped lang="scss"></style>
