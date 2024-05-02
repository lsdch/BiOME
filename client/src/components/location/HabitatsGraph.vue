<template>
  <div class="graph">
    <VueFlow
      :nodes="nodes"
      :edges="edges"
      :default-viewport="{ zoom: 10 }"
      :min-zoom="0.2"
      :max-zoom="4"
      @nodes-initialized="layoutGraph('TB')"
    >
      <Controls />
    </VueFlow>
  </div>
</template>

<script setup lang="ts">
/* these are necessary styles for vue flow */
import '@vue-flow/core/dist/style.css'
/* this contains the default theme, these are optional styles */
import '@vue-flow/core/dist/theme-default.css'

// import default controls styles
import '@vue-flow/controls/dist/style.css'

import { LocationService } from '@/api'
import { nextTick, reactive, ref, toRefs } from 'vue'
import { Edge, Element, Node, VueFlow, useVueFlow } from '@vue-flow/core'
import { Controls } from '@vue-flow/controls'
import { useLayout } from './layout'

const { onPaneReady, fitView } = useVueFlow()

const { graph, layout, previousDirection } = useLayout()

const data = await LocationService.listHabitats()
const { nodes, edges } = toRefs(
  reactive(
    data.reduce(
      (acc: { nodes: Node[]; edges: Edge[] }, { id, label, depends, incompatible }) => {
        acc.nodes.push({ id, label, position: { x: 0, y: 0 }, type: depends ? 'default' : 'input' })
        acc.edges = acc.edges.concat(
          depends
            ? depends.map((dep) => ({
                id: `e${dep.label}-${label}`,
                source: dep.id,
                target: id,
                type: 'smoothstep'
              }))
            : []
        )
        return acc
      },
      { nodes: [], edges: [] }
    )
  )
)

async function layoutGraph(direction: 'LR' | 'TB' | 'BT' | 'RL') {
  nodes.value = layout(nodes.value, edges.value, direction)

  nextTick(() => {
    fitView()
  })
}
</script>

<style scoped>
.graph {
  width: 100%;
  height: 100%;
}
</style>
