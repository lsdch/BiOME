<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <v-text-field v-model="model.name" label="Name" v-bind="schema('name')" />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <FTextField
          v-model.upper="model.code"
          label="Code"
          v-bind="schema('code')"
          class="input-font-monospace"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-textarea
          v-model="model.description"
          label="Description"
          v-bind="schema('description')"
        />
      </v-col>
    </v-row>
    <v-alert v-if="gps.error.value" density="compact" color="warning" class="mb-3">
      {{
        gps.error.value.code == gps.error.value.PERMISSION_DENIED
          ? 'Geolocation permission denied by user.'
          : gps.error.value.code == gps.error.value.POSITION_UNAVAILABLE
            ? 'Geolocation position unavailable.'
            : gps.error.value.code == gps.error.value.TIMEOUT
              ? 'Geolocation timeout.'
              : gps.error.value.message
      }}
    </v-alert>
    <v-row>
      <v-col cols="12" sm="6">
        <div class="d-flex align-center w-100">
          <div class="flex-grow-1">
            <v-number-input
              v-model.number="model.coordinates.latitude"
              label="Latitude"
              :precision="4"
              :step="0.01"
              class="input-latitude"
              v-bind="schema('coordinates', 'latitude')"
              @input="pauseGPS()"
            />
            <v-number-input
              v-model.number="model.coordinates.longitude"
              label="Longitude"
              :precision="4"
              :step="0.01"
              class="input-longitude"
              v-bind="schema('coordinates', 'longitude')"
              @input="pauseGPS()"
            />
          </div>
          <div v-if="gps.isSupported">
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
                      v-if="pendingGPS && !isHovering && !gps.error.value"
                      color="primary"
                      indeterminate
                      size="small"
                    />
                    <v-icon
                      v-else
                      :icon="pendingGPS && !gps.error.value ? 'mdi-close' : 'mdi-crosshairs-gps'"
                      :color="pendingGPS && !gps.error.value ? 'red' : undefined"
                    />
                  </template>
                </v-btn>
              </template>
            </v-hover>
            <div class="gps-link lower mb-5" />
          </div>
        </div>
        <CoordPrecisionPicker
          v-model="model.coordinates.precision"
          v-bind="schema('coordinates', 'precision')"
        />
        <v-number-input
          v-model.number="model.altitude"
          label="Altitude (m)"
          v-bind="schema('altitude')"
        />
      </v-col>
      <v-col>
        <v-card height="300" class="d-flex flex-column">
          <SiteProximityMap v-model="model.coordinates" />
        </v-card>
      </v-col>
    </v-row>
    <v-divider />

    <SiteFormLocationField
      v-model:country_code="model.country_code"
      v-model:locality="model.locality"
      v-model:user_defined_locality="model.user_defined_locality"
      :coordinates="model.coordinates"
      :schema
    />
  </v-container>
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate } from '@/api'
import { FormProps } from '@/functions/mutations'
import { SiteModel } from '@/models'
import { useGeolocation, useToggle, watchOnce } from '@vueuse/core'
import { useSchema } from '../../composables/schema'
import FTextField from '../toolkit/forms/FTextField'
import CoordPrecisionPicker from '@/components/sites/CoordPrecisionPicker.vue'
import SiteProximityMap from '@/components/sites/SiteProximityMap.vue'
import SiteFormLocationField from './SiteFormLocationField.vue'

const model = defineModel<SiteModel.SiteFormModel>({
  default: SiteModel.initialModel
})

const { mode = 'Create' } = defineProps<FormProps>()

const {
  bind: { schema }
} = useSchema(mode === 'Create' ? $SiteInput : $SiteUpdate)

/**
 * Geolocation
 */

const gps = useGeolocation({
  immediate: false,
  enableHighAccuracy: true
})

const [pendingGPS, togglePendingGPS] = useToggle(false)

function startGPS() {
  gps.error.value = null
  gps.resume()
  togglePendingGPS(true)
}

function pauseGPS() {
  gps.pause()
  togglePendingGPS(false)
}

function setCoordsFromGPS(model: SiteModel.SiteFormModel) {
  startGPS()
  watchOnce(
    () => gps.coords.value,
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
