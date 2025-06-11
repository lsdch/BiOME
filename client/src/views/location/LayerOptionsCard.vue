<template>
  <v-card :title :prepend-icon flat v-bind="$attrs">
    <template #title v-if="$slots['title']">
      <slot name="title" />
    </template>
    <template #prepend v-if="$slots['prepend']">
      <slot name="prepend" />
    </template>
    <template #append>
      <div class="d-flex ga-3 align-center">
        <slot name="before-switch" />
        <v-switch v-model="model" color="primary" hide-details />
        <slot name="after-switch" />
        <v-btn
          v-model="expanded"
          :icon="model && expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
          variant="plain"
          :rounded="100"
          color=""
          :disabled="!model"
          @click="toggleExpand()"
        />
      </div>
    </template>
    <slot name="header" />
    <v-expand-transition>
      <div v-if="expanded">
        <slot :active="model"></slot>
      </div>
    </v-expand-transition>
    <slot name="footer" />
  </v-card>
</template>

<script setup lang="ts">
import { useToggle } from '@vueuse/core'
import { watch } from 'vue'
defineProps<{
  title?: string
  prependIcon?: string
}>()
const model = defineModel<boolean>()
const [expanded, toggleExpand] = useToggle(false)
watch(model, (newValue) => {
  expanded.value = newValue ?? false
})
</script>

<style scoped lang="scss"></style>
