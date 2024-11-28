<template>
  <CRUDTable
    :items="site.events"
    :headers
    entityName="Events"
    :delete="({ id }) => EventsService.deleteEvent({ path: { id } })"
  >
    <!-- HEADERS -->
    <template #header.performed_on>
      <div class="d-flex justify-space-between align-center mr-3">
        <v-btn
          color="primary"
          variant="tonal"
          text="Add event"
          prepend-icon="mdi-plus"
          class="mb-3"
          @click="createEventDialog = true"
        />
        Date
      </div>
    </template>
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
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <!-- Event form -->
      <EventFormDialog
        :model-value="dialog || createEventDialog"
        :edit="createEventDialog ? undefined : editItem"
        :site
        @close="onClose(), (createEventDialog = false)"
        :onSuccess
      />
    </template>
  </CRUDTable>

  <!-- Event details dialog -->
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
        :disabled="focusedEvent.index >= site.events.length - 1"
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
import { DateWithPrecision, EventsService, Site, type Event } from '@/api'
import { ref } from 'vue'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import { formatDateWithPrecision } from '../toolkit/utils'
import EventCardDialog, { EventAction } from './EventCardDialog.vue'
import EventFormDialog from './EventFormDialog.vue'

const focusedEvent = ref<{ index: number; item: Event }>()
const eventTab = ref<EventAction>()
const showEventDetails = ref(false)

const props = defineProps<{ site: Site }>()

const createEventDialog = ref(false)

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
      },
      {
        title: 'Edit',
        key: 'actions',
        align: 'center',
        width: 0,
        cellProps: {
          class: 'text-no-wrap'
        }
      }
    ]
  }
]

function toggleFocus(index: number, event: Event, tab: EventAction) {
  eventTab.value = tab
  focusItem(index)
  showEventDetails.value = true
}

function focusItem(index: number) {
  focusedEvent.value = { index, item: props.site.events[index] }
}

function focusNext() {
  if (!focusedEvent.value || focusedEvent.value.index >= props.site.events.length - 1) return
  focusItem(focusedEvent.value.index + 1)
}

function focusPrev() {
  if (!focusedEvent.value || focusedEvent.value.index <= 0) return
  focusItem(focusedEvent.value.index - 1)
}
</script>

<style scoped></style>
