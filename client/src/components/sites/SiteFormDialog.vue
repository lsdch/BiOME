<template>
  <CreateUpdateForm
    v-model="item"
    :create
    :update
    :local
    @created="onCreated()"
    @updated="onUpdated()"
    @save="dialog = false"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
        v-bind="props"
        v-model="dialog"
        :title="`${mode} site`"
        :loading="loading.value"
        @submit="submit"
        :fullscreen="mdAndDown"
      >
        <!-- Expose activator slot -->
        <template #activator="slotData">
          <slot name="activator" v-bind="slotData"></slot>
        </template>

        <v-container fluid>
          <v-row>
            <v-col>
              <v-text-field v-model="model.name" label="Name" v-bind="field('name')" />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <FTextField
                v-model.upper="model.code"
                label="Code"
                v-bind="field('code')"
                class="input-font-monospace"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <v-textarea
                v-model="model.description"
                label="Description"
                v-bind="field('description')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" sm="6">
              <div class="d-flex align-center w-100">
                <div class="flex-grow-1">
                  <v-number-input
                    v-model.number="model.coordinates!.latitude"
                    label="Latitude"
                    :precision="4"
                    :step="0.01"
                    class="input-latitude"
                    v-bind="field('coordinates', 'latitude')"
                    @input="pauseGPS()"
                  />
                  <v-number-input
                    v-model.number="model.coordinates!.longitude"
                    label="Longitude"
                    :precision="4"
                    :step="0.01"
                    class="input-longitude"
                    v-bind="field('coordinates', 'longitude')"
                    @input="pauseGPS()"
                  />
                </div>
                <div v-if="isGeolocationSupported">
                  <div class="gps-link upper" />
                  <v-hover>
                    <template #default="{ isHovering, props }">
                      <v-btn
                        class="gps-btn px-2 ml-2"
                        :min-width="30"
                        :height="50"
                        text="GPS"
                        stacked
                        variant="text"
                        size="small"
                        rounded="md"
                        @click="pendingGPS ? pauseGPS() : setCoordsFromGPS(model)"
                        v-bind="props"
                      >
                        <template #prepend>
                          <v-progress-circular
                            v-if="pendingGPS && !isHovering"
                            color="primary"
                            indeterminate
                            size="small"
                          />
                          <v-icon
                            v-else
                            :icon="pendingGPS ? 'mdi-close' : 'mdi-crosshairs-gps'"
                            :color="pendingGPS ? 'red' : undefined"
                          />
                        </template>
                      </v-btn>
                    </template>
                  </v-hover>
                  <div class="gps-link lower mb-5" />
                </div>
              </div>
              <CoordPrecisionPicker
                v-model="model.precision"
                v-bind="field('coordinates', 'precision')"
              />
              <v-number-input
                v-model.number="model.altitude"
                label="Altitude (m)"
                v-bind="field('altitude')"
              />
            </v-col>
            <v-col>
              <v-card height="300">
                <SiteProximityMap v-model="model.coordinates" />
              </v-card>
            </v-col>
          </v-row>
          <v-divider />

          <!-- Handle automatic/user-defined locality and country definition -->
          <SiteFormLocationField
            v-model:country_code="model.country_code"
            v-model:locality="model.locality"
            v-model:user_defined_locality="model.user_defined_locality"
            :coordinates="model.coordinates"
            :field
          />
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate, CoordinatesPrecision, Site, SiteInput, SiteUpdate } from '@/api'
import { createSiteMutation, updateSiteMutation } from '@/api/gen/@tanstack/vue-query.gen'
import CoordPrecisionPicker from '@/components/sites/CoordPrecisionPicker.vue'
import FormDialog, { FormDialogProps } from '@/components/toolkit/forms/FormDialog.vue'
import { defineFormCreate, defineFormUpdate } from '@/functions/mutations'
import { useFeedback } from '@/stores/feedback'
import { useGeolocation, useToggle, watchOnce } from '@vueuse/core'
import { useDisplay } from 'vuetify'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FTextField from '../toolkit/forms/FTextField'
import SiteFormLocationField from './SiteFormLocationField.vue'
import SiteProximityMap from './SiteProximityMap.vue'

const { mdAndDown } = useDisplay()
const item = defineModel<Site>()
const dialog = defineModel<boolean>('dialog')

const props = defineProps<
  Omit<FormDialogProps, 'loading' | 'fullscreen'> & {
    local?: true
  }
>()

type SiteInputModel = Omit<SiteInput, 'coordinates'> & {
  coordinates: { latitude?: number; longitude?: number }
  precision: CoordinatesPrecision
}

type SiteUpdateModel = Omit<SiteUpdate, 'coordinates'> & {
  coordinates: { latitude?: number; longitude?: number }
  precision: CoordinatesPrecision
}

const initial: SiteInputModel = {
  name: '',
  code: '',
  coordinates: {},
  precision: '<100m',
  user_defined_locality: false
}

function updateTransformer({
  id,
  meta,
  $schema,
  events,
  datasets,
  country,
  coordinates,
  ...rest
}: Site): SiteUpdateModel {
  return {
    ...rest,
    country_code: country?.code,
    coordinates: {
      latitude: coordinates.latitude,
      longitude: coordinates.longitude
    },
    precision: coordinates.precision
  }
}

const create = defineFormCreate(createSiteMutation(), {
  initial,
  schema: $SiteInput,
  requestData: (model: SiteInputModel) => {
    return {
      body: {
        ...model,
        coordinates: {
          latitude: model.coordinates.latitude!,
          longitude: model.coordinates.longitude!,
          precision: model.precision
        }
      }
    }
  }
})

const update = defineFormUpdate(updateSiteMutation(), {
  itemToModel: updateTransformer,
  schema: $SiteUpdate,
  requestData: ({ code }, model: SiteUpdateModel) => {
    return {
      path: { code },
      body: {
        ...model,
        coordinates: {
          latitude: model.coordinates.latitude!,
          longitude: model.coordinates.longitude!,
          precision: model.precision
        }
      }
    }
  }
})

const { feedback } = useFeedback()

function onCreated() {
  dialog.value = false
  feedback({
    type: 'success',
    message: `Site created successfully`
  })
}

function onUpdated() {
  dialog.value = false
  feedback({
    type: 'success',
    message: `Site updated successfully`
  })
}

const {
  isSupported: isGeolocationSupported,
  coords,
  locatedAt,
  pause,
  resume,
  error
} = useGeolocation({
  immediate: false,
  enableHighAccuracy: true
})

const [pendingGPS, togglePendingGPS] = useToggle(false)

function startGPS() {
  resume()
  togglePendingGPS(true)
}

function pauseGPS() {
  pause()
  togglePendingGPS(false)
}

function setCoordsFromGPS(model: SiteInputModel | SiteUpdateModel) {
  startGPS()
  watchOnce(
    () => coords.value,
    (coords) => {
      model.coordinates.latitude = coords.latitude
      model.coordinates.longitude = coords.longitude
      pauseGPS()
    }
  )
}
</script>

<style scoped lang="scss">
.gps-link {
  height: 1rem;
  width: 58%;
  border-right: 1px solid grey;
  background-color: transparent;
  &.upper {
    border-top: 1px solid grey;
  }
  &.lower {
    border-bottom: 1px solid grey;
  }
}
</style>
