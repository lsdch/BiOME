<template>
  <v-card>
    <v-tabs v-model="tab">
      <v-tab value="sites" prepend-icon="mdi-map-marker">
        Sites
        <v-chip class="mx-1" :text="dataset.sites?.length.toString()" density="compact" />
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab">
      <v-tabs-window-item value="sites" key="sites">
        <CRUDTable :headers v-model="dataset.sites" entityName="Site" density="compact">
          <template #[`item.name`]="{ item }: { item: Site }">
            <RouterLink :to="{ name: 'site-item', params: { code: item.code } }">
              {{ item.name }}
            </RouterLink>
          </template>
        </CRUDTable>
      </v-tabs-window-item>
    </v-tabs-window>
  </v-card>
</template>

<script setup lang="ts">
import { Dataset, Site } from '@/api'
import CRUDTable from '@/components/toolkit/tables/CRUDTable.vue'
import { ref } from 'vue'

defineProps<{ dataset: Dataset }>()

type Tab = 'sites'
const tab = ref<Tab>()

const headers: CRUDTableHeader[] = [
  { key: 'name', title: 'Name' },
  { key: 'locality', title: 'Locality' },
  {
    key: 'country',
    title: 'Country',
    value({ country: { name } }) {
      return name
    }
  }
]
</script>

<style scoped lang="scss"></style>
