import { HabitatRecord, OccurrencesBySiteData } from '@/api'
import { ScaleBindingSpec } from '@/features/cartography/bindings'
import { Geocoordinates } from '@/features/cartography/coordinates'
import { brewerPalettes, withOpacity } from '@/lib/color_brewer'
import { WithRequired } from '@tanstack/vue-query'
import { UUID } from 'crypto'
import { MarkerOptions } from 'maplibre-gl'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { RGB } from 'vuetify/lib/util/colorUtils.mjs'

export type MappingFilters = WithRequired<
  Overwrite<NonNullable<OccurrencesBySiteData['body']>, { habitats?: HabitatRecord[] }>,
  'sampling_target'
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
  include_sites: SitesFilter
}

export type HexLayerSpec = BaseLayerSpec & {
  type: 'hexgrid'
  colorBinding: Required<ScaleBindingSpec>
  config: HexgridConfigSpec
  markers: MarkerLayerSpec & { minZoom: number }
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

export const markerColorPalette = [
  '#e41a1c',
  '#4daf4a',
  '#984ea3',
  '#ff7f00',
  '#377eb8',
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

export function makeMarkerLayer(name?: string): MarkerLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    name,
    active: true,
    include_sites: 'Occurrences',
    filters: {
      sampling_target: {},
      include_sites: 'Occurrences'
    }
  }
  const markerColor = nextMarkerColor()

  return {
    ...baseLayer,
    type: 'markers',
    clustered: false,
    ready: false,
    config: {
      radius: 5,
      weight: 2,
      color: withOpacity(markerColor, 0.8),
      fillColor: withOpacity(markerColor, 0.3)
    }
  }
}

export function markerLayerFromSpec<Item extends Geocoordinates>(
  spec: MarkerLayerSpec,
  data?: Item[]
): MarkerLayer<Item> {
  return {
    name: spec.name,
    active: spec.active,
    clustered: spec.clustered,
    config: spec.config,
    data
  }
}

export function makeHexLayer(): HexLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    include_sites: 'Occurrences',
    active: true,
    filters: {
      sampling_target: {},
      include_sites: 'Occurrences'
    }
  }

  return {
    ...baseLayer,
    type: 'hexgrid',
    config: {
      radius: 50,
      opacity: 0.8,
      strokeWidth: 1,
      strokeOpacity: 1,
      hover: {
        showTooltip: true,
        highlight: true
      },
      coverage: 0.95
    },
    colorBinding: { log: false, binding: 'occurrences' },
    markers: { ...makeMarkerLayer('hexgrid-markers'), minZoom: 8 }
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
export type MarkerConfig = {
  radius?: number
  fillColor?: string
  color?: string
  weight?: number
}

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
 * Type representing the configuration for a hexgrid layer in the map.
 * This configuration is used to define the appearance and behavior of the hexgrid layer.
 */
export type HexgridConfig = {
  radius: number
  radiusRange?: [number, number]
  colorRange?: RGB[]
  hover: {
    showTooltip?: boolean
    highlight?: boolean
  }
  opacity: number
  opacityRange?: [number, number]
  strokeWidth?: number
  strokeOpacity?: number
  coverage?: number
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
  colorBinding: Required<ScaleBindingSpec>
  markers: MarkerLayerSpec & { minZoom: number }
}

export type PinMarker<Item> = Geocoordinates & {
  options?: MarkerOptions
  data: Item
}
