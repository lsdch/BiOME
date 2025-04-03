import { CoordinatesPrecision, SiteInput, SiteUpdate, Site } from "@/api"
import { Reactive, reactive } from "vue"

export type SiteFormModel = Omit<SiteInput, 'coordinates'> & {
  coordinates: {
    latitude?: number
    longitude?: number
    precision: CoordinatesPrecision
  }
}

export function initialModel(): Reactive<SiteFormModel> {
  return reactive({
    name: '',
    code: '',
    coordinates: {
      precision: '<100m'
    },
    user_defined_locality: false
  })
}

export function toRequestBody({ coordinates: { latitude, longitude, precision }, ...model }: SiteFormModel): SiteInput {
  return {
    ...model,
    coordinates: {
      latitude: latitude!,
      longitude: longitude!,
      precision
    }
  } satisfies SiteUpdate
}

export function fromSite({
  id,
  meta,
  $schema,
  events,
  datasets,
  country,
  coordinates,
  ...rest
}: Site): SiteFormModel {
  return {
    ...rest,
    country_code: country?.code,
    coordinates: {
      latitude: coordinates.latitude,
      longitude: coordinates.longitude,
      precision: coordinates.precision
    }
  }
}
