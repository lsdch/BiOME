<template>
  <FormDialog
    title="Import sites"
    v-model="open"
    btn-text="Upload"
    @submit="() => (file != undefined ? importFile(file, options) : null)"
  >
    <v-row>
      <v-col>
        <v-select
          label="Existing codes"
          v-model="options.existing"
          item-value="label"
          item-title="label"
          :items="existing"
        >
          <template #item="{ item, props }">
            <v-list-item :title="item.raw.label" :subtitle="item.raw.description" v-bind="props" />
          </template>
        </v-select>
      </v-col>
      <v-col>
        <CoordPrecisionPicker
          v-model="options.defaultPrecision"
          label="Default coord precision"
          clearable
          placeholder="None"
          persistent-placeholder
        />
      </v-col>
    </v-row>

    <div class="d-flex align-center text-h7">
      <span class="text-no-wrap my-3 mr-3"> Parsing options: </span>
    </div>
    <ParserOptions v-model="options.config" />

    <div class="d-flex align-center text-h6">
      <span class="my-3 mr-3"> Preview: </span>
      <pre class="ma-3">{{ file?.name }}</pre>
      <v-divider />
      <span class="ma-3 text-subtitle-2 text-no-wrap"> First 10 lines </span>
    </div>

    <v-data-table
      class="import-preview"
      :headers="headers.map((h: DataTableHeader) => ({ ...h, sortable: false }))"
      :items="preview"
      item-value="code"
      hide-default-footer
      :expanded="preview.map(({ code }) => code ?? '')"
    >
      <template #[`item.existing`]="{ item }">
        <SiteStatusIcon :exists="item.existing" />
      </template>
      <template #expanded-row="{ item }">
        <SiteTableExpandedRow :item="item" :offset="1" />
      </template>
    </v-data-table>
  </FormDialog>
</template>

<script setup lang="ts">
import { Coordinates, CoordinatesPrecision, SiteInput } from '@/api'
import { ParseConfig, ParseError, ParseLocalConfig, parse } from 'papaparse'
import { ref, watch, watchEffect } from 'vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import ParserOptions from '../toolkit/import/ParserOptions.vue'
import SiteTableExpandedRow from './SiteTableExpandedRow.vue'

import type { Object } from 'ts-toolbelt'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'
import SiteStatusIcon from './SiteStatusIcon.vue'

export type SiteRecord = SiteInput & { existing: boolean; row: number }
export type ParsedItem = Partial<Omit<SiteRecord, 'coordinates'> & Coordinates>

type CSVTransforms<T extends object> = Partial<{ [k in keyof T]: (v: string) => T[k] }>
type Transforms = CSVTransforms<ParsedItem>

const transforms: Transforms = {
  existing(v: string) {
    return v === 'true'
  },
  code(v: string) {
    return v ?? ''
  },
  altitude(v: string): number | undefined {
    if (v == '') return undefined
    return Number(v)
  },
  latitude(v: string) {
    return v ? Number(v) : NaN
  },
  longitude(v: string) {
    return v ? Number(v) : NaN
  }
}

const existing = [
  {
    label: 'Restrict',
    description: 'Entirely disallow the use of existing codes'
  },
  {
    label: 'Omit',
    description: 'Omit existing codes from the dataset'
  },
  {
    label: 'Include',
    description: 'Include existing codes in the dataset'
  }
] as const

type Options = {
  existing: (typeof existing)[number]['label']
  defaultPrecision?: CoordinatesPrecision
  config: ParseConfig
}

const options = ref<Options>({
  existing: 'Include',
  defaultPrecision: undefined,
  config: {
    quoteChar: undefined,
    newline: undefined,
    header: true,
    skipEmptyLines: true,
    transform(value, field: keyof ParsedItem) {
      if (value === '') return undefined
      else return (transforms[field] ?? ((v) => v))(value)
    }
  }
})
const open = defineModel<boolean>('open')
const props = defineProps<{ file?: File }>()

const headers = ref<DataTableHeader[]>([
  { title: 'Status', key: 'existing', width: 0, align: 'center' },
  { title: 'Code', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Latitude', key: 'coordinates.latitude' },
  { title: 'Longitude', key: 'coordinates.longitude' },
  { title: 'Altitude (m)', key: 'altitude' }
])

const preview = ref<ParsedItem[]>([])

const emit = defineEmits<{
  ready: [file: File, config: ParseConfig]
  parseChunk: [items: ProcessedItem[], errors: ParseError[]]
}>()

watch(
  options,
  (options) => {
    console.log('options changed')
    if (props.file) loadPreview(props.file, options)
  },
  { deep: true }
)

watchEffect(() => {
  if (props.file == undefined) {
    preview.value = []
  } else {
    loadPreview(props.file, options.value)
  }
})

function loadPreview(file: File, { config }: Options) {
  const cfg: ParseLocalConfig<ParsedItem, File> = {
    ...config,
    preview: 10,
    complete(results) {
      const errors = results.errors
      console.log(errors)
      preview.value = results.data.map(processItem)
    }
  }
  parse(file, cfg)
}

function importFile(file: File, { config }: Options) {
  let cfg: ParseLocalConfig<ParsedItem, File> = {
    ...config,
    chunkSize: 1000,
    complete() {
      console.log('COMPLETE')
      open.value = false
      console.log(open.value)
    },
    chunk({ data, errors }) {
      const items = data.map(processItem)
      emit('parseChunk', items, errors)
    }
  }
  parse(file, cfg)
}

export type ProcessedItem = Object.Partial<SiteInput, 'deep'> & { exists?: boolean }

function processItem(item: ParsedItem): ProcessedItem {
  return {
    code: item.code,
    coordinates: {
      latitude: item.latitude,
      longitude: item.longitude,
      precision: item.precision ?? options.value.defaultPrecision
    },
    country_code: item.country_code,
    name: item.name,
    altitude: item.altitude,
    description: item.description,
    locality: item.locality,
    exists: item.existing
  }
}
</script>

<style lang="scss">
@use 'vuetify';

.import-preview {
  td {
    border-bottom: none !important;
    border-top: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
  }
}
</style>
