<template>
  <v-toolbar flat dense extension-height="auto" class="px-3">
    <!-- Top left icon -->
    <template v-if="icon" #prepend>
      <ClickableAvatarIcon :icon hover-icon="mdi-reload" @click="onReload?.()" />
    </template>

    <v-toolbar-title v-if="title !== undefined" style="min-width: 150px" :text="title" />

    <slot name="search" />

    <slot name="prepend-actions" />
    <slot name="actions" />
    <slot name="append-actions" />
    <!-- Expose toolbar append slot -->
    <template #append>
      <slot name="append" />
    </template>

    <!-- Search bar slot with default searchbar -->
    <template #extension>
      <v-expand-transition>
        <div class="w-100 px-3" transition="slide-y-transition">
          <slot name="extension"> </slot>
        </div>
      </v-expand-transition>
    </template>
  </v-toolbar>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { ToolbarProps } from '.'
import ClickableAvatarIcon from '../ui/ClickableAvatarIcon.vue'

defineProps<ToolbarProps>()

const emit = defineEmits<{ reload: [] }>()
</script>

<style scoped></style>
