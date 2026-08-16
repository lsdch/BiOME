import {
  CompositeDate,
  EventDatePrecision,
  H3CellWithRichness,
  ListOccurrencesData,
  ListSamplingsH3Data
} from '@/api'
import { ScaleBindingSpec } from '@/features/cartography/bindings'
import { ItemWithCoordinates } from '@/features/cartography/coordinates'
import { brewerPalettes, withOpacity } from '@/lib/color_brewer'
import { UUID } from 'crypto'
import { MarkerOptions } from 'maplibre-gl'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'
import { RGB } from 'vuetify/lib/util/colorUtils.mjs'
import { CellMarkerData } from './layer-data'
import { GlobalMarkerOptions } from '../DeckGlMap.vue'
import { WithRequired } from '@tanstack/vue-query'
import { DateTime } from 'luxon'

export interface H3Cell {
  h3_index: string
}

export type DateFilters = Overwrite<
  NonNullable<NonNullable<ListOccurrencesData['query']>['date']>,
  {
    from?: CompositeDate
    to?: CompositeDate
  }
> & {
  enabled: boolean
  is_range: boolean
  precision: EventDatePrecision
}

export type MappingFilters = Omit<NonNullable<ListOccurrencesData['query']>, 'date'> & {
  date: DateFilters
}

function compositeDateToString(
  date: CompositeDate | undefined,
  precision: EventDatePrecision
): string | undefined {
  if (!date) return undefined
  switch (precision) {
    case 'day':
      if (date.day === undefined || date.month === undefined || date.year === undefined) {
        return undefined
      }
      return DateTime.fromObject(date).toFormat('yyyy-MM-dd')
    case 'month':
      if (date.month === undefined || date.year === undefined) {
        return undefined
      }
      return DateTime.fromObject(date).toFormat('yyyy-MM')
    case 'year':
      if (date.year === undefined) {
        return undefined
      }
      return DateTime.fromObject(date).toFormat('yyyy')
  }
}

export function mappingFiltersToQuery(
  filters: MappingFilters,
  mode: 'occurrences'
): NonNullable<ListOccurrencesData['query']>

export function mappingFiltersToQuery(
  filters: MappingFilters,
  mode: 'samplings'
): NonNullable<ListSamplingsH3Data['query']>

export function mappingFiltersToQuery(
  { date, ...filters }: MappingFilters,
  mode: 'occurrences' | 'samplings'
): NonNullable<ListOccurrencesData['query'] | ListSamplingsH3Data['query']> {
  let query: NonNullable<ListOccurrencesData['query'] | ListSamplingsH3Data['query']> = {}

  const fromDate = compositeDateToString(date?.from, date.precision)
  const toDate = compositeDateToString(date?.to, date.precision)

  switch (mode) {
    case 'occurrences':
      query = {
        ...filters,
        date: date.enabled
          ? {
              from: fromDate,
              to: date.is_range ? toDate : fromDate,
              include_unknown: date.include_unknown
            }
          : undefined
      }
      break
    case 'samplings':
      const { batches, countries, limit, offset, target_taxa, target_taxa_whole_clade } =
        filters as NonNullable<ListSamplingsH3Data['query']>
      query = {
        batches,
        countries,
        date: date.enabled
          ? {
              from: fromDate,
              to: date.is_range ? toDate : fromDate,
              include_unknown: date.include_unknown
            }
          : undefined,
        limit,
        offset,
        target_taxa,
        target_taxa_whole_clade
      }
      break
  }

  if (!date.enabled || !query.date) {
    query.date = undefined
  } else if (!date.is_range) {
    query.date = { ...query.date, to: query.date.from }
  }

  console.log('Query: ', query)
  return query
}

// WithRequired<
// Overwrite<NonNullable<ListOccurrencesData['query']>, { habitats?: HabitatRecord[] }>,
// 'sampling_target'
// >
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
  include_sites?: SitesFilter
}

export type HexLayerSpec = BaseLayerSpec & {
  type: 'hexgrid'
  mode: 'occurrences' | 'samplings'
  resolutionMode: 'auto' | 'manual'
  resolution: number
  colorBinding: Required<ScaleBindingSpec>

  config: HexgridConfigSpec
  markers: MarkerLayerSpec & { minZoom: number; minZoomMode: 'auto' | 'manual' | 'never' }
}

export type MarkerLayerSpec = BaseLayerSpec & {
  type: 'markers'
  ready?: boolean
  clustered: boolean
  resolution: number
  resolutionMode: 'auto' | 'manual'
  config: MarkerConfig
}

