<template>
  <fieldset>
    <legend>{{ props.label }}</legend>
    <v-input :model-value v-bind="$attrs">
      <template #default="{ isDisabled, isValid }">
        <v-number-input
          label="Hours"
          class="flex-grow-1"
          v-model="hours"
          control-variant="stacked"
          rounded="e-0"
          @update:model-value="onInput()"
          :min="0"
          hide-details
          :disabled="isDisabled.value"
          :error="!isValid.value"
        />
        <v-number-input
          label="Minutes"
          class="flex-grow-1"
          v-model="minutes"
          control-variant="stacked"
          rounded="s-0"
          @update:model-value="onMinutesInput"
          @update:focused="onMinutesInput(minutes)"
          :min="0"
          hide-details
          :disabled="isDisabled.value"
          :error="!isValid.value"
        />
      </template>
      <template #append>
        <v-btn
          v-if="clearable && (hours || minutes)"
          @click="clear"
          color=""
          variant="tonal"
          icon="mdi-close"
          density="compact"
        />
      </template>
    </v-input>
  </fieldset>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  modelValue: number | undefined | null
  label?: string
  clearable?: boolean
}>()
const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const hours = ref()
const minutes = ref()

function clear() {
  hours.value = undefined
  minutes.value = undefined
}

function onMinutesInput(v: number) {
  if (v >= 60) {
    hours.value += Math.floor(v / 60)
    minutes.value = v % 60
  }
  onInput()
}

function onInput() {
  if (hours.value === undefined && minutes.value === undefined) {
    emit('update:modelValue', null)
  } else {
    emit('update:modelValue', (hours.value ?? 0) * 60 + (minutes.value ?? 0))
  }
}

watch(
  () => props.modelValue,
  (v) => {
    if (v === undefined || v === null) {
      hours.value = v
      minutes.value = v
    } else {
      hours.value = Math.floor(v / 60)
      minutes.value = v % 60
    }
  },
  { immediate: true }
)
</script>

<style scoped lang="scss">
fieldset {
  border: none;
}
</style>
