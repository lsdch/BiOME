<template>
  <div :class="[layout, 'node-wrapper d-flex flex-column w-auto']">
    <legend
      v-if="sideLabel"
      :class="[
        selected ? 'text-primary font-weight-bold' : '',
        'node-label text-overline align-self-center '
      ]"
    >
      {{ sideLabel }}
    </legend>
    <div :class="[layout, selected ? 'selected' : '', 'group-node flex-grow-1 bg-surface']">
      <slot :layout="layout"></slot>
      <Handle :position="Position.Left" type="target" />
      <Handle :position="Position.Bottom" type="source" style="visibility: hidden" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Handle, Position } from '@vue-flow/core'

withDefaults(
  defineProps<{
    sideLabel?: string
    layout?: 'vertical' | 'horizontal'
    selected?: boolean
  }>(),
  {
    layout: 'vertical'
  }
)
</script>

<style lang="scss">
.vue-flow__handle {
  border-radius: 1px;
  width: 8px;
}

.node-wrapper {
  min-width: 160px;
  align-items: stretch;
  &.horizontal {
    flex-direction: column;
  }
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

  &.horizontal {
    // padding: 0;
    // padding-bottom: 0;
    box-sizing: border-box;
    display: grid;
    grid-auto-columns: minmax(0, 1fr);
    grid-auto-flow: column;
  }
}

fieldset {
  box-sizing: border-box;
  border: 10px solid red;
}

legend.node-label {
  font-size: x-small;
  opacity: 0.5;
  padding-left: 10px;
  &.vertical {
    line-height: 0;
    transform: rotate(180deg);
    writing-mode: vertical-rl;
    text-orientation: mixed;
  }
}
</style>
