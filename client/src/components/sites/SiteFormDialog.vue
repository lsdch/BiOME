<template>
  <CreateUpdateForm
    v-model="item"
    :initial
    :update-transformer
    :create
    :update
    @created="onCreated()"
    @updated="onUpdated()"
  >
    <template #default="{ model, field, mode, loading, submit }">
      <FormDialog
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
              <v-number-input
                v-model.number="model.coordinates!.latitude"
                label="Latitude"
                :precision="4"
                :step="0.01"
                v-bind="field('coordinates', 'latitude')"
              />
              <v-number-input
                v-model.number="model.coordinates!.longitude"
                label="Longitude"
                :precision="4"
                :step="0.01"
                v-bind="field('coordinates', 'longitude')"
              />
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
          <v-row>
            <v-col cols="12" sm="6">
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
            <v-col cols="12" sm="6">
              <CountryPicker
                v-model="model.country_code"
                :coords="model.coordinates"
                item-value="code"
                v-bind="field('country_code')"
                clearable
              />
            </v-col>
          </v-row>
        </v-container>
      </FormDialog>
    </template>
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate, CoordinatesPrecision, Site, SiteInput, SiteUpdate } from '@/api'
import { createSiteMutation, updateSiteMutation } from '@/api/gen/@tanstack/vue-query.gen'
import CoordPrecisionPicker from '@/components/sites/CoordPrecisionPicker.vue'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useFeedback } from '@/stores/feedback'
import { useDisplay } from 'vuetify'
import CreateUpdateForm, {
  FormCreateMutation,
  FormUpdateMutation
} from '../toolkit/forms/CreateUpdateForm.vue'
import FTextField from '../toolkit/forms/FTextField'
import SiteProximityMap from './SiteProximityMap.vue'

const { mdAndDown } = useDisplay()
const item = defineModel<Site>()
const dialog = defineModel<boolean>('dialog')

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

const create: FormCreateMutation<Site, SiteInput, SiteInputModel, typeof $SiteInput> = {
  mutation: createSiteMutation,
  schema: $SiteInput,
  transformer: (model: SiteInputModel) => {
    return {
      ...model,
      coordinates: {
        latitude: model.coordinates.latitude!,
        longitude: model.coordinates.longitude!,
        precision: model.precision
      }
    }
  }
}

const update: FormUpdateMutation<Site, SiteUpdate, SiteUpdateModel, typeof $SiteUpdate> = {
  mutation: updateSiteMutation,
  schema: $SiteUpdate,
  itemID: ({ code }: Site) => ({ code }),
  transformer: (model: SiteUpdateModel) => {
    return {
      ...model,
      coordinates: {
        latitude: model.coordinates.latitude!,
        longitude: model.coordinates.longitude!,
        precision: model.precision
      }
    }
  }
}

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
</script>

<style scoped lang="scss"></style>
