<template>
  <v-card
    title="Event"
    :subtitle="
      hasID(site)
        ? 'From site'
        : site
          ? 'Register event at new site'
          : 'Waiting for site definition'
    "
    class="fill-height small-card-title"
    prepend-icon="mdi-calendar"
  >
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
      <EventFormDialog :site local @save="(ev: EventInput) => (event = ev)">
        <template #activator="{ props }">
          <v-btn
            prepend-icon="mdi-plus"
            variant="tonal"
            rounded="md"
            text="New event"
            v-bind="props"
          />
        </template>
      </EventFormDialog>
    </template>
    <v-card-text>
      <SiteEventPicker
        v-if="hasID(site)"
        @update:model-value="(ev) => (event = ev)"
        :site-code="site.code"
        clearable
        item-value="id"
      />
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Event, EventInput, Site, SiteInput } from '@/api'
import { watchEffect } from 'vue'
import SiteEventPicker from './SiteEventPicker.vue'
import { hasID } from '@/functions/db'
import EventFormDialog from '../events/EventFormDialog.vue'

const event = defineModel<Event | EventInput>()
const props = defineProps<{ site?: Site | SiteInput }>()
watchEffect(() => {
  if (!props.site) event.value = undefined
})
</script>

<style scoped lang="scss"></style>
