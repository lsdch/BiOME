<template>
  <v-input
    :variant
    class="composite-date-field v-text-field"
    clearable
    persistent-clear
    :focused="focused.day.value || focused.month.value || focused.year.value"
    :disabled="precision === 'Unknown'"
    :rules="[
      () => precision === 'Unknown' || !!model.year || `Date is required`,
      () => !precisionAtLeast('Month') || !!model.month || `Incomplete date`,
      () => precision !== 'Day' || !!model.day || `Incomplete date`,
      () =>
        DateTime.fromObject(model).isValid ||
        `Invalid date : ${DateTime.fromObject(model).invalidExplanation}`
    ]"
    validate-on="input"
    :hint
    :error-messages
    v-model="model"
  >
    <template #default="{ id, validate, isDirty }">
      <v-field
        :rounded
        :variant
        :focused="focused.day.value || focused.month.value || focused.year.value"
        active
        @click="focused.day.value = true"
        label="Date"
        :clearable="true"
        persistent-clear
        :dirty="isDirty.value"
      >
        <!-- :dirty="!!(model.day || model.month || model.year)" -->
        <template #="{}">
          <div class="v-field__input">
            <span v-show="precision === 'Unknown'" class="text-grey">Unknown</span>
            <div
              v-show="precisionAtLeast('Day')"
              :class="['padded-input', { pad: model.day && model.day < 10 }]"
            >
              <input
                v-model.number="model.day"
                :ref="inputs.day"
                class="v-field-input"
                placeholder="DD"
                type="text"
                inputmode="numeric"
                :id="`${id.value}-day`"
                :size="2"
                :max="31"
                :maxlength="2"
                @input="validate()"
                @keydown.up.prevent="increment('day', 1, 31)"
                @keydown.down.prevent="decrement('day', 1, 31)"
                @keyup.right="(e) => handleArrowRight('day')(e)"
                @keydown="(e) => handleKeyDown('day')(e)"
                @click.stop
              />
              <span class="mx-1">-</span>
            </div>

            <div
              v-show="precisionAtLeast('Month')"
              :class="['padded-input', { pad: model.month && model.month < 10 }]"
            >
              <input
                v-model.number="model.month"
                :ref="inputs.month"
                class="v-field-input"
                placeholder="MM"
                :size="2"
                type="text"
                inputmode="numeric"
                :id="`${id.value}-month`"
                :max="12"
                :min="1"
                :maxlength="2"
                @input="validate()"
                @keydown.left="(e) => handleArrowLeft('month')(e)"
                @keydown.right="(e) => handleArrowRight('month')(e)"
                @keydown.up.prevent="increment('month', 1, 12)"
                @keydown.down.prevent="decrement('month', 1, 12)"
                @keydown="(e) => handleKeyDown('month')(e)"
                @click.stop
              />
              <span class="mx-1">-</span>
            </div>
            <div v-show="precisionAtLeast('Year')" :class="['padded-input']">
              <input
                v-model.number="model.year"
                :ref="inputs.year"
                placeholder="YYYY"
                type="text"
                class="year-input v-field-input"
                inputmode="numeric"
                :id="`${id.value}-year`"
                :maxlength="4"
                :min="1500"
                :max="3000"
                @input="validate()"
                @keydown.up.prevent="
                  increment(
                    'year',
                    $CompositeDate.properties.year.minimum,
                    $CompositeDate.properties.year.maximum,
                    DateTime.now().year
                  )
                "
                @keydown.down.prevent="
                  decrement(
                    'year',
                    $CompositeDate.properties.year.minimum,
                    $CompositeDate.properties.year.maximum,
                    DateTime.now().year
                  )
                "
                @keydown.left="(e) => handleArrowLeft('year')(e)"
                @keydown="(e) => handleKeyDown('year')(e)"
                @click.stop
              />
            </div>
          </div>
        </template>
        <template #clear="{}">
          <v-icon color="" icon="mdi-close-circle" density="compact" @click="model = {}" />
        </template>
      </v-field>
    </template>
  </v-input>
</template>

