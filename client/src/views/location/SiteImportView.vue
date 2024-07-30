<template>
  <v-stepper
    v-if="!datasetCreated"
    v-model="step"
    class="fill-height d-flex flex-column"
    :mobile="mobile"
  >
    <v-stepper-header>
      <div class="d-flex align-center">
        <v-btn
          class="mx-2"
          color="secondary"
          icon="mdi-arrow-left"
          variant="text"
          :to="{ name: 'sites' }"
        />
        <v-stepper-item :value="1" title="Dataset" editable :rules="[() => validitySteps[0]]" />
      </div>
      <v-stepper-item :value="2" title="Sites coordinates" />
      <v-stepper-item :value="3" title="New sites" editable />
      <v-stepper-item :value="4" title="Review" />
    </v-stepper-header>

    <v-stepper-window class="fill-height">
      <v-stepper-window-item :value="1">
        <v-form v-model="validitySteps[0]">
          <v-container>
            <v-row>
              <v-col>
                <v-text-field
                  v-model.trim="model.label"
                  v-bind="field('label')"
                  label="Dataset name"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <PersonPicker
                  v-model.self="model.maintainers"
                  v-bind="field('maintainers')"
                  item-value="alias"
                  label="Maintainers"
                  :items="maintainers"
                  multiple
                  persistent-hint
                >
                  <template #prepend-inner>
                    <v-chip
                      :text="user?.identity.full_name"
                      color="primary"
                      prepend-icon="mdi-pin"
                    />
                  </template>
                </PersonPicker>
              </v-col>
            </v-row>
            <v-row>
              <v-col>
                <v-textarea
                  v-model.trim="model.description"
                  v-bind="field('description')"
                  label="Description"
                  variant="outlined"
                />
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </v-stepper-window-item>
      <v-stepper-window-item :value="2" class="fill-height">
        <SiteDatasetPrimer />
      </v-stepper-window-item>
      <v-stepper-window-item :value="3" class="fill-height">
        <SiteTabularImport v-model="sites" />
      </v-stepper-window-item>
      <v-stepper-window-item :value="4" class="fill-height"> SUMMARY </v-stepper-window-item>
    </v-stepper-window>

    <v-stepper-actions @click:next="step += 1" @click:prev="step -= 1">
      <template #next="{ props }">
        <v-btn v-bind="{ ...props, ...submitBtnProps }" />
      </template>
    </v-stepper-actions>
  </v-stepper>
  <v-container v-else>
    <v-row>
      <v-col> </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import {
  $SiteDatasetInput,
  Coordinates,
  LocationService,
  PeopleService,
  SiteDataset,
  SiteDatasetInput
} from '@/api'
import PersonPicker from '@/components/people/PersonPicker.vue'
import { isGranted } from '@/components/people/userRole'
import SiteDatasetPrimer from '@/components/sites/SiteDatasetPrimer.vue'
import SiteTabularImport, { RecordElement } from '@/components/sites/SiteTabularImport.vue'
import { FormProps, useForm } from '@/components/toolkit/forms/form'
import { useUserStore } from '@/stores/user'
import { computed, ref, watch } from 'vue'
import { useDisplay } from 'vuetify'

const { user } = useUserStore()
const maintainers = await PeopleService.listPersons().then(({ data, error }) => {
  if (error) {
    console.error('Failed to fetch persons: ', error)
    return []
  }
  return data.filter((p) => p.user && isGranted(p.user, 'Contributor') && p.user.id !== user?.id)
})

const submitBtnProps = computed(() => {
  if (step.value !== 3) return { text: 'Next', disabled: !validitySteps.value[step.value - 1] }
  const isReady = result.value.errorCount === 0 && result.value.validSites.length > 0
  return {
    prependIcon: 'mdi-upload',
    text: 'Submit dataset',
    disabled: !isReady,
    color: isReady ? 'success' : 'error',
    onClick: submit
  }
})

const { mobile } = useDisplay()

const step = ref(1)
const validitySteps = ref([undefined, true, undefined])

const props = defineProps<FormProps<SiteDatasetInput>>()

const initial: SiteDatasetInput = {
  label: '',
  maintainers: []
}
const { field, model } = useForm(props, $SiteDatasetInput, { initial })

const datasetCreated = ref<SiteDataset>()

async function submit() {
  const { data, error } = await LocationService.createSiteDataset({ body: model.value })
  if (error) {
    console.error(error)
    return
  }
  datasetCreated.value = data
}

const sites = ref<RecordElement[]>([])

const result = computed(() => {
  return sites.value.reduce<{ errorCount: number; validSites: RecordElement[] }>(
    (acc, site) => {
      if (site.errors) acc.errorCount += 1
      else acc.validSites.push(site)
      return acc
    },
    { errorCount: 0, validSites: [] }
  )
})

watch(
  result,
  () => {
    if (result.value.errorCount === 0 && result.value.validSites.length > 0) bindSites()
  },
  { immediate: true }
)

type ImportData = Required<Pick<SiteDatasetInput, 'new_sites' | 'sites'>>

function bindSites() {
  if (result.value.errorCount) {
    model.value.new_sites = []
    model.value.sites = []
    return
  }
  const { sites, new_sites } = result.value.validSites.reduce<ImportData>(
    (acc, site) => {
      if (site.exists) {
        acc.sites.push(site.code!)
        return acc
      }
      acc.new_sites.push({
        name: site.name!,
        code: site.code!,
        coordinates: site.coordinates as Coordinates,
        country_code: site.country_code!,
        altitude: site.altitude ? Number(site.altitude) : undefined,
        description: site.description,
        locality: site.locality
      })
      return acc
    },
    { sites: [], new_sites: [] }
  )
  model.value.new_sites = new_sites
  model.value.sites = sites
}
</script>

<style lang="scss">
.v-stepper-window {
  margin: 0;
}
</style>
