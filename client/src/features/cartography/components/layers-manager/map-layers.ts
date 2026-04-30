import { ScaleBindingSpec } from '@/composables/occurrences'
import { Geocoordinates } from '@/features/cartography/coordinates'
import { brewerPalettes, withOpacity } from '@/lib/color_brewer'
import { UUID } from 'crypto'
import { CircleMarkerOptions } from 'leaflet'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { ScaleBinding } from 'vue-leaflet-hexbin'
import { HabitatRecord, OccurrencesBySiteData } from '@/api'

export type MappingFilters = Overwrite<
  NonNullable<OccurrencesBySiteData['query']>,
  { habitats?: HabitatRecord[] }
>

/**
 * Type representing the filter options for sites in the map layers.
 * - 'All': Show all sites.
 * - 'Sampled': Show only sites that have been sampled.
 * - 'Occurrences': Show only sites that have occurrences.
 */
export type SitesFilter = 'All' | 'Sampled' | 'Occurrences'
export type LayerType = 'markers' | 'hexgrid'
export type BaseLayerSpec = {
  id: UUID
  name?: string
  active: boolean
  filters: MappingFilters
}

export type HexLayerSpec = BaseLayerSpec & {
  type: 'hexgrid'
  bindings: {
    color?: ScaleBindingSpec
    radius?: ScaleBindingSpec
    opacity?: ScaleBindingSpec
  }
  config: HexgridConfigSpec
}

export type MarkerLayerSpec = BaseLayerSpec & {
  type: 'markers'
  ready?: boolean
  clustered: boolean
  config: MarkerConfig
}

export type LayerSpec = BaseLayerSpec &
  (
    | { type: 'markers'; clustered: boolean; config: MarkerConfig }
    | {
        type: 'hexgrid'
        bindings: {
          color?: ScaleBindingSpec
          radius?: ScaleBindingSpec
          opacity?: ScaleBindingSpec
        }
        config: HexgridConfigSpec
      }
  )

const markerColorPalette = [
  '#e41a1c',
  '#377eb8',
  '#4daf4a',
  '#984ea3',
  '#ff7f00',
  '#ffff33',
  '#a65628',
  '#f781bf'
]

function createColorGenerator(palette: string[]) {
  let index = 0

  return () => {
    const color = palette[index]
    index = (index + 1) % palette.length
    return color
  }
}

const nextMarkerColor = createColorGenerator(markerColorPalette)

export function makeMarkerLayer(): MarkerLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    active: true,
    filters: {}
  }
  const markerColor = nextMarkerColor()

  return {
    ...baseLayer,
    type: 'markers',
    clustered: false,
    ready: false,
    config: {
      radius: 8,
      weight: 2,
      color: withOpacity(markerColor, 0.8),
      fillColor: withOpacity(markerColor, 0.3)
    }
  }
}

export function makeHexLayer(): HexLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    active: true,
    filters: {}
  }

  return {
    ...baseLayer,
    type: 'hexgrid',
    config: {
      radius: 10,
      opacity: 0.8,
      hover: {}
    },
    bindings: {
      color: { log: false, binding: 'occurrences' },
      opacity: { log: false },
      radius: { log: false }
    }
  }
}

export function resetLayerStyle(layer: LayerSpec) {
  const newLayer = makeMarkerLayer()
  layer.config = newLayer.config
  if (layer.type === 'hexgrid') {
    layer.bindings = {}
  } else if (layer.type === 'markers') {
    layer.clustered = false
  }
}

/**
 * Type representing the cosmetic settings for markers in the map layers.
 * Opacity can be controlled directly by 'color' and 'fill' properties.
 */
export type MarkerConfig = Overwrite<
  Omit<CircleMarkerOptions, 'opacity' | 'fillOpacity' | 'renderer'>,
  { dashArray?: string | undefined }
>

/**
 * Type representing a marker layer in the map.
 */
export type MarkerLayer<Item extends Geocoordinates> = {
  name?: string
  active: boolean
  config: MarkerConfig
  clustered: boolean
  maxClusterRadius?: number
  data?: Item[]
}

/**
 * Type representing the bindings for hexgrid scales.
 * These bindings are used to define which value is used for each hexagon
 * to determine its color, radius, and opacity.
 */
export type HexgridScaleBindings<Item> = {
  color?: ScaleBinding<Item>
  radius?: ScaleBinding<Item>
  opacity?: ScaleBinding<Item>
}

/**
 * Type representing the configuration for a hexgrid layer in the map.
 * This configuration is used to define the appearance and behavior of the hexgrid layer.
 */
export type HexgridConfig = {
  radius: number
  radiusRange?: [number, number]
  colorRange?: string[]
  hover: {
    fill?: boolean
    useScale?: boolean
    scale?: number
  }
  opacity: number
  opacityRange?: [number, number]
}

/**
 * Type representing the configuration for a hexgrid layer in the map.
 * Not intended to be used directly, but as a specification for the hexgrid layer
 * that can be converted to a HexgridConfig.
 */
export type HexgridConfigSpec = Overwrite<
  HexgridConfig,
  {
    colorRange: keyof typeof brewerPalettes
  }
>

/**
 * Type representing a hexgrid layer in the map.
 */
export type HexgridLayer<Item extends Geocoordinates> = {
  name?: string
  active: boolean
  config: HexgridConfig
  data?: Item[]
  bindings: HexgridScaleBindings<Item>
}
