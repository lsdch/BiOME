import { HabitatGroup, HabitatRecord } from "@/api";
import { computed, ref } from "vue";

export type SelectedHabitat = HabitatRecord & { group?: string }

const selection = ref<SelectedHabitat | undefined>(undefined)

export function useSelection() {
  function select(habitat: HabitatRecord, group?: HabitatGroup) {
    selection.value = { ...habitat, group: group?.label }
  }

  function isSelected(habitat: HabitatRecord) {
    return computed(() => habitat.id === selection.value?.id)
  }

  function isIncompatibleWithSelection(habitat: HabitatRecord, group?: HabitatGroup) {
    return computed(() => {
      return (selection.value?.incompatible?.find(({ id }) => id === habitat.id)) ||
        (
          selection.value?.group &&
          selection.value?.group == group?.label &&
          selection.value?.id !== habitat.id &&
          group?.exclusive_elements
        )
    })
  }

  return { selection, select, isSelected, isIncompatibleWithSelection }
}