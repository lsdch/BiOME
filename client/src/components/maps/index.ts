import { CoordinatesPrecision } from "@/api"
import { LatLngExpression } from "leaflet"

export interface Coordinates {
  latitude: number
  longitude: number
}

export interface MaybeCoordinates extends Partial<Coordinates> { }

export namespace Coordinates {
  export function isValidCoordinates(coords: MaybeCoordinates | undefined): coords is Coordinates {
    return !!coords &&
      (
        coords.latitude !== undefined &&
        coords.latitude !== null
      ) && (
        coords.longitude !== undefined &&
        coords.longitude !== null
      )
  }
}

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