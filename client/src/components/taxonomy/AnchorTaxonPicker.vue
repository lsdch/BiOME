<template>
  <v-card>
    <v-card-title class="d-flex justify-space-between">
      <span>
        <v-icon icon="mdi-sitemap" start />
        Import anchor taxon
      </span>
    </v-card-title>
    <v-card-text>
      <v-row>
        <v-col cols="12" md="5">
          <v-form
            @submit.prevent
            validate-on="submit"
            autocomplete="off"
            class="d-flex align-center h-100"
          >
            <v-row>
              <v-col cols="12" sm="5" md="12">
                <v-select :items="Object.values(Rank)" label="Rank" v-model="rank" />
              </v-col>
              <v-col cols="12" sm="7" md="12">
                <v-autocomplete
                  v-model="targetTaxon"
                  v-model:search="searchTerm"
                  item-title="canonicalName"
                  :items="autocompleteItems"
                  :loading="loading"
                  label="Taxon name"
                  return-object
                  auto-select-first
                  color="blue"
                  :no-data-text="
                    searchTerm.trim().length >= 3
                      ? 'No matching taxa found'
                      : 'At least 3 characters required'
                  "
                >
                  <template v-slot:item="{ props, item }">
                    <v-list-item
                      v-bind="props"
                      :title="item?.raw?.canonicalName"
                      :subtitle="item?.raw?.status"
                    >
                      <template v-slot:append>
                        <v-chip close>{{ item?.raw?.rank }}</v-chip>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:append-inner>
                    <v-chip v-if="targetTaxon">
                      {{ targetTaxon.rank }}
                    </v-chip>
                  </template>
                </v-autocomplete>
              </v-col>
            </v-row>
          </v-form>
        </v-col>
        <v-col md="7" class="d-flex justify-center align-center">
          <div v-if="taxonInfo != undefined">
            <v-skeleton-loader class="mx-auto" type="article, actions" v-if="loadingTaxon" />
            <div v-else>
              <v-card-title class="d-flex align-center">
                <b>{{ taxonInfo.canonicalName }}</b>
                <v-spacer />
                <v-chip color="blue" class="mr-3">{{ taxonInfo.rank }}</v-chip>
                <v-tooltip text="Go to original GBIF record" location="top">
                  <template v-slot:activator="{ props }">
                    <v-btn
                      icon
                      v-bind="props"
                      :disabled="taxonInfo == undefined"
                      :href="
                        taxonInfo ? `https://www.gbif.org/species/${taxonInfo.key}` : undefined
                      "
                      target="_blank"
                    >
                      <IconGBIF />
                    </v-btn>
                  </template>
                </v-tooltip>
              </v-card-title>
              <v-card-subtitle>
                {{ taxonInfo.authorship }}
              </v-card-subtitle>
              <v-card-text>
                <v-list>
                  <v-list-item>
                    <span v-for="{ name, id } in taxonInfo.path" :key="id" class="d-inline-block">
                      <a :href="`https://www.gbif.org/species/${id}`" target="_blank">{{ name }}</a>
                      <v-icon v-if="id !== taxonInfo.key" icon="mdi-chevron-right" />
                    </span>
                  </v-list-item>
                  <v-list-item>
                    Up to {{ countTotal(taxonInfo) }} nodes will be imported.
                  </v-list-item>
                </v-list>
              </v-card-text>
            </div>
          </div>
          <div v-else>
            <p>Pick a taxonomic group to be used as an anchor to import all its descendants.</p>
            <p>
              Parent taxons up to the kingdom will also be imported so that the full classification
              is available.
            </p>
            <p>All taxons below the root taxons will be imported recursively.</p>
          </div>
          <v-alert
            v-if="postError"
            type="error"
            title="Invalid taxon provided as anchor"
            :text="postError.message"
          />
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="orange" @click="$emit('close')">Cancel</v-btn>
      <v-btn :disabled="!taxonInfo" color="blue" @click="importAnchorTaxon(taxonInfo)">
        Validate and import
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type { Ref } from 'vue'
import { ref } from 'vue'

import axios, { AxiosError } from 'axios'

import { watch } from 'vue'

import IconGBIF from '@/components/icons/IconGBIF.vue'

import { VSkeletonLoader } from 'vuetify/labs/components'

enum Rank {
  Any = 'Any',
  Kingdom = 'Kingdom',
  Phylum = 'Phylum',
  Class = 'Class',
  Order = 'Order',
  Family = 'Family',
  Genus = 'Genus',
  Species = 'Species',
  Subspecies = 'Subspecies'
}
const rank = ref(Rank.Any)

type TaxonGBIF = {
  key: number
  canonicalName: string
  authorship: string
  rank: string
  status: string
  path?: any
}

const targetTaxon: Ref<TaxonGBIF | undefined> = ref(undefined)

const taxonInfo: Ref<TaxonGBIF | undefined> = ref(undefined)

type Item = {
  canonicalName: string
  rank: string
  status: string
}
const autocompleteItems: Ref<Item[]> = ref([])
const searchTerm = ref('')
const loading = ref(false)

const loadingTaxon = ref(false)

function endpointGBIF(suffix: string) {
  return `https://api.gbif.org/v1/species/${suffix}`
}

const emit = defineEmits<{ (event: 'close'): void }>()

const postError: Ref<AxiosError | undefined> = ref()

async function importAnchorTaxon(taxon: any) {
  try {
    await axios.post('/api/v1/taxonomy/anchors/', taxon)
    emit('close')
  } catch (error) {
    if (axios.isAxiosError(error)) {
      postError.value = error
    }
  }
}

watch(searchTerm, async (val: string) => {
  if (val.length >= 3) {
    loading.value = true
    let response = await axios.get(endpointGBIF('suggest'), {
      params: { q: val, rank: rank.value !== Rank.Any ? rank.value : undefined }
    })
    let data: Item[] = response.data
    autocompleteItems.value = data.filter(({ status }) => status !== 'DOUBTFUL')
    loading.value = false
  } else {
    autocompleteItems.value = []
  }
})

watch(targetTaxon, async (taxon) => {
  if (taxon) {
    loadingTaxon.value = true
    let response = await axios.get(endpointGBIF(taxon.key.toString()))
    let info = response.data
    info.path = Object.values(Rank)
      .filter((rank) => rank != Rank.Any && rank.toLowerCase() in info)
      .map((rank) => {
        let key = rank.toLowerCase()
        return {
          name: info[key],
          id: info[`${key}Key`]
        }
      })
    taxonInfo.value = info as TaxonGBIF
    loadingTaxon.value = false
  }
})

function countTotal(taxon: any) {
  return (
    Object.values(Rank).reduce((acc, rank) => {
      return acc + Number(rank != Rank.Any && rank.toLocaleLowerCase() in taxon)
    }, 0) + taxon.numDescendants
  )
}
</script>

<style scoped></style>
