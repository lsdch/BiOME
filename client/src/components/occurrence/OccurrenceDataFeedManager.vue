<template>
  <div class="bg-main">
    <v-list-item class="pb-3 bg-main">
      <DatasetPicker
        label="Dataset context"
        placeholder="None"
        persistent-placeholder
        prepend-icon="mdi-folder-table"
        density="compact"
        clearable
        multiple
        chips
        closable-chips
        hint="Restrict all data feeds to a one or more datasets. Individual feeds may be configured to include additional datasets."
        persistent-hint
        class="my-3"
      />
    </v-list-item>
    <v-divider class="mb-3" />
    <v-item-group selected-class="ma-3" class="d-flex flex-column ga-3 pa-2 bg-main">
      <v-item v-for="(_, i) in registry" #="{ isSelected, toggle, selectedClass }">
        <OccurrenceDataFeedCard
          v-model="registry[i]"
          :placeholder="`Feed #${i + 1}`"
          :class="selectedClass"
          :expanded="isSelected"
          @update:expanded="toggle?.()"
        >
        </OccurrenceDataFeedCard>
      </v-item>
    </v-item-group>
  </div>
  <v-divider />
  <div class="d-flex">
    <v-btn
      text="Add feed"
      size="small"
      class="flex-grow-1"
      :rounded="0"
      stacked
      prepend-icon="mdi-plus"
      variant="text"
      @click="addDataFeed()"
    />
    <v-divider vertical />
    <ConfirmDialog
      title="Reset data feeds"
      message="Are you sure you want to reset all data feeds? This will remove all custom filters and settings."
      @confirm="resetAll()"
    >
      <template #activator="{ props }">
        <v-btn
          text="Reset all"
          size="small"
          class="flex-grow-1"
          :rounded="0"
          stacked
          prepend-icon="mdi-restore"
          variant="text"
          v-bind="props"
        />
      </template>
    </ConfirmDialog>
  </div>
  <v-divider />
</template>

<script setup lang="ts">
import OccurrenceDataFeedCard from '@/components/occurrence/OccurrenceDataFeedCard.vue'
import { useDataFeeds } from './data_feeds'
import DatasetPicker from '../datasets/DatasetPicker.vue'
import ConfirmDialog from '../toolkit/ui/ConfirmDialog.vue'

const { addDataFeed, registry, resetAll } = useDataFeeds()
</script>

<style scoped lang="scss"></style>
