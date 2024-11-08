<template>
  <v-data-table :items :headers>
    <!-- HEADERS -->
    <template #header.samplings="props">
      <IconTableHeader icon="mdi-package-down" v-bind="props" />
    </template>
    <template #header.abiotic_measurements="props">
      <IconTableHeader icon="mdi-gauge" v-bind="props" />
    </template>
    <template #header.spotting="props">
      <IconTableHeader icon="mdi-binoculars" v-bind="props" />
    </template>

    <!-- ITEMS -->
    <template #item.performed_on="{ value }: { value: DateWithPrecision }">
      {{ formatDateWithPrecision(value) }}
    </template>

    <template #item.samplings="{ value, item }">
      <v-chip
        v-if="value"
        color="success"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(item, 'sampling')"
      />
    </template>
    <template #item.abiotic_measurements="{ value, item }">
      <v-chip
        v-if="value"
        color="primary"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(item, 'abiotic')"
      />
    </template>
    <template #item.spotting="{ value, item }">
      <v-chip
        v-if="value !== undefined"
        color="warning"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(item, 'spotting')"
      />
    </template>
  </v-data-table>
  <EventCardDialog :event="focusedEvent" v-model:open="showEventDetails" v-model:tab="eventTab" />
</template>

<script setup lang="ts">
import { DateWithPrecision, type Event } from '@/api'
import { ref } from 'vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import { formatDateWithPrecision } from '../toolkit/utils'
import EventCardDialog, { EventAction } from './EventCardDialog.vue'

const focusedEvent = ref<Event>()
const eventTab = ref<EventAction>()
const showEventDetails = ref(false)

function toggleFocus(event: Event, tab: EventAction) {
  eventTab.value = tab
  focusedEvent.value = event
  showEventDetails.value = true
}

defineProps<{ items: Event[] }>()

const headers: CRUDTableHeader<Event>[] = [
  {
    key: 'performed_on',
    title: 'Date',
    headerProps: { class: 'v-data-table-column--align-center' },
    cellProps: { class: 'v-data-table-column--align-end' }
  },
  {
    title: 'Actions',
    align: 'center',
    children: [
      {
        key: 'samplings',
        title: 'Samplings',
        align: 'center',
        width: 0,
        value(item: Event, fallback) {
          return item.samplings?.length
        }
      },
      {
        key: 'abiotic_measurements',
        title: 'Abiotic',
        align: 'center',
        width: 0,
        value(item: Event, fallback) {
          return item.abiotic_measurements?.length
        }
      },
      {
        key: 'spotting',
        title: 'Spotting',
        align: 'center',
        width: 0,
        value(item: Event, fallback) {
          return item.spotting?.target_taxa?.length
        }
      }
    ]
  }
]
</script>

<style scoped></style>
