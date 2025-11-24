<template>
  <div>
    <div class="d-flex justify-space-between align-center ml-2">
      Specimen quantity
      <div class="d-flex ga-1 align-center justify-end">
        <v-chip-group selected-class="text-primary" column mandatory v-model="pick">
          <v-chip
            class="justify-end"
            size="small"
            filter
            value="Exact"
            text="Exact"
            @click="model = undefined"
          />
          <v-chip
            v-for="item in predefined"
            size="small"
            filter
            :value="item.key"
            :text="item.label"
            @click="model = item.value"
          />
          <v-chip
            class="justify-end"
            size="small"
            filter
            value="Custom"
            text="Custom range"
            @click="model = typeof model === 'object' ? model : {}"
          />
          <v-chip
            class="justify-end"
            size="small"
            filter
            value="Unknown"
            text="Unknown"
            @click="model = undefined"
          />
        </v-chip-group>
      </div>
    </div>
    <div v-if="['Exact', 'Custom'].includes(pick ?? '')">
      <v-number-input
        v-if="pick === 'Exact'"
        label="Quantity"
        class="required mt-2"
        :rules="[(v) => v || 'Strictly positive quantity required']"
        @update:model-value="(v) => (model = v)"
      />
      <div v-else-if="pick === 'Custom'" class="d-flex mt-2">
        <v-number-input
          label="Lower bound"
          v-model="(model as SpecimenQuantityRangeModel).lower"
          class="required rounded-e-0"
          :step="1"
          :min="1"
          :rules="[(v) => !!v || 'Strictly positive value required']"
        />
        <v-number-input
          label="Upper bound"
          v-model="(model as SpecimenQuantityRangeModel).upper"
          class="required rounded-s-0"
          :step="1"
          :min="((model as SpecimenQuantityRangeModel).lower ?? 0) + 1"
          :rules="[
            (v) =>
              v > ((model as SpecimenQuantityRangeModel).lower ?? 0) ||
              'Upper bound must be greater than lower bound'
          ]"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { SpecimenQuantityModel, SpecimenQuantityRangeModel } from '@/models/biomat'
import { ref } from 'vue'
const model = defineModel<SpecimenQuantityModel>()
const pick = ref<'Exact' | 'Custom' | 'Unknown' | (typeof predefined)[number]['key']>()
const predefined = [
  { label: 'One', key: 'One', value: 1 },
  { label: 'Few (2-5)', key: 'Few', value: { lower: 2, upper: 5 } },
  { label: 'Several (6-20)', key: 'Several', value: { lower: 6, upper: 20 } },
  { label: 'Many (21-100)', key: 'Many', value: { lower: 21, upper: 100 } },
  { label: 'Numerous (100-1000)', key: 'Numerous', value: { lower: 101, upper: 1000 } }
] as const
</script>

<style scoped lang="scss"></style>
