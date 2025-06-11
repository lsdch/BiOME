import { SiteWithOccurrences } from "@/api"
import { HexgridConfig, HexgridScaleBindings, MarkerConfig } from "@/components/maps/SitesMap.vue"
import { UUID } from "crypto"

export type SitesFilter = 'All' | 'Sampled' | 'Occurrences'

export type MarkerLayerDefinition = {
  name?: string
  filterType: SitesFilter
  dataFeedID?: UUID
  active: boolean
  config: MarkerConfig
}

export type HexgridLayerDefinition = {
  name?: string
  active: boolean
  dataFeedID?: UUID
  filterType: SitesFilter
  bindings: HexgridScaleBindings<SiteWithOccurrences>
  config: HexgridConfig
}