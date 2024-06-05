import { HabitatGroup, HabitatRecord } from "@/api"
import { computed, ref } from "vue"

export type Dependencies = { dependencies?: (HabitatRecord & { group: HabitatGroup })[] }
export type ConnectedHabitat = HabitatRecord & Dependencies & {
  group: HabitatGroup
}
export type ConnectedGroup = HabitatGroup & Dependencies & { depends?: ConnectedHabitat, elements: ConnectedHabitat[] }
export type HabitatsGraph = {
  groups: { [k: string]: ConnectedGroup }
  habitats: Record<string, ConnectedHabitat>
}

function addGroup(group: HabitatGroup, graph: HabitatsGraph) {
  graph.groups[group.id] = {
    ...group,
    depends: group.depends ? graph.habitats[group.depends.id] : undefined,
    elements: group.elements.map((habitat) => {
      const h = { ...habitat, group }
      graph.habitats[h.id] = h
      return h
    })
  }
  return graph
}

export function registerGroup(group: HabitatGroup, graph: HabitatsGraph) {
  addGroup(group, graph)
  if (group.depends) {
    const depends = graph.habitats[group.depends.id]
    updateDependencies(graph, group.id, [depends].concat(depends.dependencies ?? []))
  }
}

function updateDependencies(
  graph: HabitatsGraph,
  groupID: string,
  dependencies: (HabitatRecord & { group: HabitatGroup })[]
) {
  graph.groups[groupID].dependencies = dependencies
  graph.groups[groupID].elements = graph.groups[groupID].elements.map((habitat) => {
    graph.habitats[habitat.id].dependencies = graph.groups[groupID].dependencies
    return graph.habitats[habitat.id]
  })
}

/**
* Indexes groups and their children habitats by UUID,
* adding references to their groups and dependencies
* so they can be used as a graph-like structure
*/
export function indexGroups(groups: HabitatGroup[]) {
  const index = groups.reduce<HabitatsGraph>(
    (acc: HabitatsGraph, group) => addGroup(group, acc),
    { groups: {}, habitats: {} }
  )
  // Index dependencies on each element up to the root node
  function collectDepends(group: HabitatGroup): (HabitatRecord & { group: HabitatGroup })[] {
    if (group.depends == undefined) return []
    else {
      const deps = collectDepends(index.habitats[group.depends.id].group)
      deps.push(index.habitats[group.depends.id])
      return deps
    }
  }
  for (const key in index.groups) {
    updateDependencies(index, key, collectDepends(index.groups[key]))
  }
  return index
}

// State management

const habitatGraph = ref<HabitatsGraph>()
const selection = ref<ConnectedHabitat>()


export function useHabitatGraph(groups?: HabitatGroup[]) {

  function select(habitat: ConnectedHabitat) {
    selection.value = habitat
  }

  function isSelected(habitat: ConnectedHabitat) {
    return computed(() => habitat.id === selection.value?.id)
  }

  function isIncompatibleWithSelection(habitat: ConnectedHabitat) {
    return computed(() => {
      return (selection.value?.incompatible?.find(({ id }) => id === habitat.id)) ||
        (
          selection.value?.group.label == habitat.group.label &&
          selection.value?.id !== habitat.id &&
          habitat.group.exclusive_elements
        )
    })
  }

  function buildGraph(groups: HabitatGroup[]) {
    habitatGraph.value = indexGroups(groups)
    return habitatGraph.value
  }

  // Graph initialization
  if (groups) {
    if (habitatGraph.value == undefined)
      habitatGraph.value = buildGraph(groups)
    else
      console.error("Graph is already initialized. Did you call useHabitatGraph with an argument multiple times ?")
  } else if (habitatGraph.value == undefined)
    console.error("Graph was never initialized, useHabitatGraph must be called with an argument")

  return { selection, select, isSelected, isIncompatibleWithSelection, addGroup, buildGraph, habitatGraph: habitatGraph.value as HabitatsGraph }
}