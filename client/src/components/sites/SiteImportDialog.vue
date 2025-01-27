<template>
  <FormDialog
    title="Import sites"
    v-model="isOpen"
    btn-text="Upload"
    :fullscreen="sm"
    @submit="() => (file != undefined ? importFile(file, options) : null)"
  >
    <SiteImportSettings v-model="options.import" />

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
import { Coordinates, SiteInput } from '@/api'
import { ParseConfig, ParseError, ParseLocalConfig, parse } from 'papaparse'
import { ref, watch, watchEffect } from 'vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import ParserOptions from '../toolkit/import/ParserOptions.vue'
import SiteTableExpandedRow from './SiteTableExpandedRow.vue'

import SiteImportSettings, { ImportSettings } from './SiteImportSettings.vue'
import SiteStatusIcon from './SiteStatusIcon'
import { useDisplay } from 'vuetify'

const { sm } = useDisplay()

export type ParsedItem = Partial<
  Record<keyof (Omit<SiteInput, 'coordinates'> & Coordinates), string>
> & { existing: boolean; row: number }

type CSVTransforms<T extends object> = {
  [k in keyof T]?: (v: string) => T[k] | string
}
type Transforms = CSVTransforms<ParsedItem>

const transforms: Transforms = {
  existing(v: string) {
    return v === 'true'
  }
}

type Options = {
  import: ImportSettings
  config: ParseConfig
}

const options = ref<Options>({
  import: {
    existing: 'Include',
    defaultPrecision: undefined
  },
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
const isOpen = defineModel<boolean>('open', { default: false })
const props = defineProps<{ file?: File }>()

const headers = ref<DataTableHeader[]>([
  { title: 'Status', key: 'existing', width: 0, align: 'center' },
  { title: 'Code', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Latitude', key: 'latitude' },
  { title: 'Longitude', key: 'longitude' },
  { title: 'Altitude (m)', key: 'altitude' }
])

const preview = ref<ParsedItem[]>([])

const emit = defineEmits<{
  ready: [file: File, config: ParseConfig]
  parseChunk: [items: SiteRecord[], errors: ParseError[]]
  complete: []
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
      preview.value = results.data
    }
  }
  parse(file, cfg)
}

function importFile(file: File, { config }: Options) {
  let cfg: ParseLocalConfig<ParsedItem, File> = {
    ...config,
    chunkSize: 1000,
    complete() {
      isOpen.value = false
      emit('complete')
    },
    chunk({ data, errors }) {
      const items = data.map(processItem)
      emit('parseChunk', items, errors)
    }
  }
  parse(file, cfg)
}

export type ProcessedItem<T extends Record<string, unknown>> = {
  [k in keyof T]: T[k] extends Record<string, unknown>
    ? ProcessedItem<T[k]>
    : T[k] | string | undefined
}
export type SiteRecord = ProcessedItem<SiteInput & { exists: boolean }>

function numberTransform(v: string | undefined): number | string | undefined {
  if (v === undefined || v === '') return undefined
  const n = Number(v)
  return isNaN(n) ? v : n
}

function processItem(item: ParsedItem): SiteRecord {
  return {
    name: item.name,
    code: item.code,
    exists: item.existing,
    coordinates: {
      latitude: numberTransform(item.latitude),
      longitude: numberTransform(item.longitude),
      precision: item.precision ?? options.value.import.defaultPrecision
    },
    altitude: numberTransform(item.altitude),
    locality: item.locality,
    country_code: item.country_code,
    description: item.description
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
