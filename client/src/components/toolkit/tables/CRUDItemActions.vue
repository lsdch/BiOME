<template>
  <!-- v-if="$slots.form" -->
  <v-icon
    v-if="show && show != 'delete'"
    color="primary"
    icon="mdi-pencil"
    @click="emit('edit', item)"
  />
  <v-icon
    v-if="show && show != 'edit'"
    color="primary"
    icon="mdi-delete"
    @click="emit('delete', item)"
  />
</template>

<script setup lang="ts" generic="ItemType extends { id?: string; meta?: Meta }">
import { Meta } from '@/api'

withDefaults(
  defineProps<{
    item: ItemType
    show?: boolean | 'edit' | 'delete'
  }>(),
  { show: true }
)

const emit = defineEmits<{
  edit: [item: ItemType]
  delete: [item: ItemType]
}>()
</script>

<style scoped></style>
