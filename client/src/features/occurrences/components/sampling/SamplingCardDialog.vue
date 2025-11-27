<template>
  <CardDialog :title="DateWithPrecision.format(sampling.performed_on)">
    <template #subtitle>
      Sampling <code>#{{ sampling.number }}</code>
    </template>
    <template #activator="props">
      <slot name="activator" v-bind="props"></slot>
    </template>
    <template #append>
      <v-btn
        class="mx-1"
        color="error"
        icon="mdi-delete"
        size="small"
        variant="tonal"
        :loading="isPending"
        @click="deleteSampling"
      />
      <v-btn
        class="mx-1"
        color="primary"
        icon="mdi-pencil"
        size="small"
        variant="tonal"
        @click="emit('edit', sampling)"
      />
    </template>

    <v-list density="compact">
      <v-divider />
      <v-list-item prepend-icon="mdi-package-variant">
        <template #append>
          <span class="text-caption text-muted">Occurrences</span>
        </template>
        <v-chip
          v-if="sampling.occurrences?.length"
          v-for="sample in sampling.occurrences"
          :text="sample.identification.taxon.name"
          class="ma-1"
        />
        <span v-else class="text-muted font-italic">None registered</span>
      </v-list-item>
      <v-divider />
      <SamplingListItems :sampling />
    </v-list>
  </CardDialog>
</template>

<script setup lang="ts">
import { DateWithPrecision, Sampling } from '@/api'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import { useAppConfirmDialog } from '@/composables/confirm_dialog'
import { useFeedback } from '@/stores/feedback'
import { useMutation } from '@tanstack/vue-query'
import { deleteSamplingMutation } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingListItems from './SamplingListItems.vue'

const { sampling } = defineProps<{ sampling: Sampling }>()

const emit = defineEmits<{
  edit: [sampling: Sampling]
  deleted: [sampling: Sampling]
}>()

const { askConfirm } = useAppConfirmDialog()
const { feedback } = useFeedback()

const { mutate, isPending } = useMutation({
  ...deleteSamplingMutation(),
  onSuccess: (deleted) => emit('deleted', deleted),
  onError: (error) => {
    if (error.status === 404) feedback({ message: 'Sampling does not exist', type: 'error' })
    else {
      feedback({ message: 'Failed to delete sampling', type: 'error' })
      console.error(error)
    }
  }
})

async function deleteSampling() {
  return askConfirm({
    title: 'Delete sampling action ?',
    message: 'All derived samples will be deleted as well for the database.'
  }).then(async ({ isCanceled }) => {
    if (isCanceled) return
    mutate({ path: { id: sampling.id } })
  })
}
</script>

<style scoped lang="scss"></style>
