<template>
  <div class="bg-main">
    <v-card
      title="Context"
      class="small-card-title"
      subtitle="Default settings for all new data feeds"
      flat
      :rounded="0"
    >
      <template #append>
        <v-switch v-model="contextEnabled" color="primary" hide-details></v-switch>
      </template>
      <v-expand-transition>
        <div v-if="contextEnabled" class="bg-main">
          <DataFeedsContextPicker />
        </div>
      </v-expand-transition>
    </v-card>
    <v-divider class="mb-3" :thickness="4" />
    <v-item-group selected-class="ma-3" class="d-flex flex-column ga-3 pa-2 bg-main">
      <v-item v-for="(_, i) in feeds" #="{ isSelected, toggle, selectedClass }">
        <OccurrenceDataFeedCard
          v-model="feeds[i]"
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
import DatasetPicker from '../datasets/DatasetPicker.vue'
import TaxonPicker from '../taxonomy/TaxonPicker.vue'
import ConfirmDialog from '../toolkit/ui/ConfirmDialog.vue'
import InlineHelp from '../toolkit/ui/InlineHelp.vue'
import { useDataFeeds } from './data_feeds'
import DataFeedsContextPicker from './DataFeedsContextPicker.vue'

const { addDataFeed, feeds, resetAll, context, contextEnabled, applyContext } = useDataFeeds()
</script>

<style scoped lang="scss"></style>
