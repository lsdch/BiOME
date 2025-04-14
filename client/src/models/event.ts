import { DateWithPrecision, DateWithPrecisionInput, Event, EventInput, OrganisationInner, PersonUser } from "@/api"
import { Reactive, reactive } from "vue"

export type EventModel = {
  performed_on: DateWithPrecisionInput
  performed_by: PersonUser[]
  performed_by_groups: OrganisationInner[]
}

export function initialModel(): Reactive<EventModel> {
  return reactive({
    performed_on: { precision: 'Day', date: {} },
    performed_by: [],
    performed_by_groups: [],
  })
}

export function fromEvent({ performed_on, performed_by, performed_by_groups }: Event): EventModel {
  return {
    performed_on: DateWithPrecision.toInput(performed_on),
    performed_by: performed_by ?? [],
    performed_by_groups: performed_by_groups ?? [],
  }
}

export function toRequestData({ performed_on, performed_by, performed_by_groups }: EventModel): EventInput {
  return {
    performed_on: performed_on,
    performed_by: performed_by.map(({ alias }) => alias),
    performed_by_groups: performed_by_groups.map(({ code }) => code),
  }
}