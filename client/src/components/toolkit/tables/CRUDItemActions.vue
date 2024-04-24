<template>
  <template v-if="canEdit">
    <!-- v-if="$slots.form" -->
    <v-icon color="primary" icon="mdi-pencil" @click="actions.edit(props.item)" />
    <v-icon color="primary" icon="mdi-delete" @click="actions.delete(props.item)" />
  </template>
</template>

<script setup lang="ts" generic="ItemType extends { id: string; meta?: Meta }">
import { Meta, User } from '@/api'
import { isGranted } from '@/components/people/userRole'
import { computed } from 'vue'
import { isOwner } from '../meta'

const props = defineProps<{
  item: ItemType
  user: User
  actions: {
    edit(item: ItemType): PromiseLike<void>
    delete(item: ItemType): PromiseLike<void>
  }
}>()

const canEdit = computed(() => {
  return isGranted(props.user, 'Maintainer') || isOwner(props.user, props.item)
})
</script>

<style scoped></style>
