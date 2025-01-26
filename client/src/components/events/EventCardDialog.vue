<template>
  <CardDialog v-model="open" :fullscreen="smAndDown" :title="event?.site.name">
    <!-- :subtitle="event ? DateWithPrecision.format(event!.performed_on) : undefined" -->
    <!-- <template #title>
      <slot name="title" />
    </template> -->
    <template #subtitle>
      <slot name="subtitle" />
    </template>

    <SamplingFormDialog
      v-if="event !== undefined"
      v-model="samplingDialog"
      :event
      :fullscreen="smAndDown"
      :edit="editingSampling"
      @created="(sampling) => event?.samplings.unshift(sampling)"
      @updated="
        (sampling) =>
          event?.samplings.map((s) => {
            return s.id === sampling.id ? sampling : s
          })
      "
    />

    <v-tabs v-model="tab" color="primary" center-active class="overflow-visible">
      <v-tab value="sampling" prepend-icon="mdi-package-down">
        <span v-if="!mobile || tab === 'sampling'"> Samplings </span>
        <v-badge color="primary" inline :content="event?.samplings.length" />
      </v-tab>
      <v-tab value="abiotic" prepend-icon="mdi-gauge">
        <span v-if="!mobile || tab === 'abiotic'"> Abiotic </span>
        <v-badge color="primary" inline :content="event?.abiotic_measurements.length" />
      </v-tab>
      <v-tab value="spotting" prepend-icon="mdi-binoculars">
        <span v-if="!mobile || tab === 'spotting'"> Spotting </span>
        <v-badge color="primary" inline :content="event?.spottings?.length ?? 0" />
      </v-tab>
    </v-tabs>
    <v-tabs-window v-model="tab" class="overflow-y-auto event-action-text">
      <v-tabs-window-item value="sampling">
        <v-container fluid>
          <v-row align-content="stretch">
            <v-col v-if="!event?.samplings?.length">
              <v-alert> No samplings reported </v-alert>
            </v-col>
            <v-col v-for="(sampling, index) in event?.samplings" cols="12" md="6">
              <SamplingCard
                :sampling
                :corner-tag="`#${index + 1} / ${event?.samplings.length}`"
                class="h-100"
                @edit="editSampling"
                @deleted="onSamplingDelete"
              />
            </v-col>
          </v-row>
        </v-container>
      </v-tabs-window-item>

      <v-tabs-window-item value="abiotic">
        <v-container v-if="!event?.abiotic_measurements?.length" fluid>
          <v-alert> No abiotic measurements reported </v-alert>
        </v-container>
        <v-list v-else :max-width="400">
          <v-list-item
            v-for="measurement in event?.abiotic_measurements"
            :title="measurement.param.label"
          >
            <template #append>
              <AbioticMeasurementChip :measurement />
            </template>
          </v-list-item>
        </v-list>
        <v-card-text v-if="!addItem">
          <v-btn
            class="ml-auto"
            text="Add measurement"
            prepend-icon="mdi-plus"
            variant="tonal"
            @click="toggleAddItem(true)"
          />
        </v-card-text>
        <v-card-text v-else>
          <v-card>
            <v-card-text class="d-flex">
              <v-row>
                <v-col class="d-flex" cols="12" sm="5">
                  <AbioticParameterPicker class="mr-3" density="compact" />
                </v-col>
                <v-col cols="12" sm="3">
                  <v-number-input label="Value" class="mr-3" density="compact" />
                </v-col>
                <v-col cols="auto" sm="4">
                  <v-btn class="mx-1" variant="tonal" color="primary" text="OK" />
                  <v-btn class="mx-1" variant="plain" text="Cancel" @click="toggleAddItem(false)" />
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-card-text>
      </v-tabs-window-item>

      <v-tabs-window-item value="spotting">
        <EventSpotting v-if="event" v-model="event" />
      </v-tabs-window-item>
    </v-tabs-window>

    <v-card-actions class="w-100 d-flex flex-column mt-auto">
      <v-divider class="w-100" />
      <v-btn
        v-if="tab == 'sampling'"
        block
        color="primary"
        prepend-icon="mdi-plus"
        text="Add sampling"
        variant="tonal"
        @click="newSampling()"
      />
      <v-list class="d-flex justify-space-between w-100">
        <v-list-item title="Performed by">
          <PersonChip v-for="p in event?.performed_by" class="ma-1" :person="p" />
        </v-list-item>
        <v-list-item title="Programs">
          <template #subtitle>
            <ProgramChip v-for="program in event?.programs" class="ma-1" :program />
          </template>
        </v-list-item>
      </v-list>
    </v-card-actions>
  </CardDialog>
</template>

<script setup lang="ts">
import { Event, Sampling } from '@/api'
import { useToggle } from '@vueuse/core'
import { ref } from 'vue'
import { useDisplay } from 'vuetify'
import PersonChip from '../people/PersonChip.vue'
import CardDialog from '../toolkit/forms/CardDialog.vue'
import AbioticParameterPicker from './AbioticParameterPicker.vue'
import EventSpotting from './EventSpotting.vue'
import ProgramChip from './ProgramChip.vue'
import SamplingCard from './SamplingCard.vue'
import SamplingFormDialog from './SamplingFormDialog.vue'
import { AbioticMeasurementChip } from './AbioticMeasurementChip'

const [samplingDialog, toggleSamplingDialog] = useToggle(false)
const editingSampling = ref<Sampling>()
function newSampling() {
  editingSampling.value = undefined
  toggleSamplingDialog(true)
}
function editSampling(sampling: Sampling) {
  editingSampling.value = sampling
  toggleSamplingDialog(true)
}

function onSamplingDelete(deleted: Sampling) {
  if (!event.value) return
  event.value.samplings = event.value.samplings.filter(({ id }) => id !== deleted.id)
}

const [addItem, toggleAddItem] = useToggle(false)

const { mobile, smAndDown } = useDisplay()

const open = defineModel<boolean>('open')
const event = defineModel<Event>()

export type EventAction = 'sampling' | 'abiotic' | 'spotting'
const tab = defineModel<EventAction>('tab', { default: 'sampling' })
</script>

<style lang="scss">
.event-action-text {
  min-height: 50vh;
}
.sampling-card {
  .v-card-item {
    padding-top: 0px;
    padding-left: 0px;
    .top-tag {
      height: 45px;
      padding: 10px;
      border-bottom-right-radius: 25%;
      font-weight: bold;
    }
  }
}
</style>
