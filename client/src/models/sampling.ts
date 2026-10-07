import {
  CoordinatesWithPrecision,
  DateWithPrecisionInput,
  Sampling,
  SamplingInput,
  SiteInput
} from '@/api'
import {
  DateWithPrecisionModel,
  toInput as DateWithPrecisionToInput,
  fromDateWithPrecision
} from '@/models/date_with_precision'

export type SamplingCoreModel = {
  coordinates: Partial<CoordinatesWithPrecision>
  altitude?: number
  performed_on?: DateWithPrecisionModel
}

export type SamplingMetadataModel = {
  access_points?: string[] | undefined
  comments?: string | undefined
  duration?: number
  fixatives?: Array<string>
  methods?: Array<string>
  performed_by?: Array<string>
  site: SiteInput
  target_taxa?: Array<string>
}

export type SamplingModel = SamplingCoreModel & SamplingMetadataModel

export function initialModel(): SamplingModel {
  return {
    site: {},
    coordinates: {},
    performed_on: { precision: 'day', date: {} }
  }
}

export function fromSampling({
  duration,
  access_points,
  comments,
  performed_on,
  performed_by,
  coordinates,
  site
}: Sampling): SamplingModel {
  return {
    site: site,
    duration,
    access_points,
    performed_on: fromDateWithPrecision(performed_on),
    coordinates,
    performed_by: performed_by ?? []
  }
}

// export function toRequestBody({
//   coordinates,
//   site,
//   duration,
//   access_points,
//   comments,
//   target_taxa,
//   fixatives,
//   methods,
//   performed_on,
//   performed_by
// }: SamplingModel): SamplingInput {
//   if ()
//   return {
//     coordinates,
//     site,
//     duration,
//     access_points,
//     comments,
//     performed_on: DateWithPrecisionToInput(performed_on),
//     performed_by: performed_by
//   } // satisfies SamplingUpdate
// }
