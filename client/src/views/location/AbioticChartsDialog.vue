<template>
  <CardDialog v-if="abiotic_data" title="Abiotic measurements" v-bind="$attrs" :fullscreen="mobile">
    <v-card-text>
      <v-select
        label="Parameter"
        v-model="item"
        :items="abiotic_data"
        item-title="param.label"
        return-object
      >
        <template #item="{ item: { title, value, raw }, props }">
          <v-list-item :title :value v-bind="props" :subtitle="raw.param.unit">
            <template #append>
              <v-chip :text="raw.points.length.toString()" />
            </template>
          </v-list-item>
        </template>
      </v-select>
      <AbioticLineChart v-if="item" :data="item" />
    </v-card-text>
  </CardDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import AbioticLineChart, { AbioticData } from './AbioticLineChart.vue'
import CardDialog from '@/components/toolkit/forms/CardDialog.vue'
import { useDisplay } from 'vuetify'

const props = defineProps<{ abiotic_data: AbioticData[] }>()

const item = ref<AbioticData>()

const { mobile } = useDisplay()

watch(
  () => props.abiotic_data,
  (items) => (item.value = items[0] ?? undefined),
  { immediate: true }
)
</script>

<style scoped></style>
