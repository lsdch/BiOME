<template>
  <FormDialog v-model="dialog" title="Create person" :loading="loading" @submit="submit">
    <v-form @submit.prevent="submit">
      <v-container fluid>
        <v-text-field v-model.trim="model.name" label="Name" v-bind="field('name')" />
        <v-text-field
          v-model.trim="model.code"
          label="Code"
          v-bind="field('code')"
          @input="() => (model.code = model.code.toUpperCase())"
        />
        <!-- <TextField label="Code" v-bind="field('code')"> </TextField> -->
        <div class="text-subtitle-1 mb-2">Coordinates (WGS84 decimal degrees)</div>
        <CoordinatesPicker v-model="model.coordinates" />
        <div class="text-subtitle-1 mb-2">Altitude</div>
        <AltitudePicker v-model="model.altitude" />
        <CountryPicker v-bind="field('country_code')" />
        <v-text-field label="Municipality" v-bind="field('municipality')" />
        <v-text-field label="Region" v-bind="field('region')" />
        <v-textarea
          v-model.trim="model.description"
          label="Description"
          variant="outlined"
          v-bind="field('description')"
        />
      </v-container>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $SiteInput, SiteDataset, SiteInput, SiteItem } from '@/api'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { FormEmits, FormProps, useForm } from '@/components/toolkit/forms/form'
import CountryPicker from '../toolkit/forms/CountryPicker.vue'
import AltitudePicker from './AltitudePicker.vue'
import CoordinatesPicker from './CoordinatesPicker.vue'

const initial: SiteInput = {
  name: '',
  code: '',
  coordinates: {
    precision: 'Unknown',
    latitude: 0,
    longitude: 0
  },
  country_code: ''
}

const dialog = defineModel<boolean>()
const props = defineProps<FormProps<SiteItem>>()
const emit = defineEmits<FormEmits<SiteDataset>>()
const { loading, field, errorHandler, model } = useForm(props, $SiteInput, {
  initial,
  transformers: {}
})

async function submit() {}
</script>

<style scoped></style>
