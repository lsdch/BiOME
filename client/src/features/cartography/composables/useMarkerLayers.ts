import { hexToArray } from '@/lib/color_brewer'
import type { Layer } from '@deck.gl/core'
import { ScatterplotLayer, TextLayer } from '@deck.gl/layers'
import { computed, MaybeRefOrGetter, toValue, type Ref } from 'vue'
import { parseHex } from 'vuetify/lib/util/colorUtils.mjs'
import { GlobalMarkerOptions } from '../components/DeckGlMap.vue'
import type { MarkerLayer } from '../components/layers-manager/map-layers'
import { ItemWithCoordinates } from '../coordinates.ts'
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

export function instanciateMarkerLayer<Item extends ItemWithCoordinates>(
  layer: MarkerLayer<Item>,
  index: number | string,
  ctx: {
    selected: Ref<MarkerSelectionInfo<Item> | undefined>
    select: (info: MarkerSelectionInfo<Item>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
    markerOptions?: GlobalMarkerOptions
    showClusterText?: boolean
  }
) {
  const radiusFn = (item: Item) => {
    const baseRadius = layer.config.radius?.(item) ?? 5
    return (
      baseRadius +
      baseRadius * (ctx.markerOptions?.cluster.radiusScaleFactor ?? DEFAULT_RADIUS_SCALE_FACTOR)
    )
  }
  // const groups = groupItemsByCoordinate(layer.data ?? [])
  // const simpleRadius = Math.max(1, Number(layer.config.radius ?? 4))
  // const clusterRadius = (count: number) =>
  //   simpleRadius +
  //   simpleRadius *
  //     (ctx.markerOptions?.cluster.radiusScaleFactor ?? DEFAULT_RADIUS_SCALE_FACTOR) *
  //     count
  const fillColor = hexToArray(parseHex(layer.config.fillColor ?? DEFAULT_HEX_COLOR))
  const strokeColor = hexToArray(parseHex(layer.config.color ?? DEFAULT_HEX_COLOR))

  // const singleItems = groups.filter((group) => group.count === 1).flatMap((group) => group.items)
  // const clusteredGroups = groups.filter((group) => group.count > 1)

  const layers: Layer[] = []

  // if (singleItems.length) {
  layers.push(
    new ScatterplotLayer<Item>({
      id: `markers-${index}`,
      data: layer.data ?? [],
      pickable: true,
      stroked: true,
      filled: true,
      radiusUnits: 'pixels',
      lineWidthUnits: 'pixels',
      lineWidthMinPixels: Number(layer.config.weight ?? 1),
      radiusMinPixels: 4,
      getRadius: radiusFn,
      getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude],
      getFillColor: () => fillColor,
      getLineColor: () => strokeColor,
      onClick: (info) => ctx.select({ type: 'item', info })
      // onHover: ctx.markerOptions?.tooltips
      //   ? (info: PickingInfo<Item>) => {
      //       const object = info.object

      //       if (!object) {
      //         ctx.hoverTooltip.value = undefined
      //         return false
      //       }

      //       ctx.hoverTooltip.value = {
      //         x: info.x,
      //         y: info.y,
      //         text: object.name ?? object.code
      //       }
      //       return true
      //     }
      //   : undefined
    })
  )

  if (ctx.showClusterText && layer.config.getText) {
    layers.push(
      new TextLayer<Item>({
        id: `markers-${index}-count`,
        data: layer.data,
        pickable: false,
        billboard: true,
        background: false,
        sizeUnits: 'pixels',
        fontWeight: 800,

        getPosition: ({ coordinates: { longitude: lon, latitude: lat } }) => [lon, lat],
        getText: layer.config.getText,
        getSize: 14,
        getColor: [255, 255, 255, 255],
        getTextAnchor: 'middle',
        getAlignmentBaseline: 'center',
        getPixelOffset: [0, 0]
      })
    )
  }

  return layers
}

export function useMarkerLayers<Item extends ItemWithCoordinates>(
  props: {
    markerLayers?: MaybeRefOrGetter<MarkerLayer<Item>[]>
    markerOptions?: MaybeRefOrGetter<GlobalMarkerOptions>
  },
  ctx: {
    selected: Ref<MarkerSelectionInfo<Item, any> | undefined>
    select: (info: MarkerSelectionInfo<Item, any>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
  }
) {
  const showClusterText = computed(() => {
    return ctx.currentZoom.value >= (toValue(props.markerOptions)?.cluster.labelZoomThreshold ?? 8)
  })

  const markerDeckLayers = computed<Layer[]>(() => {
    const markerLayers = toValue(props.markerLayers)
    if (!markerLayers) return []

    const activeLayers = markerLayers.filter(
      (layer) => layer.active && (layer.data?.length ?? 0) > 0
    )

    return activeLayers.flatMap((layer, index) => {
      return instanciateMarkerLayer(layer, index, {
        selected: ctx.selected,
        select: ctx.select,
        currentZoom: ctx.currentZoom,
        hoverTooltip: ctx.hoverTooltip,
        markerOptions: toValue(props.markerOptions),
        showClusterText: showClusterText.value
      })
    })
  })

  return { markerDeckLayers }
}
