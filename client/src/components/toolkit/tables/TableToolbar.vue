<template>
  <v-toolbar flat dense extension-height="auto">
    <!-- Top left icon -->
    <template v-if="icon" #prepend>
      <v-tooltip>
        <template #activator="{ props, isActive }">
          <v-avatar color="secondary" variant="outlined" v-bind="props">
            <v-icon dark color="secondary-darken-1" @click="emit('reload')">
              {{ isActive ? 'mdi-reload' : icon }}
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

    <v-toolbar-title style="min-width: 150px">{{ title }}</v-toolbar-title>

    <slot v-if="smAndUp && !togglableSearch" name="search" class="flex-grow-1" />

    <v-spacer />

    <!-- Toggle large searchbar component -->
    <v-tooltip v-if="xs || togglableSearch" left activator="parent" text="Toggle search">
      <template #activator="{ props }">
        <v-btn
          v-bind="props"
          size="small"
          icon="mdi-magnify"
          color="primary"
          :variant="toggleSearch ? 'flat' : 'text'"
          @click="toggleSearch = !toggleSearch"
        />
      </template>
    </v-tooltip>

    <slot name="prepend-actions" />
    <slot name="actions" />
    <slot name="append-actions" />

    <!-- Search bar slot with default searchbar -->
    <template v-if="togglableSearch || xs" #extension>
      <v-expand-transition>
        <div class="w-100 px-3" v-show="toggleSearch" transition="slide-y-transition">
          <slot name="search" class="flex-grow-1"> </slot>
        </div>
      </v-expand-transition>
    </template>
  </v-toolbar>
</template>

<script setup lang="ts" generic="ItemType extends { id: string }">
import { ref } from 'vue'
import { useDisplay } from 'vuetify'
import { ToolbarProps } from '.'

const { xs, smAndUp } = useDisplay()

const toggleSearch = ref(false)

defineProps<ToolbarProps>()

const emit = defineEmits<{ reload: [] }>()
</script>

<style scoped></style>
&
