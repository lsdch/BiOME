<template>
  <div>
    <v-form>
      <v-row>
        <v-col md="4">
          <v-text-field label="Name" v-model="filters.name" density="compact"></v-text-field>
        </v-col>
        <v-col md="4">
          <v-select
            label="Rank"
            v-model="filters.rank"
            :items="Object.values(TaxonRank)"
            density="compact"
            clearable
          />
        </v-col>
        <v-col md="4">
          <v-select
            label="Status"
            v-model="filters.status"
            :items="Object.values(TaxonStatus)"
            density="compact"
            clearable
          />
        </v-col>
      </v-row>
    </v-form>
    <v-data-table :items="items" :headers="headers" density="compact" :loading="loading">
      <template v-slot:[`item.code`]="{ item }">
        <code>{{ item.code }}</code>
      </template>
      <template v-slot:[`item.status`]="{ item }">
        <status-icon :status="item.status" size="small" />
        <LinkIconGBIF v-if="item.gbif_ID" :GBIF_ID="item.gbif_ID" variant="text" size="x-small" />
      </template>
    </v-data-table>
  </div>
</template>

<script setup lang="ts">
import type { Ref } from 'vue'
import { ref, onMounted } from 'vue'

import { TaxonRank, TaxonStatus, TaxonWithRelatives } from '@/api'
import { TaxonomyService } from '@/api'
import { computed } from 'vue'
import StatusIcon from './StatusIcon.vue'
import LinkIconGBIF from './LinkIconGBIF.vue'
import { VDataTable } from 'vuetify/components'

const taxa: Ref<TaxonWithRelatives[]> = ref([])
const loading = ref(true)

const headers: ReadonlyHeaders = [
  { title: 'Name', key: 'name' },
  { title: 'Code', key: 'code' },
  { title: 'Rank', key: 'rank' },
  { title: 'Status', key: 'status', width: 0, align: 'center' }
]

const filters = ref({
  name: undefined,
  rank: undefined,
  status: TaxonStatus.Accepted
})

const items = computed(() => {
  return taxa.value.filter(({ name, rank, status }) => {
    return (
      name.includes(filters.value.name ?? '') &&
      (filters.value.rank ? rank === filters.value.rank : true) &&
      (filters.value.status ? status === filters.value.status : true)
    )
  })
})

onMounted(async () => {
  taxa.value = await TaxonomyService.taxonomyList()
  loading.value = false
})
</script>

<style scoped></style>
