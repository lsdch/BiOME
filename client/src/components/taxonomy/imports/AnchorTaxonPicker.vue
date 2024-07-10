<template>
  <v-card>
    <v-card-title class="d-flex justify-space-between">
      <span>
        <v-icon icon="mdi-sitemap" start />
        Import anchor taxon
      </span>
    </v-card-title>
    <v-divider></v-divider>
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
                <v-select :items="rankOptions" label="Rank" v-model="rank" />
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
                    <v-switch
                      label="Import descendants"
                      v-model="importDescendants"
                      color="primary"
                      class="ml-3"
                      persistent-hint
                      :hint="
                        importDescendants
                          ? `Up to ${countTotal(taxonInfo)} nodes will be imported.`
                          : undefined
                      "
                    />
                  </v-list-item>
                </v-list>
              </v-card-text>
            </div>
          </div>
          <v-alert type="info" variant="text" title="Importing taxonomy" border v-else>
            <p>Pick a taxonomic group import from GBIF.</p>
            <p>
              Parent taxa up to the kingdom will also be imported so that the full classification is
              available.
            </p>
            <p>Descendant taxa may optionally be imported.</p>
          </v-alert>
        </v-col>
      </v-row>
      <v-alert
        v-if="postError"
        type="error"
        title="Invalid taxon provided as anchor"
        :text="postError.message"
      />
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="orange" @click="$emit('close')">Cancel</v-btn>
      <v-btn
        :disabled="taxonInfo == undefined"
        color="blue"
        @click="taxonInfo ? importAnchorTaxon(taxonInfo) : null"
      >
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

import { $TaxonRank, TaxonRank, TaxonomyGbifService } from '@/api'
import IconGBIF from '@/components/icons/IconGBIF.vue'

const importDescendants = ref(true)

type Rank = 'Any' | TaxonRank
const rank = ref<Rank>('Any')
const rankOptions: Rank[] = ['Any', ...$TaxonRank.enum]

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

async function importAnchorTaxon(taxon: TaxonGBIF) {
  try {
    await TaxonomyGbifService.importGbif({
      key: taxon.key,
      children: importDescendants.value
    })
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
      params: {
        q: val,
        rank: rank.value !== 'Any' ? rank.value : undefined
      }
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
    info.path = $TaxonRank.enum
      .filter((rank) => rank.toLowerCase() in info)
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
    $TaxonRank.enum.reduce((acc, rank) => {
      return acc + Number(rank.toLocaleLowerCase() in taxon)
    }, 0) + taxon.numDescendants
  )
}
</script>

<style scoped></style>
