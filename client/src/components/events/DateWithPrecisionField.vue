<template>
  <v-input>
    <v-select
      :label="model.precision === 'Unknown' ? 'Date' : 'Precision'"
      v-model="model.precision"
      :items="$DatePrecision.enum"
      :rounded="model.precision === 'Unknown' ? undefined : 'e-0'"
      v-bind="field('precision')"
      :max-width="150"
    />
    <v-text-field
      v-if="model.precision !== 'Unknown'"
      label="Date"
      rounded="s-0"
      v-model="model.date"
      v-maska="{
        mask: options.mask,
        eager: true
      }"
      :placeholder="options.placeholder"
      persistent-placeholder
      v-bind="field('date')"
      :rules="[
        (value: string) => {
          if (value && model.precision !== 'Unknown') {
            const parsed = DateTime.fromFormat(value, dateFormats[model.precision])
            return parsed.isValid || 'Invalid date'
          }
          return true
        }
      ]"
    >
    </v-text-field>
  </v-input>
</template>

<script setup lang="ts">
import { $DatePrecision, $DateWithPrecision, DatePrecision, DateWithPrecision } from '@/api'
import { vMaska } from 'maska/vue'
import { computed, ref } from 'vue'
import { patternRule, useSchema } from '../toolkit/forms/schema'
import { DateTime } from 'luxon'
import { parseTwoDigitYear } from 'moment'

const model = ref<DateWithPrecision>({ date: undefined, precision: 'Day' })

const { field } = useSchema($DateWithPrecision)

const dateFormats: Record<DatePrecision, string> = {
  Day: 'dd-MM-yyyy',
  Month: 'MM-yyyy',
  Year: 'yyyy',
  Unknown: ''
}

const options = computed(() => {
  switch (model.value.precision) {
    case 'Day':
      return {
        placeholder: 'DD-MM-YYYY',
        mask: '##-##-####'
      }
    case 'Month':
      return {
        placeholder: 'MM-YYYY',
        mask: '##-####'
      }
    case 'Year':
      return {
        placeholder: 'YYYY',
        mask: '####'
      }
    default:
      return {
        placeholder: 'Disabled',
        mask: undefined
      }
  }
})
</script>

<style lang="scss"></style>
