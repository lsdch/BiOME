<template>
  <v-row>
    <v-col v-for="item in anchors" :key="item.ID" cols="12" sm="6" lg="4" xl="3">
      <v-card class="h-100 d-flex flex-column">
        <v-toolbar density="compact" color="white">
          <v-toolbar-title>
            {{ item.name }}
          </v-toolbar-title>
          <LinkIconGBIF :GBIF_ID="item.GBIF_ID" />
        </v-toolbar>
        <v-card-subtitle class="d-flex">
          {{ item.rank }} <v-spacer></v-spacer> {{ item.authorship }}
        </v-card-subtitle>
        <div class="mb-5 d-flex flex-column flex-grow-1 justify-end">
          <v-card-subtitle class="flex-end">
            Last modified
            {{ moment(item.modified ?? item.created).format('DD MMM y HH:MM') }}
          </v-card-subtitle>
        </div>
      </v-card>
    </v-col>
    <v-col v-for="item in activities" :key="item.GBIF_ID" cols="12" sm="6" lg="4" xl="3">
      <ImportTaxonCard v-bind="item" />
    </v-col>
    <v-scale-transition>
      <v-col v-if="pickerActive" cols="12">
        <RootTaxonPicker :value="undefined" :finished="false" @close="pickerActive = false" />
      </v-col>
      <v-col v-else cols="12" sm="6" lg="4" xl="3">
        <v-card class="h-100" variant="tonal">
          <v-card-text class="d-flex justify-center align-center h-100">
            <!-- variant="outlined" -->
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
</template>

<script setup lang="ts">
import axios from 'axios'
import { ref, watch } from 'vue'
import type { Ref } from 'vue'

import RootTaxonPicker from '@/components/taxonomy/loading/RootTaxonPicker.vue'
import LinkIconGBIF from '@/components/taxonomy/LinkIconGBIF.vue'

import moment from 'moment'

import type { ImportProcess } from '@/components/taxonomy/loading/ImportTaxonCard.vue'
import ImportTaxonCard from '@/components/taxonomy/loading/ImportTaxonCard.vue'

import { onMounted } from 'vue'

type Taxon  = {
  ID: string
  GBIF_ID: number
  name: string
  code: string
  status: string
  anchor: boolean
  authorship?: string
  rank: string
  modified?: string
  created?: string
}

const activities: Ref<ImportProcess[]> = ref([])
const anchors: Ref<Taxon[]> = ref([])

async function updateAnchors() {
  const response = await axios.get('/api/taxonomy/anchors')
  anchors.value = response.data
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
watch(activities, (activities) => {
  updateElapsedTime()
  if (activities.filter(({ done }) => done)) {
    updateAnchors()
  }
})

const pickerActive = ref(false)

const source = new EventSource('/api/taxonomy/update/progress')
source.addEventListener('download', (event) => {
  console.log(event)
  const json: Object = JSON.parse(event.data)
  console.log('JSON', json)
  activities.value = Object.values(json).map((item) =>
    Object.assign(item, { elapsed: moment(item.started).fromNow() })
  )
})
</script>

<style></style>
