import { H3CellWithRichness } from '@/api'
import { hexToArray, paletteRGB, RGBToArray } from '@/lib/color_brewer'
import { H3HexagonLayer } from '@deck.gl/geo-layers'

import { computed, MaybeRefOrGetter, ref, Ref, toValue } from 'vue'
import { bindingLabels, getColorValue } from '../bindings'
import { GlobalMarkerOptions } from '../components/DeckGlMap.vue'
import {
  H3Cell,
  HexgridLayer,
  HexLayerSpec,
  MarkerLayer
} from '../components/layers-manager/map-layers'
import { HexagonSelection, MarkerSelection, MarkerSelectionInfo } from './marker-selection'
import { instanciateMarkerLayer } from './useMarkerLayers'

import { extent } from 'd3-array'
import { ScaleQuantize, scaleQuantize } from 'd3-scale'
import { parseHex } from 'vuetify/lib/util/colorUtils.mjs'

export function hexgridLayerFromSpec<CellData extends H3Cell>(
  spec: HexLayerSpec,
  data?: CellData[]
): HexgridLayer<CellData> {
  return {
    ...spec,
    data
  }
}

// function H3DataToResolution(
//   data: H3CellWithRichness[],
//   resolution: number
// ): Record<string, H3CellWithRichness> {
//   return data.reduce(
//     (acc, cell) => {
//       const parent = cellToParent(cell.h3_index, resolution)
//       if (!acc[parent]) {
//         acc[parent] = {
//           ...cell,
//           h3_index: parent
//         }
//       } else {
//         acc[parent].occurrences_count += cell.occurrences_count
//         acc[parent].samplings_count += cell.samplings_count
//         acc[parent].occurrence_ids = acc[parent].occurrence_ids.concat(cell.occurrence_ids)
//         acc[parent].sampling_ids = acc[parent].sampling_ids.concat(cell.sampling_ids)
//       }
//       return acc
//     },
//     {} as Record<string, H3CellWithRichness>
//   )
// }

const DEFAULT_MIN_ZOOM_THRESHOLD = 8

function showMarkers(hexgrid: HexLayerSpec, currentZoom: number) {
  return (
    !!hexgrid &&
    hexgrid.active &&
    hexgrid.markers.minZoomMode !== 'never' &&
    currentZoom >=
      (hexgrid.markers.minZoomMode === 'auto'
        ? DEFAULT_MIN_ZOOM_THRESHOLD
        : hexgrid.markers.minZoom)
  )
}

