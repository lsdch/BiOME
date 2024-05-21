import { HabitatGroup, HabitatRecord } from "@/api"

export type Dependencies = { dependencies?: (HabitatRecord & { group: HabitatGroup })[] }
export type ConnectedHabitat = HabitatRecord & Dependencies & {
  group: HabitatGroup
}
export type ConnectedGroup = HabitatGroup & Dependencies & { elements: ConnectedHabitat[] }
export type HabitatsGraph = {
  groups: { [k: string]: ConnectedGroup }
  habitats: Record<string, ConnectedHabitat>
}

export function addGroup(group: HabitatGroup, graph: HabitatsGraph) {
  graph.groups[group.id] = {
    ...group,
    elements: group.elements.map((habitat) => {
      const h = { ...habitat, group }
      graph.habitats[h.id] = h
      return h
    })
  }
  return graph
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
    index.groups[key].dependencies = collectDepends(index.groups[key])
    index.groups[key].elements = index.groups[key].elements.map((habitat) => {
      index.habitats[habitat.id].dependencies = index.groups[key].dependencies
      return index.habitats[habitat.id]
    })
  }
  return index
}