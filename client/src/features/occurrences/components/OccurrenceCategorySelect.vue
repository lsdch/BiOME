<template>
  <v-select
    :items="$OccurrenceCategory.enum"
    v-model="model"
    label="Type"
    hide-details
    placeholder="Any"
    persistent-placeholder
    density="compact"
    clearable
    persistent-clear
    :color="model ? 'primary' : undefined"
    :active="!!model"
    :prepend-inner-icon="model ? OccurrenceCategory.props[model].icon : undefined"
  >
    <template #item="{ item, props }">
      <v-list-item
        v-bind="{
          ...props,
          ...OccurrenceCategory.props[item]
        }"
        :class="`text-${OccurrenceCategory.props[item].color}`"
      ></v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $OccurrenceCategory } from '@/api'
import { OccurrenceCategory } from '@/api/adapters'

const model = defineModel<OccurrenceCategory>()
</script>

<style scoped lang="scss"></style>
