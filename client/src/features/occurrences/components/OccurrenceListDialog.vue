<template>
  <CardDialog :title v-bind="props">
    <template v-for="(_, name) in $slots" #[name]="slotData">
      <slot :name="name" v-bind="slotData ?? {}" />
    </template>

    <OccurrencesTable :with-site :occurrences />
  </CardDialog>
</template>

<script setup lang="ts" generic="WithSite extends boolean">
import CardDialog, { CardDialogProps } from '@/components/toolkit/ui/CardDialog.vue'
import OccurrencesTable, {
  OccurrenceTableItem
} from '@/features/occurrences/components/tables/OccurrencesTable.vue'
import { ComponentSlots } from 'vue-component-type-helpers'

const {
  title = 'Occurrences',
  occurrences,
  ...props
} = defineProps<
  {
    /**
     * If true, the site information is included with each occurrence.
     * If false, only the occurrence information is included.
     */
    withSite: WithSite
    occurrences?: OccurrenceTableItem<WithSite>[]
  } & CardDialogProps
>()

defineSlots<ComponentSlots<typeof CardDialog>>()
</script>

<style scoped lang="scss"></style>
