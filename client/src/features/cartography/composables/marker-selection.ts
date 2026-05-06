import { SiteItem } from '@/api'
import { HexagonLayerPickingInfo } from '@deck.gl/aggregation-layers'
import { PickingInfo } from '@deck.gl/core'
import { ColumnLayer, IconLayer } from '@deck.gl/layers'
import { computed, ref } from 'vue'
import { MarkerCluster } from './useMarkerLayers'

export type MarkerSelectionInfo<Item> =
  | { type: 'site'; info: Item }
  | {
      type: 'item'
      info: PickingInfo<Item>
    }
  | {
      type: 'cluster'
      info: PickingInfo<MarkerCluster<Item>>
    }
  | {
      type: 'hexagon'
      info: HexagonLayerPickingInfo<Item>
    }

export function useMarkerSelection<Item extends SiteItem>() {
  const selected = ref<MarkerSelectionInfo<Item>>()
  const clickHandledByLayer = ref(false)

  function select(info: MarkerSelectionInfo<Item>) {
    clickHandledByLayer.value = true
    console.debug('marker-selection: select()', info)
    selected.value = info
  }

  function clear() {
    if (clickHandledByLayer.value) {
      clickHandledByLayer.value = false
      return
    }
    selected.value = undefined
  }

  const highlightLayer = computed(() => {
    if (!selected.value || selected.value.type === 'site') return []

    // Hexagon picks expose aggregated points on the picked object.
    if (selected.value.type === 'hexagon') {
      return [makeHexagonHighlightLayer(selected.value.info as HexagonLayerPickingInfo<Item>)]
    }

    // Otherwise it's a regular marker/item pick.
    if (selected.value.info.object) {
      return [makeIconLayer(selected.value.info, MARKER_ICON_SIZE)]
    }

    return []
  })

  const THIRD_PI = Math.PI / 3
  const DIST_X = 2 * Math.sin(THIRD_PI)
  const DIST_Y = 1.5
  function getHexbinCentroid([i, j]: [number, number], radius: number): [number, number] {
    return [(i + (j & 1) / 2) * radius * DIST_X, j * radius * DIST_Y]
  }

  function makeHexagonHighlightLayer(
    info: HexagonLayerPickingInfo<Item>
  ): ColumnLayer<HexagonLayerPickingInfo<Item>> {
    console.info('Creating highlight layer for hexagon with position', info.coordinate)
    return new ColumnLayer<HexagonLayerPickingInfo<Item>>({
      id: 'marker-hex-highlight',
      data: [info],
      diskResolution: 6,
      // Use geographic coordinates and a visible radius (meters)
      coordinateSystem: 'cartesian',
      //   coordinateOrigin: [180, 90, 0],
      extruded: false,
      radius: 200000,
      filled: true,
      getLineColor: () => [255, 0, 0],
      getLineWidth: () => 2,
      getFillColor: () => [255, 0, 0, 128],
      getElevation: () => 0,
      getPosition: (d: HexagonLayerPickingInfo<Item>) => {
        return d.object?.position as [number, number]
      }
    })
  }

  function makeIconLayer(
    info: PickingInfo<Item> | PickingInfo<MarkerCluster<Item>>,
    size = MARKER_ICON_SIZE
  ): IconLayer<Item | MarkerCluster<Item>> {
    return new IconLayer<Item | MarkerCluster<Item>>({
      id: `marker-selected`,
      data: [info.object],
      pickable: true,
      billboard: true,
      sizeUnits: 'pixels',
      sizeScale: 1,
      getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude],
      getIcon: () => ({
        url: createMarkerIcon(),
        width: size,
        height: size,
        anchorX: size / 2,
        anchorY: size,
        mask: false
      }),
      getSize: () => size
    })
  }

  return {
    selected,
    select,
    clear,
    highlightLayer
  }
}

const MARKER_ICON_SIZE = 48

function preloadImage(url: string): void {
  // Preload the image to avoid lazy texture initialization in WebGL
  const img = new Image()
  img.src = url
}

function createMarkerIcon(): string {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="${MARKER_ICON_SIZE}" height="${MARKER_ICON_SIZE}">
      <path d="M12 2C8.14 2 5 5.14 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.86-3.14-7-7-7zm0 9.5A2.5 2.5 0 1 1 12 6a2.5 2.5 0 0 1 0 5.5z" stroke="white" stroke-width="0.8" stroke-linejoin="round" stroke-linecap="round" fill="#d8420b" />
    </svg>
  `
  const url = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
  // Preload the image to avoid lazy texture initialization warning in WebGL
  preloadImage(url)
  return url
}
