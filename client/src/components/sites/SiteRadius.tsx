import { CoordinatesPrecision } from '@/api'
import { LCircle } from '@vue-leaflet/vue-leaflet'
import { ComponentProps } from 'vue-component-type-helpers'
import { Geocoordinates } from '../maps'

type SiteRadiusProps = {
  site: Geocoordinates
  zoom?: number
} & Omit<ComponentProps<typeof LCircle>, 'latLng' | 'radius'>

function zoomedEnough(item: Geocoordinates, zoom?: number) {
  switch (item.coordinates.precision) {
    case undefined:
    case 'Unknown':
      return false
    case '<100m':
    case '<1km':
      return !zoom || zoom > 10
    default:
      return !zoom || zoom > 6
  }
}

export function SiteRadius({ site, interactive = false, zoom, ...props }: SiteRadiusProps) {
  return zoomedEnough(site, zoom) ? (
    <LCircle
      color="orange"
      dashArray="1px 8px"
      latLng={Geocoordinates.LatLng(site)}
      radius={CoordinatesPrecision.radius(site.coordinates.precision)}
      interactive={interactive}
      {...props}
    />
  ) : undefined
}

export default SiteRadius
