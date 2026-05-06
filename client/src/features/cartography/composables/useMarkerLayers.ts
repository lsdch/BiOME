import { SiteItem } from '@/api'
import { hexToArray } from '@/lib/color_brewer'
import type { Layer, PickingInfo } from '@deck.gl/core'
import { ScatterplotLayer, TextLayer } from '@deck.gl/layers'
import { computed, type Ref } from 'vue'
import { parseHex } from 'vuetify/lib/util/colorUtils.mjs'
import { MarkerOptions } from '../components/DeckGlMap.vue'
import type { MarkerLayer } from '../components/layers-manager/map-layers'
import { MarkerSelectionInfo } from './marker-selection'

export type MarkerCluster<Item> = {
  coordinates: { longitude: number; latitude: number }
  count: number
  items: Item[]
}

const DEFAULT_HEX_COLOR = '#FF0000'
const DEFAULT_RADIUS_SCALE_FACTOR = 0.5

function groupItemsByCoordinate<
  Item extends { coordinates: { latitude: number; longitude: number } }
>(items: Item[]) {
  const groups = new Map<string, MarkerCluster<Item>>()

  items.forEach((item) => {
    const { latitude, longitude } = item.coordinates
    const key = `${longitude}:${latitude}`
    const existing = groups.get(key)

    if (existing) {
      existing.items.push(item)
      existing.count += 1
      return
    }

    groups.set(key, {
      coordinates: { longitude, latitude },
      count: 1,
      items: [item]
    })
  })

  return [...groups.values()]
}

export function useMarkerLayers<Item extends SiteItem>(
  props: {
    markerLayers?: MarkerLayer<Item>[]
    markerOptions?: MarkerOptions
  },
  ctx: {
    selected: Ref<MarkerSelectionInfo<Item> | undefined>
    select: (info: MarkerSelectionInfo<Item>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
  }
) {
  const clusterVisible = computed(() => {
    return ctx.currentZoom.value >= (props.markerOptions?.cluster.labelZoomThreshold ?? 8)
  })

  const markerDeckLayers = computed<Layer[]>(() => {
    const activeLayers = (props.markerLayers ?? []).filter(
      (layer) => layer.active && (layer.data?.length ?? 0) > 0
    )

    const selectedCoord =
      ctx.selected.value?.type === 'item' ? ctx.selected.value?.info.coordinate : undefined

    return activeLayers.flatMap((layer, index) => {
      const groups = groupItemsByCoordinate(layer.data ?? []).filter(
        ({ coordinates: { longitude: lon, latitude: lat } }) =>
          !selectedCoord || (lon !== selectedCoord[0] && lat !== selectedCoord[1])
      ) // Exclude the selected item from the regular layers, it will be rendered in a separate layer with a different style
      const simpleRadius = Math.max(1, Number(layer.config.radius ?? 4))
      const clusterRadius = (count: number) =>
        simpleRadius +
        simpleRadius *
          (props.markerOptions?.cluster.radiusScaleFactor ?? DEFAULT_RADIUS_SCALE_FACTOR) *
          count
      const fillColor = hexToArray(parseHex(layer.config.fillColor ?? DEFAULT_HEX_COLOR))
      const strokeColor = hexToArray(parseHex(layer.config.color ?? DEFAULT_HEX_COLOR))

      const singleItems = groups
        .filter((group) => group.count === 1)
        .flatMap((group) => group.items)
      const clusteredGroups = groups.filter((group) => group.count > 1)

      const layers: Layer[] = []

      if (singleItems.length) {
        layers.push(
          new ScatterplotLayer<Item>({
            id: `markers-${index}-single`,
            data: singleItems,
            pickable: true,
            stroked: true,
            filled: true,
            radiusUnits: 'pixels',
            lineWidthUnits: 'pixels',
            lineWidthMinPixels: Number(layer.config.weight ?? 1),
            radiusMinPixels: simpleRadius,
            getRadius: () => simpleRadius,
            getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude],
            getFillColor: () => fillColor,
            getLineColor: () => strokeColor,
            onClick: (info) => ctx.select({ type: 'item', info }),
            onHover: props.markerOptions?.tooltips
              ? (info: PickingInfo<Item>) => {
                  const object = info.object

                  if (!object) {
                    ctx.hoverTooltip.value = undefined
                    return false
                  }

                  ctx.hoverTooltip.value = {
                    x: info.x,
                    y: info.y,
                    text: object.name ?? object.code
                  }
                  return true
                }
              : undefined
          })
        )
      }

      if (clusteredGroups.length) {
        layers.push(
          new ScatterplotLayer<MarkerCluster<Item>>({
            id: `markers-${index}-cluster`,
            data: clusteredGroups,
            pickable: true,
            stroked: true,
            filled: true,
            radiusUnits: 'pixels',
            lineWidthUnits: 'pixels',
            lineWidthMinPixels: Number(layer.config.weight ?? 1),
            radiusMinPixels: clusterRadius(1),
            getRadius: ({ count }) => clusterRadius(count),
            getPosition: ({ coordinates: { longitude: lon, latitude: lat } }) => [lon, lat],
            getFillColor: () => fillColor,
            getLineColor: () => strokeColor,
            onClick: (info) => ctx.select({ type: 'cluster', info })
          })
        )

        if (clusterVisible.value) {
          layers.push(
            new TextLayer<MarkerCluster<Item>>({
              id: `markers-${index}-cluster-count`,
              data: clusteredGroups,
              pickable: false,
              billboard: true,
              background: false,
              sizeUnits: 'pixels',
              fontWeight: 800,

              getPosition: ({ coordinates: { longitude: lon, latitude: lat } }) => [lon, lat],
              getText: (group) => String(group.count),
              getSize: 14,
              getColor: [255, 255, 255, 255],
              getTextAnchor: 'middle',
              getAlignmentBaseline: 'center',
              getPixelOffset: [0, 0]
            })
          )
        }
      }

      return layers
    })
  })

  return { markerDeckLayers }
}
