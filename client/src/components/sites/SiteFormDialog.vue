<template>
  <FormDialog v-model="dialog" :title="`${mode} site`" :loading @submit="submit">
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
        <v-col cols="12" sm="4">
          <v-number-input
            v-model="model.coordinates!.latitude"
            label="Latitude"
            float
            v-bind="field('coordinates', 'latitude')"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <v-number-input
            v-model="model.coordinates!.longitude"
            label="Longitude"
            float
            v-bind="field('coordinates', 'longitude')"
          />
        </v-col>
        <v-col cols="12" sm="4">
          <CoordPrecisionPicker
            v-model="model.coordinates!.precision"
            v-bind="field('coordinates', 'precision')"
          />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="8">
          <div class="d-flex">
            <v-switch v-model="model.user_defined_locality" :width="50" color="primary" />
            <v-text-field
              label="Locality"
              v-model="model.locality"
              persistent-placeholder
              :disabled="!model.user_defined_locality"
              :hint="
                model.user_defined_locality
                  ? 'User defined locality'
                  : 'Inferred from coordinates using Geoapify'
              "
              persistent-hint
              v-bind="field('locality')"
            >
            </v-text-field>
          </div>
        </v-col>
        <v-col cols="12" sm="4">
          <CountryPicker
            v-model="model.country_code"
            item-value="code"
            v-bind="field('country_code')"
          />
        </v-col>
      </v-row>
    </v-container>
  </FormDialog>
</template>

<script setup lang="ts">
import { Site, SiteInput, SiteUpdate, $SiteInput, $SiteUpdate, LocationService } from '@/api'
import CoordPrecisionPicker from '@/components/sites/CoordPrecisionPicker.vue'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import { FormEmits, FormProps, useForm, useSchema } from '@/components/toolkit/forms/form'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { reactiveComputed, useToggle } from '@vueuse/core'
import FTextField from '../toolkit/forms/FTextField'

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<Site>>()
const emit = defineEmits<FormEmits<Site>>()

const initial: SiteInput = {
  name: '',
  code: '',
  coordinates: { precision: '<100m', latitude: 0, longitude: 0 },
  country_code: '',
  user_defined_locality: false
}

const { model, mode, makeRequest } = useForm(props, {
  initial,
  updateTransformer({ id, meta, $schema, events, datasets, country, ...rest }): SiteUpdate {
    return {
      ...rest,
      country_code: country.code
    }
  }
})

const { field, errorHandler } = reactiveComputed(() =>
  useSchema(mode.value === 'Create' ? $SiteInput : $SiteUpdate)
)

const [loading, toggleLoading] = useToggle(false)

async function submit() {
  toggleLoading(true)
  return await makeRequest({
    create: LocationService.createSite,
    edit: ({ code }, body) => LocationService.updateSite({ path: { code }, body })
  })
    .then(errorHandler)
    .then((item) => emit('success', item))
    .finally(() => toggleLoading(false))
}
</script>

<style scoped lang="scss"></style>
