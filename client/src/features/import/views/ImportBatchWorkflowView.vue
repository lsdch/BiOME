<template>
  <v-card
    v-if="model"
    :title="model.batch.label"
    subtitle="Import batch"
    class="fill-height d-flex flex-column"
    flat
    rounded="0"
  >
    <template #append>
      <v-btn
        @click="materialize.mutate({ path: { id: uuid } })"
        :loading="materialize.isPending.value"
        variant="text"
        color="primary"
      >
        Materialize
      </v-btn>
    </template>
    <v-tabs class="bg-main flex-shrink-0" v-model="tab" mandatory>
      <v-tab value="overview">Overview</v-tab>
      <v-tab
        value="taxonomy"
        :prepend-icon="model.resolution_status.taxonomy ? 'mdi-check-circle' : 'mdi-alert'"
        text="Taxonomy"
      />
      <v-tab
        value="sampling-metadata"
        :prepend-icon="
          model.resolution_status.methods && model.resolution_status.fixatives
            ? 'mdi-check-circle'
            : 'mdi-alert'
        "
      >
        Sampling Metadata
      </v-tab>
      <v-tab
        value="bibliography"
        :prepend-icon="model.resolution_status.bibliography ? 'mdi-check-circle' : 'mdi-alert'"
      >
        Bibliography
      </v-tab>
    </v-tabs>
    <v-divider></v-divider>
    <v-tabs-window v-model="tab" class="flex-grow-1 overflow-y-auto">
      <v-tabs-window-item class="fill-height" value="overview">
        <v-card-text>
          <v-textarea label="Description" v-model="model.batch.description" disabled> </v-textarea>

          <AddBibliographyDialog :import_id="uuid">
            <template #activator="{ props }">
              <v-btn text="Add bibliography" prepend-icon="mdi-book-plus" v-bind="props"></v-btn>
            </template>
          </AddBibliographyDialog>
        </v-card-text>

        Model : {{ model }}

        <p class="text-error">Errors : {{ materialize.error }}</p>
      </v-tabs-window-item>
      <v-tabs-window-item class="fill-height" value="taxonomy">
        <TaxonomyResolver :uuid="uuid" :progress="model.GBIF" class="fill-height" />
      </v-tabs-window-item>
      <v-tabs-window-item class="fill-height" value="sampling-metadata">
        <SamplingMetadataResolver class="fill-height" :import_id="uuid" />
      </v-tabs-window-item>
      <v-tabs-window-item class="" value="bibliography">
        <BibliographyResolver class="fill-height" :uuid="uuid" />
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup lang="ts">
import { BatchImportsService, BatchSnapshot } from '@/api'
import {
  getImportStatusOptions,
  materializeBatchMutation
} from '@/api/gen/@tanstack/vue-query.gen.ts'
import { useMutation, useQuery } from '@tanstack/vue-query'
import { onMounted, ref } from 'vue'
import AddBibliographyDialog from '../components/AddBibliographyDialog.vue'
import BibliographyResolver from '../components/BibliographyResolver.vue'
import SamplingMetadataResolver from '../components/SamplingMetadataResolver.vue'
import TaxonomyResolver from '../components/TaxonomyResolver.vue'

const { uuid } = defineProps<{
  uuid: UUID
}>()

const model = ref<BatchSnapshot>()

type Tabs = 'overview' | 'taxonomy' | 'sampling-metadata'
const tab = ref<Tabs>('overview')

const materialize = useMutation(materializeBatchMutation())

const { data: batch, suspense } = useQuery(getImportStatusOptions({ path: { id: uuid } }))

async function monitor() {
  const { stream } = await BatchImportsService.trackImportStatus({
    path: { id: uuid },
    sseMaxRetryAttempts: 5,

    onSseEvent(event) {
      console.log('Received SSE event:')
      model.value = event.data
    }
  })
  for await (const event of stream) {
    console.log('Received event:', event)
  }
}

onMounted(async () => {
  await suspense().then((batch) => {
    if (batch.data) {
      model.value = batch.data
      monitor()
    }
  })
})
</script>

<style scoped lang="scss"></style>