export function useHexgridLayer<HexData extends H3CellWithRichness>(
  props: {
    hexgrid: () => HexgridLayer<HexData> | undefined
    // markerOptions?: () => GlobalMarkerOptions
  },
  ctx: {
    selected: Ref<MarkerSelectionInfo<any, HexData, any> | undefined>
    select: (h: HexagonSelection<HexData>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
  }
) {
  // const hexgridColorDomain = ref<{ min: number; max: number }>()
  const colorScale =
    ref<ScaleQuantize<[number, number, number] | [number, number, number, number]>>(scaleQuantize())

  const scaleExtent = ref<[number, number]>([0, 1])
  // const hexgridColorRange = computed<Array<[number, number, number]>>(() => {
  //   const hexgrid = props.hexgrid()
  //   if (!hexgrid?.config.colorRange) {
  //     return paletteRGB('Viridis').map((rgb) => RGBToArray(rgb) as [number, number, number])
  //   }
  //   return (hexgrid.config.colorRange as any[]).map((color) => {
  //     if (Array.isArray(color)) {
  //       return color as [number, number, number]
  //     }
  //     return RGBToArray(color) as [number, number, number]
  //   })
  // })

  const hexgrid = computed(() => props.hexgrid())
  const markersVisible = computed(() =>
    hexgrid.value ? showMarkers(hexgrid.value, ctx.currentZoom.value) : false
  )

  const hexgridLayer = computed(() => {
    console.debug('Computing hex layer')
    const hexgrid = props.hexgrid()
    if (!hexgrid?.active || !hexgrid.data?.length) return undefined

    const colorRange = paletteRGB(hexgrid.config.colorRange ?? 'Viridis').map((rgb) =>
      RGBToArray(rgb)
    )

    const binding = hexgrid.colorBinding.binding
    const data = hexgrid.data
    const { log: logScale } = hexgrid.colorBinding

    if (data.length) {
      scaleExtent.value = extent(data, (d) => getColorValue(d, binding)) as [number, number]
    }

    colorScale.value = scaleQuantize(
      logScale
        ? [Math.log(Math.max(1, scaleExtent.value[0])), Math.log(scaleExtent.value[1])]
        : scaleExtent.value,
      colorRange
    )

    const opacity = computed(() => {
      return markersVisible.value ? 0.25 : Number(hexgrid.config.opacity ?? 0.8)
    })

    const fillColor = (cell: HexData) => {
      let color: [number, number, number] | [number, number, number, number] = [255, 255, 255]
      if (binding === 'constant') {
        color = hexToArray(parseHex(hexgrid.config.fillColor))
      } else {
        const value = getColorValue(cell, binding)
        const scaledValue = logScale ? Math.log(Math.max(1, value)) : value
        color = colorScale.value(scaledValue)
      }
      const [r, g, b] = color
      return [r, g, b, 255 * opacity.value] as const
    }

    const layer = new H3HexagonLayer<HexData>({
      id: 'hex-grid',
      data,
      pickable: true,
      getHexagon: (item) => item.h3_index,
      extruded: false,
      coverage: hexgrid.config.coverage,
      coordinateSystem: 'lnglat',
      highlightColor: [255, 0, 0, 128],
      autoHighlight: hexgrid.config.hover.highlight && !markersVisible.value,
      highlightedObjectIndex:
        ctx.selected.value?.type === 'hexagon' ? ctx.selected.value.info.index : undefined,
      // opacity: showMarkers.value ? 0.05 : Number(hexgrid.config.opacity ?? 0.8),
      // radius: Math.max(100, Number(hexgrid.config.radius) * 1000),
      // colorRange: colorRange,
      elevationScale: 100,
      getLineColor(cell) {
        return [255, 255, 255, 50]
      },
      stroked: true,
      getLineWidth: 1,
      lineWidthUnits: 'pixels',
      getFillColor: fillColor,
      updateTriggers: {
        getFillColor: [
          logScale,
          scaleExtent,
          opacity.value,
          colorRange,
          hexgrid.config.fillColor,
          hexgrid.colorBinding
        ]
      },
      onClick: (info) => {
        ctx.select({
          type: 'hexagon',
          info,
          params: hexgrid.filters,
          resolution: hexgrid.resolution
        })
        return true
      },
      onHover:
        hexgrid.config.hover.showTooltip && binding !== 'constant'
          ? (info) => {
              const object = info.object as HexData | undefined

              if (!object || markersVisible.value) {
                ctx.hoverTooltip.value = undefined
                return false
              }

              const count = getColorValue(object, binding)

              ctx.hoverTooltip.value = {
                x: info.x,
                y: info.y,
                text: `${count} ${binding}`
              }
              return true
            }
          : undefined
      // getElevationValue: hexagonLayerColorBinding(hexgrid.colorBinding).getColorValue,
      // getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude],
      // onSetColorDomain([min, max]) {
      //   if (hexgrid.colorBinding.log) {
      //     hexgridColorDomain.value = {
      //       min: Math.round(Math.exp(min) - 1),
      //       max: Math.round(Math.exp(max) - 1)
      //     }
      //   } else {
      //     hexgridColorDomain.value = { min: min, max: max }
      //   }
      // }
    })

    return layer
  })

  // const showClusterText = computed(() => {
  //   return ctx.currentZoom.value >= (props.markerOptions?.().cluster.labelZoomThreshold ?? 8)
  // })

  // const markerLayer = computed(async () => {
  //   const hexgrid = props.hexgrid()
  //   if (!hexgrid?.active || !hexgrid.data?.length || !showMarkers(hexgrid, ctx.currentZoom.value))
  //     return undefined
  //   const { data, suspense } = useQuery(
  //     listOccurrencesH3Options({
  //       path: { resolution: 12 },
  //       query: hexgrid.filters
  //     })
  //   )
  //   await suspense()
  //   return markerLayerFromSpec(hexgrid.markers, (data.value ?? []).map(H3CellToMarkerData))
  // })

  // const hexgridMarkersLayer = computed(async () => {
  //   const layer = await markerLayer.value
  //   return layer
  //     ? instanciateMarkerLayer(layer, 'hexgrid', {
  //         ...ctx,
  //         markerOptions: props.markerOptions?.(),
  //         showClusterText: showClusterText.value
  //       })
  //     : []
  // })

  const colorDomain = computed(() => {
    const [min, max] = scaleExtent.value
    return { min, max }
  })

  return {
    hexgridLayer,
    // hexgridMarkersLayer,
    colorDomain,
    colorScale
  }
}

export function useHexgridMarkersLayer<CellData extends H3CellWithRichness>(
  props: {
    hexgrid: () => HexLayerSpec | undefined
    markers: () => MarkerLayer<CellData> | undefined
  },
  ctx: {
    select: (info: MarkerSelection<CellData>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
    markerOptions?: MaybeRefOrGetter<GlobalMarkerOptions>
  }
) {
  const showClusterText = computed(() => {
    return ctx.currentZoom.value >= (toValue(ctx.markerOptions)?.cluster.labelZoomThreshold ?? 8)
  })

  const active = computed(() => {
    const hexgrid = props.hexgrid()
    return !!hexgrid && hexgrid.active && showMarkers(hexgrid, ctx.currentZoom.value)
  })

  const markers = props.markers()

  const markerLayer = computed(() => {
    if (active.value && markers)
      return instanciateMarkerLayer(markers, 'hexgrid', {
        hoverTooltip: ctx.hoverTooltip,
        select: ctx.select,
        showClusterText: showClusterText.value
      })
    else return []
  })

  return { hexMarkersLayer: markerLayer }
}
