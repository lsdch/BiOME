<template>
  <FormDialog
    v-model="dialog"
    title="Create person"
    :loading="loading"
    @submit="submit"
    :fullscreen="smAndDown"
  >
    <v-form ref="form" @submit.prevent="submit">
      <v-container fluid>
        <v-row>
          <v-col cols="12" sm="4" md="3">
            <v-text-field
              v-model.trim="model.code"
              label="Code"
              class="input-overline"
              v-bind="field('code')"
              @input="() => (model.code = model.code?.toUpperCase())"
            />
          </v-col>
          <v-col cols="12" sm="8" md="9">
            <v-text-field v-model.trim="model.name" label="Name" v-bind="field('name')" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-text-field
              v-model.number="model.altitude"
              v-bind="field('altitude')"
              label="Altitude"
              suffix="m"
            />
          </v-col>
        </v-row>

        <div class="d-flex justify-space-between align-center mb-2">
          <span class="text-subtitle-1"> Coordinates </span>
          <span class="text-caption">WGS84 decimal degrees</span>
        </div>
        <CoordinatesPicker v-model="model.coordinates" />
        <v-row>
          <v-col cols="12" sm="4">
            <CountryPicker
              v-model="model.country_code"
              item-value="code"
              v-bind="field('country_code')"
            />
          </v-col>
          <v-col>
            <v-text-field
              label="Nearest locality"
              v-model="model.locality"
              v-bind="field('locality')"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-textarea
              v-model.trim="model.description"
              label="Description"
              variant="outlined"
              v-bind="field('description')"
            />
          </v-col>
        </v-row>
      </v-container>
    </v-form>
  </FormDialog>
</template>

<script setup lang="ts">
import { $SiteInput } from '@/api'
import FormDialog from '@/components/toolkit/forms/FormDialog.vue'
import { FormEmits, FormProps, useForm } from '@/components/toolkit/forms/form'
import { nextTick, ref, watch } from 'vue'
import { useDisplay } from 'vuetify'
import { VForm } from 'vuetify/components'
import CountryPicker from '../toolkit/forms/CountryPicker.vue'
import CoordinatesPicker from './CoordinatesPicker.vue'
import { SiteRecord } from './SiteImportDialog.vue'

const { smAndDown } = useDisplay()

const initial: SiteRecord = {
  name: '',
  code: '',
  coordinates: {
    precision: '<100m',
    latitude: undefined,
    longitude: undefined
  },
  country_code: '',
  exists: false
}

const dialog = defineModel<boolean>()
const form = ref<InstanceType<typeof VForm> | null>(null)
const props = defineProps<FormProps<SiteRecord>>()
const emit = defineEmits<FormEmits<SiteRecord>>()
const { loading, field, errorHandler, model } = useForm(props, $SiteInput, {
  initial,
  transformers: {}
})

// watch(
//   () => props.edit,
//   () => nextTick(() => form.value?.validate())
// )

watch(dialog, () => nextTick(() => form.value?.validate()))

async function submit() {
  emit('success', model.value)
}
</script>

<style scoped></style>
