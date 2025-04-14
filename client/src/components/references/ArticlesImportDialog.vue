<template>
  <FormDialog
    title="Import bibliographic references"
    v-model="isOpen"
    btn-text="Upload"
    :fullscreen="smAndDown"
    :max-width="2000"
    scrollable
  >
    <!-- @submit="() => (file != undefined ? importFile(file, options) : null)" -->
    <v-progress-linear :model-value="progress" :max="items.length" color="primary" />
    <v-divider></v-divider>
    <div class="spreadsheet-container fill-height" @mouseup="dragging = false">
      <v-data-table
        id="spreadsheet"
        ref="spreadsheet"
        class="spreadsheet fill-height"
        :headers="headers"
        :items="keyedItems"
        :items-per-page="-1"
        density="compact"
        fixed-footer
        fixed-header
        show-select
        width="100%"
      >
        <template #item.query="{ item, index, value }">
          <div
            @mousedown="select(index, 1)"
            @mouseover="dragging ? selectEnd(index, 1) : undefined"
            @dblclick="edit(index, 1)"
            style="min-height: 100%"
          >
            <v-text-field
              v-if="isEditing(index, 1)"
              autofocus
              v-model.trim="item.query"
              variant="plain"
              hide-details
              density="compact"
              @update:model-value="onQueryChange(item)"
              @blur="onEdited(item)"
            />
            <template v-else>
              <ErrorTooltip :error="item.errors?.query" class="d-flex align-center" error-class="">
                {{ value }}
              </ErrorTooltip>
            </template>
          </div>
        </template>

        <template #item.results="{ item }">
          <v-progress-circular indeterminate v-if="item.loading" color="primary" />
          <CrossRefItemSelect
            v-else-if="item.results"
            v-model="item.selectedRecord"
            :items="item.results.items"
            :total="item.results.total"
          />
        </template>

        <template #item.selectedRecord.score="{ item }">
          <v-chip
            v-if="item.selectedRecord?.score"
            size="small"
            :text="item.selectedRecord?.score?.toFixed(0).toString()"
            class="mr-2"
            :color="
              item.selectedRecord?.score === undefined
                ? ''
                : item.selectedRecord?.score > 120
                  ? 'success'
                  : item.selectedRecord?.score > 80
                    ? 'warning'
                    : 'error'
            "
          ></v-chip>
        </template>

        <template #bottom>
          <div style="position: sticky; bottom: 0px">
            <v-divider class="w-100"></v-divider>
            <div class="v-data-table-footer w-100 justify-space-between bg-surface">
              <div>
                <v-btn
                  color="primary"
                  variant="plain"
                  prepend-icon="mdi-plus"
                  text="Add row"
                  size="small"
                  @click="items.push({})"
                />
                <v-btn
                  v-if="selection"
                  color="warning"
                  variant="plain"
                  prepend-icon="mdi-close"
                  text="Delete row(s)"
                  size="small"
                  @click="items = items.filter((items, row) => !isSelected(row))"
                />
              </div>
              <v-btn
                color="warning"
                variant="plain"
                prepend-icon="mdi-close"
                text="Remove empty rows"
                size="small"
                @click="items = items.filter((item) => !isEmpty(item))"
              />
            </div>
          </div>
        </template>
      </v-data-table>
    </div>
  </FormDialog>
</template>

<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { useSpreadsheet } from '../sites'
import { Errors } from '../toolkit/validation'
import FormDialog from '../toolkit/forms/FormDialog.vue'
import ErrorTooltip from '../sites/ErrorTooltip.vue'
import { BibSearchResults, Message, ReferencesService } from '@/api'
import CrossRefItemSelect from './CrossRefItemSelect.vue'
import { computed, ref } from 'vue'
import { useDebounceFn } from '@vueuse/core'

const { smAndDown, mobile } = useDisplay()
const isOpen = defineModel<boolean>()

type Item = {
  query: string
  results?: BibSearchResults
  selectedRecord?: Message
  loading?: boolean
  errors?: Errors<'query'>
}

const {
  items,
  dragging,
  selection,
  edit,
  onEdited,
  cellHeader,
  select,
  selectEnd,
  isEditing,
  isSelected,
  isEmpty
} = useSpreadsheet<Item>(setValue, { startCol: 1 })

const onQueryChange = useDebounceFn(queryItem, 300)

const headers: DataTableHeader[] = [
  cellHeader({
    key: 'query',
    title: 'Query'
  }),
  {
    key: 'results',
    title: 'Matches',
    width: '50%',
    sortable: false
  },
  {
    key: 'selectedRecord.score',
    title: 'Score',
    width: 0
  }
]

const keyedItems = computed(() => items.value.map((it, i) => ({ ...it, id: i })))

const progress = ref(0)

async function queryItem(item: Partial<Item>) {
  if (!item.query) return
  item.loading = true
  item.results = undefined
  item.selectedRecord = undefined
  item.errors = undefined
  const { data, error } = await ReferencesService.crossRefBibSearch({ body: item.query }).finally(
    () => {
      item.loading = false
      progress.value++
    }
  )
  if (error) {
    item.errors = { query: error.detail }
    return
  }
  item.results = data
  item.selectedRecord = data.items[0]
}

function setValue(row: number, col: number, value: any) {
  if (row < items.value.length && col == 1) {
    items.value[row].query = value
    queryItem(items.value[row])
  }
}
</script>

<style scoped lang="scss">
.spreadsheet-container {
  min-width: 800px;
  overflow-x: scroll;
}
</style>
