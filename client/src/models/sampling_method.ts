import { Organisation, OrganisationInner, OrganisationInput, OrganisationUpdate, SamplingMethod, SamplingMethodInput, SamplingMethodUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type SamplingMethodFormModel = SamplingMethodInput | SamplingMethodUpdate

export function initialModel(): Reactive<SamplingMethodInput> {
  return reactive({
    code: '',
    label: ''
  })
}

export function fromSamplingMethod({ code, label, description }: SamplingMethod): SamplingMethodUpdate {
  return { code, label, description }
}