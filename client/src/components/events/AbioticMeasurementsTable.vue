<template>
  <CRUDTable :headers :items="measurements" entity-name="Abiotic measurement">
    <template #item.value="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>

    <template #item.param.unit="{ value }">
      <span class="font-monospace">
        {{ value }}
      </span>
    </template>
    <template #item.performed_on="{ value }">
      <span class="font-monospace">
        {{ DateWithPrecision.format(value) }}
      </span>
    </template>
  </CRUDTable>
</template>

<script setup lang="ts">
import { AbioticMeasurement, DateWithPrecision } from '@/api'
import CRUDTable from '../toolkit/tables/CRUDTable.vue'

defineProps<{ measurements: AbioticMeasurement[] }>()

const headers: CRUDTableHeader<AbioticMeasurement>[] = [
  { title: 'Parameter', key: 'param.label' },
  { title: 'Value', key: 'value' },
  { title: 'Unit', key: 'param.unit' },
  { title: 'Date', key: 'performed_on' }
]
</script>

<style scoped lang="scss"></style>
