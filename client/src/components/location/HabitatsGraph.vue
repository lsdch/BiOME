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
      @pane-click="onPaneClick"
      :class="{ creating }"
      @keyup.esc="console.log('PRESS'), (creating = false)"
    >
      <Background :gap="20" />
      <div v-if="creating" class="text-secondary pa-3">
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
    <ConfirmDialog
      v-model="confirmDialog.open"
      :title="confirmDialog.title"
      :message="confirmDialog.message"
      @confirm="deleteGroups(selectedGroups)"
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
import { computed, nextTick, reactive, ref, toRefs } from 'vue'
import HabitatFormDialog from './HabitatFormDialog.vue'
import HabitatGroupNode from './HabitatGroupNode.vue'
import { useSelection } from './habitats'
import { useLayout } from './layout'

import { useFeedback } from '@/stores/feedback'
import { onKeyStroke } from '@vueuse/core'
import ConfirmDialog from '../toolkit/ConfirmDialog.vue'

const confirmDialog = ref({
  open: false,
  message: '',
  title: ''
})

const selectedGroups = computed<HabitatGroup[]>(() => {
  return getSelectedNodes.value.map(({ data }) => data)
})

onKeyStroke('Escape', () => (creating.value = false))
onKeyStroke('Delete', () => {
  if (selectedGroups.value.length > 0)
    confirmDialog.value = {
      open: true,
      title:
        selectedGroups.value.length > 1
          ? `Delete ${selectedGroups.value.length} habitat groups ?`
          : `Delete habitat group ${selectedGroups.value[0].label} ?`,
      message: 'All terms in group and their references in the database will be deleted.'
    }
})

const { feedback } = useFeedback()
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

const form = ref<{ open: boolean; mode: Mode }>({
  open: false,
  mode: 'Create'
})

const { fitView, getSelectedNodes } = useVueFlow()
const { layout, blockLayout } = useLayout<HabitatGroup>()

const creating = ref(false)

const graphElement = ref(null)
const { elementX: x, elementY: y } = useMouseInElement(graphElement)
const creationPos = ref({ x: 0, y: 0 })

function onPaneClick() {
  selection.value = undefined
  if (creating.value) {
    creationPos.value = { x: x.value, y: y.value }
    form.value.open = true
  }
}

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

function addCreatedNode(group: HabitatGroup) {
  creating.value = false
  blockLayout.value = true
  nodes.value.push(createNode(group, { x: x.value, y: y.value }))
}

const { selection } = useSelection()

function createNode(group: HabitatGroup, position: { x: number; y: number }): Node<HabitatGroup> {
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
    (acc: { nodes: Node<HabitatGroup>[]; edges: Edge[] }, group) => {
      const { id, label, depends, exclusive_elements, elements } = group
      acc.nodes.push(createNode(group, { x: 0, y: 0 }))
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
