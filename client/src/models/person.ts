import { OrganisationInner, Person, PersonInput, PersonUpdate } from "@/api";
import { reactive, Reactive } from "vue";

export type PersonFormModel = Omit<PersonInput, "organisations"> & {
  organisations: OrganisationInner[]
}

export function initialModel(): Reactive<PersonFormModel> {
  return reactive({
    first_name: "",
    last_name: "",
    organisations: []
  })
}

export function fromPerson({
  alias,
  comment,
  contact,
  first_name,
  last_name,
  organisations
}: Person): PersonFormModel {
  return {
    first_name,
    last_name,
    alias,
    contact,
    comment,
    organisations
  }
}

export function toRequestBody({ organisations, ...model }: PersonFormModel): PersonInput {
  return {
    ...model,
    organisations: organisations.map(({ code }) => code)
  } satisfies PersonUpdate
}
