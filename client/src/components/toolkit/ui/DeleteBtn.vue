<template>
  <ConfirmDialog :title :message @confirm="$emit('confirm')" @cancel="$emit('cancel')">
    <template #activator="{ props }">
      <v-btn
        text="Delete"
        color="error"
        variant="tonal"
        size="small"
        class="ma-1"
        prepend-icon="mdi-delete"
        v-bind="props"
      />
    </template>
    <slot name="message" />
  </ConfirmDialog>
</template>

<script setup lang="ts">
import { useAppConfirmDialog } from '@/composables/confirm_dialog'
import ConfirmDialog from './ConfirmDialog.vue'

const { title = 'Confirm deletion', message = 'Are you sure you want to delete this item ?' } =
  defineProps<{
    title: string
    message?: string
  }>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()
</script>

<style scoped lang="scss"></style>
