import { DateWithPrecision, DateWithPrecisionInput, Event, EventInput, PersonUser, ProgramInner } from "@/api"
import { Reactive, reactive } from "vue"

export type EventModel = {
  performed_on: DateWithPrecisionInput
  performed_by: PersonUser[]
  programs: ProgramInner[]
}

export function initialModel(): Reactive<EventModel> {
  return reactive({
    performed_on: { precision: 'Day', date: {} },
    performed_by: [],
    programs: []
  })
}

export function fromEvent({ performed_on, performed_by, programs }: Event): EventModel {
  return {
    performed_on: DateWithPrecision.toInput(performed_on),
    performed_by: performed_by,
    programs: programs ?? []
  }
}

export function toRequestData({ performed_on, performed_by, programs }: EventModel): EventInput {
  return {
    performed_on: performed_on,
    performed_by: performed_by.map(({ alias }) => alias),
    programs: programs.map(({ code }) => code),
  }
}