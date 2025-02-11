<template>
  <v-menu
    location="top start"
    origin="top start"
    transition="scale-transition"
    :open-on-focus="false"
    open-on-click
  >
    <template #activator="{ props }">
      <v-chip :text="source.code" v-bind="{ ...props, ...$attrs }" />
    </template>
    <v-card
      :title="source.label ?? 'Untitled'"
      :subtitle="source.code ?? 'Unknown journal'"
      class="small-card-title bg-surface-light"
      density="compact"
      :max-width="600"
    >
      <v-card-text>
        {{ source.description }}
      </v-card-text>
      <v-card-actions v-if="source.url">
        <a :href="source.url" :text="source.url"></a>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { Article, DataSource } from '@/api/adapters'

defineProps<{ source: DataSource }>()
</script>

<style scoped lang="scss"></style>
