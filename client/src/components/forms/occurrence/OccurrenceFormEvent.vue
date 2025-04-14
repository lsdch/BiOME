<template>
  <v-card class="small-card-title" prepend-icon="mdi-calendar">
    <template #title>
      <template v-if="event && !showEdit">
        <b class="font-monospace">{{
          DateWithPrecision.format(
            event.performed_on.date instanceof Date
              ? (event.performed_on as DateWithPrecision)
              : DateWithPrecision.fromInput(event.performed_on as DateWithPrecisionInput)
          )
        }}</b>
      </template>
      <v-card-title v-else>Event</v-card-title>
    </template>
    <template #subtitle>
      <v-chip
        v-if="hasID(event)"
        label
        text="From database"
        class="font-monospace"
        color="purple"
        prepend-icon="mdi-link"
        size="x-small"
        variant="flat"
      />
      <v-chip
        v-else-if="!!event"
        label
        text="New event"
        class="font-monospace"
        color="success"
        prepend-icon="mdi-plus"
        size="x-small"
        variant="flat"
      />

      <v-card-subtitle v-else>
        {{
          hasID(site)
            ? 'Pick or register event at site'
            : !!site
              ? 'Register event at new site'
              : 'Waiting for site definition'
        }}
      </v-card-subtitle>
    </template>
    <template #append v-if="site">
      <EventFormDialog
        v-show="!hasID(event) || showEdit"
        v-model:dialog="dialog"
        title="Create event"
        btn-text="Save"
        :site
        @submit="updateEvent"
      >
        <template #activator="{ props }">
          <v-btn
            v-show="!event || showEdit"
            variant="tonal"
            rounded="md"
            v-bind="{
              ...props,
              ...(!event || hasID(event)
                ? {
                    text: 'New event',
                    prependIcon: 'mdi-plus'
                  }
                : {
                    text: 'Edit event',
                    prependIcon: 'mdi-pencil'
                  })
            }"
          />
        </template>
      </EventFormDialog>
      <v-btn
        v-show="!!event && !showEdit"
        icon="mdi-pencil"
        variant="tonal"
        size="small"
        @click="toggleEdit(true)"
      />
    </template>
    <v-card-text v-if="!!site">
      <v-expand-transition>
        <!-- :model-value="hasID(event) ? event : undefined" -->
        <SiteEventPicker
          v-if="hasID(site) && (!event || showEdit)"
          @update:model-value="updateEvent"
          :site-code="site.code"
          clearable
          return-object
        >
          <template #no-data>
            <v-alert>Selected site has no registered events yet</v-alert>
          </template>
        </SiteEventPicker>
      </v-expand-transition>
      <v-expand-transition>
        <v-card v-if="event && !showEdit" class="small-card-title" flat>
          <v-list-item prepend-icon="mdi-account-multiple">
            <template #append>
              <span class="text-caption text-muted">Participants</span>
            </template>
            <PersonChip
              v-if="event.performed_by?.length"
              v-for="person in event.performed_by"
              :person
              size="small"
            />
            <v-chip
              v-if="event.performed_by_groups?.length"
              v-for="org in event.performed_by_groups"
            >
              {{ org.code }}
            </v-chip>
            <v-list-item-title
              v-if="!event.performed_by?.length && !event.performed_by_groups?.length"
              class="text-muted font-italic"
            >
              No known participants
            </v-list-item-title>
          </v-list-item>
        </v-card>
      </v-expand-transition>
    </v-card-text>
    <template #actions v-if="showEdit && !!event">
      <v-spacer />
      <v-btn text="Cancel" @click="toggleEdit(false)" />
    </template>
  </v-card>
</template>

<script setup lang="ts">
import { DateWithPrecision, DateWithPrecisionInput, Event, Site } from '@/api'
import EventFormDialog from '@/components/forms/EventFormDialog.vue'
import SiteEventPicker from '@/components/forms/occurrence/SiteEventPicker.vue'
import PersonChip from '@/components/people/PersonChip'
import { hasID } from '@/functions/db'
import { EventModel } from '@/models'
import { SiteFormModel } from '@/models/site'
import { useToggle } from '@vueuse/core'
import { ref, watch } from 'vue'

const event = defineModel<Event | EventModel.EventModel>()
const props = defineProps<{ site?: Site | SiteFormModel }>()

const [showEdit, toggleEdit] = useToggle(false)
const dialog = ref(false)

watch(
  () => props.site,
  () => updateEvent(undefined)
)

function updateEvent(ev: Event | EventModel.EventModel | undefined) {
  event.value = ev
  dialog.value = false
  toggleEdit(!ev)
}
</script>

<style scoped lang="scss"></style>
