<template>
  <v-input :model-value v-bind="$attrs">
    <template #default="{ isDisabled, isValid }">
      <v-number-input
        label="Days"
        class="flex-grow-1"
        v-model="days"
        control-variant="stacked"
        rounded="e-0"
        @update:model-value="onInput()"
        :min="0"
        hide-details
        :disabled="isDisabled.value"
        :error="!isValid.value"
      />
      <v-number-input
        label="Hours"
        class="flex-grow-1"
        v-model="hours"
        control-variant="stacked"
        rounded="s-0"
        @update:model-value="onHoursInput"
        @update:focused="onHoursInput(hours)"
        :min="0"
        hide-details
        :disabled="isDisabled.value"
        :error="!isValid.value"
      />
    </template>
  </v-input>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

const props = defineProps<{ modelValue: number }>()
const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const days = ref()
const hours = ref()

function onHoursInput(v: number) {
  if (v >= 24) {
    days.value += Math.floor(v / 24)
    hours.value = v % 24
  }
  onInput()
}

function onInput() {
  emit('update:modelValue', (days.value ?? 0) * 24 + (hours.value ?? 0))
}

watch(
  () => props.modelValue,
  (v) => {
    days.value = Math.floor(v / 24)
    hours.value = v % 24
  },
  { immediate: true }
)
</script>

<style scoped></style>
