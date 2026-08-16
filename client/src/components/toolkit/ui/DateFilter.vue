<template>
  <v-date-input
    v-if="precision === 'day'"
    label="Date"
    density="compact"
    hide-details
    input-format="yyyy-mm-dd"
    :min="minDateTime.toJSDate()"
    :max="maxDateTime.toJSDate()"
    clearable
    :model-value="CompositeDate.toDateTime(model, precision)?.toJSDate()"
    @update:model-value="
      (v) => {
        if (!v) {
          model = {}
          return
        }
        console.log(v)
        const d = DateTime.fromJSDate(v)
        model = {
          day: d.day,
          month: d.month,
          year: d.year
        }
        console.log(model)
      }
    "
  ></v-date-input>
  <div class="d-flex align-center ga-2" v-else>
    <date-picker-year
      :model-value="model.year"
      @update:model-value="
        (year) => {
          console.log('Year updated:', year)
          model = { ...model, year }
        }
      "
      :min="minDateTime.year"
      :max="maxDateTime.year"
    />
    <!-- <v-menu>
      <template #activator="{ props }">
        <v-number-input
          label="Year"
          :model-value="model.year"
          @update:model-value="(year) => (model = { ...model, year })"
          @click:clear="
            () => (precision === 'year' ? (model = {}) : (model = { ...model, year: undefined }))
          "
          v-bind="props"
          append-inner-icon="mdi-chevron-down"
          control-variant="split"
          density="compact"
          :min="minDateTime.year"
          :max="maxDateTime.year"
          hide-details
          :step="1"
          clearable
        ></v-number-input>
      </template>
      <v-card>
        <v-date-picker-years
          :key="model.year ?? 'empty'"
          :model-value="model.year"
          @update:model-value="(year) => (model = { ...model, year })"
          :min="minDateTime.toJSDate()"
          :max="maxDateTime.toJSDate()"
          data-allow-mismatch
        />
      </v-card>
    </v-menu> -->
    <v-menu v-if="precision === 'month'">
      <template #activator="{ props }">
        <v-text-field
          v-bind="props"
          label="Month"
          type="month"
          readonly
          :model-value="
            DateTime.fromObject(
              { year: DateTime.now().year, month: model.month || 1 },
              { locale: 'en' }
            ).monthShort
          "
          append-inner-icon="mdi-chevron-down"
          density="compact"
          hide-details
          class="no-caret"
        />
      </template>
      <v-card>
        <v-date-picker-months
          :model-value="model.month ? model.month - 1 : undefined"
          @update:model-value="(v) => (model.month = v + 1)"
        />
      </v-card>
    </v-menu>
  </div>
</template>

<script setup lang="ts">
import { CompositeDate, EventDatePrecision } from '@/api'
import { DateTime } from 'luxon'
import { computed } from 'vue'
import DatePickerYear from './DatePickerYear.vue'

const { minDate, maxDate } = defineProps<{
  precision: EventDatePrecision
  minDate?: CompositeDate
  maxDate?: CompositeDate
}>()

const DEFAULT_MIN_DATE = DateTime.fromObject({ year: 1800, month: 1, day: 1 })
const DEFAULT_MAX_DATE = DateTime.now()

const minDateTime = computed(() => {
  if (!minDate || !minDate.year) return DEFAULT_MIN_DATE
  return DateTime.max(
    DEFAULT_MIN_DATE,
    DateTime.fromObject({
      year: minDate.year,
      month: minDate.month ?? 1,
      day: minDate.day ?? 1
    })
  )
})

const maxDateTime = computed(() => {
  if (!maxDate || !maxDate.year) return DEFAULT_MAX_DATE

  return DateTime.min(
    DEFAULT_MAX_DATE,
    DateTime.fromObject({
      year: maxDate.year,
      month: maxDate.month ?? 12,
      day: maxDate.day ?? 31
    })
  )
})

const model = defineModel<CompositeDate>({ default: () => ({}) })
</script>

<style lang="scss">
.no-caret input {
  caret-color: transparent;
  cursor: pointer;
}
</style>
