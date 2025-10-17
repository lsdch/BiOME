import { MapToolPreset } from "@/api"
import { HexgridLayerSpec, MarkerLayerDefinition } from "@/features/cartography/components/map-layers"
import { DataFeed } from "@/features/cartography/components/data-feeds"
import { Overwrite } from "ts-toolbelt/out/Object/Overwrite"

export type MapPresetSpec = {
  feeds: DataFeed[]
  hexgrid: HexgridLayerSpec
  markers: MarkerLayerDefinition[]
}

export type ParsedMapPreset = Overwrite<MapToolPreset, {
  spec: MapPresetSpec
}>

export function parseMapPreset(preset: MapToolPreset): ParsedMapPreset {
  return { ...preset, spec: JSON.parse(preset.spec) }
}
