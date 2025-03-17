<template>
  <v-card
    title="Location"
    class="small-card-title"
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
    <template #append>
      <div class="d-flex ga-3">
        <v-switch
          v-model="user_defined_locality"
          :true-value="false"
          :false-value="true"
          density="compact"
          color="primary"
          hide-details
          label="Auto"
        />
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
              :disabled="!user_defined_locality"
              :hint="
                user_defined_locality
                  ? 'User defined locality'
                  : 'Inferred from coordinates using Geoapify'
              "
              persistent-hint
              v-bind="field('locality')"
            />
          </v-col>
          <v-col cols="12" sm="6">
            <CountryPicker
              v-model="country_code"
              item-value="code"
              v-bind="field('country_code')"
              clearable
            />
          </v-col>
        </v-row>
      </template>
      <v-list v-else-if="Coordinates.isValidCoordinates(coordinates)">
        <v-list-item :title="locality ?? undefined">
          <template #title>
            <v-list-item-title>
              {{ locality || 'Unknown locality' }}
              <v-chip v-if="country_code" :text="country_code"></v-chip>
            </v-list-item-title>
          </template>
          <template #subtitle>
            <v-list-item-subtitle v-if="!country_code">
              Not within country boundaries
            </v-list-item-subtitle>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { Geoapify } from '@/api'
import { coordinatesToCountryOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useGeoapify } from '@/stores/geoapify'
import { useQuery } from '@tanstack/vue-query'
import { computed, watch } from 'vue'
import { Coordinates } from '../maps'
import CountryPicker from '../toolkit/forms/CountryPicker.vue'
import { FieldBinding } from '../toolkit/forms/form'
import GeoapifyStatusButton from '../toolkit/services/geoapify/GeoapifyStatusButton.vue'
const locality = defineModel<string | null | undefined>('locality', { required: true })
const country_code = defineModel<string | null | undefined>('country_code', { required: true })
const user_defined_locality = defineModel<boolean | undefined>('user_defined_locality')

const props = defineProps<{
  coordinates: Partial<Coordinates>
  field: (prop: 'country_code' | 'locality') => FieldBinding
}>()

const { reverseGeocodeQuery } = useGeoapify()

const {
  data: reverseGeocodeResult,
  error,
  isFetching: reverseGeocodeIsPending,
  refetch: refetchReverseGeoCode
} = reverseGeocodeQuery(props.coordinates, {
  enabled: computed(() => !user_defined_locality.value)
})

watch(reverseGeocodeResult, (result) => {
  if (result) {
    locality.value = result ? Geoapify.Result.toLocality(result) : null
  }
})

const {
  data: countryFromCoords,
  isFetching: countryIsPending,
  refetch: refetchCountry
} = useQuery(
  computed(() => ({
    enabled: !user_defined_locality.value && Coordinates.isValidCoordinates(props.coordinates),
    ...coordinatesToCountryOptions({
      body: {
        latitude: props.coordinates.latitude!,
        longitude: props.coordinates.longitude!
      }
    })
  }))
)

watch(countryFromCoords, (country) => {
  country_code.value = country?.code ?? null
})

function refetch() {
  refetchReverseGeoCode()
  refetchCountry()
}
</script>

<style scoped lang="scss"></style>
