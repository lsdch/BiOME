import { H3CellWithRichness, ErrorModel } from '@/api'
import { UseQueryReturnType } from '@tanstack/vue-query'
import { computed, MaybeRef, shallowReactive, unref } from 'vue'
import { BaseLayerSpec } from './map-layers'
import { ItemWithCoordinates } from '../../coordinates'
import { cellToLatLng } from 'h3-js'

export type CellMarkerData = ItemWithCoordinates & {
  h3_index: string
  occurrences_count: number
  samplings_count: number
}

const data = shallowReactive(new Map<UUID, UseQueryReturnType<H3CellWithRichness[], ErrorModel>>())

const anyLoading = computed(
  () => data.size > 0 && data.values().some((remote) => remote.isPending.value)
)

const allPending = computed(
  () => data.size > 0 && data.values().every((remote) => remote.isPending.value)
)

function registerLayer(
  dataFeed: MaybeRef<BaseLayerSpec>,
  remote: UseQueryReturnType<H3CellWithRichness[], ErrorModel>
) {
  data.set(unref(dataFeed).id, remote)
}

function deleteLayer(id: UUID) {
  data.delete(id)
}

function layerData(id: UUID): UseQueryReturnType<H3CellWithRichness[], ErrorModel> | undefined {
  return data.get(id)
}

export function H3CellToMarkerData({
  h3_index,
  occurrences_count,
  samplings_count
}: H3CellWithRichness): CellMarkerData {
  const [lat, lng] = cellToLatLng(h3_index)
  return {
    coordinates: {
      latitude: lat,
      longitude: lng
    },
    h3_index,
    occurrences_count,
    samplings_count
  }
}

export function useLayerData() {
  return {
    data,
    registerLayer,
    deleteLayer,
    anyLoading,
    allPending,
    layerData,
    H3CellToMarkerData
  }
}
