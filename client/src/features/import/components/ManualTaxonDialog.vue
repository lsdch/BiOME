<template>
  <FormDialog
    v-bind="$attrs"
    v-model="dialog"
    :title="'Add manual taxon candidate'"
    :loading="isPending"
    @submit="submit()"
  >
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-divider></v-divider>
    <v-card-text class="py-3">
      <v-text-field v-model="model.name" label="Taxon name" required v-bind="schema('name')" />
      <v-text-field
        v-model="model.authorship"
        label="Authorship (optional)"
        v-bind="schema('authorship')"
      />
      <TaxonRankPicker v-model="model.rank" label="Rank" required v-bind="schema('rank')" />
      <v-select
        v-model="model.status"
        :items="['unreferenced', 'unclassified']"
        label="Status"
        v-bind="schema('status')"
      />
      <v-combobox
        v-model="model.parent_name"
        :items="resolution_names"
        label="Parent taxon (optional)"
        v-bind="schema('parent_name')"
      />
    </v-card-text>
  </FormDialog>
</template>

<script setup lang="ts">
import { $TaxonStagingParams, TaxonRank, TaxonStatus } from '@/api'
import { createManualTaxonCandidateMutation } from '@/api/gen/@tanstack/vue-query.gen'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useSchemaBinding } from '@/composables/schema'
import TaxonRankPicker from '@/features/taxonomy/components/TaxonRankPicker'
import { useMutation } from '@tanstack/vue-query'
const dialog = defineModel<boolean>('dialog')
const model = defineModel<{
  name: string
  authorship?: string
  rank?: TaxonRank
  parent_name?: string
  status: TaxonStatus
}>({ required: true })

const { import_id, resolution_id } = defineProps<{
  import_id: UUID
  resolution_id: UUID
  resolution_names?: string[]
}>()

const emit = defineEmits<{
  created: []
}>()

const { mutateAsync, error, isPending } = useMutation(createManualTaxonCandidateMutation())

const { schema } = useSchemaBinding($TaxonStagingParams)

function submit() {
  if (!model.value.name || !model.value.rank || !model.value.parent_name) {
    throw new Error('Name, rank, and parent taxon are required')
  }

  return mutateAsync(
    {
      path: {
        id: import_id
      },
      body: {
        resolution_id,
        name: model.value.name,
        authorship: model.value.authorship,
        rank: model.value.rank,
        parent_name: model.value.parent_name,
        status: model.value.status
      }
    },
    {
      onSuccess: () => {
        emit('created')
        dialog.value = false
      }
    }
  )
}
</script>

<style scoped lang="scss"></style>
