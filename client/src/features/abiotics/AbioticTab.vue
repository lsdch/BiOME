<template>
  <AbioticMeasurementsTable :measurements="abiotic_measurements" class="mb-4" />

  <div class="d-flex justify-space-between align-center">
    <v-chip-group v-model="active_param" color="primary">
      <v-chip v-for="param in sorted_abiotic_params" :key="param.code" :value="param">
        {{ param.label }} ({{ param.unit }})
        <v-badge
          v-if="abiotic_data[param.code]?.points.length"
          :content="abiotic_data[param.code]?.points.length"
          class="ml-2"
          inline
        />
      </v-chip>
    </v-chip-group>
    <v-btn text="Add measurement" prepend-icon="mdi-plus" rounded="md"></v-btn>
  </div>
  <AbioticLineChart
    v-if="active_param && abiotic_data[active_param?.code]"
    :data="abiotic_data[active_param.code]"
  />
</template>

<script setup lang="ts">
import { AbioticMeasurement, AbioticParameter } from '@/api'
import { listAbioticParametersOptions } from '@/api/gen/@tanstack/vue-query.gen'
import AbioticMeasurementsTable from '@/features/occurrences/components/AbioticMeasurementsTable.vue'
import AbioticLineChart, {
  AbioticData,
  AbioticDataPoint
} from '@/features/site/components/AbioticLineChart.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, ref } from 'vue'

const { abiotic_measurements } = defineProps<{ abiotic_measurements: AbioticMeasurement[] }>()

const active_param = ref<AbioticParameter>()

const { data: abiotic_params, error, isPending } = useQuery(listAbioticParametersOptions())

const sorted_abiotic_params = computed(() => {
  return (abiotic_params.value ?? []).sort(
    (a, b) =>
      (abiotic_data.value[a.code]?.points.length ?? 0) -
      (abiotic_data.value[b.code]?.points.length ?? 0)
  )
})

const abiotic_data = computed(() => {
  return (
    abiotic_measurements?.reduce<Record<string, AbioticData>>(
      (acc, { performed_on, param, value }) => {
        if (performed_on?.date === undefined) return acc
        acc[param.code] = {
          param,
          points: [{ y: value, date: performed_on.date }].concat(
            acc[param.code]?.points ?? Array<AbioticDataPoint>()
          )
        }
        return acc
      },
      {}
    ) ?? {}
  )
})
</script>

<style scoped lang="scss"></style>
