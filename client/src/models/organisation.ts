import { Organisation, OrganisationInner, OrganisationInput, OrganisationUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type OrganisationFormModel = OrganisationInput | OrganisationUpdate

export function initialModel(): Reactive<OrganisationInput> {
  return reactive({
    name: '',
    code: '',
    kind: 'Lab'
  })
}

export function fromOrganisation({ name, code, kind, description }: Organisation): OrganisationUpdate {
  return {
    name,
    code,
    kind,
    description
  }
}