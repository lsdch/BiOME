<template>
  <v-tooltip bottom :open-delay="400">
    Sort by last update
    <template v-slot:activator="{ props }">
      <v-btn v-bind="{ ...sortLastUpdateIcon, ...props }" variant="plain" @click="emit('click')" />
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { VBtn } from 'vuetify/components'

const emit = defineEmits<{
  click: []
}>()

const props = defineProps<{
  sortKey: string
  sortBy?: SortItem[]
}>()

const sortLastUpdateIcon = computed(() => {
  const sortMeta = props.sortBy?.find(({ key }) => key === props.sortKey)
  if (sortMeta) {
    return sortMeta.order === 'asc'
      ? { icon: 'mdi-sort-calendar-ascending', color: 'primary' }
      : { icon: 'mdi-sort-calendar-descending', color: 'primary' }
  }
  return { icon: 'mdi-sort-calendar-ascending', color: 'secondary' }
})
</script>

<style scoped></style>
