import { BioMaterialWithDetails, Event, ExternalBioMatInput, OccurrenceCategory, Sampling, SamplingInput, Site } from "@/api"
import { SamplingFormModel } from "./sampling"
import { SiteFormModel } from "./site"
import { EventModel } from "./event"
import { BiomatModel } from "."
import { reactive, Reactive } from "vue"

export type OccurrenceModel = {
  site: SiteFormModel | Site | undefined,
  event: EventModel | Event | undefined
  sampling: SamplingFormModel | Sampling | undefined,
  biomaterial: {
    category: OccurrenceCategory | undefined,
    external: BiomatModel.ExternalBiomatModel | undefined
    // internal: BiomatModel.InternalBiomatModel | undefined
  }
  // | {
  //   category: OccurrenceCategory & "External",
  //   item: BiomatModel.ExternalBiomatModel
  // }
  // | {
  //   category: OccurrenceCategory & "Internal"
  //   item: never
  // }
}

export function initialModel(): Reactive<OccurrenceModel> {
  return reactive({
    site: undefined,
    event: undefined,
    sampling: undefined,
    biomaterial: { category: undefined, external: BiomatModel.initialModel() },
  })
}

// export function fromBioMaterial({  }: BioMaterialWithDetails)