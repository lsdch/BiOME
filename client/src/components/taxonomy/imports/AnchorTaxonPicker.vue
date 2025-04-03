<template>
  <v-card>
    <template #prepend>
      <v-icon icon="mdi-sitemap" start />
    </template>
    <template #title> Import GBIF clade </template>
    <v-divider />
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
                <AutocompleteGBIF v-model="targetTaxon" :rank="rank" />
              </v-col>
            </v-row>
          </v-form>
        </v-col>
        <v-col md="7" class="d-flex justify-center align-center">
          <div v-if="taxonInfo != undefined">
            <v-skeleton-loader class="mx-auto" type="article, actions" v-if="loadingTaxon" />
            <div v-else>
              <v-card flat>
                <template #title>
                  {{ taxonInfo.canonicalName }}
                </template>
                <template #subtitle>
                  {{ taxonInfo.authorship }}
                </template>
                <template #append>
                  <v-chip color="blue" class="mr-3">{{ taxonInfo.rank }}</v-chip>
                  <LinkIconGBIF variant="plain" :GBIF_ID="taxonInfo.key"></LinkIconGBIF>
                </template>
                <v-card-text>
                  <div class="my-3">
                    <span v-for="{ name, id } in taxonInfo.path" :key="id" class="d-inline-block">
                      <a :href="`https://www.gbif.org/species/${id}`" target="_blank">
                        {{ name }}
                      </a>
                      <v-icon v-if="id !== taxonInfo.key" icon="mdi-chevron-right" />
                    </span>
                  </div>
                  <v-switch
                    label="Import descendants"
                    v-model="importDescendants"
                    color="primary"
                    class="my-3"
                    persistent-hint
                    :hint="
                      importDescendants
                        ? `Up to ${countTotal(taxonInfo)} nodes will be imported.`
                        : undefined
                    "
                  />
                </v-card-text>
              </v-card>
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
        :text="postError.detail"
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
import { ref, watch } from 'vue'
import { $TaxonRank, ErrorModel, TaxonRank, TaxonomyGbifService } from '@/api'
import { endpointGBIF } from '.'
import LinkIconGBIF from '../LinkIconGBIF'
import AutocompleteGBIF from './AutocompleteGBIF.vue'

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
  path?: { name: string; id: number }[]
}

const targetTaxon = ref<TaxonGBIF>()

const taxonInfo = ref<TaxonGBIF>()

const loadingTaxon = ref(false)

const emit = defineEmits<{ (event: 'close'): void }>()

const postError = ref<ErrorModel>()

async function importAnchorTaxon(taxon: TaxonGBIF) {
  const { error } = await TaxonomyGbifService.importGbif({
    body: {
      key: taxon.key,
      children: importDescendants.value
    }
  })
  if (error) {
    postError.value = error
  }
  emit('close')
}

watch(targetTaxon, async (taxon) => {
  if (taxon) {
    loadingTaxon.value = true
    let response = await fetch(endpointGBIF(taxon.key.toString()))
    let info: TaxonGBIF & Record<string, unknown> = await response.json()
    info.path = $TaxonRank.enum
      .filter((rank) => rank.toLowerCase() in info)
      .map((rank) => {
        let key = rank.toLowerCase()
        return {
          name: info[key] as string,
          id: info[`${key}Key`] as number
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
