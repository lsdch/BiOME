<template>
  <div
    :class="[
      'habitat-item text-no-wrap',
      {
        connecting,
        'text-success font-weight-bold': endHandle?.handleId === habitat.id,
        'text-primary': isSelected(habitat).value
      }
    ]"
    @click="select(habitat)"
  >
    <Handle
      :id="habitat.id"
      :position="Position.Right"
      type="source"
      :connectable-start="false"
      :connectable-end="connecting"
    />
    {{ habitat.label }}
  </div>
</template>

<script setup lang="ts">
import { Handle, Position, useConnection } from '@vue-flow/core'
import { computed } from 'vue'
import { ConnectedHabitat, useHabitatGraph } from './habitat_graph'

const props = defineProps<{ habitat: ConnectedHabitat }>()

const { isSelected, select } = useHabitatGraph()
const { startHandle, endHandle } = useConnection()

const connecting = computed(() => {
  return (
    startHandle.value !== null &&
    props.habitat.group.id != startHandle.value?.nodeId &&
    props.habitat.dependencies?.find(({ group: { id } }) => id === startHandle.value?.nodeId) ===
      undefined
  )
})
</script>

<style scoped lang="scss">
@use 'vuetify';

.vue-flow__handle {
  width: 2px;
  height: 2px;
  visibility: hidden;
  &.vue-flow__handle-bottom {
    top: 100%;
  }
  &.vue-flow__handle-right {
    left: 100%;
  }
}

.habitat-item.connecting {
  .vue-flow__handle {
    visibility: visible;
    background-color: rgb(var(--v-theme-primary));
    width: 20px;
    height: 20px;
    border-radius: 100%;
    border-color: grey;
    &.connectionindicator {
      background-color: rgb(var(--v-theme-success));
    }
    &.vue-flow__handle-right {
      left: unset;
      right: -4px;
    }
  }
}

div.habitat-item {
  position: relative;
  padding: 5px;
  border-top: 1px solid grey;
  cursor: pointer;
  height: 39px;
  line-height: 30px;
  box-sizing: border-box;
}

div.habitat-item:first-child {
  border: none;
}
</style>
