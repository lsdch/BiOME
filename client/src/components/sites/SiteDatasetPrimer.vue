<template>
  <div class="fill-height" @mouseup="dragging = false">
    <v-divider />
    <v-data-table
      id="spreadsheet"
      class="spreadsheet fill-height"
      height="100"
      :headers="headers"
      :items="items"
      :items-per-page="-1"
      density="compact"
      disable-sort
      @paste.capture="handlePaste"
      fixed-footer
      fixed-header
    >
      <template
        v-for="(key, col) in ['latitude', 'longitude'] as const"
        :key="key"
        #[`item.${key}`]="{ item, index, value }"
      >
        <div
          class="fill-height"
          @mousedown="select(index, col)"
          @mouseover="dragging ? selectEnd(index, col) : undefined"
          @dblclick="edit(index, col)"
        >
          <v-text-field
            v-if="isEditing(index, col)"
            autofocus
            v-model.number="item[key]"
            variant="plain"
            hide-details
            density="compact"
            @blur="onEdited(item)"
          />
          <template v-else>
            <ErrorTooltip
              :error="item.errors?.[key]"
              class="fill-height d-flex align-center"
              error-class=""
            >
              {{ value }}
            </ErrorTooltip>
          </template>
        </div>
      </template>
      <template #[`item.precision`]="{ item, index, value }">
        <div
          class="fill-height"
          @mousedown="select(index, 2)"
          @mouseover="dragging ? selectEnd(index, 2) : undefined"
          @dblclick="edit(index, 2)"
        >
          <CoordPrecisionPicker
            v-if="isEditing(index, 2)"
            autofocus
            menu
            v-model="item.precision"
            density="compact"
            :label="undefined"
            placeholder="NA"
            variant="plain"
            hide-details
            no-label
            @update:modelValue="onEdited(item)"
            @update:menu="(open) => (open ? undefined : onEdited(item))"
          />
          <template v-else>
            <ErrorTooltip
              :error="item.errors?.precision"
              class="fill-height d-flex align-center"
              error-class=""
            >
              {{ value }}
            </ErrorTooltip>
          </template>
        </div>
      </template>

      <template #[`item.actions`]>
        <v-combobox
          v-model="s"
          :items="sites"
          density="compact"
          rounded="0"
          variant="filled"
          hide-details
          prepend-inner-icon="mdi-pencil-circle"
          @click:prepend-inner="console.log('CLICK')"
        />
      </template>

      <template #bottom>
        <div style="position: sticky; bottom: 0px">
          <v-divider class="w-100"></v-divider>
          <div class="v-data-table-footer w-100 justify-space-between">
            <div>
              <v-btn
                color="primary"
                variant="plain"
                prepend-icon="mdi-plus"
                text="Add row"
                @click="items.push({})"
              />
              <v-btn
                v-if="selection"
                color="warning"
                variant="plain"
                prepend-icon="mdi-close"
                text="Delete row(s)"
                @click="items = items.filter((items, row) => !isSelected(row))"
              />
            </div>
            <v-btn
              color="warning"
              variant="plain"
              prepend-icon="mdi-close"
              text="Remove empty rows"
              @click="items = items.filter((item) => !isEmpty(item))"
            />
          </div>
          <v-divider class="mb-3"></v-divider>
        </div>
      </template>
    </v-data-table>
  </div>
</template>

<script setup lang="ts">
import { $Coordinates, Coordinates } from '@/api'
import { useSpreadsheet } from '.'
import { Errors } from '../toolkit/validation'
import CoordPrecisionPicker from './CoordPrecisionPicker.vue'
import ErrorTooltip from './ErrorTooltip.vue'
import { ref } from 'vue'

const s = ref()
const sites = ref(['kanar'])

type Item = Partial<Coordinates> & { errors?: Errors<ObjectPaths<Coordinates>> }

const {
  items,
  dragging,
  selection,
  edit,
  onEdited,
  handlePaste,
  cellHeader,
  select,
  selectEnd,
  isEditing,
  isSelected,
  isEmpty
} = useSpreadsheet<Item>(setValue, $Coordinates)

const headers: DataTableHeader[] = [
  cellHeader({
    key: 'latitude',
    title: 'Latitude'
  }),
  cellHeader({ key: 'longitude', title: 'Longitude' }),
  cellHeader({
    key: 'precision',
    title: 'Precision',
    class: 'right-separator',
    headerProps: { class: 'border-e-lg' }
  }),
  { key: 'site_distance', title: 'Nearest site' },
  { key: 'actions', title: 'Actions' }
]

function setValue(row: number, col: number, value: any) {
  items.value[row][headers[col].key as keyof Coordinates] = value
}
</script>

<style lang="scss">
@use 'vuetify';
.spreadsheet {
  table {
    border-collapse: collapse;
    td.cell {
      border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
      user-select: none;
      box-sizing: border-box;
      &.error-cell {
        background-color: rgba(var(--v-theme-error), 0.2);
      }
      &.right-separator {
        border-right-width: 4px;
      }
      &.selected {
        outline: 1px dashed rgba(var(--v-theme-primary), 0.6);
        outline-offset: -1px;
        &.right {
          border-right: 2px solid rgb(var(--v-theme-primary));
          &.right-separator {
            border-right-width: 4px;
          }
        }
        &.left {
          border-left: 2px solid rgb(var(--v-theme-primary)) !important;
        }
        &.top {
          border-top: 2px solid rgb(var(--v-theme-primary)) !important;
        }
        &.bottom {
          border-bottom: 2px solid rgb(var(--v-theme-primary)) !important;
        }
      }
      &.editing {
        background-color: rgb(var(--v-theme-surface-light));
      }
    }
  }
}
</style>
