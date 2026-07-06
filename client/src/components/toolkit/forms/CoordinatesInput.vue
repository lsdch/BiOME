<template>
  <div class="d-flex align-center w-100">
    <div class="flex-grow-1">
      <v-number-input
        v-model.number="model.latitude"
        label="Latitude"
        :precision="4"
        :step="0.01"
        class="input-latitude"
        v-bind="schema('latitude')"
        @input="pauseGPS()"
      />
      <v-number-input
        v-model.number="model.longitude"
        label="Longitude"
        :precision="4"
        :step="0.01"
        class="input-longitude"
        v-bind="schema('longitude')"
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
</template>

<script setup lang="ts">
import { $CoordinatesWithPrecision, CoordinatesWithPrecision } from '@/api'
import { useSchema } from '@/composables/schema'
import { Coordinates } from '@/features/cartography/coordinates'
import { useGeolocation, useToggle, watchOnce } from '@vueuse/core'

const model = defineModel<Partial<Coordinates>>({
  default: { latitude: undefined, longitude: undefined }
})

const emit = defineEmits<{
  updateCoords: [value: number]
}>()

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

function setCoordsFromGPS(model: Partial<Coordinates>) {
  startGPS()
  watchOnce(
    () => gps.coords.value,
    (coords) => {
      model.latitude = coords.latitude
      model.longitude = coords.longitude
      if (coords.altitude !== null) {
        emit('updateCoords', coords.altitude)
      }
      pauseGPS()
    }
  )
}

const {
  bind: { schema }
} = useSchema($CoordinatesWithPrecision)
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
