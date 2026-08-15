import { hexToArray } from '@/lib/color_brewer'
import type { Layer } from '@deck.gl/core'
import { ScatterplotLayer, TextLayer } from '@deck.gl/layers'
import { computed, MaybeRefOrGetter, toValue, type Ref } from 'vue'
import { parseHex } from 'vuetify/lib/util/colorUtils.mjs'
import { GlobalMarkerOptions } from '../components/DeckGlMap.vue'
import type { MarkerLayer } from '../components/layers-manager/map-layers'
import { MarkerSelection } from './marker-selection'

// export type MarkerCluster<Item> = {
//   coordinates: { longitude: number; latitude: number }
//   count: number
//   items: Item[]
// }

const DEFAULT_HEX_COLOR = '#FF0000'
const DEFAULT_RADIUS_SCALE_FACTOR = 0.5

// function groupItemsByCoordinate<
//   Item extends { coordinates: { latitude: number; longitude: number } }
// >(items: Item[]) {
//   const groups = new Map<string, MarkerCluster<Item>>()

//   items.forEach((item) => {
//     const { latitude, longitude } = item.coordinates
//     const key = `${longitude}:${latitude}`
//     const existing = groups.get(key)

//     if (existing) {
//       existing.items.push(item)
//       existing.count += 1
//       return
//     }

//     groups.set(key, {
//       coordinates: { longitude, latitude },
//       count: 1,
//       items: [item]
//     })
//   })

//   return [...groups.values()]
// }

/*
  Instantiates a marker layer for the given MarkerLayer specification
  and returns an array of Deck.gl layers (ScatterplotLayer and optionally TextLayer) to be rendered on the map.
  It uses the provided context to handle selection, hover tooltips, and marker options.
*/
export function instanciateMarkerLayer<MarkerData>(
  layer: MarkerLayer<MarkerData>,
  index: number | string,
  ctx: {
    select: (info: MarkerSelection<MarkerData>) => void
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
    showClusterText?: boolean
  }
) {
  // const radiusFn = (item: MarkerData) => {
  //   const baseRadius = layer.config.baseRadius ?? 5
  //   return (
  //     baseRadius +
  //     (layer.radius?.(item) ?? 0) *
  //       (ctx.markerOptions?.cluster.radiusScaleFactor ?? DEFAULT_RADIUS_SCALE_FACTOR)
  //   )
  // }
  const radiusFn = (item: MarkerData) => {
    const baseRadius = layer.config.baseRadius ?? 5
    // const value = layer.radius?.(item) ?? 0
    const scaleFactor = layer.config.radiusScaleFactor ?? DEFAULT_RADIUS_SCALE_FACTOR

    const densityFactor = Math.pow(2, (layer.resolution - 12) * 0.5)

    const effectiveValue = (layer.radius?.(item) ?? 0) * densityFactor

    return baseRadius + Math.log1p(effectiveValue) * scaleFactor * 10
    // return baseRadius + Math.log1p(value) * scaleFactor
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

  const layers: (ScatterplotLayer<MarkerData> | TextLayer<MarkerData>)[] = []

  // if (singleItems.length) {
  layers.push(
    new ScatterplotLayer<MarkerData>({
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
      radiusMaxPixels: 50,
      getPosition: layer.getPosition,
      getFillColor: () => fillColor,
      getLineColor: () => strokeColor,
      updateTriggers: {
        getRadius: [layer.config.radiusScaleFactor, layer.config.baseRadius],
        getFillColor: [layer.config.fillColor],
        getLineColor: [layer.config.color]
      },
      onClick: (info) => {
        const [lng, lat] = layer.getPosition(info.object)
        return ctx.select({
          type: 'marker',
          info,
          params: layer.filters,
          resolution: layer.resolution,
          coordinates: { latitude: lat, longitude: lng }
        })
      }
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

  if (ctx.showClusterText && layer.getText) {
    layers.push(
      new TextLayer<MarkerData>({
        id: `markers-${index}-count`,
        data: layer.data,
        pickable: false,
        billboard: true,
        background: false,
        sizeUnits: 'pixels',
        fontWeight: 800,

        getPosition: layer.getPosition,
        getText: layer.getText,
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

export function useMarkerLayers<MarkerData>(
  props: {
    markerLayers?: MaybeRefOrGetter<MarkerLayer<MarkerData>[]>
  },
  ctx: {
    select: (info: MarkerSelection<MarkerData>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
    markerOptions?: MaybeRefOrGetter<GlobalMarkerOptions>
  }
) {
  const showClusterText = computed(() => {
    return ctx.currentZoom.value >= (toValue(ctx.markerOptions)?.cluster.labelZoomThreshold ?? 8)
  })

  const markerDeckLayers = computed<Layer[]>(() => {
    const markerLayers = toValue(props.markerLayers)
    if (!markerLayers) return []

    const activeLayers = markerLayers.filter(
      (layer) => layer.active && (layer.data?.length ?? 0) > 0
    )

    return activeLayers.flatMap((layer, index) => {
      return instanciateMarkerLayer(layer, index, {
        select: ctx.select,
        hoverTooltip: ctx.hoverTooltip,
        showClusterText: showClusterText.value
      })
    })
  })

  return { markerDeckLayers }
}
