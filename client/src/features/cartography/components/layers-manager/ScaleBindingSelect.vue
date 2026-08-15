<template>
  <v-select
    v-model="model.binding"
    :label
    :items
    item-value="value"
    item-title="title"
    item-props
    :clearable
    :density
    v-bind="$attrs"
  >
    <template #append v-if="model.binding !== 'constant'">
      <!-- :disabled="model.binding === 'constant'" -->
      <v-btn
        icon="mdi-math-log"
        :variant="model.log ? 'elevated' : 'text'"
        size="small"
        rounded="md"
        :color="model.log ? 'success' : ''"
        :active="!!model.log"
        @click="model.log = !model.log"
        v-tooltip="{ text: 'Toggle log scale', openDelay: 500 }"
      />
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { ScaleBindingSpec } from '@/features/cartography/bindings'

const { clearable } = defineProps<{
  label?: string
  clearable?: boolean
  density?: 'default' | 'compact' | 'comfortable'
}>()

const model = defineModel<ScaleBindingSpec>({
  required: true
})

const items = [
  { title: 'Constant', value: 'constant', subtitle: 'Uniform color scale' },
  {
    title: 'Sampling events',
    value: 'samplings',
    subtitle: 'Scale with the number of sampling events'
  },
  { title: 'Occurrences', value: 'occurrences', subtitle: 'Scale with the number of occurrences' },
  {
    title: 'Species richness',
    value: 'speciesRichness',
    subtitle: 'Scale with the number of species and sub-species'
  },
  {
    title: 'Genus richness',
    value: 'genusRichness',
    subtitle: 'Scale with the number of genus groups'
  },
  {
    title: 'Family richness',
    value: 'familyRichness',
    subtitle: 'Scale with the number of family groups'
  }
] as const
</script>

<style scoped lang="scss"></style>
