import {
  $CompositeDate,
  $DateWithPrecisionInput,
  $EventDatePrecision,
  CompositeDate,
  DateWithPrecision,
  DateWithPrecisionInput,
  EventDatePrecision
} from '@/api'

export type DateWithPrecisionModel = {
  date: CompositeDate
  precision: EventDatePrecision
}

export const $schema = {
  ...$DateWithPrecisionInput,
  properties: {
    date: $CompositeDate,
    precision: {
      ...$EventDatePrecision,
      enum: [...$EventDatePrecision.enum, 'Unknown'] as const
    }
  }
}

export function initialModel(): DateWithPrecisionModel {
  return {
    date: {},
    precision: 'day'
  }
}

export function fromDateWithPrecision(d?: DateWithPrecision): DateWithPrecisionModel {
  return d
    ? {
        date: CompositeDate.fromDateWithPrecision(d),
        precision: d.precision
      }
    : initialModel()
}

export function toInput({ date, precision }: DateWithPrecisionModel): DateWithPrecisionInput {
  return { date, precision }
}
