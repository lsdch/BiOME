import { HabitatGroup, HabitatRecord } from "@/api"
import { Edge, GraphEdge, MarkerType, Node, useVueFlow } from "@vue-flow/core"
import { computed, nextTick, Ref, ref } from "vue"
import { useLayout } from "./layout"


export type UpstreamData = {
  upstreamGroups: Set<UUID>
  dependencies: HabitatRecord[]
}

export type ConnectedGroup = HabitatGroup & UpstreamData & {
  depends?: ConnectedHabitat | null,
  elements: ConnectedHabitat[],
  cluster: string
}

export type ConnectedHabitat = HabitatRecord & UpstreamData & {
  group: HabitatGroup
}

type UUID = string


const { layout, blockLayout } = useLayout()
const {
  updateNodeInternals,
  updateNode,
  removeEdges,
} = useVueFlow()

export class HabitatsGraph {
  groups: Record<UUID, ConnectedGroup> = {};
  habitats: Record<UUID, ConnectedHabitat> = {};
  nodes: Ref<Node<ConnectedGroup>[]> = ref([]);
  edges: Ref<Edge[]> = ref([]);
  withNodes = true


  constructor(groupDefinitions: HabitatGroup[], withNodes = true) {
    this.withNodes = withNodes
    groupDefinitions.forEach(this._addGroup, this)
    this.updateDependencies()
    if (this.withNodes) this.buildEdges()
  }

  async layout() {
    this.nodes.value = layout(this.nodes.value, this.edges.value, 'LR')
    // Update wrt new handle positions
    nextTick(updateNodeInternals)
  }

  parentGroup(habitat: ConnectedHabitat): ConnectedGroup {
    return this.groups[habitat.group.id]
  }
  group(id: UUID): ConnectedGroup {
    return this.groups[id]
  }
  habitat(id: UUID): ConnectedHabitat {
    return this.habitats[id]
  }

  /**
   * Update the dependency chain of each group
   * This is useful to know which habitats can be connected to a group when
   * setting its requirement
   */
  updateDependencies() {
    const visited = new Set()
    const self = this
    for (const id in this.groups) {
      update(this.groups[id])
    }

    function update(group: ConnectedGroup) {
      if (!group.depends) {
        group.upstreamGroups = new Set()
        group.dependencies = []
      } else {
        const upstreamGroup = self.parentGroup(self.habitat(group.depends.id))
        update(upstreamGroup)
        if (!visited.has(group.id)) {
          group.upstreamGroups = (
            new Set([group.depends.group.id]).union(upstreamGroup.upstreamGroups)
          )
          group.dependencies = upstreamGroup.dependencies.concat([group.depends])
          visited.add(group.id)
        }
      }
      // Update upstream groups for all habitats contained in group
      group.elements.forEach(h => {
        self.habitat(h.id).upstreamGroups = group.upstreamGroups
        self.habitat(h.id).dependencies = group.dependencies
      })
    }
  }

  private _addGroup(group: HabitatGroup) {
    const createNode = this.withNodes && !this.groups[group.id]
    this.groups[group.id] = {
      ...group,
      depends: group.depends ? this.habitats[group.depends.id] : null,
      cluster: group.id,
      upstreamGroups: new Set(),
      dependencies: [],
      elements: group.elements.map((habitat) => {
        const h: ConnectedHabitat = {
          ...habitat, group,
          upstreamGroups: new Set(group.id),
          dependencies: []
        }
        this.habitats[h.id] = h
        return h
      })
    }

    if (createNode) this.nodes.value.push({
      id: group.id,
      label: group.label,
      data: this.groups[group.id],
      position: { x: 0, y: 0 },
      type: "group"
    })
  }

  addGroup(group: HabitatGroup) {
    this._addGroup(group)
    this.updateDependencies()
  }

  updateGroup(group: HabitatGroup) {
    console.log("Update node: ", group)
    this.addGroup(group)
    updateNode(group.id, { data: group, id: group.id, label: group.label })
  }

  buildEdges() {
    this.edges.value = []
    Object.values(this.groups).forEach((group) => {
      if (!!group.depends)
        this.addEdge(this.habitat(group.depends.id), group)
    })
  }

  addEdge(from: ConnectedHabitat, to: HabitatGroup) {
    const fromGroup = this.parentGroup(from)
    this.edges.value.push({
      id: `edge-${from.id}-${to.id}`,
      target: to.id,
      source: fromGroup.id,
      sourceHandle: from.id,
      markerEnd: {
        type: MarkerType.ArrowClosed,
        width: 20,
        height: 20
      }
    })
    this.updateGroup(to)
  }

  removeEdge(edge: GraphEdge) {
    // remove from vue flow internal state
    removeEdges(edge)
    // remove from edges state
    this.edges.value = this.edges.value.filter(({ id }) => id !== edge.id)
  }



  deleteGroups(groups: HabitatGroup[]) {
    this.nodes.value = this.nodes.value.filter(
      ({ id }) => !groups.find(({ id: deletedID }) => deletedID === id)
    )
    groups.forEach(group => {
      group.elements.forEach(h => {
        delete this.habitats[h.id]
      })
      delete this.groups[group.id]
    })
  }
}


const selection = ref<ConnectedHabitat>()



export function useHabitatGraphSelection() {
  function select(habitat: ConnectedHabitat) {
    selection.value = habitat
  }

  function isSelected(habitat: ConnectedHabitat) {
    return computed(() => habitat.id === selection.value?.id)
  }
  return { selection, select, isSelected }
}