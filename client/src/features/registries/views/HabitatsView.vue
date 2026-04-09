<template>
  <v-container>
    <HabitatGroupBox v-for="group in trees" :key="group.id" :group="group" open />
  </v-container>
</template>

<script setup lang="ts">
import { listHabitatGroupsOptions } from '@/api/gen/@tanstack/vue-query.gen'
import HabitatGroupBox, {
  HabitatGroupNode,
  HabitatNode
} from '@/features/registries/components/HabitatGroupBox.vue'
import { useQuery } from '@tanstack/vue-query'
import { computed } from 'vue'

const { data: groups } = useQuery({
  ...listHabitatGroupsOptions(),
  initialData: []
})

const trees = computed<HabitatGroupNode[]>(() => {
  const groupNodes = groups.value.map<HabitatGroupNode>((group) => ({
    ...group,
    type: 'group',
    children: group.elements.map<HabitatNode>((habitat) => ({
      ...habitat,
      type: 'habitat',
      children: []
    }))
  }))

  const habitatsById = new Map<string, HabitatNode>()
  groupNodes.forEach((group) => {
    group.children.forEach((habitat) => {
      habitatsById.set(habitat.id, habitat)
    })
  })

  const rootGroups: HabitatGroupNode[] = []

  groupNodes.forEach((group) => {
    const parentHabitatId = group.depends?.id
    const parentHabitat = parentHabitatId ? habitatsById.get(parentHabitatId) : undefined

    if (!parentHabitat) {
      rootGroups.push(group)
      return
    }

    const childGroups = parentHabitat.children ?? []
    parentHabitat.children = childGroups.concat(group)
  })

  return rootGroups
})
</script>

<style scoped></style>
