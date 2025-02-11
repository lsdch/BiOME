<template>
  <v-list density="compact" nav>
    <div class="d-flex justify-space-between">
      <v-list-subheader>
        {{ qualifier[0].toLocaleUpperCase() + qualifier.substring(1) }} datasets
      </v-list-subheader>
      <v-btn-toggle
        v-if="isGranted('Admin')"
        v-model="editPin"
        :rounded="true"
        variant="text"
        color="warning"
        density="compact"
      >
        <v-btn
          :value="true"
          icon="mdi-pin"
          size="small"
          style="aspect-ratio: 1"
          :width="30"
          title="Edit pinned datasets"
        ></v-btn>
      </v-btn-toggle>
    </div>
    <v-divider class="my-2" />
    <CenteredSpinner v-if="isPending" :height="100" />
    <v-alert v-else-if="error" color="error" class="text-caption" density="compact">
      Failed to load {{ qualifier }} datasets
    </v-alert>
    <v-list-item v-else-if="data?.length" v-for="dataset in data">
      <DatasetDisplayCard :dataset />
      <template v-if="editPin" #append>
        <DatasetPinButton
          class="ml-2"
          :model-value="dataset"
          @update:model-value="invalidateQueries()"
        />
      </template>
    </v-list-item>
    <v-alert v-else color="primary" class="text-caption"> No datasets to show. </v-alert>
  </v-list>
</template>

<script setup lang="ts">
import { ListDatasetsData, Options } from '@/api'
import { listDatasetsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useUserStore } from '@/stores/user'
import DatasetPinButton from '@/views/datasets/DatasetPinButton.vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { ref } from 'vue'
import CenteredSpinner from '../toolkit/ui/CenteredSpinner'
import DatasetDisplayCard from './DatasetDisplayCard.vue'

const editPin = ref(false)
const { isGranted } = useUserStore()

const { query } = defineProps<{
  qualifier: 'pinned' | 'latest'
  query: Options<ListDatasetsData>['query']
}>()

const { data, error, isPending } = useQuery(listDatasetsOptions({ query }))

const client = useQueryClient()
function invalidateQueries() {
  client.invalidateQueries({ queryKey: [{ _id: 'listDatasets' }] })
}
</script>

<style scoped lang="scss"></style>
