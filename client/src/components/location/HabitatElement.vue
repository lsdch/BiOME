<template>
  <div :class="[textClass, 'habitat-item text-no-wrap']" @click="select(habitat, group)">
    <Handle :id="habitat.id" :position="Position.Right" type="source" />
    {{ habitat.label }}
    <NodeToolbar
      v-if="isSelected(habitat).value"
      :is-visible="isSelected(habitat).value"
      :position="Position.Right"
      class="d-flex flex-column"
    >
      <BtnTooltip key="1" class="ma-1" size="small" icon="mdi-pencil" tooltip="Edit" />
      <BtnTooltip
        key="2"
        v-bind="props"
        class="ma-1"
        size="small"
        icon="mdi-arrow-right"
        tooltip="Add dependency"
      />
      <BtnTooltip
        key="3"
        class="ma-1"
        size="small"
        color="error"
        icon="mdi-arrow-right"
        tooltip="Add incompatibility"
      />
    </NodeToolbar>
  </div>
</template>

<script setup lang="ts">
import { HabitatGroup, HabitatRecord } from '@/api'
import { Handle, Position } from '@vue-flow/core'
import { NodeToolbar } from '@vue-flow/node-toolbar'
import { computed } from 'vue'
import BtnTooltip from '../toolkit/ui/BtnTooltip.vue'
import { useSelection } from './habitats'

const props = defineProps<{
  habitat: HabitatRecord
  group?: HabitatGroup
}>()

const { isSelected, isIncompatibleWithSelection, select } = useSelection()

const inRootGroup = computed(() => {
  return props.group && !props.group.depends
})

const textClass = computed(() => {
  if (isSelected(props.habitat).value) return 'text-primary'
  if (isIncompatibleWithSelection(props.habitat, props.group).value) return 'text-error'
  return ''
})
</script>

<style scoped lang="scss">
.vue-flow__handle {
  // border-radius: 2px;
  width: 2px;
  height: 2px;
  &.vue-flow__handle-bottom {
    top: 100%;
  }
  &.vue-flow__handle-right {
    left: 100%;
  }
  // visibility: hidden;
}

div.habitat-item {
  position: relative;
  padding: 5px;
  border-top: 1px solid grey;
  cursor: pointer;
  height: 39px;
  line-height: 30px;
  box-sizing: border-box;
  &.horizontal {
    padding-left: 10px;
    padding-right: 10px;
    text-align: center;
    border-top: none;
    border-left: 1px solid grey;
  }
}

div.habitat-item:first-child {
  border: none;
}
</style>
