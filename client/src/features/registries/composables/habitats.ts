import { HabitatGroupWithElements, Habitat } from '@/api'
import { getHabitatGroupsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

export type HabitatNode = Habitat & {
  children?: HabitatGroupWithElements[]
  upstream?: Habitat[]
  group: HabitatGroupWithElements
}

export type HabitatGroupNode = HabitatGroupWithElements & {
  children: Habitat[]
}

export type HabitatGraph = {
  rootGroups: HabitatGroupNode[]
  habitats: Map<UUID, HabitatNode>
  groups: Map<UUID, HabitatGroupNode>
}

export function useHabitats() {
  const { data: groups, refetch } = useQuery({
    ...getHabitatGroupsOptions(),
    initialData: []
  })

  const habitatGraph = computed<HabitatGraph>(() => {
    const graph: HabitatGraph = {
      rootGroups: [],
      habitats: new Map<UUID, HabitatNode>(),
      groups: new Map<UUID, HabitatGroupNode>()
    }

    // Initialize the graph with all groups and habitats, but without parent-child relationships
    groups.value.forEach((group) => {
      graph.groups.set(group.id, {
        ...group,
        children: group.elements
      })
      group.elements.forEach((habitat) => {
        graph.habitats.set(habitat.id, {
          ...habitat,
          group: group
        })
      })
    })

    // Build the tree by assigning group nodes to their parent habitats
    graph.groups.forEach(({ ...groupNode }) => {
      if (!groupNode.parent_id) {
        graph.rootGroups.push(groupNode)
        return
      }

      const parent = graph.habitats.get(groupNode.parent_id)!
      if (!parent.children) {
        parent.children = []
      }
      parent.children.push(groupNode)
      return
    })

    // Populate upstream habitats for each habitat in the tree
    const populateUpstream = (groupNode: HabitatGroupNode, upstreamHabitats: Habitat[] = []) => {
      groupNode.children.forEach((habitat) => {
        const habitatNode = graph.habitats.get(habitat.id)!
        habitatNode.upstream = [...upstreamHabitats]
        habitatNode.children?.forEach((childGroup) => {
          const childGroupNode = graph.groups.get(childGroup.id)!
          populateUpstream(childGroupNode, [...upstreamHabitats, habitat])
        })
      })
    }

    graph.rootGroups.forEach((rootGroup) => {
      populateUpstream(rootGroup)
    })

    return graph
  })

  function habitatDependencies(habitatId: string) {
    return habitatGraph.value.habitats.get(habitatId)?.upstream ?? []
  }

  function getGroup(id: UUID) {
    return habitatGraph.value.groups.get(id)
  }

  function getHabitat(id: UUID) {
    return habitatGraph.value.habitats.get(id)
  }

  return {
    groups,
    habitatGraph,
    getGroup,
    getHabitat,
    habitatDependencies,
    refetch
  }
}
