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
        :id="data.id"
        :class="['button', { mobile }]"
        :position="Position.Left"
        type="source"
      >
        <template #default="{ id }">
          <BtnTooltip
            v-if="selected && startHandle == null"
            :id="id"
            :size="mobile ? 'large' : 'x-small'"
            icon="mdi-arrow-left-bold"
            tooltip="Edit"
            flat
            rounded="100%"
          />
          <v-progress-circular
            v-else-if="startHandle != null"
            class="spinner"
            indeterminate
            :color="status == 'valid' ? 'success' : 'warning'"
          >
            <v-icon icon="mdi-arrow-left-bold" size="small"></v-icon>
          </v-progress-circular>
        </template>
      </Handle>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Handle, NodeProps, Position, useConnection } from '@vue-flow/core'
import { computed } from 'vue'
import BtnTooltip from '../toolkit/ui/BtnTooltip.vue'
import HabitatElement from './HabitatElement.vue'
import { ConnectedGroup, useHabitatGraph } from './habitat_graph'
import { useDisplay } from 'vuetify'
const { mobile } = useDisplay()

const props = defineProps<NodeProps<ConnectedGroup>>()

const { connection } = useHabitatGraph()
const { startHandle, status } = useConnection()

// const sideLabel = computed(() => (props.data.elements.length > 1 ? props.data.label : undefined))
const sideLabel = computed(() => props.data.label)
</script>

<style lang="scss">
.vue-flow__handle {
  border-radius: 1px;
  width: 30px;
  height: 30px;
  &.button {
    border: none;
    background: none;
    border-radius: 100%;
    left: -20px;
    &.mobile {
      width: 60px;
      height: 60px;
      left: -40px;
    }

    button,
    .spinner {
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
