import { OccurrenceCategory, Sampling, Site } from "@/api"
import { reactive, Reactive } from "vue"
import { BiomatModel } from "."
import { SamplingFormModel } from "./sampling"
import { SiteFormModel } from "./site"

export type OccurrenceModel = {
  site: SiteFormModel | Site | undefined,
  sampling: SamplingFormModel | Sampling | undefined,
  biomaterial: {
    category: OccurrenceCategory | undefined,
    external: BiomatModel.ExternalOccurrenceModel | undefined
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
    sampling: undefined,
    biomaterial: { category: undefined, external: BiomatModel.initialModel() },
  })
}

// export function fromBioMaterial({  }: BioMaterialWithDetails)