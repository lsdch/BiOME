import { CoordinatesPrecision } from '@/api'

export interface Coordinates {
  latitude: number
  longitude: number
}

export interface MaybeCoordinates extends Partial<Coordinates> {}

export namespace Coordinates {
  export function isValidCoordinates(coords: MaybeCoordinates | undefined): coords is Coordinates {
    return (
      !!coords &&
      coords.latitude !== undefined &&
      coords.latitude !== null &&
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
