<template>
  <v-card v-if="model" :title="model.workflow.label" subtitle="Import batch" class="fill-height" flat rounded="0">
    <template #append>
      <v-btn @click="materialize.mutate({ path: { id: uuid } })" :loading="materialize.isPending.value" variant="text"
        color="primary">
        Materialize
      </v-btn>
    </template>
    <v-tabs class="bg-main" v-model="tab" mandatory>
      <v-tab value="overview">Overview</v-tab>
      <v-tab value="taxonomy" :prepend-icon="model.resolution_status.taxonomy ? 'mdi-check-circle' : 'mdi-alert'"
        text="Taxonomy" />
      <v-tab value="sampling-metadata">Sampling Metadata</v-tab>
    </v-tabs>
    <v-divider></v-divider>
    <v-tabs-window v-model="tab" class="fill-height">
      <v-tabs-window-item class="fill-height" value="overview">
        <v-card-text>
          <v-textarea label="Description" v-model="model.workflow.description" disabled>
          </v-textarea>
        </v-card-text>

        Model : {{ model }}

        <p class="text-error">
          Errors : {{ materialize.error }}
        </p>
      </v-tabs-window-item>
      <v-tabs-window-item class="fill-height" value="taxonomy">
        <TaxonomyResolver :uuid="uuid" :progress="model.GBIF" class="fill-height" />
      </v-tabs-window-item>
      <v-tabs-window-item class="fill-height" value="sampling-metadata">
        <SamplingMetadataResolver class="fill-height" :import_id="uuid" />
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup lang="ts">
import { BatchImportsService, ImportEvent } from '@/api'
import { onMounted, ref } from 'vue'
import SamplingMetadataResolver from '../components/SamplingMetadataResolver.vue'
import TaxonomyResolver from '../components/TaxonomyResolver.vue'
import { useMutation } from '@tanstack/vue-query';
import { materializeBatchMutation } from '@/api/gen/@tanstack/vue-query.gen.ts';

const { uuid } = defineProps<{
  uuid: UUID
}>()

const model = ref<ImportEvent>()

type Tabs = 'overview' | 'taxonomy' | 'sampling-metadata'
const tab = ref<Tabs>('overview')

const materialize = useMutation(materializeBatchMutation())

async function monitor() {
  const { stream } = await BatchImportsService.importStatus({
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

onMounted(() => {
  monitor()
})
</script>

<style scoped lang="scss"></style>
