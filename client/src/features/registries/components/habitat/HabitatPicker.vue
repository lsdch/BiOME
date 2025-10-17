<template>
  <v-autocomplete
    v-model="model"
    v-model:search="searchTerm"
    :items
    label="Habitat tags"
    prepend-inner-icon="mdi-tag-multiple"
    color="primary"
    multiple
    chips
    clearable
    closable-chips
    item-title="label"
    item-subtitle="description"
    auto-select-first
    clear-on-select
    placeholder="Start typing to search"
    return-object
  >
    <template #chip="{ item, props }">
      <v-chip closable v-bind="props" @click:close="onDelete(item.raw)" :text="item.title" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        :title="item.title"
        :subtitle="
          graph
            .habitat(item.raw.id)
            .dependencies.map(({ label }) => label)
            .join(' › ')
        "
        v-bind="props"
      >
        <template #title="{ title }">
          <span :class="{ 'font-weight-bold': !graph.habitat(item.raw.id).dependencies?.length }">{{
            title
          }}</span>
        </template>
        <template #subtitle="{ subtitle }">
          <span class="text-caption">{{ subtitle }}</span>
        </template>
        <template #append>
          <v-chip
            class="text-overline"
            color="primary"
            :text="graph.habitat(item.raw.id).group.label"
          />
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
import { HabitatGroup, HabitatRecord, SamplingService } from '@/api'
import { useErrorHandler } from '@/api/responses'
import { computed, reactive, ref } from 'vue'
import { ConnectedGroup, ConnectedHabitat, HabitatsGraph } from './habitat_graph'

const model = defineModel<HabitatRecord[]>({ default: () => reactive([]) })
const searchTerm = ref<string | undefined>(undefined)

const habitatGroups = ref<HabitatGroup[]>(
  await SamplingService.listHabitatGroups().then(
    useErrorHandler((err) => {
      console.error('Failed to fetch habitat groups: ', err)
    })
  )
)

const graph = new HabitatsGraph(habitatGroups.value, false)

function addWithDependencies(habitat: ConnectedHabitat) {
  model.value.push(...(habitat.dependencies?.map(({ id }) => graph.habitat(id)) ?? []), habitat)
  searchTerm.value = undefined
}

function onDelete(item: HabitatRecord) {
  model.value = model.value.filter(
    ({ id }) =>
      id != item.id && graph.habitat(id).dependencies?.find(({ id }) => id == item.id) == undefined
  )
}

function compatibleHabitats(habitats: ConnectedHabitat[], selected: HabitatRecord[]) {
  return habitats.filter(
    ({ id: habitatID, incompatible }) =>
      !selected.find(({ id }) => habitatID == id || incompatible?.find((incomp) => incomp.id == id))
  )
}

function isGroupReachable(group: ConnectedGroup) {
  return model.value.find(({ id }) => group.depends?.id == id)
}

const items = computed<HabitatRecord[]>(() => {
  return Object.values(graph.groups).reduce((acc: HabitatRecord[], g) => {
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
    return Object.values(graph.habitats)
      .filter(
        ({ label, dependencies }) =>
          (dependencies.length ?? 0) > 0 && label.toLowerCase().includes(term)
      )
      .slice(0, 5)
  }
  return []
})
</script>

<style scoped></style>
