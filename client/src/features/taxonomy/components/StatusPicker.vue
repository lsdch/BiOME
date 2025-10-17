<template>
  <v-select
    v-model="model"
    :items
    label="Status"
    item-value="status"
    item-title="status"
    v-bind="$attrs"
  >
    <template #item="{ item, props }">
      <v-list-item :title="item.value" :subtitle="item.raw.description" v-bind="props">
        <template #prepend>
          <FTaxonStatusIndicator :status="item.value" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { $TaxonStatus, TaxonStatus } from '@/api'
import { FTaxonStatusIndicator, taxonStatusIndicatorProps } from './functionals'
import { computed } from 'vue'
const model = defineModel<TaxonStatus>()

const props = defineProps<{ omit?: TaxonStatus | TaxonStatus[] }>()

const items = computed(() =>
  $TaxonStatus.enum
    .filter((status) =>
      props.omit
        ? typeof props.omit == 'string'
          ? status == props.omit
          : props.omit.includes(status)
        : true
    )
    .map((status) => ({
      status,
      ...taxonStatusIndicatorProps(status)
    }))
)
</script>

<style scoped></style>
