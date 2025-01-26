<template>
  <v-card class="sampling-card" rounded="sb-0">
    <template #prepend>
      <code class="top-tag"> {{ cornerTag }} </code>
    </template>
    <template #append>
      <v-btn
        class="mx-1"
        color="error"
        icon="mdi-delete"
        size="small"
        variant="tonal"
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
      <v-divider></v-divider>
      <v-list-item title="Samples" prepend-icon="mdi-package-variant">
        <v-chip
          v-for="sample in sampling.samples"
          :text="sample.identification.taxon.name"
          :title="sample.category"
          class="ma-1"
        />
      </v-list-item>
      <v-divider></v-divider>
      <SamplingListItems :sampling />
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { Sampling, SamplingService } from '@/api'
import { useAppConfirmDialog } from '@/composables/confirm_dialog'
import { useFeedback } from '@/stores/feedback'
import { Duration } from 'luxon'
import SamplingListItems from './SamplingListItems.vue'

const { sampling } = defineProps<{ sampling: Sampling; cornerTag: string }>()
const emit = defineEmits<{
  edit: [sampling: Sampling]
  deleted: [sampling: Sampling]
}>()

const { askConfirm } = useAppConfirmDialog()
const { feedback } = useFeedback()

async function deleteSampling() {
  return askConfirm({
    title: 'Delete sampling action ?',
    message: 'All derived samples will be deleted as well for the database.'
  }).then(({ isCanceled }) => {
    if (isCanceled) return
    return SamplingService.deleteSampling({ path: { id: sampling.id } }).then(({ data, error }) => {
      if (!error) emit('deleted', data)
      else if (error.status === 404) feedback({ message: 'Sampling does not exist', type: 'error' })
      else {
        feedback({ message: 'Failed to delete sampling', type: 'error' })
        console.error(error)
      }
    })
  })
}
</script>

<style lang="scss">
@use 'vuetify';
.sampling-card {
  border-inline-start-width: 2px;
  border-inline-start-style: solid;
  border-inline-start-color: rgba(var(--v-theme-success), 0.7);

  .v-card-item {
    padding-top: 0px;
    padding-left: 0px;
    .top-tag {
      background-color: rgba(var(--v-theme-success), 0.7);
      height: 45px;
      padding: 10px;
      border-bottom-right-radius: 25%;
      font-weight: bold;
    }
  }
}
</style>
