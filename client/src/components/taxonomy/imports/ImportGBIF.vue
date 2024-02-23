<template>
  <v-container fluid>
    <v-row>
      <v-col v-for="item in anchors" :key="item.id" cols="12" sm="6" lg="4" xl="3">
        <AnchorTaxonCard v-bind="item" />
      </v-col>
      <v-col
        v-for="item in activities.filter(({ done, error }) => !done || error)"
        :key="item.GBIF_ID"
        cols="12"
        sm="6"
        lg="4"
        xl="3"
      >
        <ImportTaxonCard v-bind="item" />
      </v-col>
      <v-scale-transition>
        <v-col v-if="pickerActive" cols="12" xl="10" offset-xl="1">
          <RootTaxonPicker :value="undefined" @close="pickerActive = false" />
        </v-col>

        <v-col v-else cols="12" sm="6" lg="4" xl="3">
          <v-card class="h-100" variant="tonal">
            <v-card-text class="d-flex justify-center align-center h-100">
              <v-btn
                class="ma-2"
                color="blue"
                icon="mdi-plus"
                size="x-large"
                @click="pickerActive = true"
              />
            </v-card-text>
          </v-card>
        </v-col>
      </v-scale-transition>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import moment from 'moment'

import type { Ref } from 'vue'
import { onMounted, ref } from 'vue'

import { TaxonDB, TaxonomyService } from '@/api'
import AnchorTaxonCard from '@/components/taxonomy/imports/AnchorTaxonCard.vue'
import RootTaxonPicker from '@/components/taxonomy/imports/AnchorTaxonPicker.vue'
import type { ImportProcess } from '@/components/taxonomy/imports/ImportTaxonCard.vue'
import ImportTaxonCard from '@/components/taxonomy/imports/ImportTaxonCard.vue'

const activities: Ref<ImportProcess[]> = ref([])
const anchors: Ref<TaxonDB[]> = ref([])

async function updateAnchors() {
  const response = await TaxonomyService.taxonAnchors()
  anchors.value = response
}

onMounted(updateAnchors)

function updateElapsedTime() {
  activities.value.map((progress) => {
    progress.elapsed = moment(progress.started).fromNow()
    // elapsed.value.set(GBIF_ID, moment.duration(started).humanize())
  })
}
updateElapsedTime()
setInterval(updateElapsedTime, 5000)

const pickerActive = ref(false)

const source = new EventSource('/api/v1/taxonomy/import')
source.addEventListener('progress', (event) => {
  console.log(event)
  const json: Object = JSON.parse(event.data)
  console.log('Progress: ', json)
  activities.value = Object.values(json).map((item) =>
    Object.assign(item, { elapsed: moment(item.started).fromNow() })
  )
  updateElapsedTime()
  if (activities.value.filter(({ done }) => done).length > 0) {
    updateAnchors()
  }
})
</script>

<style scoped></style>
