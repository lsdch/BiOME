import { ErrorModel, SiteWithOccurrences } from '@/api'
import { UseQueryReturnType } from '@tanstack/vue-query'
import { computed, MaybeRef, shallowReactive, unref } from 'vue'
import { BaseLayerSpec } from './map-layers'

const data = shallowReactive(new Map<UUID, UseQueryReturnType<SiteWithOccurrences[], ErrorModel>>())

const anyLoading = computed(
  () => data.size > 0 && data.values().some((remote) => remote.isPending.value)
)

const allPending = computed(
  () => data.size > 0 && data.values().every((remote) => remote.isPending.value)
)

function registerLayer(
  dataFeed: MaybeRef<BaseLayerSpec>,
  remote: UseQueryReturnType<SiteWithOccurrences[], ErrorModel>
) {
  data.set(unref(dataFeed).id, remote)
}

function deleteLayer(id: UUID) {
  data.delete(id)
}

function layerData(id: UUID): UseQueryReturnType<SiteWithOccurrences[], ErrorModel> | undefined {
  return data.get(id)
}

export function useLayerData() {
  return {
    data,
    registerLayer,
    deleteLayer,
    anyLoading,
    allPending,
    layerData
  }
}
