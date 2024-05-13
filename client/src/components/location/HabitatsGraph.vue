<template>
  <div class="graph" ref="graphElement" @keyup.esc="console.log('PRESS'), (creating = false)">
    <VueFlow
      :nodes="nodes"
      :edges="edges"
      :default-viewport="{ zoom: 10 }"
      :min-zoom="0.2"
      :max-zoom="4"
      snap-to-grid
      :snap-grid="[20, 20]"
      @nodes-initialized="layoutGraph('LR')"
      @pane-click="selection = undefined"
      :class="{ creating }"
      @keyup.esc="console.log('PRESS'), (creating = false)"
    >
      <Background :gap="20" />
      <div class="text-secondary pa-3">
        <pre>Click / tap: create group</pre>
        <pre>Esc: cancel group creation</pre>
      </div>
      <Controls>
        <ControlButton @click="layoutGraph('LR')">
          <v-icon class="text-black">mdi-graph</v-icon>
        </ControlButton>
      </Controls>
      <v-fab
        :color="creating ? 'error' : 'primary'"
        app
        :prepend-icon="creating ? 'mdi-close' : 'mdi-plus'"
        @click="creating = !creating"
      >
        {{ creating ? 'Cancel' : 'New group' }}
      </v-fab>
      <template #node-habitat="props">
        <HabitatNode v-bind="props" />
      </template>
      <template #node-group="props">
        <HabitatGroupNode v-bind="props" />
      </template>
    </VueFlow>
    <div
      id="new-group-template"
      v-if="creating"
      class="group-node text-overline bg-primary"
      :style="{ left: `${elementX}px`, top: `${elementY}px` }"
    >
      New group
    </div>
  </div>
</template>

<script setup lang="ts">
/* these are necessary styles for vue flow */
import '@vue-flow/core/dist/style.css'
/* this contains the default theme, these are optional styles */
import '@vue-flow/core/dist/theme-default.css'
// import default controls styles
import '@vue-flow/controls/dist/style.css'

import { HabitatGroup, LocationService } from '@/api'
import { ControlButton, Controls } from '@vue-flow/controls'
import { Edge, MarkerType, Node, VueFlow, useVueFlow } from '@vue-flow/core'
import { nextTick, reactive, toRefs, ref } from 'vue'
import HabitatGroupNode from './HabitatGroupNode.vue'
import HabitatNode from './HabitatNode.vue'
import { useSelection } from './habitats'
import { useLayout } from './layout'
import HabitatNodeContainer from './HabitatNodeContainer.vue'
import { UseMouseEventExtractor, useMouse, useMouseInElement, useParentElement } from '@vueuse/core'
import { watch } from 'vue'
import { Background } from '@vue-flow/background'

import { onKeyStroke } from '@vueuse/core'

onKeyStroke('Escape', () => {
  creating.value = false
})

const { fitView } = useVueFlow()
const { layout } = useLayout<HabitatGroup>()

const creating = ref(true)
// const extractor: UseMouseEventExtractor = (event) => {
//   if (typeof Touch !== 'undefined' && event instanceof Touch) return null
//   else if (event instanceof MouseEvent) [event.offsetX, event.offsetY]
// }
const graphElement = ref(null)
const { elementX, elementY, sourceType } = useMouseInElement(graphElement)

const data = await LocationService.listHabitatGroups()

// Lookup table for element group by element ID
const nodeIndex = data.reduce(
  (
    acc: Record<string, { group: { id: string; label: string } }>,
    { id: groupId, label: groupLabel, elements }
  ) => {
    elements.forEach(({ id }) => (acc[id] = { group: { id: groupId, label: groupLabel } }))
    return acc
  },
  {}
)

const { nodes, edges } = toRefs(reactive(collectGraphElements(data)))

const { selection } = useSelection()

function collectGraphElements(groups: HabitatGroup[]) {
  return groups.reduce(
    (acc: { nodes: Node<HabitatGroup>[]; edges: Edge[] }, group) => {
      const { id, label, depends, exclusive_elements, elements } = group
      acc.nodes.push({
        id,
        label,
        position: { x: 0, y: 0 },
        data: group,
        type: 'group'
      })
      if (depends) {
        acc.edges.push({
          id: `edge-${depends.label}-${label}`,
          source: nodeIndex[depends.id].group.id ?? depends.id,
          sourceHandle: depends.id,
          target: id,
          type: 'smoothstep',
          markerEnd: {
            type: MarkerType.ArrowClosed,
            width: 30,
            height: 30
          }
        })
      }
      return acc
    },
    { nodes: [], edges: [] }
  )
}

async function layoutGraph(direction: 'LR' | 'TB' | 'BT' | 'RL') {
  nodes.value = layout(
    nodes.value.filter(({ parentNode }) => !parentNode),
    edges.value,
    direction
  )

  nextTick(() => {
    fitView()
  })
}
</script>

<style lang="scss">
#new-group-template {
  position: absolute;
  overflow: hidden;
  padding: 10px;
  display: flex;
  height: auto;
  width: auto;
  text-wrap: nowrap;
  cursor: none;
}
.graph {
  position: relative;
  width: 100%;
  height: 100%;

  .creating {
    .vue-flow__pane {
      cursor: none;
    }
  }
}
</style>
