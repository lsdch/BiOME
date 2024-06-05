<template>
  <div class="graph" ref="graphContainer" @keyup.esc="creating = false">
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
      :connection-line-options="{
        type: ConnectionLineType.SmoothStep,
        class: `${connection.status.value} animated`
      }"
      :delete-key-code="null"
    >
      <Background :gap="20" />
      <div v-if="isGranted('Admin')" id="help-pane" class="text-secondary pa-3">
        <div v-if="creating">
          <pre>Click / tap: create group</pre>
          <pre>Esc: cancel group creation</pre>
        </div>
        <div v-else-if="selectedGroups.length > 0">
          <pre>Del: delete selected groups</pre>
        </div>
        <div v-else-if="getSelectedEdges.length == 1">
          <pre>Del: delete selected edge</pre>
        </div>
      </div>
      <Controls>
        <ControlButton @click="layoutGraph()">
          <v-icon class="text-black">mdi-graph</v-icon>
        </ControlButton>
      </Controls>
      <!-- Edition controls -->
      <div v-if="isGranted('Admin')" class="vue-flow__panel bottom right">
        <v-btn
          v-if="selectedGroups.length > 0"
          class="mr-2"
          color="error"
          prepend-icon="mdi-delete"
          :text="`Delete ${selectedGroups.length} group(s)`"
          @click="askDeleteGroups(selectedGroups)"
        />
        <v-btn
          v-else-if="getSelectedEdges.length == 1"
          class="mr-2"
          color="error"
          prepend-icon="mdi-close"
          text="Drop dependency"
          @click="askDeleteEdge(getSelectedEdges[0])"
        />
        <v-btn
          :color="creating ? 'error' : 'primary'"
          :prepend-icon="creating ? 'mdi-close' : 'mdi-plus'"
          @click="creating = !creating"
          :text="creating ? 'Cancel' : 'New group'"
        />
      </div>
      <template #node-group="props">
        <HabitatGroupNode
          v-bind="props"
          @edit="(group) => ((form.open = true), (form.edit = group))"
        />
      </template>
      <template #node-habitat="props">
        {{ props.data.label }}
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
      :edit="form.edit"
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
import { HabitatGroup, LocationService } from '@/api'
import { Background } from '@vue-flow/background'
import { ControlButton, Controls } from '@vue-flow/controls'
import '@vue-flow/controls/dist/style.css'
import {
  ConnectionLineType,
  Edge,
  GraphEdge,
  MarkerType,
  Node,
  VueFlow,
  XYPosition,
  useConnection,
  useVueFlow
} from '@vue-flow/core'
import { useMouseInElement, onKeyStroke } from '@vueuse/core'
import { computed, inject, nextTick, reactive, ref, toRefs } from 'vue'
import HabitatFormDialog from './HabitatFormDialog.vue'
import HabitatGroupNode from './HabitatGroupNode.vue'
import { useLayout } from './layout'

import { ConfirmDialogKey } from '@/injection'
import { useFeedback } from '@/stores/feedback'
import { ConnectedGroup, registerGroup, useHabitatGraph } from './habitat_graph'
import { NodeData } from './layout'
import { useUserStore } from '@/stores/user'

const { isGranted } = useUserStore()

function handleDelete() {
  selectedGroups.value.length
    ? askDeleteGroups(selectedGroups.value)
    : getSelectedEdges.value.length == 1
      ? askDeleteEdge(getSelectedEdges.value[0])
      : undefined
}

onKeyStroke('Delete', handleDelete)

const data = await LocationService.listHabitatGroups()

const selectedGroups = computed<HabitatGroup[]>(() => {
  return getSelectedNodes.value.map(({ data }) => data)
})

/**
 * Creation mode: cursor displays as a node template to place on the graph
 */
const creating = ref(false)
const form = ref<{ open: boolean; edit?: ConnectedGroup }>({
  open: false,
  edit: undefined
})

const connection = useConnection()

const {
  fitView,
  getSelectedNodes,
  project,
  onConnect,
  updateNodeInternals,
  updateNodeData,
  endConnection,
  getSelectedEdges,
  removeEdges
} = useVueFlow()

const askConfirm = inject(ConfirmDialogKey)

