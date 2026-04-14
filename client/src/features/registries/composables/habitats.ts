import { HabitatGroup, HabitatRecord } from '@/api'
import { listHabitatGroupsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

export type HabitatNode = HabitatRecord & {
  children?: HabitatGroupNode[]
  upstream?: HabitatNode[]
}

export type HabitatGroupNode = HabitatGroup & {
  children: HabitatNode[]
}

export function useHabitats() {
  const { data: groups } = useQuery({
    ...listHabitatGroupsOptions(),
    initialData: []
  })

  const habitatGraph = computed<HabitatGroupNode[]>(() => {
    const habitatsById = new Map<string, HabitatNode>()

    // Create and register all habitat nodes
    groups.value.forEach((group) => {
      group.elements.forEach((habitat) => {
        const node: HabitatNode = {
          ...habitat,
          children: []
        }
        habitatsById.set(habitat.id, node)
      })
    })

    // Create group nodes and assign habitat nodes as their children
    const groupNodes: HabitatGroupNode[] = groups.value.map((group) => ({
      ...group,
      children: group.elements.map((habitat) => habitatsById.get(habitat.id)!)
    }))

    const rootGroups: HabitatGroupNode[] = []

    // Build the tree by assigning group nodes to their parent habitats
    groupNodes.forEach((groupNode) => {
      const parentHabitatId = groupNode.depends?.id
      const parentHabitat = parentHabitatId ? habitatsById.get(parentHabitatId) : undefined

      if (!parentHabitat) {
        rootGroups.push(groupNode)
        return
      }

      if (!parentHabitat.children) {
        parentHabitat.children = []
      }
      parentHabitat.children.push(groupNode)
    })

    // Populate upstream habitats for each habitat in the tree
    const populateUpstream = (
      groupNode: HabitatGroupNode,
      upstreamHabitats: HabitatNode[] = []
    ) => {
      groupNode.children.forEach((habitat) => {
        habitat.upstream = [...upstreamHabitats]
        habitat.children?.forEach((childGroup) => {
          populateUpstream(childGroup, [...upstreamHabitats, habitat])
        })
      })
    }

    rootGroups.forEach((rootGroup) => {
      populateUpstream(rootGroup)
    })

    return rootGroups
  })

  const habitatsMap = computed(() => {
    const map = new Map<string, HabitatNode>()

    const visitGroup = (group: HabitatGroupNode) => {
      group.children.forEach((habitat) => {
        map.set(habitat.id, habitat)
        habitat.children?.forEach(visitGroup)
      })
    }

    habitatGraph.value.forEach(visitGroup)

    return map
  })

  const habitatGroupByHabitatId = computed(() => {
    return groups.value.reduce((acc, group) => {
      group.elements.forEach((habitat) => {
        acc.set(habitat.id, group)
      })
      return acc
    }, new Map<string, HabitatGroup>())
  })

  function habitatDependencies(habitatId: string) {
    return habitatsMap.value.get(habitatId)?.upstream ?? []
  }

  function habitatGroupLabel(habitatId: string) {
    return habitatGroupByHabitatId.value.get(habitatId)?.label ?? ''
  }

  return {
    groups,
    habitatGraph,
    habitatsMap,
    habitatDependencies,
    habitatGroupByHabitatId,
    habitatGroupLabel
  }
}
