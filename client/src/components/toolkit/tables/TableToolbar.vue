<template>
  <v-toolbar flat dense extension-height="auto" class="px-3">
    <!-- Top left icon -->
    <template v-if="icon" #prepend>
      <v-tooltip :disabled="!onReload">
        <template #activator="{ props, isActive }">
          <v-avatar color="secondary" variant="outlined" v-bind="props">
            <v-icon
              :class="{ 'cursor-default': !onReload }"
              dark
              color="secondary-darken-1"
              @click="emit('reload')"
            >
              {{ isActive && onReload ? 'mdi-reload' : icon }}
            </v-icon>
          </v-avatar>
        </template>
        Reload items
      </v-tooltip>
    </template>

    <!-- Expose toolbar append slot -->
    <template #append>
      <slot name="append"></slot>
    </template>

    <v-toolbar-title v-if="title !== undefined" style="min-width: 150px" :text="title" />

    <slot name="search" />

    <v-spacer />

    <slot name="prepend-actions" />
    <slot name="actions" />
    <slot name="append-actions" />

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

defineProps<ToolbarProps>()

const emit = defineEmits<{ reload: [] }>()
</script>

<style scoped></style>
