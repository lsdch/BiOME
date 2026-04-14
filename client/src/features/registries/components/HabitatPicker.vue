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
      <v-chip closable v-bind="props" @click:close="onDelete(item)" :text="item.label" />
    </template>
    <template #item="{ item, props }">
      <v-list-item
        :title="item.label"
        :subtitle="
          habitatDependencies(item.id)
            .map(({ label }) => label)
            .join(' › ')
        "
        v-bind="props"
      >
        <template #title="{ title }">
          <span :class="{ 'font-weight-bold': !habitatDependencies(item.id).length }">{{
            title
          }}</span>
        </template>
        <template #subtitle="{ subtitle }">
          <span class="text-caption">{{ subtitle }}</span>
        </template>
        <template #append>
          <v-chip class="text-overline" color="primary" :text="habitatGroupLabel(item.id)" />
        </template>
      </v-list-item>
    </template>
    <template #append-item>
      <v-list v-if="quickSelect.length > 0">
        <v-list-subheader> Quick select </v-list-subheader>
        <v-list-item v-for="item in quickSelect" :key="item.id" @click="addWithDependencies(item)">
          <template v-for="dep in item.upstream ?? []" :key="dep.id">
            <v-chip class="" :text="dep.label" variant="text" />
            <v-icon icon="mdi-chevron-right" />
          </template>
          <v-chip class="ma-1" :text="item.label" color="primary" />
          <template #append>
            <v-chip class="text-overline" color="primary" :text="habitatGroupLabel(item.id)" />
          </template>
        </v-list-item>
      </v-list>
    </template>
  </v-autocomplete>
</template>

<script setup lang="ts">
import { HabitatGroup, HabitatRecord } from '@/api'
import { HabitatNode, useHabitats } from '../composables/habitats'
import { computed, reactive, ref } from 'vue'

const model = defineModel<HabitatRecord[]>({ default: () => reactive([]) })
const searchTerm = ref<string | undefined>(undefined)

const { groups: habitatGroups, habitatsMap, habitatDependencies, habitatGroupLabel } = useHabitats()

function addWithDependencies(habitat: HabitatNode) {
  const existingIds = new Set(model.value.map(({ id }) => id))
  const toAdd = [...(habitat.upstream ?? []), habitat]

  toAdd.forEach((item) => {
    if (!existingIds.has(item.id)) {
      model.value.push(item)
      existingIds.add(item.id)
    }
  })

  searchTerm.value = undefined
}

function onDelete(item: HabitatRecord) {
  model.value = model.value.filter(
    ({ id }) => id != item.id && !habitatDependencies(id).some(({ id: depId }) => depId == item.id)
  )
}

function compatibleHabitats(habitats: HabitatNode[], selected: HabitatRecord[]) {
  return habitats.filter(
    ({ id: habitatID, incompatible }) =>
      !selected.find(({ id }) => habitatID == id || incompatible?.find((incomp) => incomp.id == id))
  )
}

function isGroupReachable(group: HabitatGroup) {
  return model.value.find(({ id }) => group.depends?.id == id)
}

const items = computed<HabitatRecord[]>(() => {
  return habitatGroups.value.reduce((acc: HabitatRecord[], g) => {
    const habitats = g.elements
      .map((habitat) => habitatsMap.value.get(habitat.id))
      .filter((habitat): habitat is HabitatNode => Boolean(habitat))

    if (g.depends == undefined || isGroupReachable(g)) {
      acc = acc.concat(compatibleHabitats(habitats, model.value))
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
    return habitatsMap.value
      .values()
      .filter(
        ({ label, upstream }) => (upstream?.length ?? 0) > 0 && label.toLowerCase().includes(term)
      )
      .take(5)
      .toArray()
  }
  return []
})
</script>

<style scoped></style>
