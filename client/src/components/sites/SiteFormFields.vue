<template>
  <v-row>
    <v-col cols="12" sm="8" md="9">
      <v-text-field v-model.trim="model.name" label="Name" v-bind="field('name')" />
    </v-col>
    <v-col cols="12" sm="4" md="3">
      <v-text-field
        v-model.trim="model.code"
        label="Code"
        class="input-overline"
        v-bind="field('code')"
        @input="() => (model.code = model.code?.toUpperCase())"
      />
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
      <v-text-field label="Nearest locality" v-model="model.locality" v-bind="field('locality')" />
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
</template>

<script setup lang="ts">
import { $SiteInput, $SiteUpdate, Site, SiteInput } from '@/api'
import { useForm, useSchema } from '../toolkit/forms/form'
import { computed } from 'vue'
import { toRefs } from '@vueuse/core'
import CoordinatesPicker from './CoordinatesPicker.vue'

const props = defineProps<{ edit?: Site }>()

const initial: PartialTips<SiteInput> = { coordinates: {} }

const { model } = useForm(props, { initial, transformers: {} })

const { field } = toRefs(computed(() => useSchema(props.edit ? $SiteUpdate : $SiteInput)))
</script>

<style scoped></style>
