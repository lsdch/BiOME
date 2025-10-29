<template>
  <div class="bg-main">
    <v-card class="w-100 d-flex flex-column" :title="code" flat :rounded="0" min-height="100%">
      <template #title>
        <v-card-title class="d-flex">
          <div class="font-monospace text-wrap w-auto d-flex align-center">
            <v-menu v-if="item?.code_history" :close-on-content-click="false" target="parent">
              <template #activator="{ props }">
                {{ code }}
                <v-btn
                  icon="mdi-history"
                  variant="plain"
                  size="small"
                  color=""
                  class="ml-2"
                  v-bind="props"
                />
              </template>
              <CodeHistoryCard :codeHistory="item.code_history" />
            </v-menu>
          </div>
        </v-card-title>
      </template>
      <template #prepend>
        <ClickableAvatarIcon icon="mdi-package-variant" hover-icon="mdi-reload" @click="reload()" />
        <!-- <v-avatar variant="outlined" @click="refetch()">
          <v-icon icon="mdi-package-variant"></v-icon>
        </v-avatar> -->
      </template>
      <template #append>
        <v-btn-group size="small" divided variant="outlined">
          <v-btn icon="mdi-dna"></v-btn>
          <v-btn icon="mdi-pencil" />
        </v-btn-group>
      </template>
      <template v-if="item" #subtitle>
        <div class="d-flex align-center ga-2">
          <OccurrenceCategoryChip :category="item.category" size="small" />
          <v-chip v-if="item.is_type" prepend-icon="mdi-star-four-points" size="small" label>
            Nomenclatural type
          </v-chip>
          <v-chip
            v-if="!item.has_sequences"
            prepend-icon="mdi-dna"
            size="small"
            label
            text="No sequences"
            variant="tonal"
            color="#777"
          />
          <MetaChip :meta="item.meta" size="small" />
        </div>
      </template>

      <v-divider />

      <v-card-text class="bg-main d-flex align-center justify-center" v-if="isPending">
        <CenteredSpinner size="x-large" class="bg-main" />
      </v-card-text>
      <v-card-text v-else-if="error">
        <v-alert color="error" icon="mdi-alert">
          Failed to retrieve bio material informations
        </v-alert>
      </v-card-text>
      <template v-else-if="item">
        <v-tabs v-model="currentTab" v-if="item.external?.sequences">
          <v-tab value="occurrence" prepend-icon="mdi-text-box-outline">Occurrence</v-tab>
          <v-tab value="sequences" prepend-icon="mdi-dna" :disabled="!item.external?.sequences">
            Sequences
            <template #append>
              <v-badge
                v-if="item.external?.sequences?.length"
                :content="item.external.sequences.length"
                color="primary"
                inline
              >
              </v-badge>
            </template>
          </v-tab>
          <v-tab value="datasets" prepend-icon="mdi-folder-table" :disabled="!item.datasets">
            Datasets
            <template #append>
              <v-badge
                v-if="item.datasets?.length"
                :content="item.datasets.length"
                color="primary"
                inline
              >
              </v-badge>
            </template>
          </v-tab>
        </v-tabs>
        <v-tabs-window v-model="currentTab" crossfade>
          <v-tabs-window-item value="occurrence">
            <v-card-text class="bg-main flex-grow-1 responsive-container">
              <v-row>
                <v-col cols="12" lg="6" class="d-flex flex-column ga-4 align-stretch justify-start">
                  <v-card v-if="item.comments">
                    <v-card-text class="d-flex justify-space-between">
                      {{ item.comments }}
                      <span class="text-muted text-caption">Comments</span>
                    </v-card-text>
                  </v-card>
                  <v-card
                    title="Identification"
                    class="small-card-title"
                    prepend-icon="mdi-microscope"
                    :subtitle="DateWithPrecision.format(item.identification.identified_on)"
                  >
                    <template v-if="item.internal" #append>
                      <CongruenceChip :is_congruent="item.internal.is_congruent" />
                    </template>
                    <v-card-text>
                      <div class="d-flex align-center justify-space-between ga-1">
                        <IdentificationChip :identification="item.identification" class="my-1" />
                        <span v-if="item.identification.identified_by" class="text-no-wrap">
                          by
                          <PersonChip :person="item.identification.identified_by" />
                        </span>
                        <span class="text-muted" v-else>Curator unspecified</span>
                      </div>
                      <div v-if="item.external?.original_taxon">
                        Originally tagged as: {{ item.external.original_taxon }}
                      </div>
                    </v-card-text>
                  </v-card>
                  <v-card
                    title="Content"
                    class="small-card-title"
                    prepend-icon="mdi-hexagon-multiple"
                  >
                    <!-- <template #append>
                    <v-tooltip
                      :text="
                        item.is_homogenous
                          ? 'Sequences all identify a single taxon'
                          : 'Sequences identify different taxa'
                      "
                      open-on-click
                      location="end"
                      origin="end"
                    >
                      <template #activator="{ props }">
                        <v-chip
                          v-if="hasContentDetails"
                          size="small"
                          v-bind="{
                            ...props,
                            ...(item.is_homogenous
                              ? {
                                  color: 'success',
                                  text: 'Homogenous'
                                }
                              : {
                                  color: 'warning',
                                  text: 'Heterogenous'
                                })
                          }"
                        />
                      </template>
                    </v-tooltip>
                  </template> -->
                    <v-list v-if="item.external">
                      <v-list-item>
                        <template #title>
                          <v-chip
                            v-if="item.external.quantity"
                            :text="item.external.quantity"
                            size="small"
                          />
                          <span v-else class="text-muted text-caption">Unknown</span>
                        </template>
                        <template #append>
                          <span class="text-muted text-caption">Quantity</span>
                        </template>
                      </v-list-item>
                      <v-list-item>
                        <template #title>
                          <span
                            :class="[
                              'text-caption',
                              { 'text-muted': !item.external.content_description }
                            ]"
                          >
                            {{ item.external.content_description ?? 'None specified' }}
                          </span>
                        </template>
                        <template #append>
                          <span class="text-muted text-caption">Content description</span>
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card>
                  <v-card
                    v-if="item.external"
                    class="small-card-title"
                    title="References"
                    prepend-icon="mdi-newspaper-variant"
                  >
                    <v-divider />
                    <v-list>
                      <v-list-item>
                        <span v-if="!item.external.published_in" class="text-muted text-caption">
                          Unknown
                        </span>
                        <div v-else class="d-flex ga-2">
                          <ArticleChip v-for="article in item.external.published_in" :article />
                        </div>
                        <template #append>
                          <span class="text-muted text-caption">Publication(s)</span>
                        </template>
                      </v-list-item>
                      <v-list-item>
                        <DataSourceChip v-for="source in item.external.sources" :source />
                        <span v-if="!item.external.sources" class="text-muted text-caption">
                          Unknown
                        </span>
                        <template #append>
                          <span class="text-muted text-caption">Data sources</span>
                        </template>
                      </v-list-item>
                      <v-list-item>
                        <b v-if="item.external.archive.collection">{{
                          item.external.archive.collection
                        }}</b>
                        <span class="text-muted text-caption"> Unknown </span>
                        <template #append>
                          <span class="text-muted text-caption">Collection</span>
                        </template>
                        <template #subtitle>
                          <v-chip
                            v-for="v in item.external.archive.vouchers"
                            :text="v"
                            size="small"
                            class="ma-1"
                          />
                        </template>
                      </v-list-item>
                    </v-list>
                  </v-card>
                </v-col>

                <v-col cols="12" lg="6">
                  <div class="w-100">
                    <OccurrenceSamplingCard :item @edit="toggleSamplingEdit(true)" />
                  </div>
                </v-col>
              </v-row>
            </v-card-text>
          </v-tabs-window-item>
          <v-tabs-window-item value="sequences">
            <CRUDTable
              v-if="item.external"
              :headers="sequenceTable.headers"
              entityName="Sequence"
              :items="item.external.sequences"
            >
              <template #item.code="{ value }">
                <RouterLink :to="{ name: 'sequence', params: { code: value } }">
                  {{ CodeIdentifier.textWrap(value) }}
                </RouterLink>
              </template>
              <template #item.gene="{ value: gene }">
                <GeneChip :gene size="small" />
              </template>
              <template #item.is_identifying="{ value }">
                <v-icon
                  v-if="value"
                  icon="mdi-tag"
                  color="success"
                  v-tooltip="`Used for identification`"
                ></v-icon>
              </template>
            </CRUDTable>
          </v-tabs-window-item>
          <v-tabs-window-item value="datasets">
            <CRUDTable
              v-if="item.datasets"
              :headers="datasetTable.headers"
              entityName="Dataset"
              :items="item.datasets"
            >
              <template #item.label="{ value, item }">
                <RouterLink :to="{ name: 'occurrence-dataset-item', params: { slug: item.slug } }">
                  {{ value }}
                </RouterLink>
              </template>

              <template #item.maintainers="{ value }">
                <div class="d-flex ga-2">
                  <PersonChip
                    v-for="person in value"
                    :key="person.code"
                    :person="person"
                    size="small"
                  />
                </div>
              </template>
            </CRUDTable>
          </v-tabs-window-item>
        </v-tabs-window>
      </template>
    </v-card>
  </div>
  <SamplingFormDialogMutation
    v-if="item"
    v-model:dialog="samplingEdit"
    v-model="item.sampling"
    :site="item.sampling.site"
    @updated="
      (sampling: SamplingWithSite) => {
        item!.sampling = sampling
        toggleSamplingEdit(false)
      }
    "
  />
