<template>
  <v-card
    v-if="item"
    class="bg-surface fill-height w-100 d-flex flex-column"
    :title="item.code"
    prepend-icon="mdi-package-variant"
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
          </v-col>
          <v-col cols="12">
            <v-card
              title="Sampling"
              variant="elevated"
              prepend-icon="mdi-package-down"
              :subtitle="DateWithPrecision.format(item.event.performed_on)"
            >
              <v-card-text>
                <v-list density="compact">
                  <v-list-item
                    class="text-primary"
                    prepend-icon="mdi-map-marker-outline"
                    :title="item.event.site.name"
                    subtitle="Locality, CC"
                    :to="{ name: 'site-item', params: { code: item.event.site.code } }"
                  ></v-list-item>
                  <v-list-group value="Details" prepend-icon="mdi-text-box">
                    <template #activator="{ props }">
                      <v-list-item v-bind="props" title="Details" lines="two"></v-list-item>
                    </template>
                    <SamplingListItems :sampling="item.sampling" />
                  </v-list-group>
                  <v-divider></v-divider>
                  <v-list-item title="Other samples" prepend-icon="mdi-package-variant">
                    <v-chip
                      v-for="sample in item.sampling.samples.filter(({ id }) => id !== item!.id)"
                      :text="sample.identification.taxon.name"
                      :title="sample.category"
                      class="ma-1"
                      :to="{ name: 'biomat-item', params: { code: sample.code } }"
                      label
                    />
                  </v-list-item>
                </v-list>
              </v-card-text>
            </v-card>
          </v-col>
        </v-col>

        <v-col cols="12" lg="6">
          <v-col>
            <v-card title="Content" prepend-icon="mdi-hexagon-multiple" class="fill-height">
              <v-card-text v-if="item.external">
                Specimen quantity: <v-chip :text="item.external.quantity"></v-chip>
                <div>{{ item.external.content_description }}</div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col>
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
          </v-col>
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
  </v-card>
</template>

<script setup lang="ts">
import { SamplesService } from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import SamplingListItems from '@/components/events/SamplingListItems.vue'
import PersonChip from '@/components/people/PersonChip.vue'
import ArticleChip from '@/components/references/ArticleChip.vue'
import TaxonChip from '@/components/taxonomy/TaxonChip.vue'
import MetaChip from '@/components/toolkit/MetaChip.vue'
import { useFetchItem } from '@/composables/fetch_items'

const { code } = defineProps<{ code: string }>()

const { item } = useFetchItem(() => SamplesService.getBioMaterial({ path: { code } }))
</script>

<style scoped lang="scss"></style>
