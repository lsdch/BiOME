<template>
  <FormDialog
    v-model="dialog"
    :title="`${mode} bio material`"
    :loading
    @submit="submit"
    :fullscreen="$vuetify.display.smAndDown"
    :max-width="1200"
  >
    <!-- Expose activator slot -->
    <template #activator="slotData">
      <slot name="activator" v-bind="slotData"></slot>
    </template>
    <v-container class="bg-main overflow-y-auto" fluid min-height="100%">
      <v-row align="stretch" class="bg-main">
        <v-col>
          <SiteSelectorCard
            class="fill-height small-card-title"
            @update:model-value="(s) => (site = s?.code)"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" md="6">
          <EventFormComponent :site-code="site" v-model="event" />
        </v-col>
        <v-col cols="12" md="6">
          <v-card
            title="Sampling"
            :subtitle="event ? 'From event' : 'Waiting for event definition'"
            class="small-card-title"
            prepend-icon="mdi-package-down"
            flat
          >
            <template #append v-if="event">
              <v-btn text="New sampling" />
            </template>
            <v-card-text>
              <SamplingPicker v-if="event" :items="event.samplings ?? []" />
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card title="Bio material" prepend-icon="mdi-package-variant">
            <v-card v-if="mode === 'Create'" title="Category" flat class="small-card-title">
              <v-card-text>
                <v-btn-toggle v-model="category" mandatory divided rounded="md" variant="outlined">
                  <v-btn
                    value="Internal"
                    text="Internal"
                    :prepend-icon="OccurrenceCategory.icon('Internal')"
                    :color="OccurrenceCategory.props.Internal.color"
                  />
                  <v-btn
                    value="External"
                    text="External"
                    :prepend-icon="OccurrenceCategory.icon('External')"
                    :color="OccurrenceCategory.props.External.color"
                  />
                </v-btn-toggle>
              </v-card-text>
            </v-card>
            <v-divider v-if="category" class="mb-3" />
            <v-card-text>
              <ExternalBioMatForm v-if="category === 'External'" />
              <template v-else-if="category === 'Internal'"> [Internal bio mat form] </template>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { BioMaterial, Event, OccurrenceCategory } from '@/api'
import { FormEmits, FormProps, Mode } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useToggle } from '@vueuse/core'
import { computed, ref } from 'vue'
import ExternalBioMatForm from './ExternalBioMatForm.vue'
import SiteSelectorCard from '../sites/SiteSelectorCard.vue'
import SiteEventPicker from './SiteEventPicker.vue'
import SamplingPicker from './SamplingPicker.vue'
import EventFormComponent from './EventFormComponent.vue'

const dialog = defineModel<boolean>()

const props = defineProps<FormProps<BioMaterial>>()
const emit = defineEmits<FormEmits<BioMaterial>>()

const mode = computed<Mode>(() => (props.edit ? 'Edit' : 'Create'))
const category = ref<OccurrenceCategory>()

const site = ref<string>()
const event = ref<Event>()

const [loading, toggleLoading] = useToggle(false)

async function submit() {}
</script>

<style scoped lang="scss"></style>
