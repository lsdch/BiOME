<template>
  <tr>
    <td v-for="i in Array(offset ?? 0).keys()" :key="i"></td>
    <td colspan="1000">
      <div class="site-expanded-row">
        <div class="label">Coord. precision:</div>
        <ErrorTooltip :error="errors?.['coordinates.precision']">
          {{ item.coordinates?.precision ?? 'NA' }}
        </ErrorTooltip>
        <div class="label">Locality:</div>
        <div class="d-flex align-center">
          <ErrorTooltip :error="errors?.locality">
            {{ item.locality ?? 'NA' }}
          </ErrorTooltip>
          <ErrorTooltip :error="errors?.country_code">
            <v-chip class="mx-3 text-overline" size="small">
              {{ item.country_code || ' ? ' }}
            </v-chip>
          </ErrorTooltip>
        </div>
        <div class="label">Description:</div>
        <ErrorTooltip :error="errors?.description">
          {{ item.description ?? 'NA' }}
        </ErrorTooltip>
      </div>
    </td>
  </tr>
</template>

<script setup lang="ts">
import { CoordinatesPrecision } from '@/api'
import ErrorTooltip from './ErrorTooltip.vue'

type Item = {
  coordinates: {
    precision?: CoordinatesPrecision
  }
  locality?: string
  country_code?: string
  description?: string
}

defineProps<{
  offset?: number
  errors?: {
    [K in ObjectPaths<Item>]?: string
  }
  item: Item
}>()
</script>

<style scoped lang="scss">
.site-expanded-row {
  display: grid;
  grid-template-columns: 0fr 1fr;
  gap: 10px;
  margin-top: 10px;
  margin-bottom: 10px;
  align-items: center;
  .label {
    text-wrap: nowrap;
    text-align: right;
  }
}
</style>
