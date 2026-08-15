<template>
  <ImportBatchWorkflowView v-if="item?.status !== 'completed'" :uuid></ImportBatchWorkflowView>
  <ImportBatchItem v-else :uuid />
</template>

<script setup lang="ts">
import { getImportBatchOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'
import ImportBatchWorkflowView from './ImportBatchWorkflowView.vue'
import ImportBatchItem from '../components/ImportBatchItem.vue'

const { uuid } = defineProps<{
  uuid: UUID
}>()

const { data: item, refetch } = useQuery(
  computed(() => ({
    ...getImportBatchOptions({ path: { id: uuid } })
  }))
)
</script>

<style scoped lang="scss"></style>
