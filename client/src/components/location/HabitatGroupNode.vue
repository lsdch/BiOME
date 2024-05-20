<template>
  <div :class="['node-wrapper d-flex flex-column w-auto']">
    <legend
      v-if="sideLabel"
      :class="[
        selected ? 'text-primary font-weight-bold' : '',
        'node-label text-overline align-self-center '
      ]"
    >
      {{ sideLabel }}
    </legend>
    <div :class="[selected ? 'selected' : '', 'group-node flex-grow-1 bg-surface']">
      <HabitatElement v-for="(habitat, index) in data.elements" :key="index" :habitat="habitat" />

      <Handle :position="Position.Left" type="target" style="visibility: hidden" />

      <Handle
        v-if="selected && !data.depends"
        class="button"
        :position="Position.Left"
        type="source"
      >
        <template v-if="selected" #default="{ id }">
          <BtnTooltip :id="id" size="x-small" icon="mdi-arrow-left-bold" tooltip="Edit" flat />
        </template>
      </Handle>
    </div>
  </div>
</template>

<script setup lang="ts">
import { HabitatGroup } from '@/api'
import { Handle, NodeProps, Position } from '@vue-flow/core'
import { computed } from 'vue'
import BtnTooltip from '../toolkit/ui/BtnTooltip.vue'
import HabitatElement from './HabitatElement.vue'

const props = defineProps<NodeProps<HabitatGroup>>()

const sideLabel = computed(() => (props.data.elements.length > 1 ? props.data.label : undefined))
</script>

<style lang="scss">
.vue-flow__handle {
  border-radius: 1px;
  width: 8px;
  &.button {
    border: none;
    background: none;
    border-radius: 100%;
    width: 30px;
    height: 30px;
    left: -20px;

    button {
      pointer-events: none;
      width: 100%;
      height: 100%;
    }
  }
}

.node-wrapper {
  position: relative;
  min-width: 160px;
  align-items: stretch;
}
.group-node {
  position: relative;
  box-sizing: border-box;
  border: 2.5px solid grey;
  border-radius: 3px;
  padding: 3px;
  height: 100%;

  &.selected {
    border-color: rgb(77, 170, 187);
  }
}

fieldset {
  box-sizing: border-box;
  border: 10px solid red;
}

legend.node-label {
  position: absolute;
  top: -30px;
  font-size: x-small;
  opacity: 0.5;
}
</style>
