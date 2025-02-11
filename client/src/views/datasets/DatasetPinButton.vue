<template>
  <v-hover>
    <template #default="{ isHovering, props }">
      <v-btn
        :icon="dataset.pinned && isHovering ? 'mdi-pin-off' : 'mdi-pin'"
        :active="dataset.pinned"
        :color="dataset.pinned && isHovering ? 'error' : dataset.pinned ? 'success' : ''"
        :loading
        @click="mutate({ path: { slug: dataset.slug } })"
        :rounded="true"
        variant="text"
        size="x-small"
        :title="dataset.pinned ? 'Unpin dataset' : 'Pin dataset'"
        v-bind="{ ...$attrs, ...props }"
      >
      </v-btn>
    </template>
  </v-hover>
</template>

<script setup lang="ts">
import { togglePinDatasetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation } from '@tanstack/vue-query'

const emit = defineEmits<{
  'update:model-value': [dataset: { id: string; slug: string; pinned: boolean }]
}>()

const dataset = defineModel<{ id: string; slug: string; pinned: boolean }>({
  required: true
})

const { mutate, isPending: loading } = useMutation({
  ...togglePinDatasetMutation(),
  onError(error) {
    console.error(error)
  },
  onSuccess(updated) {
    emit('update:model-value', updated)
  }
})
</script>

<style scoped lang="scss"></style>
