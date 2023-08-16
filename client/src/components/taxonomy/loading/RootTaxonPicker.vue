<template>
  <v-card>
    <v-card-title class="d-flex justify-space-between">
      <span>
        <v-icon icon="mdi-sitemap" start :color="finished ? 'green' : ''" />
        Import anchor taxon
      </span>
      <!-- <v-spacer></v-spacer> -->
      <!-- <v-btn icon variant="flat">
        <v-icon icon="mdi-close"></v-icon>
      </v-btn> -->
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
              <v-col cols="12" sm="5" md="12" xl="5">
                <v-select :items="ranks" label="Rank" v-model="rank"></v-select>
              </v-col>
              <v-col cols="12" sm="7" md="12" xl="7">
                <v-autocomplete
                  v-model="targetTaxon"
                  v-model:search="searchTerm"
                  item-title="canonicalName"
                  :items="items_GBIF"
                  :loading="loading"
                  label="Taxon name"
                  return-object
                  auto-select-first
                  color="blue"
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
        <v-col md="7" class="justify-center">
          <div v-if="taxonInfo">
            <v-skeleton-loader class="mx-auto" type="article, actions" v-if="loadingTaxon" />
            <div v-else>
              <v-card-title class="d-flex align-center">
                <b>{{ taxonInfo.canonicalName }}</b>
                <v-spacer></v-spacer>
                <v-chip color="blue" class="mr-3">{{ taxonInfo.rank }}</v-chip>
                <v-tooltip text="Go to original GBIF record" location="top">
                  <template v-slot:activator="{ props }">
                    <v-btn
                      v-bind="props"
                      icon="mdi-plus"
                      :disabled="taxonInfo == undefined"
                      :href="
                        taxonInfo ? `https://www.gbif.org/species/${taxonInfo.key}` : undefined
                      "
                      target="_blank"
                    >
                      <IconGBIF></IconGBIF>
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
                      <v-icon v-if="id !== taxonInfo.key"> mdi-chevron-right </v-icon>
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

import axios from 'axios'

import { watch } from 'vue'

import IconGBIF from '@/components/icons/IconGBIF.vue'

import { VSkeletonLoader } from 'vuetify/labs/components'

defineProps<{
  finished: boolean
}>()

const ranks = [
  'Any',
  'Kingdom',
  'Phylum',
  'Class',
  'Order',
  'Family',
  'Genus',
  'Species',
  'Subspecies'
]
const rank = ref('Any')

const targetTaxon: any = ref(undefined)

const taxonInfo: any = ref(undefined)

type Item = {
  canonicalName: string
  rank: string
  status: string
}
const items_GBIF: Ref<Item[]> = ref([])
const searchTerm = ref('')
const loading = ref(false)

const loadingTaxon = ref(false)

function endpointGBIF(suffix: string) {
  return `https://api.gbif.org/v1/species/${suffix}`
}

async function importAnchorTaxon(taxon: any) {
  let response = await axios.post('/api/taxonomy/update/', taxon)
  console.log(response.data)
}

watch(searchTerm, async (val: string) => {
  // targetTaxon.value = undefined
  if (val.length >= 3) {
    loading.value = true
    let response = await axios.get(endpointGBIF('suggest'), {
      params: { q: val, rank: rank.value !== 'Any' ? rank.value : undefined }
    })
    console.log(response)
    let data: Item[] = response.data
    items_GBIF.value = data.filter(({ status }) => status !== 'DOUBTFUL')
    console.log(items_GBIF.value)
    loading.value = false
  } else {
    items_GBIF.value = []
  }
})

watch(targetTaxon, async (taxon) => {
  if (taxon) {
    loadingTaxon.value = true
    let response = await axios.get(endpointGBIF(taxon.key))
    console.log(response)
    taxonInfo.value = response.data
    taxonInfo.value.path = ranks
      .filter((rank) => rank != 'Any' && rank.toLowerCase() in taxonInfo.value)
      .map((rank) => {
        let key = rank.toLowerCase()
        return {
          name: taxonInfo.value[key],
          id: taxonInfo.value[`${key}Key`]
        }
      })
    loadingTaxon.value = false
  }
})

function countTotal(taxon: any) {
  return (
    ranks.reduce((acc, rank) => {
      return acc + Number(rank != 'Any' && rank.toLocaleLowerCase() in taxon)
    }, 0) + taxon.numDescendants
  )
}
</script>

<style scoped></style>
