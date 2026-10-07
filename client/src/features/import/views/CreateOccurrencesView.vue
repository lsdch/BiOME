<template>
  <v-card
    flat
    title="Create sampling and occurrences"
    class="d-flex flex-column h-100"
    prepend-icon="mdi-package-variant-plus"
  >
    <template #append>
      <v-btn color="primary" variant="tonal" text="Submit" prepend-icon="mdi-upload"></v-btn>
    </template>
    <v-divider />
    <v-container class="bg-main overflow-y-auto flex-grow-1" fluid>
      <v-row class="bg-main align-stretch">
        <v-col cols="12" lg="6">
          <div class="d-flex flex-column ga-2 justify-space-between">
            <v-card title="Sampling">
              <v-divider />
              <v-tabs v-model="samplingTab" mandatory grow>
                <v-tab text="New" prepend-icon="mdi-plus" variant="tonal" value="new"></v-tab>
                <v-tab text="Search" prepend-icon="mdi-magnify" variant="tonal" value="search" />
              </v-tabs>
              <v-card-text>
                <v-tabs-window v-model="samplingTab">
                  <v-tabs-window-item value="new" class="d-flex flex-column ga-2">
                    <SamplingCoreFields v-model="sampling" />
                  </v-tabs-window-item>
                  <v-tabs-window-item value="search"></v-tabs-window-item>
                </v-tabs-window>
              </v-card-text>
            </v-card>
          </div>
        </v-col>
        <v-col cols="12" lg="6">
          <v-card height="500" class="d-flex flex-column">
            <ItemLocationMap
              v-model:item="sampling"
              :marker-options="{
                draggable: true
              }"
              class="flex-grow-1"
            />
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card title="Sampling details" prepend-icon="mdi-package-variant">
            <v-divider />
            <v-card-text fluid class="d-flex flex-column ga-8">
              <SiteFormLocationField
                v-model:country_code="sampling.site.country_code"
                v-model:locality="sampling.site.locality"
                :coordinates="sampling.coordinates"
                variant="outlined"
                class="border-md"
              />

              <TaxonPicker
                v-model="sampling.target_taxa"
                label="Target taxa"
                item-value="name"
                multiple
                chips
                closable-chips
                clearable
              />
              <TogglableFormCard title="Protocol" subtitle="Details about the sampling protocol">
                <v-container fluid>
                  <v-row>
                    <v-col cols="12" sm="6">
                      <FixativePicker
                        label="Fixatives"
                        v-model="sampling.fixatives"
                        multiple
                        return-object
                        chips
                        closable-chips
                        clearable
                      />
                    </v-col>
                    <v-col cols="12" sm="6">
                      <SamplingMethodPicker
                        label="Sampling methods"
                        v-model="sampling.methods"
                        multiple
                        return-object
                        chips
                        closable-chips
                        clearable
                      />
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-combobox
                        label="Performed by"
                        v-model.trim="sampling.performed_by"
                        multiple
                        chips
                        closable-chips
                      />
                    </v-col>
                    <v-col cols="12" md="6">
                      <HoursMinutesInput
                        label="Duration"
                        class="mt-2"
                        v-model="sampling.duration"
                        clearable
                      />
                    </v-col>
                  </v-row>
                </v-container>
              </TogglableFormCard>

              <TogglableFormCard
                title="Habitat and access points"
                subtitle="Environmental context of the sampling site"
              >
                <v-row>
                  <v-col cols="12" sm="6">
                    <!-- <HabitatPicker
                      v-model="sampling.habitats"
                      label="Habitat tags"
                      item-value="label"
                    /> -->
                  </v-col>
                  <v-col cols="12" sm="6">
                    <AccessPointsPicker
                      v-model="sampling.access_points"
                      label="Access points"
                      hint="Pick existing terms already in use, or enter new terms"
                      persistent-hint
                      clearable
                    />
                  </v-col>
                </v-row>
              </TogglableFormCard>

              <v-textarea v-model="sampling.comments" label="Sampling comments" />

              <!-- :schema -->
              <!-- v-model:user_defined_locality="sampling.site.user_defined_locality" -->
              <!-- <v-row>
                <v-col>
                  <TaxonPicker
                    v-model="model.target_taxa"
                    v-bind="schema('target_taxa')"
                    label="Target taxa"
                    item-value="name"
                    return-object
                    multiple
                    chips
                    closable-chips
                    clearable
                  />
                </v-col>
              </v-row> -->
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-card title="Occurrences" class="" prepend-icon="mdi-package-variant">
            <v-divider />
            <v-tabs>
              <v-tab variant="tonal" text="Occurrence #1"></v-tab>
              <v-tab prepend-icon="mdi-plus" text="Add"></v-tab>
            </v-tabs>
            <v-divider />
            <OccurrenceForm v-model="occurrence" />
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-card>
</template>

<script setup lang="ts">
import { SamplingInput } from '@/api'
import OccurrenceForm, { OccurrenceModel } from '@/components/forms/occurrence/OccurrenceForm.vue'
import TogglableFormCard from '@/components/forms/occurrence/TogglableFormCard.vue'
import SamplingCoreFields from '@/components/forms/sampling/SamplingCoreFields.vue'
import SiteFormLocationField from '@/components/forms/SiteFormLocationField.vue'
import HoursMinutesInput from '@/components/toolkit/forms/HoursMinutesInput.vue'
import ItemLocationMap from '@/features/cartography/components/ItemLocationMap.vue'
import AccessPointsPicker from '@/features/occurrences/components/sampling/AccessPointsPicker.vue'
import FixativePicker from '@/features/registries/components/FixativePicker.vue'
import SamplingMethodPicker from '@/features/registries/components/SamplingMethodPicker.vue'
import TaxonPicker from '@/features/taxonomy/components/TaxonPicker.vue'
import { SamplingModel } from '@/models'
import { ref } from 'vue'

type SamplingTab = 'new' | 'search'
const samplingTab = ref<SamplingTab>('new')

const sampling = defineModel<SamplingModel.SamplingModel>('sampling', {
  default: (): SamplingModel.SamplingModel => ({
    site: {},
    coordinates: {}
  })
})

const occurrence = defineModel<OccurrenceModel>('occurrence', {
  default: (): OccurrenceModel => ({
    identification: { identified_on: { date: {}, precision: 'day' } }
  })
})
</script>

<style scoped lang="scss"></style>
