<template>
  <v-card
    class="site-location small-card-title"
    prepend-icon="mdi-town-hall"
    variant="flat"
    :subtitle="
      user_defined_locality
        ? 'User defined'
        : Coordinates.isValidCoordinates(coordinates)
          ? 'Inferred from coordinates'
          : 'Waiting for coordinates'
    "
  >
    <template #title>
      <v-card-title>
        <template v-if="user_defined_locality || !Coordinates.isValidCoordinates(coordinates)">
          Location
        </template>
        <template v-else>
          {{ locality || 'Unknown locality' }}
          <CountryChip v-if="countryFromCoords" :country="countryFromCoords" size="small" />
          <v-chip v-else text="Unknown" size="small" />
        </template>
      </v-card-title>
    </template>
    <template #append>
      <div :class="['d-flex ga-2', { 'flex-column align-end': $vuetify.display.smAndDown }]">
        <v-switch
          v-model="user_defined_locality"
          :true-value="false"
          :false-value="true"
          density="compact"
          color="primary"
          hide-details
          label="Auto"
        />
        <div class="d-flex ga-2">
          <v-btn
            icon="mdi-restore"
            title="Auto-fill from coordinates"
            variant="tonal"
            color=""
            size="small"
            :disabled="!Coordinates.isValidCoordinates(coordinates)"
            @click="refetch()"
          />
          <GeoapifyStatusButton />
        </div>
      </div>
    </template>
    <v-card-text v-if="user_defined_locality || Coordinates.isValidCoordinates(coordinates)">
      <v-progress-linear
        v-if="reverseGeocodeIsPending || countryIsPending"
        indeterminate
        color="primary"
      />
      <template v-if="user_defined_locality">
        <div class="text-caption mb-3">
          Manual input mode: please carefully review the provided informations
        </div>
        <v-row>
          <v-col cols="12" sm="6">
            <v-text-field
              label="Locality"
              v-model.trim="locality"
              persistent-placeholder
              :hint="
                user_defined_locality
                  ? 'User defined locality'
                  : 'Inferred from coordinates using Geoapify'
              "
              persistent-hint
              v-bind="schema('locality')"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <CountryPicker
              v-model="country_code"
              item-value="code"
              v-bind="schema('country_code')"
              clearable
            />
          </v-col>
        </v-row>
      </template>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Geoapify } from '@/api'
import { coordinatesToCountryOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useGeoapify } from '@/stores/geoapify'
import { useQuery } from '@tanstack/vue-query'
import { computed, unref, watch } from 'vue'
import { Coordinates } from '../maps'
import CountryPicker from '../toolkit/forms/CountryPicker.vue'
import GeoapifyStatusButton from '../toolkit/services/geoapify/GeoapifyStatusButton.vue'
import CountryChip from '../sites/CountryChip'
import { SchemaBinding } from '../../composables/schema'
const locality = defineModel<string | null | undefined>('locality', { required: true })
const country_code = defineModel<string | null | undefined>('country_code', { required: true })
const user_defined_locality = defineModel<boolean | undefined>('user_defined_locality')

const props = defineProps<{
  coordinates: Partial<Coordinates>
  schema: (...path: ['locality'] | ['country_code']) => SchemaBinding
}>()

const { reverseGeocodeQuery } = useGeoapify()

const {
  data: reverseGeocodeResult,
  error,
  isFetching: reverseGeocodeIsPending,
  refetch: refetchReverseGeoCode
} = reverseGeocodeQuery(props.coordinates, {
  staleTime: Infinity,
  enabled: computed(() => !user_defined_locality.value)
})

watch(reverseGeocodeIsPending, (isPending, wasPending) => {
  if (!isPending && wasPending) {
    locality.value = reverseGeocodeResult.value
      ? Geoapify.Result.toLocality(reverseGeocodeResult.value)
      : null
  }
})

const {
  data: countryFromCoords,
  isFetching: countryIsPending,
  refetch: refetchCountry
} = useQuery(
  computed(() => ({
    enabled: !user_defined_locality.value && Coordinates.isValidCoordinates(props.coordinates),
    staleTime: Infinity,
    ...coordinatesToCountryOptions({
      body: {
        latitude: props.coordinates.latitude!,
        longitude: props.coordinates.longitude!
      }
    })
  }))
)

function useCountryFromCoords() {
  country_code.value = unref(countryFromCoords)?.code ?? null
}

watch(
  countryIsPending,
  (isPending, wasPending) => wasPending && !isPending && useCountryFromCoords()
)

function refetch() {
  refetchReverseGeoCode()
  refetchCountry()
}
</script>

<script lang="ts">
/**
 * Handles automatic locality and country detection based on coordinates,
 * or switch to manual input mode.
 */
export default {}
</script>

<style scoped lang="scss">
.site-location {
  .v-card-item__content {
    align-self: start;
  }
}
</style>
