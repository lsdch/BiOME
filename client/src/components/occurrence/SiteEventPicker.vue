<template>
  <v-autocomplete v-model="model" :loading="isFetching" :items v-bind="$attrs" return-object>
    <template #selection="{ item: { raw: event } }">
      <v-list-item :title="DateWithPrecision.format(event.performed_on)" class="px-0">
        <template #append>
          <div class="d-flex ga-1 ml-3">
            <v-chip
              v-if="event.samplings?.length"
              :text="event.samplings.length.toString()"
              size="small"
              prepend-icon="mdi-package-down"
            />
            <v-chip
              v-if="event.abiotic_measurements?.length"
              :text="event.abiotic_measurements.length.toString()"
              size="small"
              prepend-icon="mdi-gauge"
            />
            <v-chip
              v-if="event.spottings?.length"
              :text="event.spottings.length.toString()"
              size="small"
              prepend-icon="mdi-binoculars"
            />
          </div>
        </template>
      </v-list-item>
    </template>

    <template #item="{ item: { raw: event }, props }">
      <v-list-item v-bind="props" :title="DateWithPrecision.format(event.performed_on)">
        <template #append>
          <div class="d-flex ga-1">
            <v-chip
              v-if="event.samplings?.length"
              :text="event.samplings.length.toString()"
              size="small"
              prepend-icon="mdi-package-down"
            />
            <v-chip
              v-if="event.abiotic_measurements?.length"
              :text="event.abiotic_measurements.length.toString()"
              size="small"
              prepend-icon="mdi-gauge"
            />
            <v-chip
              v-if="event.spottings?.length"
              :text="event.spottings.length.toString()"
              size="small"
              prepend-icon="mdi-binoculars"
            />
          </div>
        </template>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { DateWithPrecision, Event } from '@/api'
import { listSiteEventsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const model = defineModel<Event>()

const props = defineProps<{ siteCode: string }>()

const { data: items, isFetching } = useQuery(
  computed(() => ({
    enabled: !!props.siteCode,
    ...listSiteEventsOptions({
      path: { code: props.siteCode }
    })
  }))
)
</script>

<style scoped lang="scss"></style>
