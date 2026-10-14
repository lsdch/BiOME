<template>
  <v-alert v-if="error" type="error" variant="tonal">
    Failed to retrieve habitat resolutions.
    <v-btn text="Retry" variant="text" @click="refetch()" />
  </v-alert>
  <v-data-table
    :items="resolutions ?? []"
    :headers="headers"
    :loading="isPending"
    item-value="input_text"
  >
    <template #item.status="{ item }">
      <v-chip size="small">{{ statusLabels[item.status] }}</v-chip>
    </template>
    <template #item.habitats="{ item }">
      <v-chip v-for="id in item.habitat_ids ?? []" :key="id" class="ma-1" size="small">
        {{ habitats?.find((habitat) => habitat.id === id)?.name ?? id }}
      </v-chip>
      <v-chip
        v-for="habitat in item.create_habitats ?? []"
        :key="habitat.name"
        class="ma-1"
        size="small"
        color="success"
        prepend-icon="mdi-plus"
      >
        {{ habitat.name }}
      </v-chip>
    </template>
    <template #item.actions="{ item }">
      <v-btn text="Resolve" variant="text" size="small" @click="openResolution(item)" />
    </template>
  </v-data-table>
  <v-dialog v-model="dialog" max-width="750" :persistent="saving">
    <v-card title="Resolve habitats">
      <v-form ref="form" :disabled="saving" @submit.prevent="save">
        <v-card-text>
          <p class="mb-4 text-pre-wrap">{{ editing?.input_text }}</p>
          <v-select v-model="decision" :items="decisions" label="Decision" />
          <template v-if="decision === 'selected' || decision === 'request_creation'">
            <v-alert v-if="habitatsError" type="error" variant="tonal" class="mb-4">
              Failed to retrieve habitats.
              <v-btn text="Retry" variant="text" @click="refetchHabitats()" />
            </v-alert>
            <v-autocomplete
              v-model="habitatIds"
              :items="habitats ?? []"
              :loading="habitatsPending"
              label="Existing habitats"
              item-title="name"
              item-value="id"
              multiple
              chips
              closable-chips
              clearable
              :rules="[
                (ids: string[]) =>
                  decision !== 'selected' || ids.length > 0 || 'Select at least one habitat'
              ]"
            />
          </template>
          <template v-if="decision === 'request_creation'">
            <v-alert type="info" variant="tonal" class="mb-4">
              New habitats will be created when the batch is materialized.
            </v-alert>
            <v-card
              v-for="(habitat, index) in creations"
              :key="index"
              variant="outlined"
              class="mb-3"
            >
              <v-card-text>
                <v-text-field
                  v-model="habitat.name"
                  label="Habitat name"
                  :rules="[
                    (value: string) => !!value?.trim() || 'A habitat name is required',
                    (value: string) =>
                      creations.filter(
                        (entry) => entry.name.trim().toLowerCase() === value.trim().toLowerCase()
                      ).length === 1 || 'Names must be unique'
                  ]"
                />
                <v-text-field v-model="habitat.facet" label="Facet (optional)" />
                <v-textarea
                  v-model="habitat.description"
                  label="Description (optional)"
                  auto-grow
                  :rows="2"
                  :counter="4000"
                  :rules="[
                    (value: string) =>
                      !value || [...value].length <= 4000 || 'Maximum 4000 characters'
                  ]"
                />
                <v-btn
                  text="Remove"
                  variant="text"
                  :disabled="creations.length === 1"
                  @click="creations.splice(index, 1)"
                />
              </v-card-text>
            </v-card>
            <v-btn
              text="Add habitat"
              prepend-icon="mdi-plus"
              variant="text"
              @click="creations.push({ name: '' })"
            />
          </template>
          <v-alert v-if="saveError" type="error" variant="tonal" class="mt-4">
            {{ saveError }}
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text="Cancel" :disabled="saving" @click="dialog = false" />
          <v-btn text="Save" type="submit" color="primary" :loading="saving" />
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type { HabitatInput, SamplingHabitatResolution, VocabResolutionStatus } from '@/api'
import {
  getHabitatsResolutionOptions,
  getHabitatsResolutionQueryKey,
  listHabitatsOptions,
  resolveHabitatsMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, ref, watch } from 'vue'
import { VForm } from 'vuetify/components'

const { import_id } = defineProps<{ import_id: UUID }>()
const {
  data: resolutions,
  error,
  isPending,
  refetch
} = useQuery(computed(() => getHabitatsResolutionOptions({ path: { id: import_id } })))
const {
  data: habitats,
  error: habitatsError,
  isPending: habitatsPending,
  refetch: refetchHabitats
} = useQuery(listHabitatsOptions())
const queryClient = useQueryClient()
const { mutateAsync, isPending: saving } = useMutation(resolveHabitatsMutation())
const headers = [
  { key: 'input_text', title: 'Input Text' },
  { key: 'status', title: 'Status' },
  { key: 'habitats', title: 'Resolved Habitats', sortable: false },
  { key: 'actions', title: '', sortable: false }
]
const statusLabels: Record<VocabResolutionStatus, string> = {
  pending: 'Pending',
  auto_resolved: 'Auto-resolved',
  selected: 'Selected',
  request_creation: 'Creation requested',
  discard: 'Ignored'
}
const decisions = [
  { title: 'Select existing habitats', value: 'selected' },
  { title: 'Request new habitats', value: 'request_creation' },
  { title: 'Ignore this value', value: 'discard' },
  { title: 'Leave pending', value: 'pending' }
]
const dialog = ref(false)
const form = ref<InstanceType<typeof VForm>>()
const editing = ref<SamplingHabitatResolution>()
const decision = ref<VocabResolutionStatus>('selected')
const habitatIds = ref<string[]>([])
const creations = ref<HabitatInput[]>([])
const saveError = ref('')

watch(
  () => import_id,
  () => {
    dialog.value = false
    editing.value = undefined
  }
)

function openResolution(item: SamplingHabitatResolution) {
  editing.value = item
  decision.value = item.status === 'auto_resolved' ? 'selected' : item.status
  habitatIds.value = [...(item.habitat_ids ?? [])]
  creations.value = item.create_habitats?.length
    ? item.create_habitats.map((habitat) => ({ ...habitat }))
    : [{ name: item.input_text }]
  saveError.value = ''
  form.value?.resetValidation()
  dialog.value = true
}

async function save() {
  if (saving.value || !editing.value || !(await form.value?.validate())?.valid) return
  const item = editing.value
  saveError.value = ''
  try {
    await mutateAsync({
      path: { id: item.import_id },
      body: {
        input_text: item.input_text,
        status: decision.value,
        habitat_ids:
          decision.value === 'selected' || decision.value === 'request_creation'
            ? habitatIds.value
            : undefined,
        create_habitats:
          decision.value === 'request_creation'
            ? creations.value.map((habitat) => ({
                name: habitat.name.trim(),
                facet: habitat.facet?.trim() || undefined,
                description: habitat.description || undefined
              }))
            : undefined
      }
    })
    if (editing.value === item) dialog.value = false
    await queryClient.invalidateQueries({
      queryKey: getHabitatsResolutionQueryKey({ path: { id: item.import_id } })
    })
  } catch {
    if (editing.value === item)
      saveError.value =
        'Failed to save habitat resolution. Check the selected habitats and names, then try again.'
  }
}
</script>
