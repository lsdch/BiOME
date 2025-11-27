import { Sampling, Site } from "@/api"
import { reactive, Reactive } from "vue"
import { BiomatModel } from "."
import { SamplingFormModel } from "./sampling"
import { SiteFormModel } from "./site"



export type OccurrenceModel = {
  site: SiteFormModel | Site | undefined,
  sampling: SamplingFormModel | Sampling | undefined,
  biomaterial: BiomatModel.BiomatModel | undefined
}

export function initialModel(): Reactive<OccurrenceModel> {
  return reactive({
    site: undefined,
    sampling: undefined,
    biomaterial: BiomatModel.initialModel(),
  })
}
