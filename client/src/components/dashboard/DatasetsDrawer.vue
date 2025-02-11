<template>
  <v-navigation-drawer v-if="!mobile" location="right" permanent>
    <v-tabs v-model="tab" density="compact">
      <v-tab prepend-icon="mdi-creation-outline" value="pinned" stacked size="small" />
      <v-tab prepend-icon="mdi-update" value="recent" stacked size="small" />
    </v-tabs>
    <v-divider />
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="pinned">
        <v-list density="compact" nav>
          <v-list-subheader> Pinned datasets </v-list-subheader>
          <v-divider class="my-2" />
          <CenteredSpinner v-if="isPending" :height="100" />
          <v-alert v-else-if="error" color="error" class="text-caption" density="compact">
            Failed to load pinned datasets
          </v-alert>
          <v-list-item v-else-if="pinned?.length" v-for="dataset in pinned">
            <DatasetDisplayCard :dataset />
          </v-list-item>
          <v-alert v-else color="primary" class="text-caption">
            No datasets were pinned yet.
          </v-alert>
        </v-list>
      </v-tabs-window-item>
      <v-tabs-window-item value="recent">
        <v-list density="compact" nav>
          <v-list-subheader> Latest dataset activity </v-list-subheader>
          <v-divider class="my-2" />
          <CenteredSpinner v-if="recentPending" :height="100" />
          <v-alert v-else-if="recentError" color="error" class="text-caption" density="compact">
            Failed to load latest datasets
          </v-alert>
          <v-list-item v-else-if="recent?.length" v-for="dataset in recent">
            <DatasetDisplayCard :dataset />
          </v-list-item>
          <v-alert v-else color="primary" class="text-caption"> No registered datasets </v-alert>
        </v-list>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { ref } from 'vue'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import DatasetDisplayCard from './DatasetDisplayCard.vue'
import { useDisplay } from 'vuetify'

const { mobile } = useDisplay()

const tab = ref<'pinned' | 'recent'>('pinned')

const {
  data: pinned,
  error,
  isPending
} = useQuery(listDatasetsOptions({ query: { pinned: true, limit: 10, orderBy: 'label' } }))

const {
  data: recent,
  error: recentError,
  isPending: recentPending
} = useQuery(listDatasetsOptions({ query: { limit: 10, orderBy: 'meta.lastUpdated' } }))
</script>

<style scoped lang="scss"></style>
