<template>
  <v-card
    class="bg-surface w-100 d-flex flex-column"
    :title="code"
    flat
    :rounded="0"
    min-height="100%"
  >
    <template #prepend>
      <v-avatar variant="outlined">
        <v-icon icon="mdi-package-variant"></v-icon>
      </v-avatar>
    </template>
    <template #append>
      <v-btn color="primary" icon="mdi-pencil" variant="tonal" size="small" />
    </template>
    <template v-if="item" #subtitle>
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
        {{ item.category }} bio-material
      </v-chip>
      <v-chip
        v-if="item.is_type"
        class="mx-1"
        prepend-icon="mdi-star-four-points"
        size="small"
        label
      >
        Nomenclatural type
      </v-chip>
    </template>
    <template v-if="item" #actions>
      <v-spacer></v-spacer>
      <MetaChip :meta="item.meta"></MetaChip>
    </template>
    <v-divider></v-divider>
    <CenteredSpinner v-if="isPending" height="100%" size="x-large" />
    <v-card-text v-else-if="error">
      <v-alert color="error" icon="mdi-alert">
        Failed to retrieve bio material informations
      </v-alert>
    </v-card-text>
    <v-card-text v-else-if="item" class="flex-grow-1">
      <v-row v-if="item.comments">
        <v-col cols="12">
          <v-card flat>
            <v-card-text>
              {{ item.comments }}
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" lg="6">
          <div class="w-100 my-1">
            <v-card
              title="Identification"
              class="fill-height"
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
                <TaxonChip :taxon="item.identification.taxon" class="my-1" />
                <span class="text-no-wrap">
                  by
                  <PersonChip :person="item.identification.identified_by" />
                </span>
                <div v-if="item.external?.original_taxon">
                  Originally tagged as: {{ item.external.original_taxon }}
                </div>
              </v-card-text>
            </v-card>
          </div>
          <div class="w-100 my-1">
            <OccurrenceSamplingCard :item @edit="toggleSamplingEdit(true)" />
          </div>
        </v-col>

        <v-col cols="12" lg="6">
          <div class="w-100 my-1">
            <v-card title="Content" prepend-icon="mdi-hexagon-multiple" class="fill-height">
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
                    ></v-chip>
                  </template>
                </v-tooltip>
              </template>
              <v-card-text v-if="item.external">
                <v-list-item
                  :subtitle="item.external.content_description ?? 'No further description'"
                >
                  <template #title>
                    Specimen quantity: <v-chip :text="item.external.quantity" />
                  </template>
                  <template #subtitle>
                    <div>{{ item.external.content_description ?? 'No further description' }}</div>
                    <div v-if="item.category === 'External' && !item.external.content">
                      No sequences registered
                    </div>
                  </template>
                </v-list-item>
                <template v-if="hasContentDetails">
                  <v-divider class="my-3"></v-divider>
                  <v-expansion-panels>
                    <v-expansion-panel
                      title="Sequences by specimen"
                      :elevation="0"
                      :disabled="!item.external.content"
                    >
                      <template #text>
                        <v-treeview
                          :items="
                            item.external.content?.map(({ specimen, sequences }) => ({
                              code: specimen,
                              sequences
                            }))
                          "
                          item-children="sequences"
                          item-title="code"
                          open-on-click
                        >
                          <template #title="{ title }">
                            <code>{{ title }}</code>
                          </template>
                          <template #item="{ item }">
                            <v-treeview-item
                              :title="item.identification.taxon.name"
                              :subtitle="item.label"
                              :prepend-icon="ExtSeqOrigin.icon(item.origin)"
                              :to="{ name: 'sequence', params: { code: item.code } }"
                            >
                              <template #append>
                                <GeneChip :gene="item.gene" size="small" />
                              </template>
                            </v-treeview-item>
                          </template>
                        </v-treeview>
                      </template>
                    </v-expansion-panel>
                  </v-expansion-panels>
                </template>
              </v-card-text>
            </v-card>
          </div>
          <div class="w-100 my-1">
            <v-card v-if="item.external" title="References" prepend-icon="mdi-newspaper-variant">
              <template #append>
                <v-btn color="primary" variant="tonal" icon="mdi-link-variant" size="small"></v-btn>
              </template>
              <v-card-text>
                <v-list>
                  <v-list-item title="Published in">
                    <ArticleChip v-for="article in item.published_in" :article class="ma-1" />
                  </v-list-item>
                  <v-divider class="my-1"></v-divider>
                  <v-list-item title="In collection">
                    <b>{{ item.external.archive.collection }}</b>
                  </v-list-item>
                  <v-list-item title="Item vouchers">
                    <v-chip
                      v-for="v in item.external.archive.vouchers"
                      :text="v"
                      size="small"
                      class="ma-1"
                    />
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </div>
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
  </v-card>
  <SamplingFormDialog
    v-if="item"
    v-model="samplingEdit"
    :edit="item.sampling"
    :event="item.event"
    @updated="
      (sampling) => {
        item!.sampling = sampling
        toggleSamplingEdit(false)
      }
    "
  />
</template>

<script setup lang="ts">
import { DateWithPrecision, ExtSeqOrigin } from '@/api/adapters'
import { getBioMaterialOptions } from '@/api/gen/@tanstack/vue-query.gen'
import SamplingFormDialog from '@/components/events/SamplingFormDialog.vue'
import OccurrenceSamplingCard from '@/components/occurrence/OccurrenceSamplingCard.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import ArticleChip from '@/components/references/ArticleChip.vue'
import GeneChip from '@/components/sequences/GeneChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import CenteredSpinner from '@/components/toolkit/ui/CenteredSpinner'
import { useQuery } from '@tanstack/vue-query'
import { useToggle } from '@vueuse/core'
import { computed } from 'vue'

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
