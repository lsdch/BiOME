import { DatePrecision, DateWithPrecision } from "@/api";
import { DateTime } from "luxon";

const formats: Record<DatePrecision, string> = {
  Day: 'dd LLL yyyy',
  Month: 'LLL yyyy',
  Year: 'yyyy',
  Unknown: "'Unknown'"
}

export function formatDateWithPrecision({ date, precision }: DateWithPrecision, format?: string) {
  return DateTime.fromJSDate(date)
    .setLocale('en-gb')
    .toFormat(format ?? formats[precision])
}