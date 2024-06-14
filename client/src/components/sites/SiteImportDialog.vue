<template>
  <FormDialog
    title="Import sites"
    v-model="open"
    @submit="() => (file != undefined ? importFile(file, config) : null)"
  >
    <v-switch label="Ignore codes already in dataset"></v-switch>

    <div class="d-flex align-center text-h6">
      <span class="ma-3"> Preview: </span>
      <pre class="ma-3">{{ file?.name }}</pre>
      <v-divider />
      <span class="ma-3 text-subtitle-2 text-no-wrap"> First 10 lines </span>
    </div>

    <v-data-table
      class="import-preview"
      :headers="headers"
      :items="preview"
      item-value="code"
      hide-default-footer
      :expanded="preview.map(({ id }) => id.toString())"
    >
      <template #[`item.altitude`]="{ item }">
        {{ item.altitude?.min }}
        {{ item.altitude?.max ? ` &dash; ${item.altitude.max}` : '' }}
      </template>
      <template #expanded-row="{ item }">
        <div class="mx-3 mb-3 d-flex">
          Locality:
          {{ [item.municipality, item.region].filter((v) => v).join(',') }}
          <pre>[{{ item.country_code ?? ' ? ' }}]</pre>
        </div>
        <div>{{ item.description }}</div>
      </template>
    </v-data-table>
  </FormDialog>
</template>

<script setup lang="ts">
import { LocalFile, ParseConfig, ParseError, ParseLocalConfig, parse } from 'papaparse'
import { ref, watch, watchEffect } from 'vue'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import { AltitudeRange, Coordinates, SiteInput } from '@/api'

const open = defineModel<boolean>('open')

type CSVTransforms<T extends object> = Partial<{ [k in keyof T]: (v: string) => T[k] }>

const props = defineProps<{ file?: File }>()

const headers = ref<DataTableHeader[]>([
  { title: 'Code', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Latitude', key: 'latitude' },
  { title: 'Longitude', key: 'longitude' },
  { title: 'Altitude (m)', key: 'altitude' }
])

const preview = ref<(ParsedItem & { id: number })[]>([])

type ParsedItem = Partial<Omit<SiteInput, 'coordinates'> & Coordinates>
type Transforms = CSVTransforms<ParsedItem>

const transforms: Transforms = {
  code(v: string) {
    return v ?? ''
  },
  altitude(v: string): AltitudeRange | undefined {
    if (v == '') return undefined
    const [min, max = undefined] = v.split('-')
    return { min: parseInt(min), max: max != undefined ? parseInt(max) : max }
  },
  latitude(v: string) {
    return v ? parseFloat(v) : NaN
  },
  longitude(v: string) {
    return v ? parseFloat(v) : NaN
  }
}

const config = ref<ParseConfig>({
  header: true,
  skipEmptyLines: true,
  transform(value, field: keyof ParsedItem) {
    return (transforms[field] ?? ((v) => v))(value)
  }
})

const emit = defineEmits<{
  ready: [file: File, config: ParseConfig]
  parseChunk: [items: SiteInput[], errors: ParseError[]]
}>()

watch(
  config,
  (cfg) => {
    if (props.file) {
      loadPreview(props.file, cfg)
    }
  },
  { deep: true }
)

watchEffect(() => {
  if (props.file == undefined) {
    preview.value = []
  } else {
    loadPreview(props.file, config.value)
  }
})

function loadPreview(file: File, config: ParseConfig) {
  const cfg: ParseLocalConfig<ParsedItem, File> = {
    ...config,
    preview: 10,
    complete(results) {
      preview.value = results.data.map((item, id) => ({ id, ...item }))
    }
  }
  parse(file, cfg)
}

function importFile(file: File, config: ParseConfig) {
  let cfg: ParseLocalConfig<ParsedItem, File> = {
    ...config,
    chunkSize: 1000,
    complete() {
      console.log('COMPLETE')
      open.value = false
    },
    chunk({ data, errors, meta }, parser) {
      const items = data.map(processItem)
      emit('parseChunk', items, errors)
    }
  }
  parse(file, cfg)
}

function processItem(item: ParsedItem): SiteInput {
  const coordinates =
    item.latitude && item.longitude
      ? {
          latitude: item.latitude,
          longitude: item.longitude,
          precision: 'Km1'
        }
      : null
  return {
    code: item.code ?? '',
    coordinates,
    country_code: item.country_code ?? '',
    name: item.name ?? '',
    altitude: item.altitude,
    description: item.description,
    municipality: item.municipality,
    region: item.region
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
