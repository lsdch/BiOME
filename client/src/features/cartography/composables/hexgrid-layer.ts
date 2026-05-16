import { SiteWithOccurrences } from '@/api'
import { paletteRGB, RGBToArray } from '@/lib/color_brewer'
import { HexagonLayer } from '@deck.gl/aggregation-layers'
import { computed, ref, Ref } from 'vue'
import { bindingLabels, getBindingFn, hexagonLayerColorBinding } from '../bindings'
import {
  HexgridLayer,
  HexLayerSpec,
  markerLayerFromSpec
} from '../components/layers-manager/map-layers'
import { MarkerSelectionInfo } from './marker-selection'
import { instanciateMarkerLayer } from './useMarkerLayers'
import { GlobalMarkerOptions } from '../components/DeckGlMap.vue'

export function hexgridLayerFromSpec<Item extends SiteWithOccurrences>(
  { name, active, config, colorBinding, markers }: HexLayerSpec,
  data?: Item[]
): HexgridLayer<Item> {
  return {
    name,
    active,
    config: {
      ...config,
      colorRange: paletteRGB(config.colorRange ?? 'Viridis')
    },
    colorBinding,
    data,
    markers
  }
}

export function useHexgridLayer<Item extends SiteWithOccurrences>(
  props: {
    hexgrid: () => HexgridLayer<Item> | undefined
    markerOptions?: () => GlobalMarkerOptions
  },
  ctx: {
    selected: Ref<MarkerSelectionInfo<Item, any> | undefined>
    select: (info: MarkerSelectionInfo<Item, any>) => void
    currentZoom: Ref<number>
    hoverTooltip: Ref<{ x: number; y: number; text: string } | undefined>
  }
) {
  const showMarkers = computed(() => {
    const hexgrid = props.hexgrid()
    return (
      !!hexgrid &&
      hexgrid.active &&
      !!hexgrid.data?.length &&
      ctx.currentZoom.value >= hexgrid.markers.minZoom
    )
  })

  const hexgridColorDomain = ref<{ min: number; max: number }>()

  const hexgridColorRange = computed<Array<[number, number, number]>>(() => {
    const hexgrid = props.hexgrid()
    if (!hexgrid?.config.colorRange) {
      return paletteRGB('Viridis').map((rgb) => RGBToArray(rgb) as [number, number, number])
    }
    return (hexgrid.config.colorRange as any[]).map((color) => {
      if (Array.isArray(color)) {
        return color as [number, number, number]
      }
      return RGBToArray(color) as [number, number, number]
    })
  })

  const hexgridLayer = computed(() => {
    console.debug('Computing hex layer')
    const hexgrid = props.hexgrid()
    if (!hexgrid?.active || !hexgrid.data?.length) return undefined

    const colorRange = (hexgrid.config.colorRange ?? paletteRGB('Viridis')).map((rgb) =>
      RGBToArray(rgb)
    )

    const binding = hexgrid.colorBinding.binding

    const layer = new HexagonLayer<Item>({
      id: 'hex-grid',
      data: hexgrid.data,
      gpuAggregation: false,
      pickable: true,
      extruded: false,
      coverage: hexgrid.config.coverage,
      coordinateSystem: 'lnglat',
      highlightColor: [255, 0, 0, 128],
      autoHighlight: hexgrid.config.hover.highlight && !showMarkers.value,
      opacity: showMarkers.value ? 0.05 : Number(hexgrid.config.opacity ?? 0.8),
      radius: Math.max(100, Number(hexgrid.config.radius) * 1000),
      colorRange: colorRange,
      elevationScale: 100,
      ...hexagonLayerColorBinding(hexgrid.colorBinding),
      // getElevationValue: hexagonLayerColorBinding(hexgrid.colorBinding).getColorValue,
      updateTriggers: {
        getColorWeight: [hexgrid.colorBinding.binding, hexgrid.colorBinding.log],
        getColorValue: [hexgrid.colorBinding.binding, hexgrid.colorBinding.log]
        // getElevationValue: [hexgrid.colorBinding.binding, hexgrid.colorBinding.log]
      },
      getPosition: (item) => [item.coordinates.longitude, item.coordinates.latitude],
      onClick: (info) => {
        ctx.select({ type: 'hexagon', info })
        return true
      },
      onHover: hexgrid.config.hover.showTooltip
        ? (info) => {
            const object = info.object as { count?: number; points?: Item[] } | undefined

            if (!object) {
              ctx.hoverTooltip.value = undefined
              return false
            }

            const count = getBindingFn(binding ?? 'sites')(object.points ?? []) ?? object.count ?? 0
            if (!count) {
              ctx.hoverTooltip.value = undefined
              return false
            }

            ctx.hoverTooltip.value = {
              x: info.x,
              y: info.y,
              text: `${count} ${bindingLabels[binding ?? 'sites'] ?? binding}`
            }
            return true
          }
        : undefined,
      onSetColorDomain([min, max]) {
        if (hexgrid.colorBinding.log) {
          hexgridColorDomain.value = {
            min: Math.round(Math.exp(min) - 1),
            max: Math.round(Math.exp(max) - 1)
          }
        } else {
          hexgridColorDomain.value = { min: min, max: max }
        }
      }
    })

    return layer
  })

  const showClusterText = computed(() => {
    return ctx.currentZoom.value >= (props.markerOptions?.().cluster.labelZoomThreshold ?? 8)
  })

  const hexgridMarkersLayer = computed(() => {
    const hexgrid = props.hexgrid()
    if (!hexgrid?.active || !hexgrid.data?.length || !showMarkers.value) return []
    const markerLayer = markerLayerFromSpec(hexgrid.markers, hexgrid.data)
    return instanciateMarkerLayer(markerLayer, 'hexgrid', {
      ...ctx,
      markerOptions: props.markerOptions?.(),
      showClusterText: showClusterText.value
    })
  })

  return { hexgridLayer, hexgridMarkersLayer, hexgridColorDomain, hexgridColorRange }
}
