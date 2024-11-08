<template>
  <v-select label="Parameter" :items :loading v-model="model" item-title="label" v-bind="$attrs">
    <template #item="{ item, props }">
      <v-list-item v-bind="props" :subtitle="item.raw.description">
        <template #append>
          <v-chip :text="item.raw.unit" />
        </template>
      </v-list-item>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { AbioticParameter, EventsService, SamplingService } from '@/api'
import { handleErrors } from '@/api/responses'
import { useToggle } from '@vueuse/core'
import { onMounted, ref } from 'vue'

const model = ref<AbioticParameter>()

const items = ref<AbioticParameter[]>([])
const [loading, toggleLoading] = useToggle(true)
onMounted(async () => {
  items.value = await SamplingService.listAbioticParameters()
    .then(handleErrors((err) => console.error('Failed to fetch abiotic parameters', err)))
    .finally(() => toggleLoading(false))
})
</script>

<style lang="scss" scoped></style>
