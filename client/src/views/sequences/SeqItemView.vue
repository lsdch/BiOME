<template>
  <v-card
    v-if="item"
    class="bg-surface fill-height w-100 d-flex flex-column"
    :title="item.code"
    flat
    :rounded="0"
  >
    <template #prepend>
      <v-avatar variant="outlined">
        <v-icon icon="mdi-dna"></v-icon>
      </v-avatar>
    </template>
    <template #append>
      <v-btn color="primary" icon="mdi-pencil" variant="tonal" size="small" />
    </template>
    <template #subtitle>
      <v-chip
        class="mx-1"
        size="small"
        label
        v-bind="
          {
            Internal: {
              prependIcon: 'mdi-cube-scan',
              color: 'primary'
            },
            External: {
              prependIcon: 'mdi-arrow-collapse-all',
              color: 'warning'
            }
          }[item.category]
        "
      >
        {{ item.category }} sequence
      </v-chip>
      <GeneChip label size="small" :gene="item.gene" class="mx-1" prepend-icon="mdi-tag" />
      <v-chip
        :text="item.external?.origin"
        :prepend-icon="ExtSeqOrigin.icon(item.external!.origin)"
        :title="ExtSeqOrigin.description(item.external!.origin)"
        label
        class="mx-1"
        size="small"
      />
      <v-chip
        v-if="item.sequence"
        prepend-icon="mdi-chevron-right"
        text="ATCG"
        class="font-monospace mx-1"
        label
        size="small"
        title="Sequence available"
        @click="
          fasta?.groupItem.select(true),
            $nextTick(() => fasta?.$el.scrollIntoView({ behavior: 'smooth', block: 'start' }))
        "
      />
      <!-- @click="fasta?.scrollIntoView()" -->
      <!-- <v-chip
        v-if="item.is_type"
        class="mx-1"
        prepend-icon="mdi-star-four-points"
        size="small"
        label
      >
        Nomenclatural type
      </v-chip> -->
    </template>
    <template #actions>
      <v-spacer></v-spacer>
      <MetaChip :meta="item.meta"></MetaChip>
    </template>
    <v-divider></v-divider>
    <v-card-text class="flex-grow-1">
      <v-row>
        <v-col>
          <v-list-item :title="item.comments" prepend-icon="mdi-comment-outline"> </v-list-item>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" lg="6">
          <v-col cols="12">
            <v-card
              title="Identification"
              class="fill-height"
              prepend-icon="mdi-microscope"
              :subtitle="DateWithPrecision.format(item.identification.identified_on)"
            >
              <template #append>
                <v-tooltip
                  :text="
                    item.external?.source_sample?.identification.taxon.id ===
                    item.identification.taxon.id
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
                        ...(item.external?.source_sample?.identification.taxon.id ===
                        item.identification.taxon.id
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
                <TaxonChip :taxon="item.identification.taxon" class="my-1" />
                <span class="text-no-wrap">
                  by
                  <PersonChip :person="item.identification.identified_by" />
                </span>
              </v-card-text>
              <v-divider></v-divider>
              <v-card-actions v-if="item.external?.original_taxon">
                Originally tagged as: <code>{{ item.external.original_taxon }}</code>
              </v-card-actions>
            </v-card>
          </v-col>
          <v-col cols="12">
            <OccurrenceSamplingCard :item @edit="toggleSamplingEdit(true)" />
          </v-col>
        </v-col>
        <v-col cols="12" lg="6">
          <v-col cols="12">
            <v-card title="Origin sample" prepend-icon="mdi-package-variant">
              <template #subtitle>
                <v-chip
                  v-if="item.external?.source_sample"
                  :text="item.external.source_sample.identification.taxon.name"
                  :to="{ name: 'biomat-item', params: { code: item.external.source_sample.code } }"
                  prepend-icon="mdi-link-variant"
                  color="primary"
                  label
                  size="small"
                  class="mx-1"
                ></v-chip>
              </template>
              <v-card-text>
                <v-list-item v-if="item.external" title="Specimen identifier">
                  <template #subtitle>
                    <code>{{ item.external?.specimen_identifier }} </code>
                  </template>
                </v-list-item>
              </v-card-text>
            </v-card>
          </v-col>

          <v-col cols="12">
            <v-card v-if="item.external" title="References" prepend-icon="mdi-newspaper-variant">
              <template #append>
                <v-btn color="primary" variant="tonal" icon="mdi-link-variant" size="small"></v-btn>
              </template>
              <v-card-text>
                <v-list>
                  <v-list-item title="Databases" prepend-icon="mdi-database">
                    <SeqRefChip
                      v-for="seqRef in item.external.referenced_in"
                      :seq-ref
                      class="ma-1"
                    />
                  </v-list-item>
                  <!-- <v-list-item title="Published in">
                  <ArticleChip
                    v-if="item.external.published_in"
                    :article="item.external.published_in"
                    class="ma-1"
                  />
                </v-list-item> -->
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-expansion-panels>
            <v-expansion-panel :disabled="!item.sequence" ref="fasta-seq">
              <template #title>
                <v-chip
                  prepend-icon="mdi-chevron-right"
                  text="ATCG"
                  class="mr-2 font-monospace"
                  label
                  size="small"
                />
                {{ item.sequence ? 'Fasta sequence' : 'Fasta unavailable' }}
              </template>
              <template #text>
                <v-card
                  v-if="item.sequence"
                  variant="tonal"
                  @click="copy(`>${item.code}\n${item.sequence}`)"
                >
                  <v-snackbar
                    contained
                    class="text-center"
                    location="center"
                    origin="center"
                    color="success"
                    :width="100"
                    :timeout="500"
                    activator="parent"
                    open-on-click
                    location-strategy="connected"
                  >
                    <template #text> <v-icon icon="mdi-content-copy" /> Copied </template>
                  </v-snackbar>
                  <v-card-text>
                    <code id="fasta-seq">
                      >{{ item.code }}
                      <div>
                        {{ item.sequence.match(new RegExp('.{1,' + 10 + '}', 'g'))?.join(' ') }}
                      </div>
                    </code>
                  </v-card-text>
                </v-card>
              </template>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-col>
      </v-row>
    </v-card-text>
    <SamplingFormDialog
      v-if="item"
      v-model="samplingEdit"
      :edit="item.sampling"
      :event="item.event"
    />
    <v-divider></v-divider>
  </v-card>
</template>

<script setup lang="ts">
import { SequencesService } from '@/api'
import { DateWithPrecision, ExtSeqOrigin } from '@/api/adapters'
import SamplingFormDialog from '@/components/events/SamplingFormDialog.vue'
import OccurrenceSamplingCard from '@/components/occurrence/OccurrenceSamplingCard.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import GeneChip from '@/components/sequences/GeneChip.vue'
import SeqRefChip from '@/components/sequences/SeqRefChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import { useFetchItem } from '@/composables/fetch_items'
import { useClipboard, useScroll, useToggle } from '@vueuse/core'
import { useTemplateRef } from 'vue'

const [samplingEdit, toggleSamplingEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()

const { item, fetch } = useFetchItem(() => SequencesService.getSequence({ path: { code } }))

item.value = await fetch()

const { copy } = useClipboard()
const fasta = useTemplateRef('fasta-seq')
</script>

<style scoped lang="scss"></style>
