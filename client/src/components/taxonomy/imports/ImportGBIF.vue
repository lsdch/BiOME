<template>
  <v-progress-linear indeterminate v-if="loading" />
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

import { ref, watch } from 'vue'

import { listAnchorsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import AnchorTaxonCard from '@/components/taxonomy/imports/AnchorTaxonCard.vue'
import RootTaxonPicker from '@/components/taxonomy/imports/AnchorTaxonPicker.vue'
import type { ImportProcess } from '@/components/taxonomy/imports/ImportTaxonCard.vue'
import ImportTaxonCard from '@/components/taxonomy/imports/ImportTaxonCard.vue'
import { useQuery } from '@tanstack/vue-query'
import { useEventSource } from '@vueuse/core'
import { onBeforeRouteLeave } from 'vue-router'

const activities = ref<ImportProcess[]>([])

const {
  data: anchors,
  error,
  isFetching: loading,
  refetch
} = useQuery({
  ...listAnchorsOptions(),
  initialData: []
})

function updateElapsedTime() {
  activities.value.map((progress) => {
    progress.elapsed = moment(progress.started).fromNow()
  })
}
updateElapsedTime()
setInterval(updateElapsedTime, 5000)

const pickerActive = ref(false)

const { data, status, close } = useEventSource('/api/v1/import/taxonomy/monitor', [
  'state'
] as const)
watch(data, (data) => {
  if (!data) return
  const json: Object = JSON.parse(data)
  console.log('Monitor: ', status, json)
  activities.value = Object.values(json).map((item) =>
    Object.assign(item, { elapsed: moment(item.started).fromNow() })
  )
  updateElapsedTime()
  if (activities.value.filter(({ done, error }) => done && !error).length > 0) {
    refetch()
  }
})

onBeforeRouteLeave(() => close())
</script>

<style scoped></style>
