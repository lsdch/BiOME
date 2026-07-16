<template>
  <v-data-table :headers="headers" :items="taxonResolutions ?? []" :sort-by="[{ key: 'input_name', order: 'asc' }]">
    <template #item.candidates="{ item }">
      <v-select :model-value="item.resolved_to" @update:model-value="
        (candidate_id) => {
          console.log('Resolving taxon for item.id:', item.id, 'to candidate_id:', candidate_id)
          resolveTaxon(item.id, candidate_id)
        }
      " :items="item.candidates ?? []" item-title="name" item-value="id" hide-details density="compact">
        <template #selection="{ item: candidate }">
          <div class="d-flex ga-3">
            <v-chip :text="candidate.source"></v-chip>
            {{ candidate.name }}
            <span v-if="candidate.authorship" class="text-muted">
              {{ candidate.authorship }}
            </span>
          </div>
        </template>
        <template #item="{ item: candidate, props }">
          <v-list-item :title="candidate.name" :subtitle="candidate.authorship ?? 'No authorship'" v-bind="props">
            <template #append>
              <v-chip size="small" :text="candidate.rank"></v-chip>
              <v-chip size="small" :text="candidate.status"></v-chip>
              <v-chip size="small" :text="candidate.priority"></v-chip>
              <v-chip size="small" :text="candidate.source"></v-chip>
              <v-chip size="small" :text="candidate.gbif_id"></v-chip>
            </template>
          </v-list-item>
        </template>
      </v-select>
    </template>
  </v-data-table>
  <span class="text-muted">
    {{ taxonResolutions }}
  </span>
</template>

<script setup lang="ts">
import { ProgressSnapshot } from '@/api'
import { getTaxonResolutionsOptions, resolveTaxonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const { uuid, progress } = defineProps<{
  uuid: UUID
  progress?: ProgressSnapshot
}>()

const headers: DataTableHeader[] = [
  {
    key: 'input_name',
    title: 'Input Name',
    width: 0
  },
  {
    key: 'status',
    title: 'Status',
    width: 0
  },
  {
    key: 'candidates',
    title: 'Candidates',
    sortable: false
  }
]

const { data: taxonResolutions, refetch } = useQuery(
  computed(() => ({
    enabled: progress?.Status === 'completed',
    ...getTaxonResolutionsOptions({ path: { id: uuid } })
  }))
)

const resolutionErrors = ref<{
  [resolutionID: string]: string
}>({})

const { mutateAsync: resolveAsync } = useMutation(resolveTaxonMutation())

async function resolveTaxon(resolutionID: UUID, candidateID: UUID) {
  console.log('Resolving taxon for resolutionID:', resolutionID, 'to candidateID:', candidateID)
  await resolveAsync({
    path: { id: uuid },
    body: {
      resolution_id: resolutionID,
      candidate_id: candidateID
    }
  })
    .then(() => {
      console.log('Resolved taxon for resolutionID:', resolutionID, 'to candidateID:', candidateID)
      delete resolutionErrors.value[resolutionID]
      refetch()
    })
    .catch((error) => {
      resolutionErrors.value[resolutionID] = error.message
    })
}
</script>

<style scoped lang="scss"></style>
