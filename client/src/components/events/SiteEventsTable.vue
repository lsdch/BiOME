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

    <template #item.samplings="{ value, item, index }">
      <v-chip
        v-if="value"
        color="success"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(index, item, 'sampling')"
      />
    </template>
    <template #item.abiotic_measurements="{ value, item, index }">
      <v-chip
        v-if="value"
        color="primary"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(index, item, 'abiotic')"
      />
    </template>
    <template #item.spotting="{ value, item, index }">
      <v-chip
        v-if="value !== undefined"
        color="warning"
        :text="value.toString()"
        size="small"
        density="compact"
        @click="toggleFocus(index, item, 'spotting')"
      />
    </template>
  </v-data-table>
  <EventCardDialog
    :event="focusedEvent?.item"
    v-model:open="showEventDetails"
    v-model:tab="eventTab"
    @next="focusedEvent && focusItem(focusedEvent.index + 1)"
    @prev="focusedEvent && focusItem(focusedEvent.index - 1)"
  >
    <template v-if="focusedEvent" #title>
      <v-btn
        color="primary"
        icon="mdi-arrow-left"
        @click="focusNext"
        :disabled="focusedEvent.index >= items.length - 1"
      />
      {{ focusedEvent.item.site.name }} |
      {{ formatDateWithPrecision(focusedEvent.item.performed_on) }}
      <v-btn
        color="primary"
        icon="mdi-arrow-right"
        @click="focusPrev"
        :disabled="focusedEvent.index <= 0"
      />
    </template>
  </EventCardDialog>
</template>

<script setup lang="ts">
import { DateWithPrecision, type Event } from '@/api'
import { ref } from 'vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import { formatDateWithPrecision } from '../toolkit/utils'
import EventCardDialog, { EventAction } from './EventCardDialog.vue'

const focusedEvent = ref<{ index: number; item: Event }>()
const eventTab = ref<EventAction>()
const showEventDetails = ref(false)

function toggleFocus(index: number, event: Event, tab: EventAction) {
  eventTab.value = tab
  focusItem(index)
  showEventDetails.value = true
}

function focusItem(index: number) {
  focusedEvent.value = { index, item: props.items[index] }
}

function focusNext() {
  if (!focusedEvent.value || focusedEvent.value.index >= props.items.length - 1) return
  focusItem(focusedEvent.value.index + 1)
}

function focusPrev() {
  if (!focusedEvent.value || focusedEvent.value.index <= 0) return
  focusItem(focusedEvent.value.index - 1)
}

const props = defineProps<{ items: Event[] }>()

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
