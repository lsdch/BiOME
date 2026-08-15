<template>
  <v-card
    class="w-100 d-flex flex-column wrap-card-title small-card-title"
    :title="CodeIdentifier.textWrap(code)"
    flat
    :rounded="0"
    min-height="100%"
  >
    <template #prepend>
      <ClickableAvatarIcon icon="mdi-dna" hover-icon="mdi-reload" @click="refetch()" />
    </template>
    <template #append>
      <v-btn v-if="canEdit" color="primary" icon="mdi-pencil" variant="tonal" size="small" />
    </template>
    <template v-if="item" #subtitle>
      <div class="d-flex align-center ga-2">
        <OccurrenceCategoryChip :category="item.category" size="small" />
        <GeneChip label size="small" :gene="item.gene" prepend-icon="mdi-tag" />
        <v-chip
          v-if="item.sequence"
          prepend-icon="mdi-chevron-right"
          text="ATCG"
          class="font-monospace"
          label
          size="small"
          v-tooltip="`Sequence available`"
          @click="
            () => {
              fasta?.groupItem.select(true)
              $nextTick(() => fasta?.$el.scrollIntoView({ behavior: 'smooth', block: 'start' }))
            }
          "
        />
        <MetaChip :meta="item.meta" size="small" />
      </div>
    </template>
    <v-divider />
    <v-card-text v-if="isPending" class="bg-main d-flex justify-center align-center">
      <CenteredSpinner class="bg-main" size="x-large" />
    </v-card-text>
    <PageErrors v-else-if="error" :error />
    <v-card-text v-else-if="item" class="flex-grow-1 bg-main responsive-container">
      <v-row>
        <v-col cols="12" lg="6" class="d-flex flex-column ga-4 align-stretch justify-start">
          <v-card v-if="item.label || item.comments">
            <v-list>
              <v-list-item v-if="item.label" prepend-icon="mdi-tag">
                {{ item.label }}
              </v-list-item>
              <v-list-item
                v-if="item.comments"
                prepend-icon="mdi-dots-horizontal"
                :class="['text-small', { 'font-italic': !item.comments }]"
              >
                {{ item.comments }}
              </v-list-item>
            </v-list>
          </v-card>
          <v-card
            v-if="item.category === 'External'"
            title="Occurrence"
            prepend-icon="mdi-package-variant"
          >
            <v-list>
              <v-list-item
                :title="item.occurrence.code"
                :subtitle="item.occurrence.identification.taxon.name"
                :to="{ name: 'occurrence-item', params: { code: item.occurrence.code } }"
                link
              />
            </v-list>
          </v-card>
          <v-card
            v-else-if="item.category === 'Internal'"
            title="Identification"
            prepend-icon="mdi-microscope"
            :subtitle="DateWithPrecision.format(item.identification.identified_on)"
          >
            <template #append>
              <v-tooltip
                :text="
                  item.occurrence.identification.taxon.id === item.identification.taxon.id
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
                      ...(item.occurrence.identification.taxon.id === item.identification.taxon.id
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
              <div class="d-flex justify-space-between align-center">
                <IdentificationChip :identification="item.identification" class="my-1" />
                <span v-if="item.identification.identified_by" class="text-no-wrap">
                  by
                  <v-chip
                    v-for="person in item.identification.identified_by"
                    :text="person"
                    :key="person"
                  />
                </span>
                <span class="text-muted" v-else>Curator unspecified</span>
              </div>
            </v-card-text>
            <v-divider />
            <v-list-item v-if="item.occurrence.verbatim_identification">
              <code>{{ item.occurrence.verbatim_identification }}</code>
              <template #append>
                <span class="text-muted text-caption">Original ident.</span>
              </template>
            </v-list-item>
            <v-list-item>
              <v-chip class="font-monospace" label size="small">
                {{ item.specimen_identifier }}
              </v-chip>
              <template #append>
                <span class="text-muted text-caption"> Specimen identifier </span>
              </template>
            </v-list-item>
            <v-list-item>
              <OccurrenceLinkChip
                v-if="item.occurrence"
                :biomat="item.occurrence"
                size="small"
                class="mx-1"
              />
              <span class="text-muted font-italic" v-else>No attached bio-material</span>
              <template #append>
                <span class="text-muted text-caption"> Source sample </span>
              </template>
            </v-list-item>
          </v-card>

          <v-card title="References" prepend-icon="mdi-newspaper-variant">
            <v-list>
              <v-list-item prepend-icon="mdi-database">
                <div class="d-flex ga-1 align-center">
                  <SeqRefChip
                    v-if="item.referenced_in"
                    v-for="seqRef in item.referenced_in"
                    :seqRef
                    size="small"
                  />
                  <span v-else class="text-muted">None registered</span>
                </div>
                <template #append>
                  <span class="text-muted text-caption">Repositories</span>
                </template>
              </v-list-item>
              <v-list-item
                :subtitle="item.occurrence.published_in ? undefined : 'No registered references'"
                prepend-icon="mdi-newspaper-variant"
              >
                <div class="d-flex align-center ga-2">
                  <ArticleChip
                    v-for="article in item.occurrence.published_in"
                    :article
                    size="small"
                  />
                </div>
                <template #append>
                  <span class="text-muted text-caption"> Publication(s) </span>
                </template>
              </v-list-item>
            </v-list>
          </v-card>
        </v-col>
        <v-col cols="12" lg="6">
          <OccurrenceSamplingCard :item="item.occurrence" @edit="toggleSamplingEdit(true)" />
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
                {{
                  item.sequence
                    ? `Fasta sequence (${item.sequence.length} bp)`
                    : 'Fasta unavailable'
                }}
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
    <SamplingFormDialogMutation
      v-if="item"
      v-model:dialog="samplingEdit"
      v-model="item.occurrence.sampling"
      :site="item.occurrence.sampling.site"
    />
    <v-divider />
  </v-card>
</template>

<script setup lang="ts">
import { CodeIdentifier, DateWithPrecision } from '@/api/adapters'
import { getSequenceOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingFormDialogMutation from '@/components/forms/SamplingFormDialogMutation.vue'
import OccurrenceCategoryChip from '@/features/occurrences/components/OccurrenceCategoryChip'
import OccurrenceLinkChip from '@/features/occurrences/components/OccurrenceLinkChip'
// import SamplingFormDialog from '@/components/events/SamplingFormDialog.vue'
import MetaChip from '@/components/toolkit/MetaChip'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import ClickableAvatarIcon from '@/components/toolkit/ui/ClickableAvatarIcon.vue'
import PageErrors from '@/components/toolkit/ui/PageErrors.vue'
import OccurrenceSamplingCard from '@/features/occurrences/components/OccurrenceSamplingCard.vue'
import ArticleChip from '@/features/registries/components/ArticleChip'
import GeneChip from '@/features/sequences/components/GeneChip'
import SeqRefChip from '@/features/sequences/components/SeqRefChip'
import IdentificationChip from '@/features/taxonomy/components/IdentificationChip'
import { useUserStore } from '@/stores/user'
import { useQuery } from '@tanstack/vue-query'
import { useClipboard, useToggle } from '@vueuse/core'
import { computed, nextTick, useTemplateRef } from 'vue'

const [samplingEdit, toggleSamplingEdit] = useToggle(false)

const { code } = defineProps<{ code: string }>()
nextTick(() => {
  document.title = code
})

const { data: item, error, isPending, refetch } = useQuery(getSequenceOptions({ path: { code } }))

const { copy } = useClipboard()
const fasta = useTemplateRef('fasta-seq')

const { isGranted, isOwner } = useUserStore()

const canEdit = computed(() => item.value && (isOwner(item.value) || isGranted('maintainer')))
</script>

<style scoped lang="scss"></style>
