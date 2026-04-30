<template>
  <v-card class="small-card-title" v-bind="$attrs">
    <template #prepend>
      <v-icon icon="mdi-hexagon-multiple" class="text-muted" />
    </template>
    <template #title>
      <v-text-field
        v-model="layer.name"
        density="compact"
        hide-details
        :placeholder="title"
        variant="underlined"
      />
    </template>
    <template #append>
      <v-progress-circular
        v-if="remote.isPending.value"
        indeterminate
        size="small"
        color="warning"
      ></v-progress-circular>
      <v-switch v-model="layer.active" color="primary" hide-details></v-switch>
      <v-btn
        :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
        color=""
        variant="text"
        :rounded="100"
        size="small"
        @click="() => (expanded = !expanded)"
      />
    </template>
    <v-expand-transition>
      <div v-show="expanded" class="">
        <v-divider />
        <v-tabs v-model="tab" density="compact">
          <v-tab prepend-icon="mdi-database" value="data">Data</v-tab>
          <v-tab prepend-icon="mdi-hexagon-multiple-outline" value="layer">Style</v-tab>
          <v-tab prepend-icon="mdi-chart-bar" value="stats" :disabled="remote.isFetching.value">
            Stats
          </v-tab>
        </v-tabs>
        <v-divider></v-divider>
        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="data">
            <v-confirm-edit v-model="layer.filters" @save="register()">
              <template #default="{ model, actions, isPristine }">
                <LayerDataFeed v-model="model.value" class="pa-2" />
                <div v-if="!isPristine" class="d-flex justify-end">
                  <component :is="actions"></component>
                </div>
              </template>
            </v-confirm-edit>
          </v-tabs-window-item>
          <v-tabs-window-item value="layer">
            <HexgridLayerStylePanel v-model="layer" />
          </v-tabs-window-item>
          <v-tabs-window-item value="stats">
            <OccurrencesStats :sites="remote.data.value" />
          </v-tabs-window-item>
        </v-tabs-window>
        <!-- <v-list-item prepend-icon="mdi-shape-polygon-plus" title="Polygon">
          <template #append>
            <v-btn icon="mdi-shape-polygon-plus" variant="text" />
            <v-btn icon="mdi-eye-outline" variant="text" />
            <v-switch color="primary" hide-details></v-switch>
          </template>
        </v-list-item> -->
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script setup lang="ts">
import { occurrencesBySiteOptions } from '@/api/gen/@tanstack/vue-query.gen'
import OccurrencesStats from '@/features/occurrences/components/OccurrencesStats.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed, onMounted, ref } from 'vue'
import HexgridLayerStylePanel from './HexgridLayerStylePanel.vue'
import { useLayerData } from './layer-data'
import LayerDataFeed from './LayerDataFeed.vue'
import { HexLayerSpec } from './map-layers'

const layer = defineModel<HexLayerSpec>('layer', { required: true })

const title = computed(() => (layer.value.name?.length ? layer.value.name : 'Hexgrid layer'))

const tab = ref<'data' | 'layer' | 'stats'>('data')

const expanded = defineModel<boolean>('expanded', { default: false })

const { layerData, registerLayer } = useLayerData()

const emit = defineEmits<{
  delete: []
  'draghandle-down': []
}>()

const remote = useQuery(
  computed(() =>
    occurrencesBySiteOptions({
      query: {
        ...layer.value.filters,
        habitats: layer.value.filters.habitats?.map(({ label }) => label)
      }
    })
  )
)

function register() {
  registerLayer(layer.value, remote)
}

onMounted(() => {
  register()
})
</script>

<style scoped lang="scss"></style>
