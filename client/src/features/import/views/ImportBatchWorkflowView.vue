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
        prepend-icon="mdi-file-download"
        text="Download raw data"
        variant="tonal"
        @click="downloadRawFile()"
      ></v-btn>
      <v-btn
        @click="
          materialize.mutate(
            { path: { id: uuid } },
            {
              onSuccess: () => {
                router?.replace(router.currentRoute.value)
              }
            }
          )
        "
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

  <v-dialog
    :model-value="
      model?.status === 'materializing' && !model?.materialization_steps.materialization_complete
    "
    persistent
  >
    <v-card title="Materializing batch...">
      <v-list>
        <v-list-item title="Filling GBIF dependencies">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.fill_gbif_dependencies" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
        <v-list-item title="Materializing taxa">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.materialize_taxa" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
        <v-list-item title="Materializing samplings">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.materialize_samplings" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
        <v-list-item title="Materializing occurrences">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.materialize_occurrences" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
        <v-list-item title="Materializing bibliography">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.materialize_bibliography" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
        <v-list-item title="Generating occurrence codes">
          <template #prepend>
            <v-icon v-if="model?.materialization_steps.refresh_occurrence_codes" color="success"
              >mdi-check-circle</v-icon
            >
            <v-progress-circular
              v-else
              indeterminate
              color="primary"
              size="20"
            ></v-progress-circular>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { BatchImportsService, BatchSnapshot, DownloadRawFileData } from '@/api'
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
import { client } from '@/api/gen/client.gen.ts'
import { useRouter } from 'vuetify/lib/composables/router.mjs'

const { uuid } = defineProps<{
  uuid: UUID
}>()

const model = ref<BatchSnapshot>()

const router = useRouter()

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

function downloadRawFile() {
  const url = client.buildUrl<DownloadRawFileData>({
    path: { id: uuid },
    url: '/import-batches/{id}/raw'
  })
  window.open(url, '_blank')
}
</script>

<style scoped lang="scss"></style>