</template>

<script setup lang="ts">
import {
  CodeIdentifier,
  Dataset,
  DateWithPrecision,
  ExternalSequence,
  Gene,
  SamplingWithSite
} from '@/api/adapters'
import { getBioMaterialOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingFormDialogMutation from '@/components/forms/SamplingFormDialogMutation.vue'
import CongruenceChip from '@/features/occurrences/components/CongruenceChip'
import OccurrenceCategoryChip from '@/features/occurrences/components/OccurrenceCategoryChip'
import OccurrenceSamplingCard from '@/features/occurrences/components/OccurrenceSamplingCard.vue'
import PersonChip from '@/features/people/components/PersonChip'
import ArticleChip from '@/features/registries/components/ArticleChip'
import DataSourceChip from '@/features/registries/components/DataSourceChip'
import GeneChip from '@/features/sequences/components/GeneChip'
import TaxonChip from '@/features/taxonomy/components/TaxonChip'
import MetaChip from '@/components/toolkit/MetaChip'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import ClickableAvatarIcon from '@/components/toolkit/ui/ClickableAvatarIcon.vue'
import { useFeedback } from '@/stores/feedback'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { ref } from 'vue'
import CodeHistoryCard from '../components/CodeHistoryCard.vue'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'

const [samplingEdit, toggleSamplingEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()

const {
  data: item,
  error,
  isPending,
  refetch
} = useQuery(getBioMaterialOptions({ path: { code } }))

const { feedback } = useFeedback()

async function reload() {
  await refetch()
  feedback({ message: 'Data reloaded', type: 'info' })
}

type Tabs = 'occurrence' | 'sequences'
const currentTab = ref<Tabs>('occurrence')

const sequenceTable: { headers: CRUDTableHeader<ExternalSequence>[] } = {
  headers: [
    { key: 'code', title: 'Code', cellProps: { class: 'font-monospace' } },
    { key: 'gene', title: 'Gene', sort: (a: Gene, b: Gene) => a.code.localeCompare(b.code) },
    { key: 'specimen_identifier', title: 'Specimen', cellProps: { class: 'font-monospace' } },
    { key: 'is_identifying', title: '', width: 0, align: 'end' }
  ]
}
const datasetTable: { headers: CRUDTableHeader<Dataset>[] } = {
  headers: [
    { key: 'label', title: 'Label' },
    { key: 'maintainers', title: 'Maintainers' }
  ]
}
</script>

<style scoped lang="scss"></style>
