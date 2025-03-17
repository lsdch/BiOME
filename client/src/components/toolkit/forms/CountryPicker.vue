<template>
  <v-autocomplete
    v-model="model"
    v-bind="$attrs"
    label="Country"
    :item-value
    :return-object
    :items="items"
    item-title="name"
    filter-mode="some"
    :loading="isPending || coordsIsPending"
    :error-messages="error?.detail"
    :custom-filter="
      (_: any, q: string, item: any) => {
        const { code, name }: Country = item?.raw
        if (q == '') return true
        return (
          code.toLowerCase().includes(q.toLowerCase()) ||
          name.toLowerCase().includes(q.toLowerCase())
        )
      }
    "
  >
    <template #item="{ item, props }">
      <v-list-item v-bind="props" :title="item.raw.name">
        <template #append>
          <span class="text-overline">
            {{ item.raw.code }}
          </span>
        </template>
      </v-list-item>
    </template>
    <template #append>
      <v-btn
        icon="mdi-restore"
        @click.stop="refetchCountry()"
        title="Reset country"
        variant="tonal"
        color=""
        size="small"
        :disabled="!Coordinates.isValidCoordinates(props.coords)"
      />
    </template>
    <template #append-inner>
      <v-chip v-if="model" :text="typeof model == 'string' ? model : model.code"></v-chip>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { CoordinatesToCountryResponse, Country } from '@/api'
import { coordinatesToCountryOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { Coordinates } from '@/components/maps'
import { useCountries } from '@/stores/countries'
import { useQuery } from '@tanstack/vue-query'
import { storeToRefs } from 'pinia'
import { Ref, unref } from 'vue'
import { computed, MaybeRef, watch } from 'vue'

const model = defineModel<Country | string | undefined | null>({ required: true })

const { countries: items, isPending, error } = storeToRefs(useCountries())

const props = withDefaults(
  defineProps<{
    coords?: {
      latitude?: number
      longitude?: number
    }
    returnObject?: boolean
    itemValue?: 'code' | 'name'
  }>(),
  { itemValue: 'code' }
)

const {
  data: countryFromCoords,
  refetch,
  isFetching: coordsIsPending
} = useQuery(
  computed(() => ({
    enabled: () => Coordinates.isValidCoordinates(props.coords),
    ...coordinatesToCountryOptions({
      body: {
        latitude: props.coords?.latitude ?? 0,
        longitude: props.coords?.longitude ?? 0
      }
    })
  }))
)

watch(countryFromCoords, setCountry)

// Reset country if incomplete coordinates
watch(
  () => props.coords,
  (coords) => {
    if (!!coords && !Coordinates.isValidCoordinates(coords)) {
      model.value = null
    }
  },
  { deep: true }
)

// Updates model value with inferred country from coordinates
function setCountry(country: MaybeRef<CoordinatesToCountryResponse> | Ref<undefined>) {
  const value = unref(country)
  if (value) {
    model.value = props.returnObject ? value : value[props.itemValue]
  } else {
    model.value = null
  }
}

// Manually refetch country from coordinates
async function refetchCountry() {
  setCountry((await refetch()).data)
}
</script>

<style scoped></style>