export type MarkerLayerGetters<Item> = {
  getPosition: (item: Item) => [number, number]
  radius?: (item: Item) => number
  getText?: (item: Item) => string
}

export type LayerSpec = HexLayerSpec | MarkerLayerSpec

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

export type MarkerLayerParams = Partial<
  Pick<MarkerLayerSpec, 'include_sites' | 'filters' | 'ready' | 'config'>
>
export function makeMarkerLayer(name?: string, params?: MarkerLayerParams): MarkerLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    name,
    active: true,
    include_sites: 'Occurrences',
    filters: {
      ...{
        sampling_target: {},
        include_sites: 'Occurrences',
        date: { precision: 'year', enabled: false, is_range: false }
      },
      ...params?.filters
    }
  }
  const markerColor = nextMarkerColor()

  return {
    ...baseLayer,
    type: 'markers',
    clustered: false,
    ready: true,
    resolutionMode: 'auto',
    resolution: 12,
    config: {
      ...{
        baseRadius: 5,
        radiusScaleFactor: 1,
        weight: 2,
        color: withOpacity(markerColor, 0.8),
        fillColor: withOpacity(markerColor, 0.3)
      },
      ...params?.config
    }
  }
}

export function markerLayerFromSpec<Item>(
  spec: MarkerLayerSpec,
  data: Item[],
  getters: MarkerLayerGetters<Item>
): MarkerLayer<Item> {
  return {
    ...spec,
    ...getters,
    data
  }
}

export function makeHexLayer(markersParams?: MarkerLayerParams): HexLayerSpec {
  const baseLayer: BaseLayerSpec = {
    id: crypto.randomUUID(),
    include_sites: 'Occurrences',
    active: true,
    filters: {
      date: { precision: 'year', enabled: false, is_range: false }
      // sampling_target: {},
      // include_sites: 'Occurrences'
    }
  }

  console.log('Creating new hex layer')

  return {
    ...baseLayer,
    type: 'hexgrid' as const,
    mode: 'occurrences' as const,
    resolutionMode: 'auto' as const,
    resolution: 8,
    config: {
      radius: 50,
      opacity: 0.8,
      strokeWidth: 1,
      strokeOpacity: 1,
      hover: {
        showTooltip: true,
        highlight: true
      },
      coverage: 0.95,
      fillColor: '#FF7F00'
    },
    colorBinding: { log: true, binding: 'occurrences' as const },
    markers: {
      ...makeMarkerLayer('hexgrid-markers', markersParams),
      minZoom: 8,
      minZoomMode: 'auto'
    }
  }
}

export function resetLayerStyle(layer: LayerSpec, globalOpts?: GlobalMarkerOptions): void {
  if (layer.type === 'hexgrid') {
    const newLayer = makeHexLayer()
    layer.config = newLayer.config
    layer.colorBinding = newLayer.colorBinding
  } else if (layer.type === 'markers') {
    const newLayer = makeMarkerLayer(undefined, {
      config: { radiusScaleFactor: globalOpts?.cluster.radiusScaleFactor ?? 1 }
    })
    layer.config = newLayer.config
  }
}

/**
 * Type representing the cosmetic settings for markers in the map layers.
 * Opacity can be controlled directly by 'color' and 'fill' properties.
 */
export type MarkerConfig = {
  baseRadius?: number
  radiusScaleFactor?: number
  fillColor?: string
  color?: string
  weight?: number
}

/**
 * Type representing a marker layer in the map.
 */
export type MarkerLayer<Item> = MarkerLayerSpec &
  MarkerLayerGetters<Item> & {
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
  fillColor: string
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
export type HexgridLayer<CellData extends H3Cell> = HexLayerSpec & {
  data?: CellData[]
}

export type PinMarker<Item> = ItemWithCoordinates & {
  options?: MarkerOptions
  data: Item
}

export function automaticResolution(layer: LayerSpec, zoom: number): number {
  switch (layer.type) {
    case 'markers':
      // if (zoom < 3) return 3
      if (layer.resolutionMode === 'auto') {
        if (zoom < 4) return 4
        if (zoom < 6) return 5
        if (zoom < 8) return 8
        if (zoom < 10) return 10
        return 12
      }
    case 'hexgrid':
      if (zoom <= 2) return 2
      if (zoom <= 3) return 3
      if (zoom <= 5) return 4
      if (zoom <= 7) return 5
      if (zoom <= 9) return 6
      if (zoom <= 11) return 7
      if (zoom <= 13) return 8
      if (zoom <= 15) return 9
      if (zoom <= 17) return 10
      if (zoom <= 19) return 11
      return 12
  }
}
