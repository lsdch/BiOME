<template>
  <CardDialog :title v-bind="props">
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>

    <SamplingWithOccurrencesTable :with-site :samplings />
  </CardDialog>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import { SamplingDateWithOccurrences, SiteItem } from '@/api'
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import SamplingWithOccurrencesTable, {
  SamplingTableItem
} from '@/features/occurrences/components/tables/SamplingWithOccurrencesTable.vue'
import { ComponentSlots } from 'vue-component-type-helpers'

const {
  title = 'Sampling events',
  samplings,
  ...props
} = defineProps<
  {
    /**
     * If true, the site information is included with each occurrence.
     * If false, only the occurrence information is included.
     */
    withSite: WithSite
    samplings?: SamplingTableItem<SamplingDateWithOccurrences, WithSite, SiteItem>[]
  } & CardDialogProps
>()

defineSlots<ComponentSlots<typeof CardDialog>>()
</script>

<style scoped lang="scss"></style>
