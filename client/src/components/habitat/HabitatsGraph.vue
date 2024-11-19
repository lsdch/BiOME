<template>
  <div class="graph" ref="graphContainer" @keyup.esc="toggleCreating(false)">
    <VueFlow
      :nodes="graph.nodes.value"
      :edges="graph.edges.value"
      :default-viewport="{ zoom: 10 }"
      :min-zoom="0.4"
      :max-zoom="2"
      snap-to-grid
      :snap-grid="[20, 20]"
      @nodes-initialized="layout"
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
        <ControlButton @click="layout">
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
      v-if="creating"
      id="new-group-template"
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
      @close="toggleCreating(false)"
    />
  </div>
</template>

<script setup lang="ts">
/* necessary styles for vue flow */
import '@vue-flow/core/dist/style.css'
/* optional styles: default theme for vue flow  */
import '@vue-flow/core/dist/theme-default.css'

import { HabitatGroup, LocationService } from '@/api'
import { Background } from '@vue-flow/background'
import { ControlButton, Controls } from '@vue-flow/controls'
import '@vue-flow/controls/dist/style.css'
import { ConnectionLineType, GraphEdge, VueFlow, useConnection, useVueFlow } from '@vue-flow/core'
import { onKeyStroke, useMouseInElement, useToggle } from '@vueuse/core'
import { computed, nextTick, ref } from 'vue'
import HabitatFormDialog from './HabitatFormDialog.vue'
import HabitatGroupNode from './HabitatGroupNode.vue'
import { useLayout } from './layout'

import { mergeResponses } from '@/api/responses'
import { useAppConfirmDialog } from '@/composables'
import { useFeedback } from '@/stores/feedback'
import { useUserStore } from '@/stores/user'
import {
  ConnectedGroup,
  ConnectedHabitat,
  HabitatsGraph,
  useHabitatGraphSelection
} from './habitat_graph'

const { isGranted } = useUserStore()

function handleDelete() {
  if (selectedGroups.value.length) askDeleteGroups(selectedGroups.value)
  else if (getSelectedEdges.value.length == 1) askDeleteEdge(getSelectedEdges.value[0])
}

onKeyStroke('Delete', handleDelete)
onKeyStroke('Escape', () => {
  if (!form.value.open) creating.value = false
})

const data = await LocationService.listHabitatGroups().then(({ data, error }) => {
  if (error !== undefined) {
    console.error('Failed to fetch habitat groups: ', error)
    return []
  }
  return data
})

const graph = new HabitatsGraph(data)

const graphContainer = ref(null)
const { elementX: x, elementY: y } = useMouseInElement(graphContainer)
const creationPos = ref({ x: 0, y: 0 })

const { selection } = useHabitatGraphSelection()

const selectedGroups = computed<HabitatGroup[]>(() => {
  return getSelectedNodes.value.map(({ data }) => data)
})

/**
 * Creation mode: cursor displays as a node template to place on the graph
 */
const [creating, toggleCreating] = useToggle(false)
const form = ref<{ open: boolean; edit?: ConnectedGroup }>({
  open: false,
  edit: undefined
})

const connection = useConnection()
const { askConfirm } = useAppConfirmDialog()
const { feedback } = useFeedback()
const { fitView, getSelectedNodes, onConnect, endConnection, getSelectedEdges } = useVueFlow()

onConnect(({ source, target: dependGroupID, targetHandle }) => {
  connectGroupHabitat(graph.group(source), graph.habitat(targetHandle!))
})

function connectGroupHabitat(group: ConnectedGroup, habitat: ConnectedHabitat) {
  askConfirm({
    title: 'Confirm connection',
    message: `Set group ${group.label} as a refinement of ${habitat.label} ?`
  })
    .then(async ({ isCanceled }) => {
      if (isCanceled) return
      await LocationService.updateHabitatGroup({
        path: { code: group.label },
        body: { depends: habitat.label }
      }).then(({ data: updated, error }) => {
        if (error !== undefined) {
          console.error('Error:', error)
          return
        }
        graph.addEdge(habitat, updated)
      })
    })
    .finally(() => endConnection()) ??
    console.error('Failed to inject confirmation dialog in component')
}

function onPaneClick({ layerX, layerY }: MouseEvent) {
  selection.value = undefined
  form.value.edit = undefined
  if (creating.value) {
    creationPos.value = { x: layerX, y: layerY }
    form.value.open = true
  }
}

function askDeleteEdge(edge: GraphEdge) {
  const group: string = edge.targetNode.data.label
  askConfirm({ title: `Drop dependency of '${group}'?` }).then(async ({ isCanceled }) => {
    if (isCanceled) return console.info('Dependency drop canceled')
    else {
      const { data: updated, error } = await LocationService.updateHabitatGroup({
        path: { code: group },
        body: { depends: null }
      })
      if (error !== undefined) {
        feedback({ type: 'error', message: 'Failed to remove dependency.' })
        console.error(error)
        return
      }
      graph.updateGroup(updated)
      graph.removeEdge(edge)
      feedback({ type: 'success', message: `Dropped dependency for '${group}'` })
    }
  })
}

async function askDeleteGroups(groups: HabitatGroup[]) {
  if (groups.length === 0) return
  const title =
    groups.length === 1
      ? `Delete habitat group ${groups[0].label} ?`
      : `Delete ${groups.length} habitat groups ?`
  askConfirm({
    title,
    message: 'All terms in group and their references in the database will be deleted.'
  }).then(async ({ isCanceled }) => {
    if (isCanceled) {
      console.info('Group deletion cancelled')
      return
    }
    return await deleteHandler(groups)
  }) ?? console.error('No confirm dialog provider')

  async function deleteHandler(groups: HabitatGroup[]) {
    await Promise.all(
      groups.map((group) => LocationService.deleteHabitatGroup({ path: { code: group.label } }))
    )
      .then(mergeResponses)
      .then(({ data, error }) => {
        if (error !== undefined) {
          feedback({
            type: 'error',
            message: `Failed to delete habitat group(s)`
          })
          console.error(error)
          return
        }
        blockLayout.value = true
        graph.deleteGroups(data)
        feedback({
          type: 'success',
          message:
            data.length > 1
              ? `Deleted ${data.length} habitat group`
              : `Deleted habitat group ${data[0].label}`
        })
      })
  }
}

const { blockLayout } = useLayout()
function addCreatedNode(group: HabitatGroup) {
  toggleCreating(false)
  blockLayout.value = true
  graph.addGroup(group)
  feedback({ type: 'success', message: `Created habitat group: ${group.label}` })
}

async function layout() {
  graph.layout()
  nextTick(fitView)
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
