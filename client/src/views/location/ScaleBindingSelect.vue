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
    <template #append>
      <v-btn-toggle
        v-model="model.log"
        :color="model.binding ? 'success' : ''"
        size="small"
        :density
        mandatory
        divided
        rounded="sm"
        border="sm"
        :disabled="!model.binding"
      >
        <v-btn class="text-caption" :density size="small" text="Linear" :value="false" />
        <v-btn class="text-caption" :density size="small" text="Log" :value="true" />
        <v-btn class="text-caption" :density size="small" :value="10"> Log<sub>10</sub> </v-btn>
      </v-btn-toggle>
    </template>
  </v-select>
</template>

<script setup lang="ts">
import { ScaleBindingSpec } from '@/composables/occurrences'
import { onMounted, reactive } from 'vue'

const { clearable } = defineProps<{
  label?: string
  clearable?: boolean
  density?: 'default' | 'compact' | 'comfortable'
}>()

const model = defineModel<ScaleBindingSpec>({
  default: () =>
    reactive({
      log: false,
      binding: undefined
    })
})

onMounted(() => {
  if (!clearable) {
    model.value.binding = items[0].value
  }
})

const items = [
  { title: 'Sites', value: 'sites' },
  { title: 'Sampling events', value: 'samplings' },
  {
    title: 'Occurrences',
    value: 'occurrences'
  },
  {
    title: 'Species richness',
    value: 'speciesRichness',
    subtitle: 'Includes species and sub-species'
  },
  { title: 'Genus richness', value: 'genusRichness' },
  { title: 'Family richness', value: 'familyRichness' }
] as const
</script>

<style scoped lang="scss"></style>
