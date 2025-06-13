import { Geocoordinates } from "@/components/maps"
import { ScaleBindingSpec } from "@/composables/occurrences"
import { brewerPalettes } from "@/functions/color_brewer"
import { UUID } from "crypto"
import { CircleMarkerOptions } from "leaflet"
import { Overwrite } from "ts-toolbelt/out/Object/Overwrite"
import { ScaleBinding } from "vue-leaflet-hexbin"

/**
 * Type representing the filter options for sites in the map layers.
 * - 'All': Show all sites.
 * - 'Sampled': Show only sites that have been sampled.
 * - 'Occurrences': Show only sites that have occurrences.
 */
export type SitesFilter = 'All' | 'Sampled' | 'Occurrences'


/**
 * Type representing the cosmetic settings for markers in the map layers.
 * Opacity can be controlled directly by 'color' and 'fill' properties.
 */
export type MarkerConfig = Overwrite<
  Omit<CircleMarkerOptions, 'opacity' | 'fillOpacity' | 'renderer'>,
  { dashArray?: string | undefined }
>

/**
 * Type representing the definition of a marker layer in the map,
 * to be used with data feeds.
 */
export type MarkerLayerDefinition = {
  name?: string
  filterType: SitesFilter
  dataFeedID?: UUID
  active: boolean
  clustered: boolean
  config: MarkerConfig
}

/**
 * Type representing a marker layer in the map.
 */
export type MarkerLayer<Item extends Geocoordinates> = {
  name?: string
  active: boolean
  config: MarkerConfig
  clustered: boolean
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
export type HexgridConfigSpec = Overwrite<HexgridConfig, {
  colorRange: keyof typeof brewerPalettes
}>

/**
 * Type representing the definition of a hexgrid layer in the map,
 * to be used with data feeds.
 */
export type HexgridLayerSpec = {
  name?: string
  active: boolean
  dataFeedID?: UUID
  filterType: SitesFilter
  bindings: {
    color?: ScaleBindingSpec
    radius?: ScaleBindingSpec
    opacity?: ScaleBindingSpec
  }
  config: HexgridConfigSpec
}


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