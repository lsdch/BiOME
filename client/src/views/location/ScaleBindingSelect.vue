<template>
  <v-select
    v-model="model.binding"
    :label
    :items
    item-value="value"
    item-title="title"
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
import { SiteWithOccurrences } from '@/api'
import { ScaleBindingSpec, useScaleBinding } from '@/composables/occurrences'
import { onMounted, reactive, watch } from 'vue'
import { ScaleBinding } from 'vue-leaflet-hexbin'

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

const emit = defineEmits<{
  updateFn: [binding?: ScaleBinding<SiteWithOccurrences>]
}>()

watch(
  () => model.value,
  (newValue) => {
    if (!newValue.binding) model.value.log = false
    console.log('update')
    emit('updateFn', useScaleBinding(newValue))
  },
  { deep: true }
)

const items = [
  { title: 'Sites', value: 'sites' },
  {
    title: 'Occurrences',
    value: 'occurrences'
  }
  // {
  //   title: 'Species',
  //   value: (d) =>
  //     Object.values(
  //       d.reduce<SampledTaxa>(
  //         (acc, { data }) => ({ ...acc, ...occurringTaxa(data.occurrences) }),
  //         occurringTaxa([])
  //       )
  //     ).reduce((acc, taxon) => ({...acc, ...taxon})).length
  // }
] as const
</script>

<style scoped lang="scss"></style>
