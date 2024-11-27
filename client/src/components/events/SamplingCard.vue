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
      <v-list-item title="Targets" prepend-icon="mdi-bullseye">
        <v-chip
          v-if="sampling.target.kind === 'Unknown'"
          class="ma-1"
          :text="sampling.target.kind"
          prepend-icon="mdi-help-circle"
        />
        <v-chip
          v-else-if="sampling.target.kind === 'Community'"
          class="ma-1"
          :text="sampling.target.kind"
          prepend-icon="mdi-grid"
        />
        <v-chip
          v-else-if="sampling.target.kind === 'Taxa'"
          class="ma-1"
          v-for="taxon in sampling.target.target_taxa"
          :text="taxon.name"
          rounded
        />
      </v-list-item>
      <v-list-item title="Duration" prepend-icon="mdi-update">
        <code>
          {{ Duration.fromObject({ minutes: sampling.duration }).toFormat("hh'h' mm'm'") }}
        </code>
      </v-list-item>
      <v-list-item title="Fixatives" prepend-icon="mdi-snowflake">
        <v-chip v-for="f in sampling.fixatives" :text="f.label" />
        <v-list-item-subtitle v-if="sampling.fixatives.length == 0" class="font-italic">
          Unknown
        </v-list-item-subtitle>
      </v-list-item>
      <v-list-item title="Methods" prepend-icon="mdi-hook">
        <v-chip v-for="m in sampling.methods" :text="m.label" />
        <v-list-item-subtitle v-if="sampling.methods.length == 0" class="font-italic">
          Unknown
        </v-list-item-subtitle>
      </v-list-item>
      <v-list-item prepend-icon="mdi-image-filter-hdr-outline">
        <v-list-item title="Habitat" class="px-0">
          <v-list-item-subtitle v-if="sampling.habitats.length == 0" class="font-italic">
            Unknown
          </v-list-item-subtitle>
          <v-chip
            v-for="h in sampling.habitats"
            class="ma-1"
            :text="h.label"
            :title="h.description"
          />
        </v-list-item>
        <v-list-item title="Access points" class="px-0">
          <v-list-item-subtitle v-if="sampling.access_points.length == 0" class="font-italic">
            Unknown
          </v-list-item-subtitle>
          <v-chip v-for="access in sampling.access_points" class="ma-1" :text="access" />
        </v-list-item>
      </v-list-item>
      <v-list-item
        title="Comments"
        v-if="sampling.comments"
        prepend-icon="mdi-note-edit-outline"
        :subtitle="sampling.comments"
      />
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import { Sampling, SamplingService } from '@/api'
import { useAppConfirmDialog } from '@/composables/confirm_dialog'
import { useFeedback } from '@/stores/feedback'
import { Duration } from 'luxon'

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
  .v-list-item__prepend > .v-icon ~ .v-list-item__spacer {
    width: 16px;
  }
}
</style>
