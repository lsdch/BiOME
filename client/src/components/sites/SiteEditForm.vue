<template>
  <v-form>
    <v-container>
      <v-row>
        <v-col>
          <v-text-field v-model="model.name" label="Name" v-bind="field('name')" />
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-text-field v-model="model.code" label="Code" v-bind="field('code')" />
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
        <v-col cols="12" sm="4">
          <NumberInput
            v-model="model.coordinates.latitude"
            label="Latitude"
            float
            v-bind="field('coordinates', 'latitude')"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <NumberInput
            v-model="model.coordinates.longitude"
            label="Longitude"
            float
            v-bind="field('coordinates', 'longitude')"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <CoordPrecisionPicker
            v-model="model.coordinates.precision"
            v-bind="field('coordinates', 'precision')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="8">
          <v-text-field label="Locality" v-model="model.locality" v-bind="field('locality')" />
        </v-col>
        <v-col cols="12" sm="4">
          <CountryPicker v-model="model.country" v-bind="field('country_code')" />
        </v-col>
      </v-row>
    </v-container>
  </v-form>
</template>

<script setup lang="ts">
import { $SiteInput, Site, SiteUpdate } from '@/api'
import CountryPicker from '../toolkit/forms/CountryPicker.vue'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'
import NumberInput from '../toolkit/ui/NumberInput.vue'
import { ref, watch } from 'vue'
import { useSchema } from '../toolkit/forms/schema'

const props = defineProps<{ site: Site }>()

const model = ref(props.site)

watch(
  () => props.site,
  (s) => (model.value = s)
)

const { field } = useSchema($SiteInput)

function updateInput(site: Site): SiteUpdate {
  const { name, code, coordinates, description, altitude, locality } = site
  return {
    name,
    code,
    coordinates,
    description,
    altitude,
    locality,
    country_code: site.country.code
  }
}
</script>

<style scoped></style>
