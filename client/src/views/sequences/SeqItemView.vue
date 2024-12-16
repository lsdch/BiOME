<template>
  <v-card
    v-if="item"
    class="bg-surface fill-height w-100 d-flex flex-column"
    :title="item.code"
    prepend-icon="mdi-dna"
    flat
    :rounded="0"
  >
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
      <GeneChip label size="small" :gene="item.gene" prepend-icon="mdi-tag" />
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
        <v-col>
          <v-card title="Origin sample" prepend-icon="mdi-package-variant">
            <template #subtitle>
              <v-chip
                :text="item.external?.origin"
                :prepend-icon="ExtSeqOrigin.icon(item.external!.origin)"
                label
                class="mx-1"
                size="small"
              ></v-chip>
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
        <v-col>
          <OccurrenceSamplingCard :item @edit="toggleSamplingEdit(true)" />
        </v-col>
        <v-col>
          <v-card v-if="item.external" title="References" prepend-icon="mdi-newspaper-variant">
            <template #append>
              <v-btn color="primary" variant="tonal" icon="mdi-link-variant" size="small"></v-btn>
            </template>
            <v-card-text>
              <v-list>
                <v-list-item title="Databases" prepend-icon="mdi-database">
                  <SeqRefChip v-for="seqRef in item.external.referenced_in" :seq-ref class="ma-1" />
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
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { SamplesService, SequencesService } from '@/api'
import { ExtSeqOrigin } from '@/api/adapters'
import OccurrenceSamplingCard from '@/components/occurrence/OccurrenceSamplingCard.vue'
import ArticleChip from '@/components/references/ArticleChip.vue'
import GeneChip from '@/components/sequences/GeneChip.vue'
import SeqRefChip from '@/components/sequences/SeqRefChip.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import { useFetchItem } from '@/composables/fetch_items'
import { useToggle } from '@vueuse/core'

const [samplingEdit, toggleSamplingEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()

const { item } = useFetchItem(() => SequencesService.getSequence({ path: { code } }))
</script>

<style scoped lang="scss"></style>
