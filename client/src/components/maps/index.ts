import { CoordinatesPrecision } from "@/api"
import { LatLngExpression } from "leaflet"

export interface Geocoordinates {
  coordinates: {
    latitude: number
    longitude: number
    precision?: CoordinatesPrecision
  }
}

export namespace Geocoordinates {
  export function LatLng({ coordinates: { latitude, longitude } }: Geocoordinates): LatLngExpression {
    return [latitude, longitude]
  }
}