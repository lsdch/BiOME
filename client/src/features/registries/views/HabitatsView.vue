<template>
  <v-card flat class="fill-height d-flex flex-column" rounded="0">
    <v-toolbar flat class="flex-shrink-0">
      <v-toolbar-title>Habitats</v-toolbar-title>
      <v-text-field
        v-model="search"
        label="Search habitats"
        prepend-inner-icon="mdi-magnify"
        clearable
        hide-details
        density="compact"
      />
      <v-btn
        icon="mdi-refresh"
        aria-label="Reload habitats"
        :loading="isFetching"
        @click="refetch()"
      />
      <v-btn v-if="canManage" prepend-icon="mdi-plus" text="New habitat" @click="openForm()" />
    </v-toolbar>
    <v-card-text v-if="error" class="flex-shrink-0 flex-grow-0">
      <v-alert type="error" variant="tonal" class="mt-4">Failed to retrieve habitats</v-alert>
    </v-card-text>
    <v-data-table
      :headers="headers"
      :items="tableItems"
      :loading="isPending"
      :search="search ?? ''"
      :filter-keys="['name', 'facet', 'description']"
      :group-by="[{ key: 'facet', order: 'asc' }]"
      :sort-by="[{ key: 'name', order: 'asc' }]"
      open-all
      item-value="id"
      hover
      class="flex-grow-1 justify-content-start overflow-y-auto"
    >
      <template #group-header="{ item, columns, toggleGroup, isGroupOpen }">
        <tr>
          <td :colspan="columns.length">
            <v-btn
              :icon="isGroupOpen(item) ? 'mdi-chevron-down' : 'mdi-chevron-right'"
              :aria-label="`${isGroupOpen(item) ? 'Collapse' : 'Expand'} ${item.value || 'No facet'}`"
              variant="text"
              size="small"
              @click="toggleGroup(item)"
            />
            <v-chip label class="font-weight-medium">{{ item.value || 'No facet' }}</v-chip>
            <v-chip size="small" class="ml-2">{{ item.items.length }}</v-chip>
          </td>
        </tr>
      </template>
      <template #item.description="{ value }">
        <span class="text-wrap">{{ value || '—' }}</span>
      </template>
      <template #item.created_at="{ item }">{{ item.created_at.toLocaleString() }}</template>
      <template #item.actions="{ item }">
        <v-btn
          icon="mdi-pencil"
          aria-label="Edit habitat"
          variant="text"
          size="small"
          @click="openForm(item)"
        />
        <v-btn
          icon="mdi-delete"
          aria-label="Delete habitat"
          variant="text"
          size="small"
          color="error"
          :disabled="deleteMutation.isPending.value"
          @click="openDeleteDialog(item)"
        />
      </template>
    </v-data-table>
  </v-card>
  <v-dialog v-model="formDialog" max-width="600" :persistent="saving">
    <v-card :title="editingHabitat ? 'Edit habitat' : 'Create habitat'">
      <v-form ref="form" :disabled="saving" @submit.prevent="submit">
        <v-card-text>
          <v-text-field
            v-model="model.name"
            label="Habitat name"
            autofocus
            :rules="[(value: string) => !!value?.trim() || 'A habitat name is required']"
          />
          <v-combobox v-model="model.facet" :items="facets" label="Facet (optional)" clearable />
          <v-textarea
            v-model="model.description"
            label="Description (optional)"
            :counter="4000"
            auto-grow
            :rules="[
              (value: string) => !value || value.length <= 4000 || 'Maximum 4000 characters'
            ]"
          />
          <v-alert v-if="formError" type="error" variant="tonal">{{ formError }}</v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn text="Cancel" :disabled="saving" @click="formDialog = false" />
          <v-btn type="submit" color="primary" text="Save" :loading="saving" />
        </v-card-actions>
      </v-form>
    </v-card>
  </v-dialog>
  <ConfirmDialog
    v-model="deleteDialog"
    :title="`Delete habitat ${habitatToDelete?.name ?? ''}?`"
    message="This habitat will also be removed from associated samplings."
    @cancel="habitatToDelete = undefined"
    @confirm="deleteHabitat"
  />
