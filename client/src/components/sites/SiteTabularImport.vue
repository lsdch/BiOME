<template>
  <v-data-table :items="model" :headers="headers" show-select select-strategy="all" show-expand>
    <!-- <template v-for="header in headers" :key="header.key" #[headerSlotName(header)]>
      {{ header.title }} <br />
      {{ header.key }}
    </template> -->

    <template #top>
      <v-toolbar>
        <v-spacer />
        <v-btn prepend-icon="mdi-plus" color="primary" text="New site" @click="addSite" />
      </v-toolbar>
    </template>
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
    <template #[`footer.prepend`]> </template>
  </v-data-table>
  <SiteImportDialog
    v-model="importDialog"
    :file="importedFile"
    @parse-chunk="
      (items, errors) => {
        console.log('ITEMS:', items)
        model = model.concat(model, items)
        console.log('ERRORS:', errors)
      }
    "
  />
  <SiteFormDialog title="New site" v-model="formDialog"></SiteFormDialog>
</template>

<script setup lang="ts">
import { SiteDatasetInput, SiteInput } from '@/api'
import DropZone from '@/components/toolkit/import/DropZone.vue'
import { reactive, ref } from 'vue'
import SiteFormDialog from './SiteFormDialog.vue'
import SiteImportDialog from './SiteImportDialog.vue'
import { ParseConfig, ParseLocalConfig, parse } from 'papaparse'

const model = reactive<SiteInput[]>([])

const headers: DataTableHeader[] = [
  { title: undefined },
  { title: 'Code', key: 'code' },
  { title: 'Name', key: 'name' },
  { title: 'Latitude', key: 'latitude' },
  { title: 'Longitude', key: 'longitude' },
  { title: 'Altitude (m)', key: 'altitude' },
  { title: '', key: 'data-table-expand' }
  // { title: 'Region' },
  // { title: 'Municipality' },
  // { title: 'Country' },
  // { title: 'Description' }
]

const formDialog = ref(false)
function addSite() {
  formDialog.value = true
}

const importDialog = ref(false)
const importedFile = ref<File>()

function onUpload(files: File[]) {
  importedFile.value = files.pop()
  importDialog.value = true
}
</script>

<style lang="scss"></style>
