<template>
  <CRUDTable
    v-model="model"
    v-model:selected="selected"
    :headers="headers"
    density="compact"
    select-strategy="all"
    entityName="Site"
    item-value="id"
    :itemRepr="({ code, id }) => code ?? `#${id}`"
    :toolbar="{
      title: 'Import sites',
      icon: 'mdi-map-marker-radius',
      noSort: true,
      noFilters: true
    }"
    show-actions="edit"
    show-select
    show-expand
    :row-props="
      ({ item }: Record<'item', Item> & {}) => ({
        class: Object.keys(item.errors ?? {}).length > 0 ? 'error-row' : undefined
      })
    "
  >
    <template #[`toolbar-append-actions`]>
      <v-btn
        color="primary"
        variant="plain"
        icon="mdi-tune-vertical"
        @click="settingsDialog = true"
      />
      <SiteImportSettingsDialog v-model="settingsDialog" />
    </template>

    <!-- Columns -->
    <template
      v-for="(header, i) in headers.filter(({ key }) => !['exists', 'errors'].includes(key!))"
      :key="i"
      #[`item.${header.key}`]="{ value, item, column }"
    >
      <v-tooltip v-if="column.key !== null && item.errors[column.key] !== undefined">
        <template #activator="{ props }">
          <span class="text-error" v-bind="props">{{ value }}</span>
        </template>
        {{ item.errors[column.key] }}
      </v-tooltip>
      <span v-else>{{ value }}</span>
    </template>

    <!-- Status column -->
    <template #[`header.exists`]="props">
      <IconTableHeader v-bind="props" icon="mdi-link-variant-plus" />
    </template>
    <template #[`item.exists`]="{ value }">
      <SiteStatusIcon :exists="value" />
    </template>

    <!-- Errors column -->
    <template #[`header.errors`]="props">
      <IconTableHeader v-bind="props" icon="mdi-alert-circle-outline" color="error" />
    </template>
    <template #[`item.errors`]="{ value, item }">
      <v-chip
        v-if="value > 0"
        color="error"
        size="small"
        rounded
        @click="debug ? console.log(item.errors) : null"
        >{{ value }}</v-chip
      >
      <v-icon v-else color="success">mdi-check</v-icon>
    </template>

    <template #expanded-row="{ item }">
      <SiteTableExpandedRow :offset="2" :item="item" :errors="item.errors" />
      <!-- <tr> {{ item.errors }} </tr> -->
    </template>

    <!-- Empty table -->
    <template #[`no-data`]>
      <div class="d-flex justify-center align-center">
        <DropZone
          class="ma-5"
          width="100%"
          :height="200"
          :datatypes="['text/csv', 'text/tab-separated-values']"
          hint="Accepted formats: CSV, TSV"
          @upload="onUpload"
        />
      </div>
    </template>

    <!-- Footer -->
    <template #[`footer.prepend`]>
      <v-btn
        v-if="selected.length > 0"
        color="warning"
        variant="plain"
        :text="`Remove ${selected.length} selected items`"
        @click="removeItems(selected)"
      />
      <v-spacer />
      <v-btn
        color="primary"
        text="preview"
        prepend-icon="mdi-map"
        variant="plain"
        @click="showPreview = true"
      />
      <v-spacer />
    </template>

    <!-- Item form -->
    <template #form="{ dialog, onClose, onSuccess, editItem }">
      <SiteFormDialog
        title="New site"
        :edit="editItem"
        :model-value="dialog"
        @success="
          (item) => {
            console.log(item)
            return onSuccess(validateItem(item))
          }
        "
        @close="onClose"
      />
    </template>
  </CRUDTable>
  <SitesMapPreview v-model:open="showPreview" :sites="itemsPreview" />
  <SiteImportDialog v-model:open="importDialog" :file="importedFile" @parse-chunk="addItems" />
</template>

<script setup lang="ts">
import { $SiteInput } from '@/api'
import DropZone from '@/components/toolkit/import/DropZone.vue'
import { ParseError } from 'papaparse'
import { computed, ref } from 'vue'
import { ImportItem } from '.'
import { useSchema } from '../toolkit/forms/schema'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import { indexErrors } from '../toolkit/validation'
import SiteFormDialog from './SiteFormDialog.vue'
import SiteImportDialog, { ProcessedItem } from './SiteImportDialog.vue'
import SiteImportSettingsDialog from './SiteImportSettingsDialog.vue'
import SiteStatusIcon from './SiteStatusIcon.vue'
import SiteTableExpandedRow from './SiteTableExpandedRow.vue'
import SitesMapPreview from './SitesMapPreview.vue'

const showPreview = ref(false)

type Item = ImportItem

const model = ref<Item[]>([])
const itemsPreview = computed(() => {
  return model.value.filter(({ coordinates }) => {
    return (
      coordinates !== undefined &&
      validate('coordinates', 'latitude')(coordinates.latitude) === true &&
      validate('coordinates', 'longitude')(coordinates.longitude) === true
    )
  })
})

const selected = ref<string[]>([])
const debug = ref(false)
const settingsDialog = ref(false)

const emit = defineEmits<{
  ready: [sites: Item[]]
}>()

/**
 * Item ID generator
 */
const id = (function* genID() {
  let n = 0
  while (true) {
    yield `${(n += 1)}`
  }
})()

function removeItems(items: string[]) {
  model.value = model.value.filter(({ id }) => {
    return !items.includes(id ?? '')
  })
  selected.value = []
}

const headers: CRUDTableHeaders = [
  {
    title: 'Site status',
    key: 'exists',
    width: 0,
    align: 'center'
  },
  { title: 'Code', key: 'code', cellProps: { class: 'text-overline' } },
  { title: 'Name', key: 'name' },
  { title: 'Latitude', key: 'coordinates.latitude' },
  { title: 'Longitude', key: 'coordinates.longitude' },
  { title: 'Altitude (m)', key: 'altitude' },
  {
    title: 'Errors',
    key: 'errors',
    width: 0,
    align: 'center',
    value(item) {
      return Object.keys(item.errors).length ?? 0
    }
  }
]

const { schema, validate, paths, validateAll } = useSchema($SiteInput)

function validateItem(item: Item) {
  item.errors = indexErrors(validateAll(item))
  return item
}

function addItems(items: ProcessedItem[], parseErrors: ParseError[]) {
  const toAdd = items.map<Item>((item: ProcessedItem): Item => {
    return validateItem({
      id: id.next().value,
      code: item.code,
      name: item.name,
      altitude: item.altitude,
      description: item.description,
      locality: item.locality,
      coordinates: {
        latitude: item.coordinates?.latitude,
        longitude: item.coordinates?.longitude,
        precision: item.coordinates?.precision
      },
      country_code: item.country_code,
      exists: item.exists
    })
  })
  model.value.unshift(...toAdd)
}

const importDialog = ref(false)
const importedFile = ref<File>()

function onUpload(files: File[]) {
  importedFile.value = files.pop()
  importDialog.value = true
}
</script>

<style lang="scss">
@use 'vuetify';
table tr {
  > td:first-child,
  th:first-child {
    border-left: 2px solid transparent;
  }
  &.error-row > td:first-child {
    border-left: 2px solid rgb(var(--v-theme-error));
  }
}
</style>