</template>

<script setup lang="ts">
import type { Habitat, HabitatInput } from '@/api'
import {
  createHabitatMutation,
  deleteHabitatMutation,
  getHabitatsOptions,
  getHabitatsQueryKey,
  updateHabitatMutation
} from '@/api/gen/@tanstack/vue-query.gen'
import ConfirmDialog from '@/components/toolkit/ui/ConfirmDialog.vue'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import { useMutation, useQuery, useQueryClient } from '@tanstack/vue-query'
import { computed, reactive, ref } from 'vue'
import { VForm } from 'vuetify/components'

const { data: habitats, error, isPending, isFetching, refetch } = useQuery(getHabitatsOptions())
const queryClient = useQueryClient()
const { feedback } = useFeedback()
const { isGranted } = useUserStore()
const canManage = computed(() => isGranted('admin'))
const search = ref('')
const tableItems = computed(() =>
  (habitats.value ?? []).map((habitat) => ({ ...habitat, facet: habitat.facet ?? '' }))
)
const headers = computed(() => [
  { key: 'data-table-group', title: 'Facet' },
  { key: 'name', title: 'Name' },
  { key: 'description', title: 'Description', cellProps: { class: 'text-muted' } },
  { key: 'created_at', title: 'Created' },
  ...(canManage.value ? [{ key: 'actions', title: 'Actions', sortable: false }] : [])
])
const facets = computed(() =>
  [...new Set((habitats.value ?? []).flatMap(({ facet }) => (facet ? [facet] : [])))].sort()
)
const formDialog = ref(false)
const form = ref<InstanceType<typeof VForm>>()
const editingHabitat = ref<Habitat>()
const habitatToDelete = ref<Habitat>()
const deleteDialog = ref(false)
const model = reactive({ name: '', facet: null as string | null, description: '' })
const formError = ref('')
const createMutation = useMutation(createHabitatMutation())
const updateMutation = useMutation(updateHabitatMutation())
const saving = computed(() => createMutation.isPending.value || updateMutation.isPending.value)
const deleteMutation = useMutation({
  ...deleteHabitatMutation(),
  onSuccess: async () => {
    feedback({ type: 'success', message: 'Habitat deleted' })
    await refreshHabitats()
  },
  onError: () => feedback({ type: 'error', message: 'Failed to delete habitat' })
})
function refreshHabitats() {
  return queryClient.invalidateQueries({ queryKey: getHabitatsQueryKey() })
}
function openForm(habitat?: Habitat) {
  editingHabitat.value = habitat
  model.name = habitat?.name ?? ''
  model.facet = habitat?.facet ?? null
  model.description = habitat?.description ?? ''
  formError.value = ''
  form.value?.resetValidation()
  formDialog.value = true
}
async function submit() {
  if (saving.value || !canManage.value || !(await form.value?.validate())?.valid) return
  formError.value = ''
  const body: HabitatInput = {
    name: model.name.trim(),
    facet: model.facet?.trim() || undefined,
    description: model.description || undefined
  }
  try {
    if (editingHabitat.value) {
      await updateMutation.mutateAsync({
        path: { id: editingHabitat.value.id },
        body: { ...body, facet: body.facet ?? null, description: body.description ?? null }
      })
    } else {
      await createMutation.mutateAsync({ body })
    }
    formDialog.value = false
    feedback({
      type: 'success',
      message: editingHabitat.value ? 'Habitat updated' : 'Habitat created'
    })
    await refreshHabitats()
  } catch {
    formError.value = 'Failed to save habitat. Check that the name is unique and try again.'
  }
}
function openDeleteDialog(habitat: Habitat) {
  habitatToDelete.value = habitat
  deleteDialog.value = true
}
function deleteHabitat() {
  if (!habitatToDelete.value || !canManage.value) return
  deleteMutation.mutate({ path: { id: habitatToDelete.value.id } })
  habitatToDelete.value = undefined
}
</script>
