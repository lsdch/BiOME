<template>
  <CardDialog title="Save preset" v-model="dialog">
    <template #activator="props">
      <slot name="activator" v-bind="props" />
    </template>
    <template #append>
      <v-btn
        variant="tonal"
        text="Save"
        @click="mutate({ body: { ...preset, spec: JSON.stringify(specs) } })"
      />
    </template>
    <v-card-text class="pt-1">
      <v-alert class="text-caption mb-3">
        Saving a preset allows you to quickly restore your current map configuration, including data
        feeds, hexgrid settings, and marker layers. You can share presets with others or keep them
        private.
      </v-alert>
      <v-text-field label="Name" v-model="preset.name" v-bind="schema('name')" />
      <v-textarea label="Description" v-model="preset.description" />
      <v-switch
        v-if="isGranted('Contributor')"
        label="Publish preset"
        v-model="preset.is_public"
        v-bind="schema('is_public')"
        persistent-hint
      />
      <v-switch
        v-if="isGranted('Maintainer')"
        label="Global preset"
        v-model="preset.is_global"
        persistent-hint
        v-bind="schema('is_global')"
      />
    </v-card-text>
  </CardDialog>
</template>

<script setup lang="ts">
import { $MapToolPresetInput, MapToolPresetInput } from '@/api'
import { useSchema } from '@/composables/schema'
import { reactive } from 'vue'
import CardDialog from '../toolkit/ui/CardDialog.vue'
import { DataFeed } from '../occurrence/data_feeds'
import { HexgridLayerSpec, MarkerLayerDefinition } from './map-layers'
import { useMutation } from '@tanstack/vue-query'
import { useFeedback } from '@/stores/feedback'
import { createUpdateMapPresetMutation } from '@/api/gen/@tanstack/vue-query.gen'
import { useUserStore } from '@/stores/user'

export type MapPreset = {
  feeds: DataFeed[]
  hexgrid: HexgridLayerSpec
  markers: MarkerLayerDefinition[]
}

const dialog = defineModel<boolean>('dialog', {
  default: false
})

const preset = defineModel<MapToolPresetInput>({
  default: () => reactive({ name: '', is_global: false, is_public: false })
})

const { specs } = defineProps<{
  specs: MapPreset
}>()

const {
  bind: { schema }
} = useSchema($MapToolPresetInput)

const { isGranted } = useUserStore()

const { feedback } = useFeedback()

const { mutate } = useMutation({
  ...createUpdateMapPresetMutation(),
  onSuccess: () => {
    feedback({
      type: 'success',
      message: `Preset "${preset.value.name}" saved successfully.`
    })
    dialog.value = false
  }
})
</script>

<style scoped lang="scss"></style>
