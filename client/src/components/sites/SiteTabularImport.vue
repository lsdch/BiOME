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
    appendActions="edit"
    show-select
    show-expand
    :row-props="
      ({ item }: Record<'item', RecordElement> & {}) => ({
        class: item.errors ? 'error-row' : undefined
      })
    "
  >
    <template #[`toolbar-append-actions`]>
      <SiteImportSettingsDialog>
        <template #activator="{ open }">
          <v-btn color="primary" variant="plain" icon="mdi-tune-vertical" @click="open" />
        </template>
      </SiteImportSettingsDialog>
    </template>

    <!-- Columns -->
    <template
      v-for="(header, i) in headers.filter(({ key }) => !['exists', 'errors'].includes(key!))"
      :key="i"
      #[`item.${header.key}`]="{ value, item, column }"
    >
      <v-tooltip v-if="column.key !== null && item.errors?.[column.key] !== undefined">
        <template #activator="{ props }">
          <span v-if="value" class="text-error" v-bind="props">{{ value }}</span>
          <v-icon v-else icon="mdi-alert-outline" v-bind="props" color="error"></v-icon>
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
    <template #[`item.errors`]="{ value, toggleExpand, internalItem }">
      <v-chip
        v-if="value > 0"
        color="error"
        size="small"
        rounded
        :text="`${value}`"
        @click="toggleExpand(internalItem)"
      />
      <v-icon v-else color="success" icon="mdi-check" />
    </template>

    <template #expanded-row="{ item }">
      <SiteTableExpandedRow :offset="2" :item="item" :errors="item.errors" />
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

      <SitesMapPreview :sites="itemsPreview">
        <template #activator="{ open }">
          <v-btn
            color="primary"
            text="preview"
            prepend-icon="mdi-map"
            variant="plain"
            @click="open"
          />
        </template>
      </SitesMapPreview>

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
            return onSuccess(validateItem(item))
          }
        "
        @close="onClose"
      />
    </template>
  </CRUDTable>

  <div>Errors:{{ result.errorCount }} Valid sites : {{ result.validSites.length }}</div>

  <SiteImportDialog
    v-model:open="importDialog"
    :file="importedFile"
    @parse-chunk="addItems"
    @complete="importedFile = undefined"
  />
</template>

<script setup lang="ts">
import { $SiteInput } from '@/api'
import DropZone from '@/components/toolkit/import/DropZone.vue'
import { ParseError } from 'papaparse'
import { computed, ref } from 'vue'
import { useSchema } from '../toolkit/forms/schema'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'
import IconTableHeader from '../toolkit/tables/IconTableHeader.vue'
import { Errors, indexErrors } from '../toolkit/validation'
import SiteFormDialog from './SiteFormDialog.vue'
import SiteImportDialog, { SiteRecord } from './SiteImportDialog.vue'
import SiteImportSettingsDialog from './SiteImportSettingsDialog.vue'
import SiteStatusIcon from './SiteStatusIcon.vue'
import SiteTableExpandedRow from './SiteTableExpandedRow.vue'
import SitesMapPreview from './SitesMapPreview.vue'

/**
 * RecordElement is a parsed CSV line before any validation is applied
 */
export type RecordElement = SiteRecord & {
  id: string
  errors?: Errors<ObjectPaths<SiteRecord>>
}

const model = defineModel<RecordElement[]>({ default: [] })

/**
 * Selected sites codes
 */
const selected = ref<string[]>([])

/**
 * Subset of records that can be previewed, i.e. have valid geospatial coordinates
 */
const itemsPreview = computed(() => {
  return model.value.filter(({ coordinates }) => {
    return (
      validate('coordinates', 'latitude')(coordinates.latitude) === true &&
      validate('coordinates', 'longitude')(coordinates.longitude) === true
    )
  })
})

const result = computed(() => {
  return model.value.reduce<{ errorCount: number; validSites: RecordElement[] }>(
    (acc, site) => {
      if (site.errors) acc.errorCount += 1
      else acc.validSites.push(site)
      return acc
    },
    { errorCount: 0, validSites: [] }
  )
})

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
      return item.errors ? (Object.keys(item.errors).length ?? 0) : 0
    }
  }
]

const { schema, validate, paths, validateAll } = useSchema($SiteInput)

function validateItem(item: SiteRecord): RecordElement {
  return {
    id: id.next().value,
    ...item,
    errors: indexErrors(validateAll(item))
  }
}

function addItems(items: SiteRecord[], parseErrors: ParseError[]) {
  const toAdd = items.map<RecordElement>(validateItem)
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
