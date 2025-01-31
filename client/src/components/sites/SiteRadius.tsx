import { CoordinatesPrecision } from '@/api'
import { LCircle } from '@vue-leaflet/vue-leaflet'
import { ComponentProps } from 'vue-component-type-helpers'
import { Geocoordinates } from '../maps'

type SiteRadiusProps = {
  site: Geocoordinates
} & Omit<ComponentProps<typeof LCircle>, 'latLng' | 'radius'>

export function SiteRadius({ site, interactive = false, ...props }: SiteRadiusProps) {
  return (
    <LCircle
      latLng={Geocoordinates.LatLng(site)}
      radius={CoordinatesPrecision.radius(site.coordinates.precision)}
      interactive={interactive}
      {...props}
    />
  )
}

export default SiteRadius
