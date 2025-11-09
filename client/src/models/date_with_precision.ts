import { $CompositeDate, $DatePrecision, $DateWithPrecisionInput, CompositeDate, DatePrecision, DateWithPrecision, DateWithPrecisionInput } from "@/api";
import { reactive, Reactive } from "vue";

export type DatePrecisionModel = DatePrecision | 'Unknown'
export type DateWithPrecisionModel = {
  date: CompositeDate,
  precision: DatePrecisionModel
}

export const $schema = {
  ...$DateWithPrecisionInput,
  properties: {
    date: $CompositeDate,
    precision: {
      ...$DatePrecision,
      enum: [...$DatePrecision.enum, 'Unknown'] as const
    }
  }
}

export function initialModel(): Reactive<DateWithPrecisionModel> {
  return reactive({
    date: {},
    precision: 'Unknown'
  })
}

export function fromDateWithPrecision(d?: DateWithPrecision): DateWithPrecisionModel {
  return d ? {
    date: CompositeDate.fromDateWithPrecision(d),
    precision: d.precision
  } : initialModel()
}

export function toInput({ date, precision }: DateWithPrecisionModel): DateWithPrecisionInput | undefined {
  if (precision === 'Unknown') {
    return undefined
  }
  return { date, precision }
}