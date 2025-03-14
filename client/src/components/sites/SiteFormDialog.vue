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
  </CreateUpdateForm>
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate, Site, SiteInput, SiteUpdate } from '@/api'
import { createSiteMutation, updateSiteMutation } from '@/api/gen/@tanstack/vue-query.gen'
import CoordPrecisionPicker from '@/components/sites/CoordPrecisionPicker.vue'
import CountryPicker from '@/components/toolkit/forms/CountryPicker.vue'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { useFeedback } from '@/stores/feedback'
import CreateUpdateForm from '../toolkit/forms/CreateUpdateForm.vue'
import FTextField from '../toolkit/forms/FTextField'

const item = defineModel<Site>()
const dialog = defineModel<boolean>('dialog')

const initial: SiteInput = {
  name: '',
  code: '',
  coordinates: { precision: '<100m', latitude: 0, longitude: 0 },
  user_defined_locality: false
}

const { feedback } = useFeedback()

function updateTransformer({
  id,
  meta,
  $schema,
  events,
  datasets,
  country,
  ...rest
}: Site): SiteUpdate {
  return {
    ...rest,
    country_code: country?.code
  }
}

const create = {
  mutation: createSiteMutation,
  schema: $SiteInput
}

const update = {
  mutation: updateSiteMutation,
  schema: $SiteUpdate,
  itemID: ({ code }: Site) => ({ code })
}

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
