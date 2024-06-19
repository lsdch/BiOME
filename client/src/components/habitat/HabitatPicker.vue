<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerm"
    label="Habitat tags"
    prepend-inner-icon="mdi-tag-multiple"
    color="primary"
    multiple
    chips
    clearable
    closable-chips
    :items="items"
    item-title="label"
    item-value="label"
    return-object
    density="compact"
    auto-select-first
    clear-on-select
    placeholder="Start typing to search"
  >
    <template #chip="{ item, props }">
      <v-chip closable v-bind="props" @click:close="onDelete(item.raw)" :text="item.title" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        :title="item.title"
        :subtitle="item.raw.dependencies?.map(({ label }) => label).join(' › ')"
        v-bind="props"
      >
        <template #title="{ title }">
          <span :class="{ 'font-weight-bold': !item.raw.dependencies?.length }">{{ title }}</span>
        </template>
        <template #subtitle="{ subtitle }">
          <span class="text-caption">{{ subtitle }}</span>
        </template>
        <template #append>
          <v-chip class="text-overline" color="primary" :text="item.raw.group.label" />
        </template>
      </v-list-item>
    </template>
    <template #append-item>
      <v-list v-if="quickSelect.length > 0">
        <v-list-subheader> Quick select </v-list-subheader>
        <v-list-item v-for="item in quickSelect" :key="item.id" @click="addWithDependencies(item)">
          <template v-for="dep in item.dependencies" :key="dep.id">
            <v-chip class="" :text="dep.label" variant="text" />
            <v-icon icon="mdi-chevron-right" />
          </template>
          <v-chip class="ma-1" :text="item.label" color="primary" />
          <template #append>
            <v-chip class="text-overline" color="primary" :text="item.group.label" />
          </template>
        </v-list-item>
      </v-list>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { HabitatGroup, HabitatRecord, LocationService } from '@/api'
import { computed, ref } from 'vue'
import { ConnectedGroup, ConnectedHabitat, Dependencies, indexGroups } from './habitat_graph'
import { handleErrors } from '@/api/responses'

function addWithDependencies(habitat: ConnectedHabitat) {
  model.value.push(
    ...(habitat.dependencies?.map(({ id }) => habitatsGraph.value.habitats[id]) ?? []),
    habitat
  )
  searchTerm.value = undefined
}

function onDelete(item: HabitatRecord & Dependencies) {
  model.value = model.value.filter(
    ({ dependencies }) => dependencies?.find(({ id }) => id == item.id) == undefined
  )
}

const habitatGroups = ref<HabitatGroup[]>(
  await LocationService.listHabitatGroups().then(
    handleErrors((err) => {
      console.error('Failed to fetch habitat groups: ', err)
    })
  )
)

const habitatsGraph = ref(indexGroups(habitatGroups.value))

const model = ref<ConnectedHabitat[]>([])
const searchTerm = ref<string | undefined>(undefined)

function compatibleHabitats(habitats: ConnectedHabitat[], selected: ConnectedHabitat[]) {
  return habitats.filter(
    ({ id: habitatID, incompatible }) =>
      !selected.find(({ id }) => habitatID == id || incompatible?.find((incomp) => incomp.id == id))
  )
}

function isGroupReachable(group: ConnectedGroup) {
  return model.value.find(({ id }) => group.depends?.id == id)
}

const items = computed<ConnectedHabitat[]>(() => {
  return Object.values(habitatsGraph.value.groups).reduce((acc: ConnectedHabitat[], g) => {
    if (g.depends == undefined || isGroupReachable(g)) {
      acc = acc.concat(compatibleHabitats(g.elements, model.value))
    }
    return acc
  }, [])
})

/**
 * List of deep habitat tags that are not immediately accessible given the current tags selection
 */
const quickSelect = computed(() => {
  if (searchTerm.value != undefined && searchTerm.value.length > 0) {
    const term = searchTerm.value.toLowerCase()
    return Object.values(habitatsGraph.value.habitats)
      .filter(
        ({ label, dependencies }) =>
          (dependencies?.length ?? 0) > 0 && label.toLowerCase().includes(term)
      )
      .slice(0, 5)
  }
  return []
})
</script>

<style scoped></style>