onConnect(({ source: groupID, target: dependGroupID, targetHandle: dependHabitatID }) => {
  const group = habitatGraph.groups[groupID]
  const habitat = habitatGraph.habitats[dependHabitatID!]
  askConfirm?.({
    title: 'Confirm connection',
    message: `Set group ${group.label} as a refinement of ${habitat.label} ?`
  })
    .then(async ({ isCanceled }) => {
      if (isCanceled) return
      await LocationService.updateHabitatGroup({
        code: group.label,
        requestBody: { depends: habitat.label }
      }).then((updated) => {
        edges.value.push({
          id: `e-${dependHabitatID}-${groupID}`,
          target: groupID,
          source: dependGroupID,
          sourceHandle: dependHabitatID,
          markerEnd: {
            type: MarkerType.ArrowClosed,
            width: 20,
            height: 20
          }
        })
        registerGroup(updated, habitatGraph)
        updateNodeData(groupID, { depends: habitatGraph.habitats[dependHabitatID!] })
      })
    })
    .finally(() => endConnection()) ??
    console.error('Failed to inject confirmation dialog in component')
})

const { feedback } = useFeedback()
const confirmDeletion = inject(ConfirmDialogKey)

function askDeleteEdge(edge: GraphEdge) {
  const group: string = edge.targetNode.data.label
  confirmDeletion?.({ title: `Drop dependency of '${group}'?` }).then(async ({ isCanceled }) => {
    if (isCanceled) return console.info('Dependency drop canceled')
    else {
      const updated = await LocationService.updateHabitatGroup({
        code: group,
        requestBody: { depends: null }
      })
      registerGroup(updated, habitatGraph)
      updateNodeData(updated.id, { depends: null })
      removeEdges(edge)
      edges.value = edges.value.filter(({ id }) => id !== edge.id)
      feedback({ type: 'success', message: `Dropped dependency for '${group}'` })
    }
  })
}

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

const graphContainer = ref(null)
const { elementX: x, elementY: y } = useMouseInElement(graphContainer)
const creationPos = ref({ x: 0, y: 0 })

function onPaneClick({ layerX, layerY }: MouseEvent) {
  selection.value = undefined
  form.value.edit = undefined
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
  nodes.value.push(
    createNode(habitatGraph.groups[group.id], {
      cluster: group.id,
      position: project(creationPos.value),
      type: 'group'
    })
  )
  feedback({ type: 'success', message: `Created habitat group: ${group.label}` })
}

function createNode(
  group: ConnectedGroup,
  options: Partial<Node<NodeData>> & {
    cluster: string
    position: XYPosition
    type: 'group' | 'habitat'
  }
): Node<NodeData> {
  const { id, label } = group
  return {
    id,
    label,
    data:
      options.type === 'group'
        ? { ...group, cluster: options.cluster }
        : { ...group, cluster: options.cluster },
    ...options
  }
}

function edgeId({ label, depends }: HabitatGroup) {
  return `edge-${depends.label}-${label}`
}

function collectGraphElements(groups: HabitatGroup[]) {
  return groups.reduce<{ nodes: Node<NodeData>[]; edges: Edge[] }>(
    (acc, group) => {
      const { id, depends } = group
      acc.nodes.push(
        createNode(habitatGraph.groups[group.id], {
          cluster: id,
          type: 'group',
          position: { x: 0, y: 0 }
        })
      )
      if (depends) {
        acc.edges.push({
          id: edgeId(group),
          source: habitatGraph.habitats[depends.id].group.id,
          sourceHandle: depends.id,
          target: id,
          markerEnd: {
            type: MarkerType.ArrowClosed,
            width: 20,
            height: 20
          }
        })
      }
      return acc
    },
    { nodes: [], edges: [] }
  )
}

const { layout, blockLayout } = useLayout()
async function layoutGraph() {
  nodes.value = layout(nodes.value, edges.value, 'LR')
  nextTick(() => {
    // Update wrt new handle positions
    updateNodeInternals()
    fitView()
  })
}
</script>

<style lang="scss">
@use 'vuetify';
.vue-flow__edge-path {
  stroke-width: 2px;
}

.vue-flow__connection-path {
  stroke-width: 3px;
  stroke: rgb(var(--v-theme-primary));
  &.valid {
    stroke: rgb(var(--v-theme-success));
  }
}

.vue-flow__edge.selected {
  .vue-flow__edge-path {
    stroke: rgb(var(--v-theme-primary));
  }
}

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
