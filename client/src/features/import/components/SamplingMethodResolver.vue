<template>
  <v-data-table :items="resolutions" :headers>
    <template #item.resolved_method_id="{ item }">
      <SamplingMethodPicker :prepend-icon="undefined" :model-value="item.resolved_method_id"
        @update:model-value="(resolved_method_id: UUID) => resolve(item, resolved_method_id)" item-value="id"
        hide-details density="compact" />
    </template>
  </v-data-table>
</template>

<script setup lang="ts">
import { SamplingMethodResolution } from '@/api';
import { getMethodsResolutionOptions, resolveMethodMutation } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingMethodPicker from '@/features/registries/components/SamplingMethodPicker.vue'
import { useMutation, useQuery } from '@tanstack/vue-query'

const { import_id } = defineProps<{
  import_id: UUID
}>()
const { data: resolutions, refetch } = useQuery(getMethodsResolutionOptions({ path: { id: import_id } }))
const { mutateAsync: resolveMethod, error } = useMutation(resolveMethodMutation())

function resolve(item: SamplingMethodResolution, resolved_method_id: UUID) {
  resolveMethod({
    path: { id: item.import_id },
    body: { input_text: item.input_text, resolved_method_id: resolved_method_id, status: 'selected' }
  }, {
    onSuccess() {
      refetch()
    }
  })
}

const headers = [
  {
    key: 'input_text',
    title: 'Input Text'
  },
  {
    key: 'resolved_method_id',
    title: 'Resolved Method',
    sortable: false
  }
]
</script>

<style scoped lang="scss"></style>
