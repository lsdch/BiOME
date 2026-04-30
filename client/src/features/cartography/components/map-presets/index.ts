import { MapToolPreset } from '@/api'
import {
  HexLayerSpec,
  MarkerLayerSpec
} from '@/features/cartography/components/layers-manager/map-layers'
import { Overwrite } from 'ts-toolbelt/out/Object/Overwrite'

export type MapPresetSpec = {
  hexgrid: HexLayerSpec
  markers: MarkerLayerSpec[]
}

export type ParsedMapPreset = Overwrite<
  MapToolPreset,
  {
    spec: MapPresetSpec
  }
>

export function parseMapPreset(preset: MapToolPreset): ParsedMapPreset {
  return { ...preset, spec: JSON.parse(preset.spec) }
}
