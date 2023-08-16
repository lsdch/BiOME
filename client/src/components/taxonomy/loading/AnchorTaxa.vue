<template>
  <div>
    <div v-if="items.length">
      <v-card v-for="item in items" :key="item"> {{ item }} </v-card>
    </div>
    <div v-if="inProgress.length">
      <v-row>
        <v-col v-for="item in inProgress" :key="item.GBIF_ID" cols="12" sm="6" lg="4">
          <ImportTaxonCard v-bind="item" />
        </v-col>
        <v-expand-transition>
          <v-col v-if="pickerActive" cols="12">
            <RootTaxonPicker :value="undefined" :finished="false"></RootTaxonPicker>
          </v-col>
          <v-col v-else cols="12" sm="6" lg="4">
            <v-card class="h-100" variant="tonal">
              <v-card-text class="d-flex justify-center align-center h-100">
                <!-- variant="outlined" -->
                <v-btn
                  class="ma-2"
                  color="blue"
                  icon="mdi-plus"
                  size="x-large"
                  @click="pickerActive = true"
                ></v-btn>
              </v-card-text>
            </v-card>
          </v-col>
        </v-expand-transition>
      </v-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import { ref } from 'vue'
import moment from 'moment'
import RootTaxonPicker from './RootTaxonPicker.vue'
import type { ImportProcess } from './ImportTaxonCard.vue'
import ImportTaxonCard from './ImportTaxonCard.vue'
import { watch } from 'vue'

const pickerActive = ref(false)

const props = defineProps<{ inProgress: Array<ImportProcess> }>()

// const elapsed = ref(new Map<number, string>())

const response = await axios.get('/api/taxonomy/anchors')
const items = ref(response.data)

function updateElapsedTime() {
  props.inProgress.map((progress) => {
    progress.elapsed = moment(progress.started).fromNow()
    // elapsed.value.set(GBIF_ID, moment.duration(started).humanize())
  })
}
updateElapsedTime()
setInterval(updateElapsedTime, 5000)
watch(items, updateElapsedTime)
</script>

<style scoped less></style>
