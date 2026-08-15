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
      bibResolutions?.filter((item) => !item.status || statusFilter.includes(item.status)) ?? []
    "
    :sort-by="[{ key: 'input_name', order: 'asc' }]"
    v-bind="$attrs"
  >
    <template #item.verbatim="{ item }">
      <span v-if="item.verbatim">{{ item.verbatim }}</span>
      <v-chip v-else-if="item.doi">
        <div class="d-flex align-end">
          ?.filter((item) => !item.status || statusFilter.includes(item.status)) ?? []
          <span class="text-overline text-muted text-label-medium">DOI:&nbsp;</span>
          <span>{{ item.doi }}</span>
        </div>
      </v-chip>
    </template>
    <template #item.candidates="{ item }">
      <v-select
        :model-value="item.resolved_id"
        @update:model-value="
          (candidate_id) => {
            console.log(
              'Resolving publication for item.id:',
              item.id,
              'to candidate_id:',
              candidate_id
            )
            resolvePublication(item.id, candidate_id)
          }
        "
        :items="item.candidates ?? []"
        item-value="id"
        item-title="authors"
        hide-details
        density="compact"
      >
        <template #selection="{ item: candidate }">
          <div class="d-flex ga-3">
            <v-chip :text="candidate.source"></v-chip>
            {{ shortAuthors(candidate.authors) }} ({{ candidate.year }}) - {{ candidate.title }}
          </div>
        </template>
        <template #item="{ item: candidate, props }">
          <v-list-item
            v-bind="props"
            :title="candidate.authors?.join(', ') ?? 'No authors'"
            :subtitle="candidate.title ?? 'No title'"
          >
            <template #subtitle>
              <div class="d-flex ga-2 align-center">
                <v-chip size="small" :text="candidate.year"></v-chip>
                {{ candidate.title ?? candidate.verbatim }}
              </div>
            </template>
            <template #append>
              <v-chip size="small" :text="candidate.source"></v-chip>
              <v-chip v-if="candidate.doi" size="small" :text="candidate.doi"></v-chip>
              <v-chip
                v-if="candidate.score >= 0"
                size="small"
                :text="candidate.score.toFixed(2)"
              ></v-chip>
            </template>
          </v-list-item>
        </template>
      </v-select>
    </template>
  </v-data-table>
  <span class="text-muted">
    <!-- {{ bibResolutions }} -->
  </span>
</template>

<script setup lang="ts">
import { $ResolutionStatus, ProgressSnapshot, ResolutionStatus } from '@/api'
import {
  getBibliographyResolutionsOptions,
  resolvePublicationMutation,
  resolveTaxonMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const { uuid, progress } = defineProps<{
  uuid: UUID
  progress?: ProgressSnapshot
}>()

const statusFilter = ref<ResolutionStatus[]>([...$ResolutionStatus.enum])

const searchTerm = ref<string>()

function shortAuthors(authors: string[] | undefined): string {
  if (!authors || authors.length === 0) {
    return ''
  }
  const firstAuthor = authors[0]
  // get the last name of the first author
  if (authors.length === 1) {
    return firstAuthor
  }
  return `${firstAuthor} et al.`
}

const headers: DataTableHeader[] = [
  {
    key: 'verbatim',
    title: 'Verbatim',
    value(item, fallback) {
      return item.verbatim ?? item.doi
    }
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

const { data: bibResolutions, refetch } = useQuery(
  computed(() => ({
    ...getBibliographyResolutionsOptions({ path: { id: uuid } })
  }))
)

const resolutionErrors = ref<{
  [resolutionID: string]: string
}>({})

const { mutateAsync: resolveAsync } = useMutation(resolvePublicationMutation())

async function resolvePublication(resolutionID: UUID, candidateID: UUID) {
  console.log(
    'Resolving publication for resolutionID:',
    resolutionID,
    'to candidateID:',
    candidateID
  )
  await resolveAsync({
    path: { id: uuid },
    body: {
      resolution_id: resolutionID,
      candidate_id: candidateID
    }
  })
    .then(() => {
      console.log(
        'Resolved publication for resolutionID:',
        resolutionID,
        'to candidateID:',
        candidateID
      )
      delete resolutionErrors.value[resolutionID]
      refetch()
    })
    .catch((error) => {
      resolutionErrors.value[resolutionID] = error.message
    })
}
</script>

<style scoped lang="scss"></style>
