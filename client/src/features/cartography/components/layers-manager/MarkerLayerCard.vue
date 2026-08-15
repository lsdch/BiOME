<template>
  <v-card class="small-card-title" v-bind="$attrs">
    <template #prepend>
      <v-tooltip open-on-click :open-on-hover="false" :close-delay="500">
        <template #activator="{ props }">
          <v-icon
            icon="mdi-swap-vertical"
            class="handle cursor-grab text-muted"
            v-bind="props"
            @mousedown="emit('draghandle-down')"
          ></v-icon>
        </template>
        <span>Drag to reorder</span>
      </v-tooltip>
      <SvgCircle
        :size="20"
        :fill-color="layer.config.fillColor"
        :stroke-color="withOpacity(layer.config.color)"
        class="cursor-pointer"
        @click="((expanded = true), (tab = 'layer'))"
      />
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
      <div class="d-flex ga-1 align-center">
        <v-menu>
          <template #activator="{ props }">
            <v-btn
              icon="mdi-dots-vertical"
              variant="plain"
              :rounded="100"
              size="small"
              v-bind="props"
            />
          </template>
          <v-card>
            <v-list density="compact">
              <v-list-item
                prepend-icon="mdi-refresh"
                title="Reload data"
                @click="remote.refetch()"
              />
              <v-list-item
                prepend-icon="mdi-content-copy"
                title="Duplicate"
                @click="emit('duplicate')"
              />
              <v-list-item prepend-icon="mdi-delete" title="Delete" @click="emit('delete')" />
            </v-list>
          </v-card>
        </v-menu>
        <v-progress-circular
          v-if="layer.ready && remote.isPending.value"
          indeterminate
          size="small"
          color="warning"
        ></v-progress-circular>
        <v-icon
          v-if="!layer.ready"
          icon="mdi-database-alert"
          size="small"
          color="warning"
          v-tooltip="{ text: 'Waiting for data definition' }"
          @click="((expanded = true), (tab = 'data'))"
        >
        </v-icon>
        <v-switch v-else v-model="layer.active" color="primary" hide-details></v-switch>
        <v-btn
          :icon="expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
          color=""
          variant="text"
          :rounded="100"
          size="small"
          @click="() => (expanded = !expanded)"
        />
      </div>
    </template>
    <v-expand-transition>
      <div v-show="expanded" class="">
        <v-divider />
        <v-tabs v-model="tab" density="compact">
          <v-tab prepend-icon="mdi-database" value="data">Data</v-tab>
          <v-tab prepend-icon="mdi-circle-multiple-outline" value="layer">Style</v-tab>
          <v-tab prepend-icon="mdi-chart-bar" value="stats" :disabled="remote.isFetching.value">
            Stats
          </v-tab>
        </v-tabs>
        <v-divider></v-divider>
        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="data">
            <v-confirm-edit v-model="layer.filters">
              <template #default="{ model, actions: _, isPristine, cancel, save }">
                <LayerDataFeed v-model="model.value" class="pa-2" />
                <template v-if="!layer.ready || !isPristine">
                  <v-divider></v-divider>
                  <div class="d-flex justify-end pa-2 ga-1">
                    <v-btn text="Cancel" @click="cancel()" variant="text"></v-btn>
                    <v-btn text="Ok" @click="(save(), setReady())" rounded="md"></v-btn>
                  </div>
                </template>
              </template>
            </v-confirm-edit>
          </v-tabs-window-item>
          <v-tabs-window-item value="layer">
            <MarkerLayerStylePanel v-model="layer" />
          </v-tabs-window-item>
          <v-tabs-window-item value="stats">
            <OccurrencesStats :cells="remote.data.value" />
          </v-tabs-window-item>
        </v-tabs-window>
      </div>
    </v-expand-transition>
  </v-card>
</template>

<script setup lang="ts">
import { listOccurrencesH3Options } from '@/api/gen/@tanstack/vue-query.gen.ts'
import SvgCircle from '@/components/toolkit/ui/SvgCircle.vue'
import OccurrencesStats from '@/features/occurrences/components/OccurrencesStats.vue'
import { withOpacity } from '@/lib/color_brewer'
import { useQuery } from '@tanstack/vue-query'
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useLayerData } from './layer-data'
import LayerDataFeed from './LayerDataFeed.vue'
import { automaticResolution, MarkerLayerSpec } from './map-layers'
import MarkerLayerStylePanel from './MarkerLayerStylePanel.vue'

const { index, zoom } = defineProps<{
  index: number
  zoom: number
}>()

const layer = defineModel<MarkerLayerSpec>('layer', { required: true })

const title = computed(() =>
  layer.value.name?.length ? layer.value.name : `Markers #${index + 1}`
)

const tab = ref<'data' | 'layer' | 'stats'>('data')

const expanded = defineModel<boolean>('expanded', { default: false })

const emit = defineEmits<{
  delete: []
  'draghandle-down': []
  duplicate: []
}>()

const { registerLayer } = useLayerData()

watch(
  () => zoom,
  (newZoom) => {
    if (layer.value.resolutionMode === 'auto') {
      console.log('Updating marker resolution based on zoom', newZoom)
      layer.value.resolution = automaticResolution(layer.value, newZoom)
    }
  },
  { immediate: true }
)

const remote = useQuery(
  computed(() => {
    return {
      enabled: layer.value.ready && layer.value.active,
      initialData: [],
      refetchOnMount: true,
      gcTime: 1000 * 60 * 10, // 10 minutes
      ...listOccurrencesH3Options({
        path: { resolution: layer.value.resolution },
        query: {
          ...layer.value.filters
        }
      })
    }
  })
)

onMounted(() => {
  if (layer.value.ready) {
    registerLayer(layer.value, remote)
  }
})

function setReady() {
  if (!layer.value.ready) {
    nextTick(() => {
      registerLayer(layer.value, remote)
      layer.value.ready = true
    })
  }
}
</script>

<style scoped lang="scss"></style>
