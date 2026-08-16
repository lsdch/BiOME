<template>
  <v-menu :width="400">
    <template #activator="{ props }">
      <v-number-input
        label="Year"
        v-model="model"
        v-bind="props"
        append-inner-icon="mdi-chevron-down"
        control-variant="split"
        density="compact"
        :min="minYear"
        :max="maxYear"
        hide-details
        :step="1"
        clearable
        :loading="isLoading"
      ></v-number-input>
    </template>
    <v-card :max-height="400">
      <div class="d-flex ga-3 align-center flex-wrap">
        <v-btn
          v-for="year in range(minYear, maxYear)"
          :key="year"
          :variant="year === model ? 'elevated' : 'text'"
          @click="model = year"
        >
          {{ year }}
        </v-btn>
      </div>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { listSamplingYearsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { DateTime } from 'luxon'
import { computed } from 'vue'

const model = defineModel<number>()

const { min = 1800, max = DateTime.now().year } = defineProps<{
  min?: number
  max?: number
}>()

function* range(from: number, to: number) {
  for (let i = from; i <= to; i++) {
    yield i
  }
}

const minYear = computed(() => Math.max(years.value?.[0] ?? min, min))
const maxYear = computed(() => Math.min(years.value?.[years.value.length - 1] ?? max, max))

const { data: years, isPending: isLoading } = useQuery(listSamplingYearsOptions())
</script>

<style scoped lang="scss"></style>
