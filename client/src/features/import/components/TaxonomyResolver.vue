<template>
  <v-chip-group v-model="statusFilter" filter multiple mandatory color="primary">
    <v-chip text="Pending" value="pending" key="pending"></v-chip>
    <v-chip text="Needs decision" value="needs_decision" key="needs_decision"></v-chip>
    <v-chip text="User-resolved" value="user_resolved" key="user_resolved"></v-chip>
    <v-chip text="Auto-resolved" value="auto_resolved" key="auto_resolved"></v-chip>
  </v-chip-group>
  <v-text-field v-model="searchTerm" label="Search" clearable></v-text-field>
  <v-data-table
    :headers="headers"
    :items="
      taxonResolutions?.filter((item) => !item.status || statusFilter.includes(item.status)) ?? []
    "
    :sort-by="[{ key: 'input_name', order: 'asc' }]"
    :search="searchTerm"
  >
    <template #item.input_name="{ item }">
      <div class="d-flex flex-column">
        {{ item.input_name }}
        <span v-if="item.from_resolution_name" class="text-label-small text-muted"
          >From: {{ item.from_resolution_name }}</span
        >
        <span v-else-if="item.sampling_target" class="text-label-small text-muted"
          >From sampling target</span
        >
      </div>
    </template>
    <template #item.candidates="{ item }">
      <div class="d-flex ga-2 align-center">
        <v-select
          :model-value="item.resolved_to"
          @update:model-value="
            (candidate_id) => {
              console.log('Resolving taxon for item.id:', item.id, 'to candidate_id:', candidate_id)
              resolveTaxon(item.id, candidate_id)
            }
          "
          :items="item.candidates ?? []"
          item-title="name"
          item-value="id"
          hide-details
          density="compact"
        >
          <template #selection="{ item: candidate }">
            <div class="d-flex ga-3">
              <v-chip :text="candidate.source"></v-chip>
              {{ candidate.name }}
              <span v-if="candidate.authorship" class="text-muted">
                {{ candidate.authorship }}
              </span>
            </div>
          </template>
          <template #append-item>
            <v-divider></v-divider>

            <manual-taxon-dialog
              :model-value="{
                name: item.input_name,
                authorship: item.input_authorship,
                rank: item.input_rank ? TaxonRank.fromString(item.input_rank) : undefined,
                status: item.input_authorship ? 'unreferenced' : 'unclassified'
              }"
              :import_id="uuid"
              :resolution_id="item.id"
              :resolution_names="taxonResolutions?.map((r) => r.input_name) ?? []"
              @created="refetch()"
            >
              <template #activator="{ props }">
                <v-list-item
                  title="Create manual taxon"
                  prepend-icon="mdi-plus"
                  v-bind="props"
                ></v-list-item>
              </template>
            </manual-taxon-dialog>
          </template>
          <template #item="{ item: candidate, props }">
            <v-list-item
              :title="candidate.name"
              :subtitle="candidate.authorship ?? 'No authorship'"
              v-bind="props"
            >
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
      </div>
    </template>
  </v-data-table>
  <span class="text-muted">
    {{ taxonResolutions }}
  </span>
</template>

<script setup lang="ts">
import { $ResolutionStatus, ProgressSnapshot, ResolutionStatus, TaxonRank } from '@/api'
import { getTaxonResolutionsOptions, resolveTaxonMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import ManualTaxonDialog from './ManualTaxonDialog.vue'

const { uuid, progress } = defineProps<{
  uuid: UUID
  progress?: ProgressSnapshot
}>()

const statusFilter = ref<ResolutionStatus[]>([...$ResolutionStatus.enum])

const searchTerm = ref<string>()

const headers: DataTableHeader[] = [
  {
    key: 'input_name',
    title: 'Input Name'
  },
  {
    key: 'input_rank',
    title: 'Rank',
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
    ...getTaxonResolutionsOptions({ path: { id: uuid } })
  }))
)

watch(
  () => progress,
  (newProgress) => {
    // if (newProgress?.Status === 'completed') {
    refetch()
    // }
  }
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
