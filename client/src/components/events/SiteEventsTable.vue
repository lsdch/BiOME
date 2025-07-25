<template>
  <CRUDTable
    ref="events-table"
    :items="site.events"
    :headers
    entityName="Events"
    :delete="{
      mutation: deleteEventMutation,
      params: ({ id }: Event) => ({ path: { id } })
    }"
    @item-created="(item, index) => toggleFocus(index, item, 'sampling')"
    @click:row="
      (_: PointerEvent, { index, item }: RowClick) => toggleFocus(index, item, 'sampling')
    "
  >
    <!-- HEADERS -->
    <template #header.corner>
      <div class="d-flex justify-space-between align-center mr-3">
        <v-btn
          color="primary"
          variant="tonal"
          text="Add event"
          prepend-icon="mdi-plus"
          class="mb-3"
          @click="eventsTable?.actions.create()"
        />
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
    <template #item.performed_on="{ item }">
      {{ DateWithPrecision.format(item.performed_on) }}
    </template>

    <template #item.samples="{ value }: { value: BioMaterial[] }">
      <v-chip
        v-for="s in value"
        class="ma-1"
        :text="s.identification.taxon.name"
        label
        :to="{ name: 'biomat-item', params: { code: s.code } }"
        size="small"
      />
    </template>

    <template #item.samplings="{ value, item, index }">
      <v-chip
        v-if="value"
        color="success"
        :text="value.toString()"
        size="small"
        density="compact"
        @click.stop="toggleFocus(index, item, 'sampling')"
      />
    </template>
    <template #item.abiotic_measurements="{ value, item, index }">
      <v-chip
        v-if="value"
        color="primary"
        :text="value.toString()"
        size="small"
        density="compact"
        @click.stop="toggleFocus(index, item, 'abiotic')"
      />
    </template>
    <template #item.spotting="{ value, item, index }">
      <v-chip
        v-if="value !== undefined"
        color="warning"
        :text="value.toString()"
        size="small"
        density="compact"
        @click.stop="toggleFocus(index, item, 'spotting')"
      />
    </template>

    <!-- Event form -->
    <template #form="{ dialog, mode, onClose, onSuccess, editItem }">
      <EventMutationFormDialogMutation
        :dialog
        @update:dialog="(v) => !v && onClose()"
        :model-value="editItem"
        :site
        @close="onClose()"
        @success="onSuccess"
      />
    </template>
  </CRUDTable>

  <!-- Event details dialog -->
  <EventCardDialog
    v-if="focusedEvent"
    v-model="focusedEvent.item"
    v-model:open="showEventDetails"
    v-model:tab="eventTab"
    @next="focusedEvent && focusItem(focusedEvent.index + 1)"
    @prev="focusedEvent && focusItem(focusedEvent.index - 1)"
  >
    <template #subtitle>
      <v-btn
        color="primary"
        icon="mdi-arrow-left"
        @click="focusNext"
        size="small"
        variant="text"
        density="compact"
        :disabled="focusedEvent.index >= (site.events?.length ?? 0) - 1"
      />
      {{ DateWithPrecision.format(focusedEvent.item.performed_on) }}
      <!-- {{ focusedEvent.item.site.name }} |
      {{ DateWithPrecision.format(focusedEvent.item.performed_on) }} -->
      <v-btn
        color="primary"
        icon="mdi-arrow-right"
        @click="focusPrev"
        size="small"
        variant="text"
        density="compact"
        :disabled="focusedEvent.index <= 0"
      />
    </template>
  </EventCardDialog>
</template>

<script setup lang="ts">
import { BioMaterial, Site, type Event } from '@/api'
import { DateWithPrecision } from '@/api/adapters'
import { deleteEventMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { ref, useTemplateRef } from 'vue'
import { ComponentExposed } from 'vue-component-type-helpers'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import EventCardDialog, { EventAction } from './EventCardDialog.vue'
import EventMutationFormDialogMutation from '../forms/EventFormDialogMutation.vue'

type RowClick = {
  index: number
  item: Event
}

const focusedEvent = ref<{ index: number; item: Event }>()
const eventTab = ref<EventAction>()
const showEventDetails = ref(false)

const eventsTable = useTemplateRef<ComponentExposed<typeof CRUDTable>>('events-table')

const props = defineProps<{ site: Site }>()

const headers: CRUDTableHeader<Event>[] = [
  {
    title: 'Events',
    key: 'corner',
    children: [
      {
        key: 'performed_on',
        title: 'Date',
        value: 'performed_on.date',
        align: 'end'
      },
      {
        key: 'samples',
        title: 'Samples',
        align: 'end',
        value(item: Event) {
          return item.samplings?.flatMap(({ samples }) => samples ?? []) ?? []
        }
      }
    ]
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
          return item.spottings?.length
        }
      },
      {
        key: 'actions',
        align: 'center',
        sortable: false,
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
  focusedEvent.value = { index, item: props.site.events![index] }
}

function focusNext() {
  if (!focusedEvent.value || focusedEvent.value.index >= (props.site.events?.length ?? 0) - 1)
    return
  focusItem(focusedEvent.value.index + 1)
}

function focusPrev() {
  if (!focusedEvent.value || focusedEvent.value.index <= 0) return
  focusItem(focusedEvent.value.index - 1)
}
</script>

<style scoped></style>
