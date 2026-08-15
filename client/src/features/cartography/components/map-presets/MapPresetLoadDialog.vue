<template>
  <CardDialog title="Load preset" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <v-card-text class="bg-main">
      <v-chip-group v-model="include" multiple mandatory>
        <v-chip
          v-if="isGranted('contributor')"
          value="personal"
          text="My presets"
          color="primary"
        />
        <v-chip value="globals" text="App. level presets" color="warning" />
        <v-chip value="public" text="Public presets" color="success" />
      </v-chip-group>
      <MapToolPresetPicker
        v-model="model"
        :include
        clearable
        label="Presets"
        placeholder="Search by name or author"
      />
      <v-expand-transition>
        <v-card v-if="model" :title="`Preset: ${model.name}`">
          <template #append>
            <v-btn variant="tonal" text="Apply" rounded="sm" @click="apply(model)" />
          </template>
          <template #subtitle>
            <div class="d-flex ga-1">
              <v-chip text="App level" color="warning" size="small" />
              <v-chip :text="model.meta.created_by?.name" size="small" />
            </div>
          </template>
          <v-card-text v-if="model.description" class="text-caption">{{
            model.description
          }}</v-card-text>
          <v-list>
            <v-divider />
            <v-list-item title="Hexgrid layer" prepend-icon="mdi-hexagon-multiple">
              <template v-if="model.spec.hexgrid.config.colorRange" #append>
                <ColorPalettePreview :gradient="model.spec.hexgrid.config.colorRange" />
              </template>
            </v-list-item>
            <v-divider />
            <v-list-group value="markers">
              <template #activator="{ props }">
                <v-list-item
                  title="Data feeds"
                  prepend-icon="mdi-circle-multiple-outline"
                  v-bind="props"
                >
                  <template #title>
                    Marker layers
                    <v-badge color="success" inline :content="model.spec.markers.length" />
                  </template>
                </v-list-item>
              </template>
              <v-list-item
                v-for="(layer, i) in model.spec.markers"
                :title="layer.name ?? `Unnamed layer #${i + 1}`"
              >
                <template #prepend>
                  <SvgCircle
                    class="mr-2"
                    :size="20"
                    :fill-color="layer.config.fillColor"
                    :stroke-color="withOpacity(layer.config.color)"
                  />
                </template>
              </v-list-item>
            </v-list-group>
          </v-list>
        </v-card>
      </v-expand-transition>
    </v-card-text>
  </CardDialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CardDialog from '@/components/toolkit/ui/CardDialog.vue'
import { ParsedMapPreset } from '.'
import MapToolPresetPicker from './MapToolPresetPicker.vue'
import SvgCircle from '@/components/toolkit/ui/SvgCircle.vue'
import { withOpacity } from '@/lib/color_brewer'
import { useUserStore } from '@/stores/user'
import ColorPalettePreview from '@/components/toolkit/ui/ColorPalettePreview.vue'

const include = ref<('personal' | 'globals' | 'public')[]>(['personal', 'globals', 'public'])

const dialog = defineModel<boolean>('dialog', { default: false })

const model = defineModel<ParsedMapPreset>()

const { isGranted } = useUserStore()

const emit = defineEmits<{
  apply: [preset: ParsedMapPreset]
}>()

function apply(preset: ParsedMapPreset) {
  emit('apply', preset)
  dialog.value = false
}
</script>

<style scoped lang="scss"></style>