<script setup lang="ts">
import { $CompositeDate, $DatePrecision, DatePrecision } from '@/api'
import { CompositeDate } from '@/api/adapters'
import { useFocus } from '@vueuse/core'
import { DateTime } from 'luxon'
import { ref } from 'vue'
import { VTextField } from 'vuetify/components'

const model = defineModel<CompositeDate>({ required: true })

const props = withDefaults(
  defineProps<{
    precision?: DatePrecision
    clearable?: boolean
    variant?: VTextField['$props']['variant']
    rules?: VTextField['$props']['rules']
    hint?: VTextField['$props']['hint']
    errorMessages?: VTextField['$props']['errorMessages']
    rounded?: boolean | string | number
  }>(),
  { variant: 'outlined', precision: 'Day', hint: '⇅ : increment/decrement' }
)

const inputs = {
  day: ref<HTMLInputElement>(),
  month: ref<HTMLInputElement>(),
  year: ref<HTMLInputElement>()
}
const focused = {
  day: useFocus(inputs.day, { initialValue: false }).focused,
  month: useFocus(inputs.month, { initialValue: false }).focused,
  year: useFocus(inputs.year, { initialValue: false }).focused
}

function precisionAtLeast(precision: DatePrecision) {
  return $DatePrecision.enum.indexOf(precision) >= $DatePrecision.enum.indexOf(props.precision)
}

type Subfield = 'day' | 'month' | 'year'

function increment(field: Subfield, min: number, max: number, initial?: number) {
  model.value[field] = Math.min((model.value[field] || (initial ?? min) - 1) + 1, max)
}

function decrement(field: Subfield, min: number, max: number, initial?: number) {
  model.value[field] = Math.max((model.value[field] || (initial ?? max) + 1) - 1, min)
}

function focusNext() {
  if (focused.day.value) {
    focused.month.value = true
    inputs.month.value!.selectionStart = 0
  } else if (focused.month.value) {
    focused.year.value = true
    inputs.year.value!.selectionStart = 0
  }
}

function focusPrev() {
  if (focused.month.value) {
    focused.day.value = true
    inputs.day.value!.selectionStart = inputs.day.value!.value.length
  }
  if (focused.year.value) {
    focused.month.value = true
    inputs.month.value!.selectionStart = inputs.month.value!.value.length
  }
}

function isCursorAtStart(input: HTMLInputElement) {
  return input.selectionStart === 0
}

function isCursorAtEnd(input: HTMLInputElement) {
  return input.selectionStart === input.value.length
}

function handleArrowRight(field: 'day' | 'month' | 'year') {
  return (event: KeyboardEvent) => {
    if (isCursorAtEnd(inputs[field].value!)) {
      focusNext()
      event.preventDefault()
      event.stopImmediatePropagation()
    }
  }
}

function handleArrowLeft(field: 'day' | 'month' | 'year') {
  return (event: KeyboardEvent) => {
    if (isCursorAtStart(inputs[field].value!)) {
      focusPrev()
      event.preventDefault()
      event.stopImmediatePropagation()
    }
  }
}

function handleKeyDown(field: 'day' | 'month' | 'year') {
  return (event: KeyboardEvent) => {
    if (event.key === 'ArrowLeft') return
    if (event.altKey || event.shiftKey || event.ctrlKey || event.metaKey) return
    if (event.key === '-' && focused[field].value && model.value[field]) focusNext()
    if (event.key === 'Backspace' && !model.value[field]) {
      event.preventDefault()
      focusPrev()
    }
    if (event.key.length > 1) return
    if (event.key > '9' || event.key < '0') {
      event.preventDefault()
    }
  }
}
</script>

<style scoped lang="scss">
.composite-date-field {
  .padded-input {
    flex-grow: 0;
    flex-shrink: 0;
    flex-basis: content;
    position: relative;
    &.pad::before {
      content: '0';
      position: absolute;
      left: 0;
      pointer-events: none;
    }
  }
  &.v-theme--dark {
    input:focus {
      background-color: rgb(68, 68, 68);
      border-radius: 3px;
    }
  }
  &.v-theme--light {
    input:focus {
      background-color: rgb(214, 214, 214);
      border-radius: 3px;
    }
  }
  input {
    text-align: right;
    font-family: monospace;

    &.year-input {
      text-align: left;
    }
  }
}
</style>
