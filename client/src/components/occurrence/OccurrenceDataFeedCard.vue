<template>
  <v-card v-bind="$attrs">
    <template #title>
      <v-text-field
        v-model="feed.name"
        :label
        :placeholder
        density="compact"
        hide-details
        clearable
      />
    </template>
    <template #append>
      <v-menu>
        <template #activator="{ props }">
          <v-btn icon="mdi-dots-vertical" variant="plain" :rounded="100" v-bind="props" />
        </template>
        <v-card>
          <v-list density="compact">
            <MapStatsDialog
              :sites="remote.data.value"
              :title="`${feed.name} statistics`"
              :loading="remote.isPending.value"
            >
              <template #activator="{ props }">
                <v-list-item v-bind="props" title="Stats" prepend-icon="mdi-poll" />
              </template>
            </MapStatsDialog>
            <v-list-item
              prepend-icon="mdi-content-copy"
              title="Duplicate"
              @click="duplicateFeed(feed)"
            />
            <v-list-item prepend-icon="mdi-restore" title="Reset" @click="feed.filters = {}" />
            <v-list-item
              v-if="registry.length > 1"
              prepend-icon="mdi-delete"
              title="Delete"
              @click="deleteFeed(feed.id)"
            />
          </v-list>
        </v-card>
      </v-menu>
      <v-progress-circular
        v-if="remote.isPending.value"
        indeterminate
        size="small"
        color="warning"
      />
      <CardDialog
        v-else-if="remote.error.value"
        :title="`Error: ${feed.name}`"
        prepend-icon="mdi-alert"
        :max-width="500"
      >
        <template #activator="{ props }">
          <v-btn v-bind="props" icon="mdi-alert" color="error" :rounded="100" variant="text" />
        </template>
        <v-card-text>
          <v-alert color="error">
            An error occurred while fetching data for this feed:
            <strong>{{ remote.error.value?.detail }}</strong>
          </v-alert>
          <span class="text-caption">
            This may be due to a server error, please file an issue or try again later.
          </span>
        </v-card-text>
      </CardDialog>
      <v-icon
        v-else
        icon="mdi-circle"
        color="success"
        size="small"
        v-tooltip="{ text: 'Data is up to date', openDelay: 500, openOnClick: true }"
        @click="() => {}"
      />
      <v-btn
        :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
        color=""
        variant="text"
        :rounded="100"
        size="small"
        @click="() => (expanded = !expanded)"
      />
    </template>
    <v-expand-transition>
      <div v-show="expanded" class="">
        <v-divider />
        <OccurrenceDataFeedFilters v-model="feed.filters" class="pa-2" />
        <!-- <v-list-item prepend-icon="mdi-shape-polygon-plus" title="Polygon">
          <template #append>
            <v-btn icon="mdi-shape-polygon-plus" variant="text" />
            <v-btn icon="mdi-eye-outline" variant="text" />
            <v-switch color="primary" hide-details></v-switch>
          </template>
        </v-list-item> -->
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import OccurrenceDataFeedFilters from './OccurrenceDataFeedFilters.vue'
import { DataFeed, useDataFeeds } from './data_feeds'
import { useQuery } from '@tanstack/vue-query'
import { occurrencesBySiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import MapStatsDialog from '@/components/occurrence/OccurrenceStatsDialog.vue'
import CardDialog from '../toolkit/ui/CardDialog.vue'

defineProps<{
  placeholder?: string
  label?: string
}>()

const feed = defineModel<DataFeed>({
  default: () => reactive({ filters: {} })
})

const expanded = defineModel<boolean>('expanded', { default: false })

const { registry, register, deleteFeed, duplicateFeed } = useDataFeeds()

const remote = useQuery(
  computed(() =>
    occurrencesBySiteOptions({
      query: {
        ...feed.value.filters,
        habitats: feed.value.filters.habitats?.map(({ label }) => label)
      }
    })
  )
)

register(feed, remote)
</script>

<style scoped lang="scss"></style>
