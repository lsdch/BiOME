<template>
  <div
    class="graph"
    ref="graphElement"
    @keyup.esc="creating = false"
    @keyup.delete="askDeleteGroups(selectedGroups)"
  >
    <VueFlow
      :nodes="nodes"
      :edges="edges"
      :default-viewport="{ zoom: 10 }"
      :min-zoom="0.4"
      :max-zoom="2"
      snap-to-grid
      :snap-grid="[20, 20]"
      @nodes-initialized="layoutGraph()"
      @pane-click="onPaneClick"
      connect-on-click
      :class="{ creating }"
    >
      <Background :gap="20" />
      <div v-if="creating" class="text-secondary pa-3">
        <pre>Click / tap: create group</pre>
        <pre>Esc: cancel group creation</pre>
      </div>
      <Controls>
        <ControlButton @click="layoutGraph()">
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
      <template #node-group="props">
        <HabitatGroupNode v-bind="props" />
      </template>
    </VueFlow>
    <div
      id="new-group-template"
      v-if="creating"
      class="group-node text-overline bg-primary"
      :style="
        form.open
          ? { left: `${creationPos.x}px`, top: `${creationPos.y}px` }
          : { left: `${x}px`, top: `${y}px` }
      "
    >
      New group
    </div>
    <HabitatFormDialog
      v-model="form.open"
      mode="Create"
      @success="addCreatedNode"
      @close="creating = false"
    />
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
import { type Mode } from '@/components/toolkit/forms/form'
import { Background } from '@vue-flow/background'
import { ControlButton, Controls } from '@vue-flow/controls'
import { Edge, MarkerType, Node, VueFlow, useVueFlow } from '@vue-flow/core'
import { useMouseInElement } from '@vueuse/core'
import { computed, inject, nextTick, reactive, ref, toRefs } from 'vue'
import HabitatFormDialog from './HabitatFormDialog.vue'
import HabitatGroupNode from './HabitatGroupNode.vue'
import { useLayout } from './layout'

import { ConfirmDialogKey } from '@/injection'
import { useFeedback } from '@/stores/feedback'
import { ConnectedGroup, useHabitatGraph } from './habitat_graph'

const data = await LocationService.listHabitatGroups()

const selectedGroups = computed<HabitatGroup[]>(() => {
  return getSelectedNodes.value.map(({ data }) => data)
})

/**
 * Creation mode: cursor displays as a node template to place on the graph
 */
const creating = ref(false)
const form = ref<{ open: boolean; mode: Mode }>({
  open: false,
  mode: 'Create'
})

const { fitView, getSelectedNodes, project, onConnect } = useVueFlow()

onConnect(({ source, target, targetHandle }) => {
  edges.value.push({
    id: `edge-${target}-${source}`,
    target: source,
    source: target,
    sourceHandle: targetHandle,
    type: 'smoothstep',
    markerEnd: {
      type: MarkerType.ArrowClosed,
      width: 30,
      height: 30
    }
  })
})

const { feedback } = useFeedback()
const confirmDeletion = inject(ConfirmDialogKey)

function askDeleteGroups(groups: HabitatGroup[]) {
  const title =
    groups.length === 1
      ? `Delete habitat group ${groups[0].label} ?`
      : groups.length > 0
        ? `Delete ${groups.length} habitat groups ?`
        : null
  if (title !== null)
    confirmDeletion?.({
      title,
      message: 'All terms in group and their references in the database will be deleted.'
    }).then(({ isCanceled }) =>
      isCanceled ? console.info('Group deletion cancelled') : deleteGroups(groups)
    ) ?? console.error('No confirm dialog provider')
}

async function deleteGroups(groups: HabitatGroup[]) {
  await Promise.all(
    groups.map((group) => LocationService.deleteHabitatGroup({ code: group.label }))
  )
    .then((deleted) => {
      blockLayout.value = true
      nodes.value = nodes.value.filter(
        ({ id }) => !deleted.find(({ id: deletedID }) => deletedID === id)
      )
      feedback({
        type: 'success',
        message:
          deleted.length > 1
            ? `Deleted ${deleted.length} habitat group`
            : `Deleted habitat group ${deleted[0].label}`
      })
    })
    .catch(() => {
      feedback({
        type: 'error',
        message: `Failed to delete habitat group(s)`
      })
    })
}

const graphElement = ref(null)
const { elementX: x, elementY: y } = useMouseInElement(graphElement)
const creationPos = ref({ x: 0, y: 0 })

function onPaneClick({ layerX, layerY }: MouseEvent) {
  selection.value = undefined
  if (creating.value) {
    creationPos.value = { x: layerX, y: layerY }
    form.value.open = true
  }
}

const { selection, habitatGraph, addGroup } = useHabitatGraph(data)

const { nodes, edges } = toRefs(reactive(collectGraphElements(data)))

function addCreatedNode(group: HabitatGroup) {
  creating.value = false
  blockLayout.value = true
  addGroup(group, habitatGraph)
  nodes.value.push(createNode(habitatGraph.groups[group.id], project(creationPos.value)))
  feedback({ type: 'success', message: `Created node group: ${group.label}` })
}

function createNode(
  group: ConnectedGroup,
  position: { x: number; y: number }
): Node<ConnectedGroup> {
  const { id, label } = group
  return {
    id,
    label,
    position,
    data: group,
    type: 'group'
  }
}

function collectGraphElements(groups: HabitatGroup[]) {
  return groups.reduce(
    (acc: { nodes: Node<ConnectedGroup>[]; edges: Edge[] }, group) => {
      const { id, label, depends } = group
      acc.nodes.push(createNode(habitatGraph.groups[group.id], { x: 0, y: 0 }))
      if (depends) {
        acc.edges.push({
          id: `edge-${depends.label}-${label}`,
          source: habitatGraph.habitats[depends.id].group.id ?? depends.id,
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

const { layout, blockLayout } = useLayout<ConnectedGroup>()
async function layoutGraph() {
  nodes.value = layout(
    nodes.value.filter(({ parentNode }) => !parentNode),
    edges.value,
    'LR'
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
