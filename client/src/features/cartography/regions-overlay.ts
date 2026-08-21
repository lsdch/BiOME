import { Map, StyleSpecification } from 'maplibre-gl'
import regionsLayerStyle from './data/regions-layer.json'
import { MaybeRef, unref } from 'vue'
const regionLayerIds = [
  'boundary_2',
  'boundary_disputed',
  'waterway_line_label',
  'water_name_point_label',
  'water_name_line_label',
  'label_city',
  'label_city_capital',
  'label_town',
  'label_village',
  'label_country_1',
  'label_country_2',
  'label_country_3',
  'label_state'
] as const

const roadsLayerIds = [
  'highway-motorway',
  'highway-trunk',
  'highway-primary',
  'highway-secondary-tertiary',
  'highway-minor',
  'highway-path',
  'tunnel-motorway',
  'tunnel-trunk-primary',
  'tunnel-secondary-tertiary',
  'tunnel-minor',
  'tunnel-path',
  'bridge-motorway',
  'bridge-trunk-primary',
  'bridge-secondary-tertiary',
  'bridge-minor',
  'bridge-path',
  'highway-name-major',
  'highway-name-minor',
  'highway-name-path'
] as const

export function createRegionLayers(visible: boolean): StyleSpecification['layers'] {
  return (regionsLayerStyle.layers ?? [])
    .filter((layer) => regionLayerIds.includes(layer.id as (typeof regionLayerIds)[number]))
    .map((layer) => ({
      ...layer,
      layout: {
        ...(layer.layout ?? {}),
        visibility: visible ? 'visible' : 'none'
      }
    })) as StyleSpecification['layers']
}

export function createRoadsLayers(visible: boolean): StyleSpecification['layers'] {
  return (regionsLayerStyle.layers ?? [])
    .filter((layer) => roadsLayerIds.includes(layer.id as (typeof roadsLayerIds)[number]))
    .map((layer) => ({
      ...layer,
      layout: {
        ...(layer.layout ?? {}),
        visibility: visible ? 'visible' : 'none'
      },
      paint: {
        ...(layer.paint ?? {})
      }
    })) as StyleSpecification['layers']
}

export { regionsLayerStyle }

export function updateRegionsLayerVisibility(map: MaybeRef<Map | undefined>, value: boolean) {
  const instance = unref(map)
  if (!instance) return
  const visibility = value ? 'visible' : 'none'

  regionLayerIds.forEach((layerId) => {
    if (instance.getLayer(layerId)) {
      instance.setLayoutProperty(layerId, 'visibility', visibility)
    }
  })
}

export function updateRoadsLayerVisibility(map: MaybeRef<Map | undefined>, value: boolean) {
  const instance = unref(map)
  if (!instance) return
  const visibility = value ? 'visible' : 'none'

  roadsLayerIds.forEach((layerId) => {
    if (instance.getLayer(layerId)) {
      instance.setLayoutProperty(layerId, 'visibility', visibility)
    }
  })
}
