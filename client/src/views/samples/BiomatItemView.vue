<template>
  <v-card class="w-100 d-flex flex-column" :title="code" flat :rounded="0" min-height="100%">
    <template #title>
      <v-card-title class="font-monospace text-wrap">
        {{ CodeIdentifier.textWrap(code) }}
      </v-card-title>
    </template>
    <template #prepend>
      <v-avatar variant="outlined">
        <v-icon icon="mdi-package-variant"></v-icon>
      </v-avatar>
    </template>
    <template #append>
      <v-btn color="primary" icon="mdi-pencil" variant="tonal" size="small" />
    </template>
    <template v-if="item" #subtitle>
      <div class="d-flex align-center ga-2">
        <OccurrenceCategoryChip :category="item.category" element="BioMaterial" size="small" />
        <v-chip v-if="item.is_type" prepend-icon="mdi-star-four-points" size="small" label>
          Nomenclatural type
        </v-chip>
        <v-chip prepend-icon="mdi-dna" size="small" label>
          {{
            item.has_sequences
              ? `${item.external?.content?.reduce((acc, { sequences }) => acc + sequences.length, 0)} sequences`
              : 'No attached sequences'
          }}
        </v-chip>
      </div>
    </template>
    <template v-if="item" #actions>
      <v-spacer />
      <MetaChip :meta="item.meta"></MetaChip>
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
    <v-card-text v-else-if="item" class="bg-main flex-grow-1 responsive-container">
      <v-row v-if="item.comments">
        <v-col cols="12">
          <v-card flat>
            <v-card-text class="d-flex justify-space-between">
              {{ item.comments }}
              <span class="text-muted text-caption">Comments</span>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" lg="6" class="d-flex flex-column ga-4 align-stretch justify-start">
          <v-card
            title="Identification"
            class="small-card-title"
            prepend-icon="mdi-microscope"
            :subtitle="DateWithPrecision.format(item.identification.identified_on)"
          >
            <template #append>
              <v-tooltip
                :text="
                  item.is_congruent
                    ? 'Bio material identification matches its sequences identification'
                    : 'Bio material identification contradicted by its sequences identification'
                "
                open-on-click
                location="end"
                origin="center"
              >
                <template #activator="{ props }">
                  <v-chip
                    v-bind="{
                      ...props,
                      ...(item?.is_congruent
                        ? {
                            color: 'success',
                            text: 'Congruent'
                          }
                        : {
                            color: 'warning',
                            text: 'Incongruent'
                          })
                    }"
                    size="small"
                  >
                  </v-chip>
                </template>
              </v-tooltip>
            </template>
            <v-card-text>
              <div class="d-flex align-center justify-space-between ga-1">
                <TaxonChip :taxon="item.identification.taxon" class="my-1" />
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
          <v-card title="Content" class="small-card-title" prepend-icon="mdi-hexagon-multiple">
            <template #append>
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
            </template>
            <v-list v-if="item.external">
              <v-list-item lines="three">
                <template #title>
                  <v-chip :text="item.external.quantity" size="small" />
                </template>
                <template #subtitle>
                  <div>{{ item.external.content_description ?? 'No further description' }}</div>
                  <div v-if="item.category === 'External' && !item.external.content">
                    No sequences registered
                  </div>
                </template>
                <template #append>
                  <span class="text-muted text-caption">Quantity</span>
                </template>
              </v-list-item>
              <template v-if="hasContentDetails">
                <v-divider class="my-3" />
                <v-expansion-panels>
                  <v-expansion-panel
                    title="Sequences by specimen"
                    :elevation="0"
                    :disabled="!item.external.content"
                  >
                    <template #text>
                      <v-treeview
                        :items="
                          item.external.content?.map<TreeViewSpecimenItem | ExternalBioMatSequence>(
                            ({ specimen, sequences }) => ({
                              code: specimen,
                              sequences: sequences.map((seq) => ({
                                ...seq,
                                props: {
                                  to: { name: 'sequence', params: { code: seq.code } }
                                }
                              }))
                            })
                          )
                        "
                        indent-lines
                        item-children="sequences"
                        item-title="code"
                        open-on-click
                      >
                        <template #title="{ title, item }">
                          <code v-if="'sequences' in item">{{ title }}</code>
                          <v-list-item-title v-else>{{
                            item.identification.taxon.name
                          }}</v-list-item-title>
                        </template>
                        <template #subtitle="{ item }">
                          <template v-if="!('sequences' in item)">
                            {{ item.label }}
                          </template>
                        </template>
                        <template #append="{ item }">
                          <v-chip
                            v-if="'sequences' in item"
                            :text="item.sequences.length"
                            size="small"
                          />
                          <GeneChip v-else :gene="item.gene" size="small" class="ma-1" />
                        </template>
                      </v-treeview>
                    </template>
                  </v-expansion-panel>
                </v-expansion-panels>
              </template>
            </v-list>
          </v-card>
          <v-card
            v-if="item.external"
            class="small-card-title"
            title="References"
            prepend-icon="mdi-newspaper-variant"
          >
            <template #append>
              <v-btn color="primary" variant="tonal" icon="mdi-link-variant" size="small" />
            </template>

            <v-divider />
            <v-list>
              <v-list-item>
                <ArticleChip v-for="article in item.published_in" :article class="ma-1" />
                <template #append>
                  <span class="text-muted text-caption">Publication(s)</span>
                </template>
              </v-list-item>
              <v-list-item>
                <DataSourceChip
                  v-if="item.external.original_source"
                  :source="item.external.original_source"
                />
                <span v-else class="text-muted text-caption"> Unknown </span>
                <template #append>
                  <span class="text-muted text-caption">Data source</span>
                </template>
              </v-list-item>
              <v-list-item>
                <b v-if="item.external.archive.collection">{{
                  item.external.archive.collection
                }}</b>
                <span class="text-muted"> Unknown </span>
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

    <v-divider />
  </v-card>
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
  DateWithPrecision,
  ExternalBioMatSequence,
  ExtSeqOrigin,
  SamplingWithSite
} from '@/api/adapters'
import { getBioMaterialOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingFormDialogMutation from '@/components/forms/SamplingFormDialogMutation.vue'
import OccurrenceCategoryChip from '@/components/occurrence/OccurrenceCategoryChip'
import OccurrenceSamplingCard from '@/components/occurrence/OccurrenceSamplingCard.vue'
import PersonChip from '@/components/people/PersonChip'
import ArticleChip from '@/components/references/ArticleChip'
import DataSourceChip from '@/components/references/DataSourceChip'
import GeneChip from '@/components/sequences/GeneChip'
import TaxonChip from '@/components/taxonomy/TaxonChip'
import MetaChip from '@/components/toolkit/MetaChip'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { computed } from 'vue'

type TreeViewSpecimenItem = {
  code: string
  sequences: ExternalBioMatSequence[]
}

const [samplingEdit, toggleSamplingEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()

const { data: item, error, isPending } = useQuery(getBioMaterialOptions({ path: { code } }))

const hasContentDetails = computed(() => {
  switch (item.value?.category) {
    case 'External':
      return !!item.value.external?.content
    case 'Internal':
      return true // TODO
    default:
      return false
  }
})
</script>

<style scoped lang="scss"></style>
